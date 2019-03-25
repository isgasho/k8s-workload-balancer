package eviction

import (
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodEvictReasonss string

func EvictPod(client kubernetes.Interface, pod *corev1.Pod, reason PodEvictReasonss) error {
	eviction := &policyv1beta1.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: &metav1.DeleteOptions{},
	}
	err := client.Policy().Evictions(eviction.Namespace).Evict(eviction)
	if err != nil && apierrors.IsNotFound(err) {
		glog.Errorf("evict pod %s/%s error: %v", pod.Namespace, pod.Name, err)
		return err
	}
	if apierrors.IsNotFound(err) {
		glog.Errorf("pod %s/%s not found", pod.Namespace, pod.Name)
		return nil
	}
	glog.Infof("evict pod %s/%s cased by: %s", pod.Namespace, pod.Name, reason)
	return nil
}
