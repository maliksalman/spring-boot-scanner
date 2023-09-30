package cf

import "os"

type TestActor struct{}

func NewTestActor(_ any, _ any, _ any) *TestActor {
	return &TestActor{}
}

func (a *TestActor) DownloadCurrentDropletByAppName(app string, spaceGUID string) ([]byte, string, string, error) {
	file, err := os.ReadFile("test-droplet.tgz")
	return file, "", "", err
}
