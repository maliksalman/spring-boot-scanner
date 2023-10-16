package k8s

import (
	"context"
	"github.com/google/go-containerregistry/cmd/crane/cmd"
	"github.com/google/go-containerregistry/pkg/crane"
	"log"
	"os"
)

type PullSecretsDataProvider interface {
	GetPullSecretServerDataForAppImage(ns string, app string, image string) []PullSecretServerData
}

type OciImageDownloader struct {
	secretsData PullSecretsDataProvider
	ctx         context.Context
	cacheDir    string
}

func NewOciImageDownloader(ctx context.Context, provider PullSecretsDataProvider) *OciImageDownloader {
	return &OciImageDownloader{
		secretsData: provider,
		ctx:         ctx,
		cacheDir:    "./.image-layers-cache",
	}
}

func (d *OciImageDownloader) DeleteDir(image string, ociDir string) error {
	log.Printf("Deleting Image: Name=%s, OciDir=%s\n", image, ociDir)
	return os.RemoveAll(ociDir)
}

func (d *OciImageDownloader) DownloadAsDir(ns string, app string, image string) (string, error) {

	// create temp directory where oci image will be saved
	ociDir, err := os.MkdirTemp("", "oci-image-*")
	if err != nil {
		return "", err
	}

	// login to the registry if there is a pull-secret
	secretsForAppImage := d.secretsData.GetPullSecretServerDataForAppImage(ns, app, image)
	if len(secretsForAppImage) > 0 {
		authLogin := cmd.NewCmdAuthLogin()
		authLogin.SetArgs([]string{
			secretsForAppImage[0].Server,
			"--username", secretsForAppImage[0].Username,
			"--password", secretsForAppImage[0].Password,
		})
		err := authLogin.ExecuteContext(d.ctx)
		if err != nil {
			return "", err
		}
	}

	// make the cacheDir if it doesn't exist
	err = os.MkdirAll(d.cacheDir, 0755)
	if err != nil {
		return "", err
	}

	// pull the image and store it as oci directory
	pullCmd := cmd.NewCmdPull(&[]crane.Option{})
	pullCmd.SetArgs([]string{
		image, ociDir,
		"--format", "oci",
		"--cache_path", d.cacheDir})

	err = pullCmd.ExecuteContext(d.ctx)
	if err != nil {
		return "", err
	}

	log.Printf("Pulled Image: Name=%s, OciDir=%s\n", image, ociDir)
	return ociDir, nil
}
