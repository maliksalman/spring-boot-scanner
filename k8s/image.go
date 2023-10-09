package k8s

import (
	"context"
	"github.com/google/go-containerregistry/cmd/crane/cmd"
	"github.com/google/go-containerregistry/pkg/crane"
	"log"
	"os"
)

func DownloadImageAsOciDir(ctx context.Context, image string) (string, error) {

	// create temp directory where oci image will be saved
	ociDir, err := os.MkdirTemp("", "oci-image-*")
	if err != nil {
		return "", err
	}

	// pull the image and store it as oci directory
	pullCmd := cmd.NewCmdPull(&[]crane.Option{})
	pullCmd.SetArgs([]string{image, ociDir, "--format", "oci"})

	err = pullCmd.ExecuteContext(ctx)
	if err != nil {
		return "", err
	}

	log.Printf("Pulled Image: Name=%s, OciDir=%s\n", image, ociDir)
	return ociDir, nil
}
