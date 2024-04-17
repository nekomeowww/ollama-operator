/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kollama

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericiooptions"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
)

var (
	schemaGroupVersion = ollamav1.GroupVersion

	modelSchemaResourceName         = "models"
	modelSchemaGroupVersionResource = ollamav1.GroupVersion.WithResource(modelSchemaResourceName)
)

// NewCmd provides a cobra command wrapping NamespaceOptions
func NewCmd(streams genericiooptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kollama [cmd] [args] [flags]",
		Short: "CLI for Ollama Operator",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(NewCmdDeploy(streams))

	return cmd
}
