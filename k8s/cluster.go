package k8s

import (
	"context"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type ClusterClientSet struct {
	clientSet *kubernetes.Clientset
}

func NewClusterClientSet() (*ClusterClientSet, error) {
	cs, err := createClientSet()
	if err != nil {
		return nil, err
	}

	return &ClusterClientSet{cs}, nil
}

func createClientSet() (*kubernetes.Clientset, error) {

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

func (ccs *ClusterClientSet) GetSecret(ctx context.Context, namespace string, name string) (*coreV1.Secret, error) {
	return ccs.clientSet.CoreV1().
		Secrets(namespace).Get(ctx, name, metaV1.GetOptions{})
}

func (ccs *ClusterClientSet) ListDeployments(ctx context.Context, ns string, opts metaV1.ListOptions) (*appsV1.DeploymentList, error) {
	return ccs.clientSet.AppsV1().
		Deployments(ns).List(ctx, opts)
}

func (ccs *ClusterClientSet) ListNamespaces(ctx context.Context, opts metaV1.ListOptions) (*coreV1.NamespaceList, error) {
	return ccs.clientSet.CoreV1().
		Namespaces().List(ctx, opts)
}

func (ccs *ClusterClientSet) GetServiceAccount(ctx context.Context, namespace string, accountName string) (*coreV1.ServiceAccount, error) {
	return ccs.clientSet.CoreV1().
		ServiceAccounts(namespace).Get(ctx, accountName, metaV1.GetOptions{})
}
