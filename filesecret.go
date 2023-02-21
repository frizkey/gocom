package gocom

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func FileSecret(dsn string, key string) (string, error) {

	targetFile := dsn

	if strings.HasSuffix(targetFile, "/") {
		targetFile += key
	} else {
		targetFile += "/" + key
	}

	byteA, err := ioutil.ReadFile(targetFile)

	if err != nil {
		return "", fmt.Errorf("read secret file error : %w", err)
	}

	return string(byteA), nil
}

func init() {

	RegSecretFunc("file", FileSecret)
}
