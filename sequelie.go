package sequelie

import (
	"log"
)

type Map map[string]interface{}

var Settings = Options{AllowComments: false, Logger: log.Default()}

func ReadFile(file string) error {
	return readFile(file, &Settings)
}

func ReadDirectories(directories ...string) error {
	return readDirs(directories, &Settings)
}

func ReadDirectory(dir string) error {
	return readDir(dir, &Settings)
}

func GetAndTransform(address string, transformers Map) string {
	return transform(Get(address), transformers, &Settings)
}

func Get(address string) string {
	q, ex := store[address]
	if !ex {
		panic("cannot find any query with the address " + address)
	}
	return q
}
