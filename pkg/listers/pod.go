package listers

import (
	"github.com/li-ang/k8s-workload-balancer/pkg/filters"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
)

// =================================================================================================

func ListAllPods(plister listerscorev1.PodLister) ([]*corev1.Pod, error) {
	pods, err := plister.List(labels.Everything())
	if err != nil {
		glog.Errorf("list pods error: %v", err)
		return nil, err
	}
	return pods, nil
}

// =================================================================================================

func ListAllPodsByDefaultFiltes(plister listerscorev1.PodLister, fs *filters.FilterStore) ([]*corev1.Pod, error) {
	pods, err := plister.List(labels.Everything())
	if err != nil {
		glog.Errorf("list all pods error: %v", err)
		return nil, err
	}
	return fs.FilterUnEvictablePods(pods), nil
}

// =================================================================================================

func ListPodsByDefaultFiltes(plister listerscorev1.PodLister, fs *filters.FilterStore, node *corev1.Node) ([]*corev1.Pod, error) {
	ret := make([]*corev1.Pod, 0)
	pods, err := ListAllPodsByDefaultFiltes(plister, fs)
	if err != nil {
		glog.Errorf("ListAllPodsByDefaultFiltes error: %v", err)
		return nil, err
	}
	for i, _ := range pods {
		if pods[i].Spec.NodeName == node.Name {
			ret = append(ret, pods[i])
		}
	}
	return ret, nil
}
