package balancer

import (
	"github.com/golang/glog"
	"github.com/li-ang/k8s-workload-balancer/pkg/filters"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
)

type Balancer struct {
	NodeLister  listerscorev1.NodeLister
	PodLister   listerscorev1.PodLister
	Client      kubernetes.Interface
	FilterStore *filters.FilterStore
}

func New(clientset kubernetes.Interface, factory informers.SharedInformerFactory) *Balancer {
	return &Balancer{
		NodeLister:  factory.Core().V1().Nodes().Lister(),
		PodLister:   factory.Core().V1().Pods().Lister(),
		Client:      clientset,
		FilterStore: filters.New(clientset, factory),
	}

}

func (b *Balancer) Run(stopCh <-chan struct{}) {
	var err error
	err = b.FilterStore.Start(stopCh)
	if err != nil {
		glog.Fatal("FilterStore store started failed, error: %v", err)
		return
	}

	err = b.EvictPodByNotMatchNode()
	if err != nil {
		glog.Error(err)
	}
	err = b.EvictBestEffortPods()
	if err != nil {
		glog.Error(err)
	}
	err = b.EvictBurstablePods()
	if err != nil {
		glog.Error(err)
	}
	err = b.EvictGuaranteedPods()
	if err != nil {
		glog.Error(err)
	}

}
