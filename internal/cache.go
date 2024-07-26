package internal

import (
	"errors"
	"io"
	"os"
	"path"
	"strconv"
)

var UserCacheDir, _ = os.UserCacheDir()
var CacheDir = path.Join(UserCacheDir, "aocutil")

func InitCache() {
	os.MkdirAll(CacheDir, 0600)
}

func LoadUserInput(year int, day int) []byte {
	fileDir := path.Join(CacheDir, strconv.Itoa(year))
	filePath := path.Join(fileDir, strconv.Itoa(day)+".input")

	var file *os.File
	var err error

	file, err = os.Open(filePath)
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	return data
}

func SaveUserInput(year int, day int, input []byte) error {
	fileDir := path.Join(CacheDir, strconv.Itoa(year))
	filePath := path.Join(fileDir, strconv.Itoa(day)+".input")

	var file *os.File
	var err error

	os.MkdirAll(fileDir, 0600)
	file, err = os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(filePath)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	defer file.Close()

	file.Write(input)
	return nil
}
