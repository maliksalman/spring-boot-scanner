package cmd

import (
	"encoding/json"
	"github.com/maliksalman/spring-boot-scanner/k8s"
	"github.com/spf13/cobra"
	"os"
)

type K8sAppInfo struct {
	Namespace    string
	Name         string
	Container    string
	Image        string
	Labels       map[string]string
	Info         *k8s.JavaInfo
	IsSpringBoot bool
}

func NewCmdK8s() *cobra.Command {
	return &cobra.Command{
		Use:   "k8s",
		Short: "Scans a Kubernetes cluster",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {

			clientSet, err := k8s.NewClusterClientSet()
			if err != nil {
				return err
			}

			apps, err := k8s.FindApps(cmd.Context(), clientSet)
			if err != nil {
				return err
			}

			secretsData := k8s.BuildPullSecretsCache(cmd.Context(), apps, clientSet)
			imageDownloader := k8s.NewOciImageDownloader(cmd.Context(), secretsData)

			infos := make([]K8sAppInfo, 0)
			for _, app := range apps {
				for container, image := range app.Images {

					extractedDir, err := imageDownloader.DownloadAsDir(app.Namespace, app.Name, image)
					if err != nil {
						return err
					}
					isSpringBoot, info, err := k8s.FindJavaInfoFromExtractedImage(extractedDir)
					if err != nil {
						return err
					}
					infos = append(infos, K8sAppInfo{
						Namespace:    app.Namespace,
						Name:         app.Name,
						Container:    container,
						Image:        image,
						Labels:       app.Labels,
						Info:         info,
						IsSpringBoot: isSpringBoot,
					})
					imageDownloader.DeleteDir(image, extractedDir)
				}
			}

			jsonBytes, _ := json.Marshal(infos)
			os.WriteFile("output.json", jsonBytes, 0644)
			return nil
		},
	}
}
