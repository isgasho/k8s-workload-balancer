package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var CPUManagerFilterKey = "cpumanager"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &CPUManagerFilter{}
	}
	Register(CPUManagerFilterKey, newFunc)
}

// TODO: 暂未实现
type CPUManagerFilter struct{}

func (c *CPUManagerFilter) IsEvictablePod(pod *corev1.Pod) bool { return true }
