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
			apps, err := k8s.FindApps(cmd.Context())
			if err != nil {
				return err
			}

			infos := make([]K8sAppInfo, 0)
			for _, app := range apps {
				for container, image := range app.Images {

					extractedDir, err := k8s.DownloadImageAsOciDir(cmd.Context(), image)
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
				}
			}

			jsonBytes, _ := json.Marshal(infos)
			os.WriteFile("output.json", jsonBytes, 0644)
			return nil
		},
	}
}
