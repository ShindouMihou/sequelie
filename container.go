package sequelie

import (
	"errors"
	"os"
	"strings"
	"sync"
)

var store = make(map[string]string)

func readDirs(directories []string, options *Options) error {
	var waitGroup sync.WaitGroup
	var errs []error
	for _, dir := range directories {
		waitGroup.Add(1)
		dir := dir
		go func() {
			if err := readDir(dir, options); err != nil {
				errs = append(errs, err)
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func readDir(directory string, options *Options) error {
	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	var waitGroup sync.WaitGroup
	var errs []error
	for _, file := range dir {
		file := file
		if file.IsDir() {
			waitGroup.Add(1)
			go func() {
				if err := readDir(directory+file.Name(), options); err != nil {
					errs = append(errs, err)
				}
				waitGroup.Done()
			}()
			continue
		}
		if strings.HasSuffix(file.Name(), ".sql") {
			waitGroup.Add(1)
			go func() {
				if err := reader.read(directory+file.Name(), store, options); err != nil {
					errs = append(errs, err)
				}
				waitGroup.Done()
			}()
		}
	}
	waitGroup.Wait()
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func readFile(file string, options *Options) error {
	return reader.read(file, store, options)
}
