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

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
)

const (
	deployExample = `# Deploy a phi model

kollama deploy phi

or

kubectl ollama deploy phi

# Deploy a phi model in a specific namespace

kollama deploy phi -n phi-namespace

or

kubectl ollama deploy phi -n phi-namespace`
)

// CmdDeployOptions provides information required to deploy a model
type CmdDeployOptions struct {
	configFlags     *genericclioptions.ConfigFlags
	clientConfig    clientcmd.ClientConfig
	kubeConfig      *rest.Config
	dynamicClient   dynamic.Interface
	discoveryClient discovery.DiscoveryInterface

	userSpecifiedNamespace string

	genericiooptions.IOStreams
}

// NewCmdDeployOptions provides an instance of NamespaceOptions with default values
func NewCmdDeployOptions(streams genericiooptions.IOStreams) *CmdDeployOptions {
	return &CmdDeployOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

// NewCmdNamespace provides a cobra command wrapping NamespaceOptions
func NewCmdDeploy(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewCmdDeployOptions(streams)

	cmd := &cobra.Command{
		Use:          "deploy [model name] [flags]",
		Short:        "Deploy a model with the given name by using Ollama Operator",
		Example:      deployExample,
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

func (o *CmdDeployOptions) runE(cmd *cobra.Command, args []string) error {
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
	if model != nil {
		if model.Spec.Image == modelImage {
			fmt.Println(modelImage, "deploy")
			return nil
		}

		model.Spec.Image = modelImage

		unstructuredObj, err := Unstructured(model)
		if err != nil {
			return err
		}

		_, err = o.dynamicClient.
			Resource(modelSchemaGroupVersionResource).
			Namespace(o.userSpecifiedNamespace).
			Update(ctx, unstructuredObj, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		fmt.Println(modelImage, "updated")

		return nil
	}

	model = &ollamav1.Model{
		TypeMeta: metav1.TypeMeta{
			APIVersion: schemaGroupVersion.String(),
			Kind:       "Model",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: modelImage,
		},
		Spec: ollamav1.ModelSpec{
			Image: modelImage,
		},
	}

	unstructuredObj, err := Unstructured(model)
	if err != nil {
		return err
	}

	_, err = o.dynamicClient.
		Resource(modelSchemaGroupVersionResource).
		Namespace(o.userSpecifiedNamespace).
		Create(ctx, unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Println(modelImage, "deployed")

	return nil
}
