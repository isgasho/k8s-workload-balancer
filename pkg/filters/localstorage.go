package filters

import (
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
)

var LocalStorageFilterKey = "localstorage"

func init() {
	newFunc := func(clientset kubernetes.Interface, factory informers.SharedInformerFactory) Filter {
		return &LocalStorageFilter{
			PVCLister: factory.Core().V1().PersistentVolumeClaims().Lister(),
			PVLister:  factory.Core().V1().PersistentVolumes().Lister(),
		}
	}
	Register(LocalStorageFilterKey, newFunc)
}

type LocalStorageFilter struct {
	PVCLister listerscorev1.PersistentVolumeClaimLister
	PVLister  listerscorev1.PersistentVolumeLister
}

func (l *LocalStorageFilter) IsEvictablePod(pod *corev1.Pod) bool {
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim != nil {
			b, err := l.IsLocalStoragePVC(pod.Namespace, v.Name)
			if err != nil {
				glog.Error(err)
				return false
			}
			if b {
				return false
			}
		}
	}
	return true
}

func (l *LocalStorageFilter) IsLocalStoragePVC(ns string, pvcName string) (bool, error) {
	pvc, err := l.PVCLister.PersistentVolumeClaims(ns).Get(pvcName)
	if err != nil {
		glog.Errorf("get pvc %s/%s error: %v", ns, pvcName, err)
		return false, err
	}
	pvName := pvc.Spec.VolumeName
	pv, err := l.PVLister.Get(pvName)
	if err != nil {
		glog.Errorf("get pv %s error: %v", pvName, err)
		return false, err
	}
	if pv.Spec.Local != nil {
		return true, nil
	}
	return false, nil
}
