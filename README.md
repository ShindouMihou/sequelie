### Sequelie

A Golang library for off-loading SQL queries from code; Sequelie is a library for Golang applications to off-load their 
SQL queries from their code to another file in an effort to make code cleaner.

To learn by code, you can review our examples:
- [books.sql](examples/books.sql)
- [using `books.get`](examples/get)
- [using `books.get_with_field`](examples/get-with-field)

#### 📦 Installation
```shell
go install github.com/ShindouMihou/sequelie
```

#### 💭 Creating Sequelie Files

If you are into the long text, then let's delve into Sequelie files. You can create Sequelie-supported files by 
creating a regular SQL (`.sql`) file, but with some quirks:
1. It can contain more than one query, delimited by an address.
2. It can contain a namespace.
3. All queries must have one address, not more, not less.

Once you understand those rules, let us understand what the terms "address" and "namespace" are. In Sequelie, we need 
a way to uniquely address queries and that is addresses come to play. Addresses are simply keys to the query, we can use 
it to reference the query, such as "books.get" or simply "get".

Addresses are the most important part into Sequelie. These are automatically lowercase, does not use spaces and 
should be unique to one query. You can set the address of a query by using the `-- sequelie:query {address}` 
syntax, such as in the example:
```sql
-- sequelie:query get
SELECT * FROM books WHERE id = $1 
```

As a general rule, you should categorize your addresses to prevent overriding each other, an example of how to do this 
is, as follows:
```sql
-- sequelie:query books.get
SELECT * FROM books WHERE id = $1 
```

But adding that category can get tedious fast, and that is where namespaces come to place. You can define them by using 
the `-- sequelie:namespace {namespace}` syntax, and it can be used to prepend the namespace into each query, such as 
in the example:
```sql
-- sequelie:namespace books

-- sequelie:query get
SELECT * FROM books WHERE id = $1 

-- sequelie:query get_with_name
SELECT * FROM books WHERE name = $1
```

Once you have created your Sequelie file, we can then load them into our code.  We recommend separating non-Sequelie files,
i.e. literal SQL files, from the Sequelie files since Sequelie will still assume regular `.sql` files as a potential 
Sequelie file and will throw errors if it doesn't find specific rules it likes.

To load Sequelie files, you can do either of the following:
```go
sequelie.ReadDirectory("<directory here>")
sequelie.ReadDirectories("<directory>", "<directory two>")
sequelie.ReadFile("<file path here>")
```

To fetch a query, you can use `sequelie.Get` method, such as in the example:
```go
sequelie.Get("books.get")
```

You can do more with Sequelie such as using declarations and literal interpolations. To learn more about them, you can
scroll down a tiny bit from here.

##### 🔬 Declarations

Declarations are variables in Sequelie that are interpolated directly upon loading of the file. They are used to 
remove magic values, or magic numbers, and you can define them by using the 
`-- sequelie:define {name} {value, spaces are accepted without the need for quoting}`.

There are a few rules into declarations and those are:
1. They can be overridden when re-declared again, i.e mutable variables.
2. You must declare them before actually being used, kind of common sense, at this point.
3. Declarations have case-sensitive names.
4. Declarations cannot have spaces in the names, although they support symbols.

To use a declaration, you can use the `{$$key}` placeholder, such as the following example:
```sql
-- sequelie:namespace books
-- sequelie:declare TABLE books

-- sequelie:query get
SELECT * FROM {$$TABLE} WHERE id = $1 
```

##### 🔬 Literal Interpolation

Literal Interpolation, similar to declarations, are a way to interpolate data into queries, but unlike declarations, you can 
use literal interpolations to interpolate data from the Golang code. Although, they are not escaped or anything, as the name 
implies, i.e literal. 

> **Warning**
> 
> We do not recommend IN ANY WAY using literal interpolation to interpolate user-inputted values, or even values, in general, 
> especially without proper escaping since literal interpolation interpolates the data as indicated.

> **Note**
> 
> Literal Interpolation will be so much slower than using `sequelie.Get` since this has to do some interpolation which 
> is simply `marshal-replace` and will definitely take some processing time, but it shouldn't take horrifically long.

To use literal interpolation, you can use the `{&key}` placeholder, such as the following example:
```sql
-- sequelie:namespace books
-- sequelie:declare TABLE books

-- sequelie:query get_with_field
SELECT * FROM {$$TABLE} WHERE {&field} = $1 
```
```go
query := sequelie.GetAndTransform("books.get_with_field", sequelie.Map{"field":"id"})
```

Additionally, Sequelie can automatically handle marshaling the data into JSON by adding the following 
to the struct's declaration:
```go
type ExampleStruct struct {}
func (example *ExampleStruct) SequelieJson() bool {
	return true
}
```

You can also force Sequelie to custom marshal the data by adding either `MarshalSequelie()`, `MarshalString()` `String()`, or even `MarshalJson()` 
with the priority being as how we laid them out in this text, i.e:
1. `MarshalSequence()`
2. `MarshalString()`
3. `String()`
4. `MarshalJson()`
5. `fmt.Sprint()`

As an example of using `MarshalSequelie()`:
```go
func (example *ExampleStruct) MarshalSequelie() string {
    return "hi"
}
```

If Sequelie cannot find any marshaling method for the structure, it will use `fmt.Sprint` to marshal the data into a String.

##### 🔬 Reusing Queries

Sequelie supports a moderate amount of query reuse through an operator named `Insert Operator` which enables you to insert 
existing queries into the new query, albeit, this is a simple addition that requires the query to be:
1. Initialized **AFTER** the queries being inserted, this means that all queries to be reused must be initialized first.
2. `Local Operators` has to be enabled for the query.

To get started with reusing queries, you have to remember the operator which is:
```sql
{$INSERT:namespace.query}
```

Once you know the operator, we can start with enabling `Local Operators` for the specific query. The reason that this mode exists 
is to prevent additional overhead when reading queries or files that do not use operators such as `Insert Operator`, to enable this, 
you have to add the line:
```sql
-- sequelie:enable operators
```

And now, you can then reuse the queries such in the example of:
```sql
-- sequelie:query reuse
-- sequelie:enable operators

SELECT * FROM {$$TABLE} AND {$INSERT:articles.get}
```

You can view the full example of this in:
- [`examples/articles.sql`](examples/articles.sql)
