package listers

import (
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
)

// =================================================================================================

func ListReadyNodes(nlister listerscorev1.NodeLister) ([]*corev1.Node, error) {
	nodes, err := nlister.List(labels.Everything())
	if err != nil {
		glog.Errorf("list nodes error: %v", err)
		return nil, err
	}

	ret := make([]*corev1.Node, 0)
	for i, n := range nodes {
		if !isReadyNode(n) {
			glog.Infof("node %s is not in ready, skip it", n.Name)
			continue
		}
		ret = append(ret, nodes[i])
	}
	return ret, nil
}

func isReadyNode(node *corev1.Node) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady &&
			cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// =================================================================================================

func ListReadyNodesMap(nlister listerscorev1.NodeLister) (map[string]*corev1.Node, error) {
	nodes, err := ListReadyNodes(nlister)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*corev1.Node)
	for i, n := range nodes {
		ret[n.Name] = nodes[i]
	}
	return ret, nil
}
