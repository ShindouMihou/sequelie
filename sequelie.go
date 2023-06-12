package sequelie

import (
	"log"
)

type Map map[string]interface{}

var Settings = Options{AllowComments: false, Logger: log.Default()}

func ReadFile(file string) error {
	return readFile(file, &Settings)
}

func ReadFileWithSettings(file string, settings *Options) error {
	return readFile(file, settings)
}

func ReadDirectories(directories ...string) error {
	return readDirs(directories, &Settings)
}

func ReadDirectoriesWithSettings(settings *Options, directories ...string) error {
	return readDirs(directories, settings)
}

func ReadDirectory(dir string) error {
	return readDir(dir, &Settings)
}

func ReadDirectoryWithSettings(dir string, settings *Options) error {
	return readDir(dir, settings)
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
