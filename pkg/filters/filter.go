package filters

import (
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type Filter interface {
	IsEvictablePod(*corev1.Pod) bool
}

var (
	mu         sync.Mutex
	newfuncMap = make(map[string]NewFunc)
)

// NewFunc is new filter func.
type NewFunc func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter

func Register(filterKey string, newFunc NewFunc) error {
	mu.Lock()
	defer mu.Unlock()
	_, ok := newfuncMap[filterKey]
	if ok {
		return fmt.Errorf("filter: %s exists", filterKey)
	}
	newfuncMap[filterKey] = newFunc
	return nil
}

// =================================================================================================

type FilterStore struct {
	mu      sync.Mutex
	m       map[string]Filter
	factory informers.SharedInformerFactory
}

func New(clientset kubernetes.Interface, factory informers.SharedInformerFactory) *FilterStore {
	mu.Lock()
	defer mu.Unlock()
	m := make(map[string]Filter)
	for k, newFunc := range newfuncMap {
		filter := newFunc(clientset, factory)
		m[k] = filter
	}
	ps := &FilterStore{
		m:       m,
		factory: factory,
	}
	return ps
}

// =================================================================================================

func (fs *FilterStore) Start(stopCh <-chan struct{}) error {
	informTypes := fs.factory.WaitForCacheSync(stopCh)
	for informType, HasSynced := range informTypes {
		if !HasSynced {
			name := informType.Name()
			return fmt.Errorf("filter init failure, unable to sync caches for %s", name)
		}
	}
	return nil
}

func (fs *FilterStore) JudgePodEvictable(pod *corev1.Pod) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, f := range fs.m {
		if !f.IsEvictablePod(pod) {
			return false
		}
	}
	return true
}

func (fs *FilterStore) FilterUnEvictablePods(pods []*corev1.Pod) []*corev1.Pod {
	ret := make([]*corev1.Pod, 0)
	for i, _ := range pods {
		if fs.JudgePodEvictable(pods[i]) {
			ret = append(ret, pods[i])
		}
	}
	return ret
}
