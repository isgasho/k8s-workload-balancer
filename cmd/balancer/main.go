package main

import (
	"flag"

	"github.com/li-ang/k8s-workload-balancer/pkg/balancer"
	"github.com/li-ang/k8s-workload-balancer/pkg/signals"

	"github.com/golang/glog"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Fatal(err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)
	factory := informers.NewSharedInformerFactory(clientset, 0)

	go factory.Start(stopCh)
	balancer.New(clientset, factory).Run(stopCh)
}
