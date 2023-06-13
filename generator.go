package sequelie

import (
	"bytes"
	"errors"
	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
	"sync"
)

var titleCaser = cases.Title(language.AmericanEnglish, cases.NoLower)
var importSequelie = []byte("import \"sequelie\"")

func renderQuery(key string, q *Query) string {
	separatorIndex := strings.Index(key, ".")
	if separatorIndex == -1 {
		panic("cannot find the separator for the key")
	}
	token := key[separatorIndex+1:]
	token = strings.ReplaceAll(token, "_", " ")
	token = titleCaser.String(token)
	token = strings.ReplaceAll(token, " ", "")

	return jen.Var().
		Id(token+"Query").
		Op("=").
		Qual("sequelie", "Query").
		Call(jen.Lit(q.String())).
		GoString()
}

func getPackage(key string) string {
	separatorIndex := strings.Index(key, ".")
	if separatorIndex == -1 {
		panic("cannot find the separator for the key")
	}
	return strings.ToLower(key[:separatorIndex])
}

func flatten() map[string][][]byte {
	var t = make(map[string][][]byte)
	for k, v := range store {
		pkg := getPackage(k)
		t[pkg] = append(t[pkg], []byte(renderQuery(k, v)))
	}
	return t
}

func render(pkg string, declarations [][]byte) []byte {
	t := [][]byte{
		[]byte("package " + pkg),
		{},
		importSequelie,
		{},
	}
	t = append(t, declarations...)
	return bytes.Join(t, newLineBytes)
}

func Generate() error {
	t := flatten()

	var waitGroup sync.WaitGroup
	var errs []error

	for pkg, declarations := range t {
		pkg := pkg
		declarations := declarations

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			rendered := render(pkg, declarations)
			err := os.MkdirAll(".sequelie/"+pkg, os.ModePerm)
			if err != nil {
				errs = append(errs, err)
				return
			}
			f, err := os.Create(".sequelie/" + pkg + "/queries.go")
			if err != nil {
				errs = append(errs, err)
				return
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					errs = append(errs, err)
					return
				}
			}(f)
			_, err = f.Write(rendered)
			if err != nil {
				errs = append(errs, err)
				return
			}
		}()
	}
	waitGroup.Wait()
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
