package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var HostPathFilterKey = "hostpath"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &HostPathFilter{}
	}
	Register(HostPathFilterKey, newFunc)
}

type HostPathFilter struct{}

func (c *HostPathFilter) IsEvictablePod(pod *corev1.Pod) bool {
	for _, volume := range pod.Spec.Volumes {
		if volume.HostPath != nil {
			return false
		}
	}
	return true
}
