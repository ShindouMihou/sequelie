package sequelie

import (
	"fmt"

	"testing"
)

type marshalStringStruct struct {
	Hello string
}

type stringStruct struct {
	Hello string
}

type sequelieMarshalStruct struct {
	Hello string
}

type marshalJsonStruct struct {
	Hello string `json:"hello"`
}

func (test *marshalJsonStruct) SequelieJson() bool {
	return true
}

func (test *marshalStringStruct) MarshalString() string {
	return fmt.Sprint("\"", test.Hello, "\"")
}

func (test *stringStruct) String() string {
	return test.Hello
}

func (test *sequelieMarshalStruct) MarshalSequelie() string {
	return fmt.Sprint("'", test.Hello, "'")
}

func TestReadDir(t *testing.T) {
	err := ReadDirectory("examples/")
	if err != nil {
		t.Fatal("failed to read examples dir: ", err)
		return
	}
}

func TestGet(t *testing.T) {
	v := Get("books.get")
	if v != "SELECT * FROM books WHERE id = $1" {
		t.Error("books.get does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformMarshalString(t *testing.T) {
	v := GetAndTransform("books.test", Map{
		"field": "id",
		"value": &marshalStringStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = \"world\"" {
		t.Error("transform books.test with MarshalString() does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformString(t *testing.T) {
	v := GetAndTransform("books.test", Map{
		"field": "id",
		"value": &stringStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = world" {
		t.Error("transform books.test with String() does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformMarshalSequelie(t *testing.T) {
	v := GetAndTransform("books.test", Map{
		"field": "id",
		"value": &sequelieMarshalStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = 'world'" {
		t.Error("transform books.test with MarshalSequelie() does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformMarshalJson(t *testing.T) {
	v := GetAndTransform("books.test", Map{
		"field": "id",
		"value": &marshalJsonStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = {\"hello\":\"world\"}" {
		t.Error("transform books.test with MarshalJson() does not match desired value, instead got:", v)
	}
}

func TestGetDeclaration(t *testing.T) {
	v := Get("books.get_romance_books")
	if v != "SELECT * FROM books WHERE category = 'romance'" {
		t.Error("books.get_romance_books does not match desired value, instead got: ", v)
	}
}
