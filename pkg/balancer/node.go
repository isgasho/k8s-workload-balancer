package balancer

import (
	"fmt"

	"github.com/li-ang/k8s-workload-balancer/pkg/eviction"
	"github.com/li-ang/k8s-workload-balancer/pkg/listers"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	api_errors "k8s.io/apimachinery/pkg/util/errors"
	corev1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
)

// =================================================================================================

func (b *Balancer) EvictPodByNotMatchNode() error {
	readyNodeMap, err := listers.ListReadyNodesMap(b.NodeLister)
	if err != nil {
		return err
	}
	pods, err := listers.ListAllPods(b.PodLister)
	if err != nil {
		return err
	}
	errList := make([]error, 0)
	for _, pod := range pods {
		node, ok := readyNodeMap[pod.Spec.NodeName]
		if !ok {
			glog.Infof("pod %s/%s is not running in ready node now, skip it", node.Namespace, node.Name)
			continue
		}
		if !checkPodMatchNode(pod, labels.Set(node.Labels)) {
			err = eviction.EvictPod(b.Client, pod, "")
			if err != nil {
				err := fmt.Errorf("evict pod %s/%s error: %v", pod.Namespace, pod.Name, err)
				glog.Error(err)
				errList = append(errList, err)
			}
		}
	}
	if len(errList) > 0 {
		return api_errors.NewAggregate(errList)
	}
	return nil
}

func checkPodMatchNode(pod *corev1.Pod, nlabels labels.Labels) bool {
	if len(pod.Spec.NodeSelector) > 0 {
		selector := labels.SelectorFromSet(pod.Spec.NodeSelector)
		if !selector.Matches(nlabels) {
			return false
		}
	}
	if pod.Spec.Affinity != nil &&
		pod.Spec.Affinity.NodeAffinity != nil &&
		pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		pterms := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
		if len(pterms) == 0 {
			return true
		}
		for _, term := range pterms {
			nodeSelector, err := corev1helper.NodeSelectorRequirementsAsSelector(term.MatchExpressions)
			if err != nil {
				glog.Errorf("terms: %s error: %v", term.MatchExpressions, err)
				return false
			}
			if nodeSelector.Matches(nlabels) {
				return true
			}
		}
	}
	return true
}

// =================================================================================================
