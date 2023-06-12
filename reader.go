package sequelie

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

type iReader struct{}

var reader = iReader{}

// CACHING BYTES
// These are to cache the following strings to byte conversions.
// It may be redundant, but still better to have.

var prefixBytes = []byte("-- sequelie:")
var prefixBytesLen = len(prefixBytes)
var commentBytes = []byte("--")

var declareBytes = []byte("declare")
var queryBytes = []byte("query")
var namespaceBytes = []byte("namespace")

var separatorBytes = []byte(".")
var newLineBytes = []byte("\n")

var declarationPrefixBytes = []byte("{$$")
var encloserBytes = []byte{'}'}

func (reader *iReader) read(file string, m map[string]string, options *Options) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			options.Logger.Println("ERR sequelie failed to close file ", file, ": ", err)
		}
	}(f)
	buf := bufio.NewScanner(f)

	var builder strings.Builder
	var namespace, address *[]byte
	var declarations [][][]byte
	var push = func() error {
		if address != nil {
			m[string(*address)] = strings.TrimSpace(builder.String())
			builder.Reset()
			return nil
		}
		return errors.New("[no_address] file " + file + " has a query with no address")
	}

	for buf.Scan() {
		line := buf.Bytes()
		if !insensitiveHasPrefix(line, prefixBytes) {
			if address != nil {
				if bytes.HasPrefix(line, commentBytes) {
					if !options.AllowComments {
						continue
					}
				} else {
					for _, declaration := range declarations {
						line = bytes.ReplaceAll(line, declaration[0], declaration[1])
					}
				}
				builder.Write(newLineBytes)
				builder.Write(line)
			}
			continue
		}
		args := bytes.Fields(line[prefixBytesLen:])
		if len(args) > 1 {
			if bytes.EqualFold(args[0], namespaceBytes) {
				namespace = ptr(bytes.ToLower(args[1]))
			} else if bytes.EqualFold(args[0], queryBytes) {
				if namespace == nil {
					return errors.New("[no_namespace] file " + file + " has no namespace")
				}
				if address != nil {
					if err := push(); err != nil {
						return err
					}
				}
				address = ptr(bytes.Join([][]byte{*namespace, bytes.ToLower(args[1])}, separatorBytes))
			} else if bytes.EqualFold(args[0], declareBytes) {
				declarations = append(declarations, [][]byte{
					bytes.Join([][]byte{declarationPrefixBytes, encloserBytes}, args[1]),
					line[(prefixBytesLen + 1 + len(args[0]) + 1 + len(args[1])):],
				})
			}
		}
	}
	return push()
}
