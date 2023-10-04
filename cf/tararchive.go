package cf

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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

func findJavaCompilerVersion(dropletBytes []byte) string {

	gzipReader, _ := gzip.NewReader(bytes.NewReader(dropletBytes))
	defer gzipReader.Close()

	javaCompilerVersion := ""
	tarReader := tar.NewReader(gzipReader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if hdr.Typeflag == tar.TypeReg && hdr.Size >= 8 && (strings.HasPrefix(hdr.Name, "./app/BOOT-INF/classes/") || strings.HasPrefix(hdr.Name, "./app/WEB-INF/classes/")) && strings.HasSuffix(hdr.Name, ".class") {
			header := make([]byte, 8)
			tarReader.Read(header)
			foundVersion := int(header[7])
			if foundVersion >= 45 && foundVersion <= 48 {
				javaCompilerVersion = legacyJavaVersion[foundVersion]
				log.Printf("Found class file: Name=%s, 8thByte=%d, Version=%s\n", hdr.Name, foundVersion, javaCompilerVersion)
				break
			} else if foundVersion >= 49 {
				javaCompilerVersion = fmt.Sprintf("%d", foundVersion-44)
				log.Printf("Found class file: Name=%s, 8thByte=%d, Version=%s\n", hdr.Name, foundVersion, javaCompilerVersion)
				break
			}
		}
	}
	return javaCompilerVersion
}
