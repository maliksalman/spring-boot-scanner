package k8s

import (
	"context"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sApp struct {
	Name             string
	ImagePullSecrets []string
	Namespace        string
	Images           map[string]string
	Labels           map[string]string
}

//go:generate counterfeiter . AppsProvider

type AppsProvider interface {
	ListNamespaces(ctx context.Context, opts metaV1.ListOptions) (*coreV1.NamespaceList, error)
	ListDeployments(ctx context.Context, namespace string, opts metaV1.ListOptions) (*appsV1.DeploymentList, error)
	GetServiceAccount(ctx context.Context, namespace string, accountName string) (*coreV1.ServiceAccount, error)
}

func FindApps(ctx context.Context, provider AppsProvider) ([]K8sApp, error) {

	apps := make([]K8sApp, 0)

	namespaceList, err := provider.ListNamespaces(ctx, metaV1.ListOptions{Limit: 1000})
	if err != nil {
		return nil, err
	}

	for _, ns := range namespaceList.Items {

		deploymentsList, err := provider.
			ListDeployments(ctx, ns.Name, metaV1.ListOptions{Limit: 1000})
		if err != nil {
			return nil, err
		}

		for _, app := range deploymentsList.Items {

			images := make(map[string]string)

			for _, cont := range app.Spec.Template.Spec.Containers {
				images[cont.Name] = cont.Image
			}

			secrets := findPullSecretNames(ctx, app, provider)

			apps = append(apps, K8sApp{
				Name:             app.Name,
				Namespace:        app.Namespace,
				Labels:           app.Labels,
				Images:           images,
				ImagePullSecrets: secrets,
			})
		}
	}

	return apps, nil
}

func findPullSecretNames(ctx context.Context, app appsV1.Deployment, clientSet AppsProvider) []string {

	names := make([]string, 0)

	// find direct references to names
	for _, secret := range app.Spec.Template.Spec.ImagePullSecrets {
		names = append(names, secret.Name)
	}

	// if there is any service account look into that
	saName := app.Spec.Template.Spec.ServiceAccountName
	if saName == "" {
		saName = "default"
	}

	// find the SA and look for names mentioned there
	account, _ := clientSet.GetServiceAccount(ctx, app.Namespace, saName)
	for _, secret := range account.ImagePullSecrets {
		names = append(names, secret.Name)
	}

	return names
}
