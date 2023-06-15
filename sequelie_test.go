package sequelie

import (
	"fmt"
	"strings"
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
	v := Get("books.get").String()
	if v != "SELECT * FROM books WHERE id = $1" {
		t.Error("books.get does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformMarshalString(t *testing.T) {
	v := Get("books.test").Interpolate(Map{
		"field": "id",
		"value": &marshalStringStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = \"world\"" {
		t.Error("transform books.test with MarshalString() does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformString(t *testing.T) {
	v := Get("books.test").Interpolate(Map{
		"field": "id",
		"value": &stringStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = world" {
		t.Error("transform books.test with String() does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformMarshalSequelie(t *testing.T) {
	v := Get("books.test").Interpolate(Map{
		"field": "id",
		"value": &sequelieMarshalStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = 'world'" {
		t.Error("transform books.test with MarshalSequelie() does not match desired value, instead got:", v)
	}
}

func TestGetAndTransformMarshalJson(t *testing.T) {
	v := Get("books.test").Interpolate(Map{
		"field": "id",
		"value": &marshalJsonStruct{"world"},
	})
	if v != "SELECT * FROM books WHERE id = {\"hello\":\"world\"}" {
		t.Error("transform books.test with MarshalJson() does not match desired value, instead got:", v)
	}
}

func TestGetDeclaration(t *testing.T) {
	v := Get("books.get_romance_books").String()
	if v != "SELECT * FROM books WHERE category = 'romance'" {
		t.Error("books.get_romance_books does not match desired value, instead got: ", v)
	}
	v = Get("articles.reuse").String()
	if v != "SELECT * FROM articles AND "+Get("articles.get").String() {
		t.Error("articles.reuse does not match desired value, instead got: ", v)
	}
}

func TestGenerate(t *testing.T) {
	err := Generate()
	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkReadDirectory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ReadDirectory("examples/")
		if err != nil {
			b.Fatal("failed to read examples dir: ", err)
			return
		}
	}
}

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ReadFile("examples/books.sql")
		if err != nil {
			b.Fatal("failed to read books.sql: ", err)
			return
		}
	}
}

func BenchmarkScan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := reader.scan(strings.NewReader(`
-- sequelie:namespace articles
-- sequelie:declare TABLE articles

-- sequelie:query get

SELECT * FROM {$$TABLE} WHERE id = $1

-- sequelie:query reuse
-- You can reuse existing queries by using the {$INSERT:query} operator, but this feature can only be enabled
-- per query due to performance consequences. To enable the query, you have to write:

-- sequelie:enable operators

-- Additionally, the query being depended upon must be initialized before the query. As such, this isn't the
-- best to use with queries that you cannot guarantee will be initialized before this.
SELECT * FROM {$$TABLE} AND {$INSERT:articles.get}`), store, &Settings)
		if err != nil {
			b.Fatal("failed to scan: ", err)
			return
		}
	}
}

func BenchmarkReadFileWithOperator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ReadFile("examples/articles.sql")
		if err != nil {
			b.Fatal("failed to read books.sql: ", err)
			return
		}
	}
}
