package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute)
	eventInformer := informerFactory.Core().V1().Events()
	eventInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, eventInformer.Informer().HasSynced) {
		return
	}
	fmt.Println("for")
	<- stopCh
}

func onAdd(obj interface{}) {
	event := obj.(*v1.Event)
	// events add Type: Normal, Name: cloud-tool-686c6d8c8f-q49zj.170ef3910d867272, Message: Created container cloud-tool
	fmt.Printf("events add Type: %s, Name: %s, Message: %s\n", event.Type, event.Name, event.Message)
}

func onUpdate(oldObj, newObj interface{}) {
	fmt.Printf("events update\n")
}

func onDelete(obj interface{}) {
	fmt.Printf("events delete\n")
}
