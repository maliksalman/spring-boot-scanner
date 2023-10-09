package scan

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

var legacyJavaVersion = map[int]string{
	45: "1.1",
	46: "1.2",
	47: "1.3",
	48: "1.4",
}

func FindJavaCompilerVersionFromBootJar(bootJar string, bootContentPrefix string) (string, error) {
	reader, err := zip.OpenReader(bootJar)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	for _, f := range reader.File {
		if !f.FileInfo().IsDir() &&
			strings.Contains(f.Name, bootContentPrefix+"/classes/") &&
			strings.HasSuffix(f.Name, ".class") {

			classFileReader, err := f.Open()
			if err != nil {
				return "", err
			}
			defer classFileReader.Close()

			found, ver := getClassFileVersion(classFileReader, f.Name)
			if found {
				return ver, nil
			}
		}
	}

	return "", errors.New("Can't find a class file in BootJar")
}

func FindJavaCompilerVersionFromReader(tarGzippedReader io.Reader, bootContentPrefix string) (bool, string) {
	gzipReader, _ := gzip.NewReader(tarGzippedReader)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if hdr.Typeflag == tar.TypeReg &&
			hdr.Size >= 8 &&
			strings.Contains(hdr.Name, bootContentPrefix+"/classes/") &&
			strings.HasSuffix(hdr.Name, ".class") {

			found, ver := getClassFileVersion(tarReader, hdr.Name)
			if found {
				return true, ver
			}
		}
	}

	return false, ""
}

func getClassFileVersion(reader io.Reader, filename string) (bool, string) {

	header := make([]byte, 8)
	reader.Read(header)

	foundVersion := int(header[7])
	if foundVersion >= 45 && foundVersion <= 48 {
		ver := legacyJavaVersion[foundVersion]
		log.Printf("Found class file: Name=%s, 8thByte=%d, Version=%s\n", filename, foundVersion, ver)
		return true, ver
	} else if foundVersion >= 49 {
		ver := fmt.Sprintf("%d", foundVersion-44)
		log.Printf("Found class file: Name=%s, 8thByte=%d, Version=%s\n", filename, foundVersion, ver)
		return true, ver
	}
	return false, ""
}

func FindJavaCompilerVersionFromContent(tarGzippedContent []byte, bootContentPrefix string) (bool, string) {
	return FindJavaCompilerVersionFromReader(bytes.NewReader(tarGzippedContent), bootContentPrefix)
}
