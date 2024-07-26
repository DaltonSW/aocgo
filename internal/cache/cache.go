package cache

import (
	"errors"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/charmbracelet/log"
)

var UserCacheDir, _ = os.UserCacheDir()
var CacheDir = path.Join(UserCacheDir, "aocutil")
var InputCacheDir = path.Join(CacheDir, "inputs")

func InitCache() {
	os.MkdirAll(CacheDir, 0600)
}

func LoadUserInput(year int, day int, userSession string) []byte {
	log.Infof("Loading user puzzle input for Day %v (%v) for user %v", day, year, userSession)

	fileDir := path.Join(InputCacheDir, userSession, strconv.Itoa(year))
	filePath := path.Join(fileDir, strconv.Itoa(day)+".input")

	var file *os.File
	var err error

	file, err = os.Open(filePath)
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	log.Infof("Success")
	return data
}

func SaveUserInput(year int, day int, userSession string, input []byte) error {
	log.Infof("Saving user puzzle input for Day %v (%v) for user %v", day, year, userSession)

	fileDir := path.Join(InputCacheDir, userSession, strconv.Itoa(year))
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
	log.Infof("Success")
	return nil
}
