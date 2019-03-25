package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var StaticPodFilterKey = "staticpod"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &StaticPodFilter{}
	}
	Register(StaticPodFilterKey, newFunc)
}

type StaticPodFilter struct{}

func (s *StaticPodFilter) IsEvictablePod(pod *corev1.Pod) bool {
	if pod.Annotations != nil && pod.Annotations["kubernetes.io/config.source"] == "file" {
		return false
	}
	return true
}
