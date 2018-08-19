package main

import (
	"go-k8s/handler"
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

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true})

	// Output to stdout instead of the default stderr, could also be a file.
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)

}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Kube config path for outside of cluster access",
			Destination: &pathToConfig,
		},

		cli.StringFlag{
			Name:        "namespace, n",
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
		go pollDeployments()
		for {
			//keep the application alive.
			time.Sleep(5 * time.Second)
			logrus.Debug("Doing nothing")
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
			handler.HandleService(service)

		}
		time.Sleep(10 * time.Second)
	}
}

func pollDeployments() error {
	for {
		deploymentClient := clientset.AppsV1().Deployments(namespace)
		deployments, err := deploymentClient.List(metav1.ListOptions{})
		if err != nil {
			logrus.Warnf("Failed to poll the services: %v", err)
			continue
		}
		for _, deployment := range deployments.Items {
			handler.HandleDeployment(deployment, deploymentClient, "stop")

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
