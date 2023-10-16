package k8s_test

import (
	"context"
	"github.com/maliksalman/spring-boot-scanner/k8s"
	"github.com/maliksalman/spring-boot-scanner/k8s/k8sfakes"
	"github.com/stretchr/testify/assert"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestFindApps_VerifyAppProperties(t *testing.T) {

	// GIVEN
	appsProvider := new(k8sfakes.FakeAppsProvider)
	appsProvider.ListNamespacesReturns(&coreV1.NamespaceList{Items: []coreV1.Namespace{
		{ObjectMeta: metaV1.ObjectMeta{Name: "ns-1"}},
	}}, nil)
	appsProvider.ListDeploymentsReturns(&appsV1.DeploymentList{Items: []appsV1.Deployment{createDeployment("ns-1", "app-1", "")}}, nil)
	appsProvider.GetServiceAccountReturns(createServiceAccount("ns-1", "default"), nil)

	// WHEN
	apps, err := k8s.FindApps(context.Background(), appsProvider)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, "app-1", apps[0].Name)
	assert.Equal(t, "ns-1", apps[0].Namespace)
	assert.Equal(t, map[string]string{"key-1": "val-1", "key-2": "val-2"}, apps[0].Labels)
	assert.Equal(t, map[string]string{"cont-1": "img-1", "cont-2": "img-2"}, apps[0].Images)
}

func TestFindApps_VerifyNamedServiceAccount(t *testing.T) {

	// GIVEN
	appsProvider := new(k8sfakes.FakeAppsProvider)
	appsProvider.ListNamespacesReturns(&coreV1.NamespaceList{Items: []coreV1.Namespace{
		{ObjectMeta: metaV1.ObjectMeta{Name: "ns-1"}},
	}}, nil)
	appsProvider.ListDeploymentsReturns(&appsV1.DeploymentList{Items: []appsV1.Deployment{createDeployment("ns-1", "app-1", "sa-1")}}, nil)
	appsProvider.GetServiceAccountReturns(createServiceAccount("ns-1", "sa-1"), nil)

	// WHEN
	apps, err := k8s.FindApps(context.Background(), appsProvider)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, []string{"secret-1", "secret-2", "secret-3", "secret-4"}, apps[0].ImagePullSecrets)

	// verify the call to get-service-account
	assert.Equal(t, 1, appsProvider.GetServiceAccountCallCount())
	_, saNs, saAcctName := appsProvider.GetServiceAccountArgsForCall(0)
	assert.Equal(t, "ns-1", saNs)
	assert.Equal(t, "sa-1", saAcctName)
}

func TestFindApps_VerifyDefaultServiceAccount(t *testing.T) {

	// GIVEN
	appsProvider := new(k8sfakes.FakeAppsProvider)
	appsProvider.ListNamespacesReturns(&coreV1.NamespaceList{Items: []coreV1.Namespace{
		{ObjectMeta: metaV1.ObjectMeta{Name: "ns-1"}},
	}}, nil)
	appsProvider.ListDeploymentsReturns(&appsV1.DeploymentList{Items: []appsV1.Deployment{createDeployment("ns-1", "app-1", "")}}, nil)
	appsProvider.GetServiceAccountReturns(createServiceAccount("ns-1", "default"), nil)

	// WHEN
	apps, err := k8s.FindApps(context.Background(), appsProvider)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, []string{"secret-1", "secret-2", "secret-3", "secret-4"}, apps[0].ImagePullSecrets)

	// verify the call to get-service-account
	assert.Equal(t, 1, appsProvider.GetServiceAccountCallCount())
	_, saNs, saAcctName := appsProvider.GetServiceAccountArgsForCall(0)
	assert.Equal(t, "ns-1", saNs)
	assert.Equal(t, "default", saAcctName)
}

func createServiceAccount(ns string, name string) *coreV1.ServiceAccount {
	return &coreV1.ServiceAccount{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		ImagePullSecrets: []coreV1.LocalObjectReference{
			{Name: "secret-3"},
			{Name: "secret-4"},
		},
	}
}

func createDeployment(ns string, name string, sa string) appsV1.Deployment {
	return appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"key-1": "val-1",
				"key-2": "val-2",
			}},
		Spec: appsV1.DeploymentSpec{
			Template: coreV1.PodTemplateSpec{
				Spec: coreV1.PodSpec{
					ServiceAccountName: sa,
					Containers: []coreV1.Container{
						{Name: "cont-1", Image: "img-1"},
						{Name: "cont-2", Image: "img-2"},
					},
					ImagePullSecrets: []coreV1.LocalObjectReference{
						{Name: "secret-1"},
						{Name: "secret-2"},
					},
				},
			},
		},
	}
}
