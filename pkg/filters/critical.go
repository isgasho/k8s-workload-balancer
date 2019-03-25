package filters

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var CriticalPodFilterKey = "criticalpod"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &CriticalPodFilter{}
	}
	Register(CriticalPodFilterKey, newFunc)
}

type CriticalPodFilter struct{}

func (c *CriticalPodFilter) IsEvictablePod(pod *corev1.Pod) bool {
	return !IsCriticalPod(pod)
}

func IsCriticalPod(pod *corev1.Pod) bool {
	if IsCritical(pod.Namespace, pod.Annotations) {
		return true
	}
	return false
}

func IsCritical(ns string, annotations map[string]string) bool {
	if ns != metav1.NamespaceSystem {
		return false
	}
	val, ok := annotations["scheduler.alpha.kubernetes.io/critical-pod"]
	if ok && val == "" {
		return true
	}
	return false
}
