// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mesh

import (
	"time"

	"github.com/spf13/cobra"

	"istio.io/istio/operator/pkg/util/clog"
	buildversion "istio.io/pkg/version"
)

type operatorInitArgs struct {
	// common is shared operator args
	common operatorCommonArgs

	// inFilenames is the path to the input IstioOperator CR.
	inFilename string

	// kubeConfigPath is the path to kube config file.
	kubeConfigPath string
	// context is the cluster context in the kube config.
	context string
	// readinessTimeout is maximum time to wait for all Istio resources to be ready.
	readinessTimeout time.Duration
	// wait is flag that indicates whether to wait resources ready before exiting.
	wait bool
}

const (
	istioControllerComponentName = "Operator"
	istioOperatorCRComponentName = "OperatorCustomResource"
)

func addOperatorInitFlags(cmd *cobra.Command, args *operatorInitArgs) {
	hub, tag := buildversion.DockerInfo.Hub, buildversion.DockerInfo.Tag
	if hub == "" {
		hub = "gcr.io/istio-testing"
	}
	if tag == "" {
		tag = "latest"
	}
	cmd.PersistentFlags().StringVarP(&args.inFilename, "filename", "f", "", "Path to file containing IstioOperator custom resource")
	cmd.PersistentFlags().StringVarP(&args.kubeConfigPath, "kubeconfig", "c", "", "Path to kube config")
	cmd.PersistentFlags().StringVar(&args.context, "context", "", "The name of the kubeconfig context to use")
	cmd.PersistentFlags().DurationVar(&args.readinessTimeout, "readiness-timeout", 300*time.Second, "Maximum seconds to wait for the Istio operator to be ready."+
		" The --wait flag must be set for this flag to apply")
	cmd.PersistentFlags().BoolVarP(&args.wait, "wait", "w", false, "Wait, if set will wait until all Pods, Services, and minimum number of Pods "+
		"of a Deployment are in a ready state before the command exits. It will wait for a maximum duration of --readiness-timeout seconds")

	cmd.PersistentFlags().StringVar(&args.common.hub, "hub", hub, "The hub for the operator controller image")
	cmd.PersistentFlags().StringVar(&args.common.tag, "tag", tag, "The tag for the operator controller image")
	cmd.PersistentFlags().StringVar(&args.common.operatorNamespace, "operatorNamespace", "istio-operator",
		"The namespace the operator controller is installed into")
	cmd.PersistentFlags().StringVar(&args.common.istioNamespace, "istioNamespace", "istio-system",
		"The namespace Istio is installed into")
	cmd.PersistentFlags().StringVarP(&args.common.charts, "charts", "d", "", chartsFlagHelpStr)
}

func operatorInitCmd(rootArgs *rootArgs, oiArgs *operatorInitArgs) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Installs the Istio operator controller in the cluster.",
		Long:  "The init subcommand installs the Istio operator controller in the cluster.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			l := clog.NewConsoleLogger(cmd.OutOrStdout(), cmd.ErrOrStderr(), installerScope)
			operatorInit(rootArgs, oiArgs, l)
		}}
}

// operatorInit installs the Istio operator controller into the cluster.
func operatorInit(args *rootArgs, oiArgs *operatorInitArgs, l clog.Logger) {
	initLogsOrExit(args)

	restConfig, clientset, client, err := K8sConfig(oiArgs.kubeConfigPath, oiArgs.context)
	if err != nil {
		l.LogAndFatal(err)
	}
	// Error here likely indicates Deployment is missing. If some other K8s error, we will hit it again later.
	already, _ := isControllerInstalled(clientset, oiArgs.common.operatorNamespace)
	if already {
		l.LogAndPrintf("Operator controller is already installed in %s namespace, updating.", oiArgs.common.operatorNamespace)
	}

	l.LogAndPrintf("Using operator Deployment image: %s/operator:%s", oiArgs.common.hub, oiArgs.common.tag)

	vals, mstr, err := renderOperatorManifest(args, &oiArgs.common, l)
	if err != nil {
		l.LogAndFatal(err)
	}

	installerScope.Debugf("Installing operator charts with the following values:\n%s", vals)
	installerScope.Debugf("Using the following manifest to install operator:\n%s\n", mstr)

	opts := &Options{
		DryRun:      args.dryRun,
		Wait:        oiArgs.wait,
		WaitTimeout: oiArgs.readinessTimeout,
		Kubeconfig:  oiArgs.kubeConfigPath,
		Context:     oiArgs.context,
	}

	// If CR was passed, we must create a namespace for it and install CR into it.
	customResource, istioNamespace, err := getCRAndNamespaceFromFile(oiArgs.inFilename, l)
	if err != nil {
		l.LogAndFatal(err)
	}

	if err := applyManifest(restConfig, client, mstr, istioControllerComponentName, opts, l); err != nil {
		l.LogAndFatal(err)
	}

	if customResource != "" {
		if err := CreateNamespace(clientset, istioNamespace); err != nil {
			l.LogAndFatal(err)

		}
		if err := applyManifest(restConfig, client, customResource, istioOperatorCRComponentName, opts, l); err != nil {
			l.LogAndFatal(err)
		}
	}

	l.LogAndPrint("\n*** Success. ***\n")
}
