package utils

import (
	"bytes"
	"io/ioutil"
	"math/rand"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdef0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func ReplaceImageTag(transID string, blobPath string, oldTag string, newTag string) bool {
	input, err := ioutil.ReadFile(blobPath + ".tmp")
	if err != nil {
		ConfigZap().Errorf("[%s] Replacing old tag [%s] with new tag [%s]...FAILED: %v", transID, oldTag, newTag, err)
		return false
	}

	output := bytes.Replace(input, []byte(oldTag), []byte(newTag), -1)

	if err = ioutil.WriteFile(blobPath, output, 0666); err != nil {
		ConfigZap().Errorf("[%s] Writing new temporary file to local...OK: %v", transID, err)
		return false
	}

	return true
}