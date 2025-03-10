package helpers

import (
	"encoding/base64"
	"os"
)

func DecodeBase64ToBytes(base64Data string) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func EncodeBytesToBase64(bytesData []byte) string {
	base64Str := base64.StdEncoding.EncodeToString(bytesData)
	return base64Str
}

func WriteBytesToFile(bytes []byte, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// write data to file
	file.Write(bytes)
	return nil
}
