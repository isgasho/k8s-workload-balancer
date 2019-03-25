package filters

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

var JobFilterKey = "job"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &JobFilter{}
	}
	Register(JobFilterKey, newFunc)
}

type JobFilter struct{}

func (j *JobFilter) IsEvictablePod(pod *corev1.Pod) bool {
	for _, o := range pod.OwnerReferences {
		if o.Kind == "Job" ||
			pod.Status.Phase == corev1.PodSucceeded ||
			pod.Status.Phase == corev1.PodFailed {
			return false
		}
	}
	return true
}
