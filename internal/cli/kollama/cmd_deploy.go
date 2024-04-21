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
	deployExample = `
  # Deploy a phi model
  $ %s deploy phi

  # Deploy a model with exposed through NodePort service
  $ %s deploy phi --expose

  # Deploy a model with a specific image
  $ %s deploy phi --image=phi-image

  # Deploy a phi model in a specific namespace
  $ %s deploy phi -n phi-namespace

`

	deployedAlreadyMessage = `%s has been deployed already.

To undeploy it, use

  %s undeploy %s

`

	deployedNonExposedMessage = `🎉 Successfully deployed %s.

💡 Currently the deployed model has not yet exposed. If this is unintentional, you can expose the model through

  %s expose %s

Or create with a exposed port with

  %s deploy %s --expose

next time.

To expose manually, use the following command:

  kubectl expose deployment %s --name=%s-nodeport --type=NodePort --port 11434

`

	deployedExposedMessage = `🎉 Successfully deployed %s.

🌐 The model has been exposed through a service over %s.

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
  }'

	`
)

// CmdDeployOptions provides information required to deploy a model
type CmdDeployOptions struct {
	configFlags     *genericclioptions.ConfigFlags
	clientConfig    clientcmd.ClientConfig
	kubeConfig      *rest.Config
	kubeClient      client.Client
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface

	modelImage   string
	expose       bool
	serviceType  string
	serviceName  string
	nodePort     int32
	storageClass string
	pvAccessMode string

	genericiooptions.IOStreams
}

// NewCmdDeployOptions provides an instance of CmdDeployOptions with default values
func NewCmdDeployOptions(streams genericiooptions.IOStreams) *CmdDeployOptions {
	return &CmdDeployOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

// NewCmdNamespace provides a cobra command wrapping CmdDeployOptions
func NewCmdDeploy(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewCmdDeployOptions(streams)

	cmd := &cobra.Command{
		Use:     "deploy [model name] [flags]",
		Short:   "Deploy a model with the given name by using Ollama Operator",
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

	cmd.Flags().StringVar(&o.modelImage, "image", "", ""+
		"Model image to deploy. If not specified, the model name will be used as the "+
		"image name (will be pulled from registry.ollama.ai/library/<model name> by "+
		"default if no registry is specified), the tag will be latest.",
	)

	cmd.Flags().StringVarP(&o.storageClass, "storage-class", "", "", ""+
		"StorageClass to use for the model's associated PersistentVolumeClaim. If not specified, "+
		"the default StorageClass will be used.",
	)

	cmd.Flags().StringVarP(&o.pvAccessMode, "pv-access-mode", "", "", ""+
		"Access mode for Ollama Operator created image store (to cache pulled images)'s StatefulSet "+
		"resource associated PersistentVolume. If not specified, the access mode will be ReadWriteOnce. "+
		"If you are deploying models into default deployed kind and k3s clusters, you should keep "+
		"it as ReadWriteOnce. If you are deploying models into a custom cluster, you can set it to "+
		"ReadWriteMany if StorageClass supports it.",
	)

	cmd.Flags().BoolVar(&o.expose, "expose", false, ""+
		"Whether to expose the model through a service for external access and makes it "+
		"easy to interact with the model. By default, --expose will create a NodePort "+
		"service. Use --service-type=LoadBalancer to create a LoadBalancer service",
	)

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

func (o *CmdDeployOptions) runE(cmd *cobra.Command, args []string) error {
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

	modelImage, err := getImage(cmd, args)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	model, err := getOllama(ctx, o.dynamicClient, namespace, modelName)
	if err != nil {
		return err
	}
	if model != nil {
		fmt.Printf(deployedAlreadyMessage, modelName, command(), modelName)
		os.Exit(1)
	}

	createdModel, err := createOllamaModel(ctx, o.dynamicClient, namespace, modelName, modelImage, o.storageClass, o.pvAccessMode)
	if err != nil {
		return err
	}
	if !o.expose {
		fmt.Printf(deployedNonExposedMessage, modelName, command(), modelName, command(), modelName, createdModel.Name, createdModel.Name)
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
	fmt.Printf(deployedExposedMessage, modelName, ollamaHost, ollamaHost, modelName, ollamaHost, modelName)

	return nil
}
