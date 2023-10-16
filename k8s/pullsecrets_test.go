package k8s_test

import (
	"context"
	"fmt"
	"github.com/maliksalman/spring-boot-scanner/k8s"
	"github.com/maliksalman/spring-boot-scanner/k8s/k8sfakes"
	"github.com/stretchr/testify/assert"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"testing"
)

func TestBuildPullSecretsCache_SingleAppNoSecrets(t *testing.T) {

	// GIVEN
	secretsProvider := new(k8sfakes.FakeSecretsProvider)
	secretsProvider.GetSecretReturnsOnCall(0, &coreV1.Secret{Data: secretData(t, 5, "")}, nil)

	// WHEN
	data := k8s.BuildPullSecretsCache(context.Background(), []k8s.K8sApp{k8sApp("ns-1", "app-0", []string{}, "reg-2")}, secretsProvider)
	pullSecretServerData := data.GetPullSecretServerDataForAppImage("ns-1", "app-0", "reg-2/img")

	// THEN
	assert.Len(t, pullSecretServerData, 0)
}
func TestBuildPullSecretsCache_SingleAppAndRegistry(t *testing.T) {

	// GIVEN
	secretsProvider := new(k8sfakes.FakeSecretsProvider)
	secretsProvider.GetSecretReturnsOnCall(0, &coreV1.Secret{Data: secretData(t, 1, "")}, nil)

	// WHEN
	data := k8s.BuildPullSecretsCache(context.Background(), []k8s.K8sApp{k8sApp("ns-1", "app-0", []string{"sec-0"}, "reg-0")}, secretsProvider)
	pullSecretServerData := data.GetPullSecretServerDataForAppImage("ns-1", "app-0", "reg-0/img")

	// THEN
	assert.Len(t, pullSecretServerData, 1)
	assert.Equal(t, "reg-0", pullSecretServerData[0].Server)
	assert.Equal(t, "user-0", pullSecretServerData[0].Username)
	assert.Equal(t, "pass-0", pullSecretServerData[0].Password)
}

func TestBuildPullSecretsCache_MultipleAppsAndRegistries(t *testing.T) {

	// GIVEN
	secretsProvider := new(k8sfakes.FakeSecretsProvider)
	secretsProvider.GetSecretReturnsOnCall(0, &coreV1.Secret{Data: secretData(t, 1, "0-")}, nil)
	secretsProvider.GetSecretReturnsOnCall(1, &coreV1.Secret{Data: secretData(t, 1, "1-")}, nil)

	// WHEN
	data := k8s.BuildPullSecretsCache(context.Background(), []k8s.K8sApp{
		k8sApp("ns-1", "app-0", []string{"sec-0"}, "0-reg-0"),
		k8sApp("ns-1", "app-1", []string{"sec-1"}, "1-reg-1"),
	}, secretsProvider)

	// THEN
	pullSecretServerData := data.GetPullSecretServerDataForAppImage("ns-1", "app-1", "1-reg-0/img")
	assert.Len(t, pullSecretServerData, 1)
	assert.Equal(t, "1-reg-0", pullSecretServerData[0].Server)
	assert.Equal(t, "1-user-0", pullSecretServerData[0].Username)
	assert.Equal(t, "1-pass-0", pullSecretServerData[0].Password)

	pullSecretServerData = data.GetPullSecretServerDataForAppImage("ns-1", "app-0", "0-reg-0/img")
	assert.Len(t, pullSecretServerData, 1)
	assert.Equal(t, "0-reg-0", pullSecretServerData[0].Server)
	assert.Equal(t, "0-user-0", pullSecretServerData[0].Username)
	assert.Equal(t, "0-pass-0", pullSecretServerData[0].Password)
}

func TestBuildPullSecretsCache_SingleAppMultipleRegistries(t *testing.T) {

	// GIVEN
	secretsProvider := new(k8sfakes.FakeSecretsProvider)
	secretsProvider.GetSecretReturnsOnCall(0, &coreV1.Secret{Data: secretData(t, 3, "")}, nil)

	// WHEN
	data := k8s.BuildPullSecretsCache(context.Background(), []k8s.K8sApp{
		k8sApp("ns-1", "app-0", []string{"sec-0"}, "reg-1"),
	}, secretsProvider)

	// THEN
	pullSecretServerData := data.GetPullSecretServerDataForAppImage("ns-1", "app-0", "reg-1/img")
	assert.Len(t, pullSecretServerData, 1)
	assert.Equal(t, "reg-1", pullSecretServerData[0].Server)
	assert.Equal(t, "user-1", pullSecretServerData[0].Username)
	assert.Equal(t, "pass-1", pullSecretServerData[0].Password)
}

func TestBuildPullSecretsCache_SingleAppMultipleSecrets(t *testing.T) {

	// GIVEN
	secretsProvider := new(k8sfakes.FakeSecretsProvider)
	secretsProvider.GetSecretReturnsOnCall(0, &coreV1.Secret{Data: secretData(t, 6, "")}, nil)
	secretsProvider.GetSecretReturnsOnCall(1, &coreV1.Secret{Data: secretData(t, 6, "")}, nil)

	// WHEN
	data := k8s.BuildPullSecretsCache(context.Background(), []k8s.K8sApp{
		k8sApp("ns-1", "app-0", []string{"sec-0", "sec-1"}, "reg-5"),
	}, secretsProvider)

	// THEN
	pullSecretServerData := data.GetPullSecretServerDataForAppImage("ns-1", "app-0", "reg-5/img")
	assert.Len(t, pullSecretServerData, 2)
	assert.Equal(t, "reg-5", pullSecretServerData[0].Server)
	assert.Equal(t, "reg-5", pullSecretServerData[1].Server)
}

func k8sApp(ns string, name string, secrets []string, registry string) k8s.K8sApp {
	return k8s.K8sApp{
		Name:             name,
		Namespace:        ns,
		Images:           map[string]string{"app": registry + "/img"},
		Labels:           map[string]string{},
		ImagePullSecrets: secrets,
	}
}

func secretData(t *testing.T, numRegistries int, prefix string) map[string][]byte {

	content := k8s.DockerConfig{Auths: make(map[string]k8s.PullSecretData)}
	for i := 0; i < numRegistries; i++ {
		content.Auths[fmt.Sprintf("%sreg-%d", prefix, i)] = k8s.PullSecretData{
			Username: fmt.Sprintf("%suser-%d", prefix, i),
			Password: fmt.Sprintf("%spass-%d", prefix, i),
		}
	}

	jsonStringAsBytes, err := json.Marshal(content)
	assert.NoError(t, err)

	return map[string][]byte{".dockerconfigjson": jsonStringAsBytes}
}
