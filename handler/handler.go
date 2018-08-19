package handler

import (
	"github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
)

const ctmJobAnnotation string = "sample.com/job-orchestrator"

//Handle the services with such annotations defined in ctmJobAnnotation
func HandleService(service v1.Service) {
	a := service.ObjectMeta.GetAnnotations()
	if a[ctmJobAnnotation] != "" {
		logrus.Infof("Service (%s) has the annotation %s set to %s", service.ObjectMeta.Name, ctmJobAnnotation, a[ctmJobAnnotation])
	} else {
		logrus.Infof("The service (%s) does not have the annotation", service.ObjectMeta.Name)
	}

}
