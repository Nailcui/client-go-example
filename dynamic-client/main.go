package main

import (
	"context"
	"flag"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
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

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	unStructList, err := dynamicClient.
		Resource(gvr).
		Namespace("default").
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	pods := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructList.UnstructuredContent(), pods)
	if err != nil {
		panic(err.Error())
	}

	for _, item := range pods.Items {
		fmt.Printf("ns: %s, name: %s, status: %s\n", item.Namespace, item.Name, item.Status.Phase)
	}
}
