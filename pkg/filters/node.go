package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var NodeFilterKey = "node"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &NodeFilter{}
	}
	Register(NodeFilterKey, newFunc)
}

type NodeFilter struct{}

func (n *NodeFilter) IsEvictablePod(pod *corev1.Pod) bool {
	if pod.Spec.NodeName == "" {
		return false
	}
	return true
}
