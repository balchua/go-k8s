/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

var namespace string

var pathToConfig string

const ctmJobAnnotation string = "sample.com/job-orchestrator"

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Kube config path for outside of cluster access",
			Destination: &pathToConfig,
		},

		cli.StringFlag{
			Name:        "namespace, ns",
			Value:       "default",
			Usage:       "the namespace where the application will poll the service.",
			Destination: &namespace,
		},
	}

	app.Action = func(c *cli.Context) error {
		var err error
		clientset, err = getClient()
		if err != nil {
			logrus.Error(err)
			return err
		}
		go pollServices()
		// watchNodes()
		for {
			time.Sleep(5 * time.Second)
			logrus.Infof("Doing nothing")
		}
	}
	app.Run(os.Args)
}

func pollServices() error {
	for {
		services, err := clientset.Core().Services(namespace).List(metav1.ListOptions{})
		if err != nil {
			logrus.Warnf("Failed to poll the services: %v", err)
			continue
		}
		for _, service := range services.Items {
			a := service.ObjectMeta.GetAnnotations()
			if a[ctmJobAnnotation] != "" {
				logrus.Infof("Service (%s) has the annotation %s set to %s", service.ObjectMeta.Name, ctmJobAnnotation, a[ctmJobAnnotation])
			} else {
				logrus.Infof("The service (%s) does not have the annotation", service.ObjectMeta.Name)
			}

		}
		time.Sleep(10 * time.Second)
	}
}

func getClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if pathToConfig == "" {
		logrus.Info("Using in cluster config")
		config, err = rest.InClusterConfig()
		// in cluster access
	} else {
		logrus.Info("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", pathToConfig)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
