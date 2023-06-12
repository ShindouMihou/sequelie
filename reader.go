package sequelie

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
	"sync"
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

var enableBytes = []byte("enable")
var operatorsBytes = []byte("operators")

var separatorBytes = []byte(".")
var newLineBytes = []byte("\n")

var declarationPrefixBytes = []byte("{$$")
var encloserBytes = []byte{'}'}

var insertPrefixBytes = []byte("{$INSERT:")

var mutex = sync.RWMutex{}

type localOptions struct {
	Operators bool
}

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
	local := localOptions{Operators: false}

	var push = func() error {
		if address != nil {
			mutex.Lock()
			m[string(*address)] = strings.TrimSpace(builder.String())
			mutex.Unlock()

			builder.Reset()
			local = localOptions{Operators: false}
			return nil
		}
		return errors.New("[no_address] file " + file + " has a query with no address")
	}

	for buf.Scan() {
		line := buf.Bytes()
		if !insensitiveHasPrefix(line, prefixBytes) {
			// FUN FACT: This if-statement exists to prevent white spaces while address is null.
			if address != nil {
				if bytes.HasPrefix(line, commentBytes) && !options.AllowComments {
					continue
				} else {
					for _, declaration := range declarations {
						line = bytes.ReplaceAll(line, declaration[0], declaration[1])
					}
					if local.Operators && namespace != nil {
						clauses := bytes.Fields(line)
						for _, clause := range clauses {
							if bytes.HasPrefix(clause, insertPrefixBytes) && bytes.HasSuffix(clause, encloserBytes) {
								name := clause[len(insertPrefixBytes):]
								name = name[:len(name)-1]

								mutex.RLock()
								v, e := m[string(name)]
								mutex.RUnlock()
								if e {
									line = bytes.Replace(line, clause, []byte(v), 1)
								} else {
									options.Logger.Println(
										"ERR sequelie couldn't insert ", string(name), " into ", string(*address),
										", it is likely because the former hasn't been initialized.")
								}
							}
						}
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
			} else if bytes.EqualFold(args[0], enableBytes) {
				if bytes.EqualFold(args[1], operatorsBytes) {
					local.Operators = true
				}
			}
		}
	}
	return push()
}
