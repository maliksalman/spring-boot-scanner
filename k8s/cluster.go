package k8s

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type K8sApp struct {
	Name      string
	Namespace string
	Images    map[string]string
	Labels    map[string]string
}

func getClientSet() (*kubernetes.Clientset, error) {

	var kubeConfig string
	if envVar := os.Getenv("KUBECONFIG"); envVar != "" {
		kubeConfig = envVar
	} else {
		kubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func FindApps(ctx context.Context) ([]K8sApp, error) {

	clientSet, err := getClientSet()
	if err != nil {
		return nil, err
	}

	apps := make([]K8sApp, 0)
	namespaceList, err := clientSet.CoreV1().
		Namespaces().
		List(ctx, v1.ListOptions{Limit: 1000})
	if err != nil {
		return nil, err
	}
	for _, ns := range namespaceList.Items {
		deploymentsList, err := clientSet.AppsV1().
			Deployments(ns.Name).List(ctx, v1.ListOptions{Limit: 1000})
		if err != nil {
			return nil, err
		}
		for _, app := range deploymentsList.Items {
			images := make(map[string]string)
			for _, cont := range app.Spec.Template.Spec.Containers {
				images[cont.Name] = cont.Image
			}
			apps = append(apps, K8sApp{
				Name:      app.Name,
				Namespace: app.Namespace,
				Labels:    app.Labels,
				Images:    images,
			})
		}
	}
	return apps, nil
}
