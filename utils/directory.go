package utils

import (
	"io/ioutil"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}

func CopyFile(sourceContent, desFilePath string) error {
	dExist, err := PathExists(desFilePath)

	if err != nil {
		return err
	}
	if !dExist {
		_, cErr := os.Create(desFilePath)
		if cErr != nil {
			return cErr
		}
	} else {
		return nil
	}

	wErr := ioutil.WriteFile(desFilePath, []byte(sourceContent), 0755)
	if wErr != nil {
		return wErr
	}
	return nil
}
