package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var PreemptorPodFilterKey = "preemptorpod"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &PreemptorPodFilter{}
	}
	Register(PreemptorPodFilterKey, newFunc)
}

type PreemptorPodFilter struct{}

func (c *PreemptorPodFilter) IsEvictablePod(pod *corev1.Pod) bool {
	return pod.Status.NominatedNodeName == ""
}
