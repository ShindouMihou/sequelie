package sequelie

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type iReader struct{}

var reader = iReader{}
var prefix = "-- sequelie:"
var prefixLen = len(prefix)

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

	content := ""
	namespace := ""
	var address *string
	declarations := make(map[string]string)
	var push = func() error {
		if address != nil {
			m[strings.ToLower(namespace+"."+*address)] = strings.TrimSpace(content)
			content = ""
			return nil
		}
		return errors.New("[no_address] file " + file + " has a query with no address")
	}
	for buf.Scan() {
		text := buf.Text()
		if !strings.HasPrefix(text, prefix) {
			if address != nil {
				if strings.HasPrefix(text, "--") {
					if !options.AllowComments {
						continue
					}
				} else {
					for k, v := range declarations {
						text = strings.ReplaceAll(text, k, v)
					}
				}
				content += "\n" + text
			}
			continue
		}
		args := strings.Fields(text[12:])
		if len(args) > 1 {
			switch args[0] {
			case "namespace":
				namespace = args[1]
			case "query":
				if address != nil {
					if err := push(); err != nil {
						return err
					}
				}
				address = &args[1]
			case "declare":
				declarations["{$$"+args[1]+"}"] = text[(prefixLen + 1 + len(args[0]) + 1 + len(args[1])):]
			}
		}
	}
	return push()
}
