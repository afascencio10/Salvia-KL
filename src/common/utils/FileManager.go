package utils

import (
	"bufio"
	"errors"
	"io"
	"mime/multipart"
	"os"
)

func ReadStringFileHeader(file *multipart.FileHeader) (string, error) {
	uploadedFile, err := file.Open()
	if err != nil {

		return "", errors.New("file error")
	}
	defer uploadedFile.Close()
	buffer := make([]byte, file.Size)
	_, err = uploadedFile.Read(buffer)
	if err != nil {
		return "", errors.New("file error")
	}

	return string(buffer), nil
}

func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
