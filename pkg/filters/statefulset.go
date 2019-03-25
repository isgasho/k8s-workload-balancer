package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var StatefulSetFilterKey = "statefulset"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &StatefulSetFilter{}
	}
	Register(StatefulSetFilterKey, newFunc)
}

type StatefulSetFilter struct{}

func (d *StatefulSetFilter) IsEvictablePod(pod *corev1.Pod) bool {
	for _, o := range pod.OwnerReferences {
		if o.Kind == "StatefulSet" {
			return false
		}
	}
	return true
}
