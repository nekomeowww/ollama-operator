package kollama

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

const (
	unDeployExample = `# Undeploy a phi model

kollama undeploy phi

or

kubectl ollama undeploy phi

# Undeploy a phi model in a specific namespace

kollama undeploy phi -n phi-namespace

or

kubectl ollama undeploy phi -n phi-namespace`
)

type CmdUndeployOptions struct {
	configFlags     *genericclioptions.ConfigFlags
	clientConfig    clientcmd.ClientConfig
	kubeConfig      *rest.Config
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface

	userSpecifiedNamespace string

	genericiooptions.IOStreams
}

func NewCmdUndeployOptions(streams genericiooptions.IOStreams) *CmdUndeployOptions {
	return &CmdUndeployOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

func NewCmdUndeploy(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewCmdUndeployOptions(streams)

	cmd := &cobra.Command{
		Use:          "undeploy [model name] [flags]",
		Short:        "Undeploy a model with the given name by using Ollama Operator",
		Example:      unDeployExample,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			return o.runE(c, args)
		},
	}

	o.configFlags.AddFlags(cmd.Flags())
	o.clientConfig = o.configFlags.ToRawKubeConfigLoader()
	o.kubeConfig = lo.Must(o.clientConfig.ClientConfig())
	o.dynamicClient = lo.Must(dynamic.NewForConfig(o.kubeConfig))
	o.discoveryClient = lo.Must(discovery.NewDiscoveryClientForConfig(o.kubeConfig))

	return cmd
}

func (o *CmdUndeployOptions) runE(cmd *cobra.Command, args []string) error {
	var err error

	o.userSpecifiedNamespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}
	if o.userSpecifiedNamespace == "" {
		var ok bool

		o.userSpecifiedNamespace, ok, err = o.clientConfig.Namespace()
		if err != nil {
			return err
		}
		if !ok {
			o.userSpecifiedNamespace = "default"
		}
	}

	supported, err := IsOllamaOperatorCRDSupported(o.discoveryClient, modelSchemaResourceName)
	if err != nil {
		return err
	}
	if !supported {
		return ErrOllamaModelNotSupported
	}

	modelImage := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	model, err := getOllama(ctx, o.dynamicClient, o.userSpecifiedNamespace, modelImage)
	if err != nil {
		return err
	}
	if model == nil {
		fmt.Println(modelImage, "undeployed")
		return nil
	}

	err = o.dynamicClient.
		Resource(modelSchemaGroupVersionResource).
		Namespace(o.userSpecifiedNamespace).
		Delete(ctx, modelImage, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	fmt.Println(modelImage, "undeployed")

	return nil
}
