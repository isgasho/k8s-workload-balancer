package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var DaemonSetFilterKey = "daemonset"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &DaemonSetFilter{}
	}
	Register(DaemonSetFilterKey, newFunc)
}

type DaemonSetFilter struct{}

func (d *DaemonSetFilter) IsEvictablePod(pod *corev1.Pod) bool {
	for _, o := range pod.OwnerReferences {
		if o.Kind == "DaemonSet" {
			return false
		}
	}
	return true
}
