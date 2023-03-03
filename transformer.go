package sequelie

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type hasMarshalString interface {
	MarshalString() string
}
type hasString interface {
	String() string
}
type Marshal interface {
	MarshalSequelie() string
}
type marshalIntoJson interface {
	SequelieJson() bool
}

func transform(contents string, transformers map[string]interface{}, options *Options) string {
	for k, v := range transformers {
		var t string
		if rv, ok := v.(string); ok {
			t = rv
		} else {
			r := reflect.ValueOf(v)
			if m, ok := r.Interface().(Marshal); ok {
				t = m.MarshalSequelie()
			} else if m, ok := r.Interface().(hasMarshalString); ok {
				t = m.MarshalString()
			} else if m, ok := r.Interface().(hasString); ok {
				t = m.String()
			} else if m, ok := r.Interface().(marshalIntoJson); ok {
				if m.SequelieJson() {
					j, e := json.Marshal(v)
					if e != nil {
						options.Logger.Print("ERR sequelie noticed value ", v, " should be marshaled to JSON, "+
							"but cannot marshal to JSON, therefore defaulting to other means: ", e)
					} else {
						t = string(j)
					}
				}
			} else if m, ok := r.Interface().(json.Marshaler); ok {
				j, e := m.MarshalJSON()
				if e != nil {
					options.Logger.Print("ERR sequelie noticed value ", v, " has JSON marshaler, "+
						"but cannot marshal to JSON, therefore defaulting to other means: ", e)
				} else {
					t = string(j)
				}
			} else {
				t = fmt.Sprint(v)
			}
		}
		contents = strings.ReplaceAll(contents, "{&"+k+"}", t)
	}
	return contents
}
