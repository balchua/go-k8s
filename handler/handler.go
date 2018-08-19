package handler

import (
	"github.com/Sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
)

const ctmJobAnnotation string = "sample.com/job-orchestrator"

//Handle the services with such annotations defined in ctmJobAnnotation
func HandleService(service corev1.Service) {
	a := service.ObjectMeta.GetAnnotations()
	if a[ctmJobAnnotation] != "" {
		logrus.Infof("Service (%s) has the annotation %s set to %s", service.ObjectMeta.Name, ctmJobAnnotation, a[ctmJobAnnotation])
	} else {
		logrus.Infof("The service (%s) does not have the annotation", service.ObjectMeta.Name)
	}

}

//Verifies if the deployment contains the right annotation, it will scale the deployment to stop or start
func HandleDeployment(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface, action string) {
	if action == "stop" {
		scaleToZero(deployment, deploymentClient)
	} else {

	}

}

func scaleToZero(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface) {
	a := deployment.ObjectMeta.GetAnnotations()

	if a["intraday-enabled"] == "true" {
		if *deployment.Spec.Replicas > int32(0) {
			logrus.Infof("Deployment (%s) has the annotation scaling to zero", deployment.ObjectMeta.Name)
			deployment.Spec.Replicas = int32Ptr(0)
			deploymentClient.Update(&deployment)
		} else {
			logrus.Infof("Deployment (%s) is already scaled to zero", deployment.ObjectMeta.Name)
		}
	}

}

func int32Ptr(i int32) *int32 {
	return &i
}
