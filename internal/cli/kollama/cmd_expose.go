package kollama

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

const (
	exposedMessage = `🎉 The model has been exposed through a service over %s.

To start a chat with ollama:

  OLLAMA_HOST=%s ollama run %s

To integrate with your OpenAI API compatible client:

  curl http://%s/v1/chat/completions -H "Content-Type: application/json" -d '{
    "model": "%s",
    "messages": [
      {
        "role": "user",
        "content": "Hello!"
      }
    ]
  }'.
`
)

// CmdDeployOptions provides information required to expose a model
type CmdExposeOptions struct {
	configFlags     *genericclioptions.ConfigFlags
	clientConfig    clientcmd.ClientConfig
	kubeConfig      *rest.Config
	kubeClient      client.Client
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface

	serviceType string
	serviceName string
	nodePort    int32

	genericiooptions.IOStreams
}

// NewCmdExposeOptions provides an instance of CmdExposeOptions with default values
func NewCmdExposeOptions(streams genericiooptions.IOStreams) *CmdExposeOptions {
	return &CmdExposeOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

// NewCmdExpose provides a cobra command wrapping NamespaceOptions
func NewCmdExpose(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewCmdExposeOptions(streams)

	cmd := &cobra.Command{
		Use:     "expose [model name] [flags]",
		Short:   "Expose a model with the given name as NodePort, or LoadBalancer service for external access",
		Example: fmt.Sprintf(deployExample, command(), command(), command(), command()),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("model name is required")
			}
			if args[0] == "" {
				return fmt.Errorf("model name cannot be empty")
			}

			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			return o.runE(c, args)
		},
	}

	cmd.Flags().StringVar(&o.serviceType, "service-type", "", ""+
		"Type of the Service to expose the model. If not specified, the service will be "+
		"exposed as NodePort. Use LoadBalancer to expose the service as LoadBalancer.",
	)
	cmd.Flags().StringVar(&o.serviceName, "service-name", "", ""+
		"Name of the Service to expose the model. If not specified, the model name will "+
		"be used as the service name with -nodeport as the suffix for NodePort.",
	)
	cmd.Flags().Int32Var(&o.nodePort, "node-port", 0, ""+
		"NodePort to expose the model. If not specified, a random port will be assigned."+
		"Only valid when --expose is specified, and --service-type is set to NodePort.",
	)

	o.configFlags.AddFlags(cmd.Flags())
	o.clientConfig = o.configFlags.ToRawKubeConfigLoader()
	o.kubeConfig = lo.Must(o.clientConfig.ClientConfig())
	o.kubeClient = lo.Must(client.New(o.kubeConfig, client.Options{}))
	o.dynamicClient = lo.Must(dynamic.NewForConfig(o.kubeConfig))
	o.discoveryClient = lo.Must(discovery.NewDiscoveryClientForConfig(o.kubeConfig))

	return cmd
}

func (o *CmdExposeOptions) runE(cmd *cobra.Command, args []string) error {
	var err error

	namespace, err := getNamespace(o.clientConfig, cmd)
	if err != nil {
		return err
	}

	supported, err := IsOllamaOperatorCRDSupported(o.discoveryClient, modelSchemaResourceName)
	if err != nil {
		return err
	}
	if !supported {
		return ErrOllamaModelNotSupported
	}

	modelName := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	model, err := getOllama(ctx, o.dynamicClient, namespace, modelName)
	if err != nil {
		return err
	}
	if model == nil {
		fmt.Println("Ollama Model", modelName, "not found, did you deploy it?")
		os.Exit(1)

		return nil
	}

	svc, err := exposeOllamaModel(
		ctx,
		o.kubeClient,
		namespace,
		modelName,
		lo.Ternary(o.serviceType == "", corev1.ServiceTypeNodePort, corev1.ServiceType(o.serviceType)),
		o.serviceName,
		o.nodePort,
	)
	if err != nil {
		return err
	}

	parsedHost, err := url.Parse(o.kubeConfig.Host)
	if err != nil {
		return err
	}

	ollamaHost := fmt.Sprintf("%s:%d", parsedHost.Hostname(), svc.Spec.Ports[0].NodePort)
	fmt.Printf(exposedMessage, ollamaHost, ollamaHost, modelName, ollamaHost, modelName)

	return nil
}
