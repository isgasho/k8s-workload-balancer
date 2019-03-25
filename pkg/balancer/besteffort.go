package balancer

import (
	"fmt"

	"github.com/li-ang/k8s-workload-balancer/pkg/eviction"
	"github.com/li-ang/k8s-workload-balancer/pkg/listers"
	"github.com/li-ang/k8s-workload-balancer/pkg/qos"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	api_errors "k8s.io/apimachinery/pkg/util/errors"
)

func (b *Balancer) EvictBestEffortPods() error {
	rnodeMap, err := listers.ListReadyNodesMap(b.NodeLister)
	if err != nil {
		return err
	}
	pods, err := listers.ListAllPodsByDefaultFiltes(b.PodLister, b.FilterStore)
	if err != nil {
		return err
	}

	pods = FilterPods(pods, IsBestEffortPod)
	errList := make([]error, 0)
	success := 0
	for _, pod := range pods {
		if (success - len(pods)/2) > 0 {
			break
		}
		_, ok := rnodeMap[pod.Spec.NodeName]
		if !ok {
			glog.Infof("pod %s/%s is not running in ready node now, skip it", pod.Namespace, pod.Name)
			continue
		}
		err = eviction.EvictPod(b.Client, pod, "")
		if err != nil {
			err := fmt.Errorf("evict pod %s/%s error: %v", pod.Namespace, pod.Name, err)
			glog.Error(err)
			errList = append(errList, err)
			continue
		}
		success++
	}
	if len(errList) > 0 {
		return api_errors.NewAggregate(errList)
	}
	return nil
}

func IsBestEffortPod(pod *corev1.Pod) bool {
	return qos.GetPodQOS(pod) == corev1.PodQOSBestEffort
}
