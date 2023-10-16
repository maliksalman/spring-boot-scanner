package k8s

import (
	"context"
	"encoding/json"
	coreV1 "k8s.io/api/core/v1"
	"log"
	"strings"
)

type PullSecretApp struct {
	Namespace string
	Name      string
}

type PullSecretName struct {
	Namespace string
	Name      string
}

type PullSecretServer struct {
	Namespace string
	Server    string
}

type PullSecretData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PullSecretServerData struct {
	Username string
	Password string
	Server   string
}

type DockerConfig struct {
	Auths map[string]PullSecretData `json:"auths"`
}

type PullSecretsDataCache struct {
	appCache map[PullSecretApp][]map[string]PullSecretData
}

func (c *PullSecretsDataCache) Add(ns string, app string, data []map[string]PullSecretData) {
	c.appCache[PullSecretApp{Namespace: ns, Name: app}] = data
}

func (c *PullSecretsDataCache) GetPullSecretServerDataForAppImage(ns string, app string, image string) []PullSecretServerData {

	data := make([]PullSecretServerData, 0)
	possible := c.appCache[PullSecretApp{
		Namespace: ns,
		Name:      app,
	}]

	for _, auths := range possible {
		for serverName, secretData := range auths {
			if strings.HasPrefix(image, serverName) {
				data = append(data, PullSecretServerData{
					Username: secretData.Username,
					Password: secretData.Password,
					Server:   serverName,
				})
			}
		}
	}

	return data
}

//go:generate counterfeiter . SecretsProvider
type SecretsProvider interface {
	GetSecret(ctx context.Context, namespace string, name string) (*coreV1.Secret, error)
}

func BuildPullSecretsCache(ctx context.Context, apps []K8sApp, provider SecretsProvider) *PullSecretsDataCache {

	dataCache := &PullSecretsDataCache{make(map[PullSecretApp][]map[string]PullSecretData)}

	var nameCache = make(map[PullSecretName]map[string]PullSecretData)
	for _, app := range apps {

		appSecretsDataList := make([]map[string]PullSecretData, 0)
		for _, secretName := range app.ImagePullSecrets {

			psName := PullSecretName{Namespace: app.Namespace, Name: secretName}
			if nameCache[psName] != nil {
				appSecretsDataList = append(appSecretsDataList, nameCache[psName])
				continue
			}

			// get the secret
			secret, err := provider.GetSecret(ctx, app.Namespace, secretName)
			if err != nil {
				log.Printf("Can't get secret: Namespace=%s, Secret=%s, Error=%s",
					app.Namespace,
					secretName,
					err.Error())
				continue
			}

			// decode the encoded data
			jsonData := secret.Data[".dockerconfigjson"]

			// unmarshall the json data
			var data DockerConfig
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				log.Printf("Can't unmarshal secret json data: Namespace=%s, Secret=%s, JsonData=%s, Error=%s",
					app.Namespace,
					secretName,
					jsonData,
					err.Error())
				continue
			}

			// add data to name-cache
			nameCache[psName] = data.Auths
			appSecretsDataList = append(appSecretsDataList, nameCache[psName])
		}

		// save the data for easy lookup by app
		dataCache.Add(app.Namespace, app.Name, appSecretsDataList)
	}

	return dataCache
}
