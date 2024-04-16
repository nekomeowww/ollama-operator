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

package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/nekomeowww/ollama-operator/internal/cli/kollama"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	// Import authentication plugins (Go)
	// By default, plugins that use client-go cannot authenticate to Kubernetes clusters on many cloud providers. To address this, include the following import in your plugin:
	// from - Best practices · Krew
	// https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	flags := pflag.NewFlagSet("kollama", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := kollama.NewCmd(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
