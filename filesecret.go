package gocom

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var fileSecretCache map[string]string = map[string]string{}

func FileSecret(dsn string, key string) (string, error) {

	ret, ok := fileSecretCache[dsn+"_"+key]

	if !ok {
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

		ret = string(byteA)
		fileSecretCache[dsn+"_"+key] = ret
	}

	return ret, nil
}

func init() {

	RegSecretFunc("file", FileSecret)
}
