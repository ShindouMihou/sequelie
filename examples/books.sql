-- Note: Sequelie is highly case-sensitive from declaration names.
-- Note: Sequelie namespace and address are transformed into lowercase automatically.
-- Note: This file assumes that the SQL Dialect (hi Postgres) uses $N instead of other placeholders.

-- A simple example of how to declare sequelie files.
-- Sequelie files have two important parts, the namespace (optional) and the address (from sequelie:query).

-- To define a namespace, which can be used to prepend a category e.g. "books.", you can use the "-- sequelie:namespace [namespace]"
-- sequelie:namespace books

-- Declarations are variables that exists globally in the file. You can use it to declare common values and
-- imprint it into the query on import (a.k.a when the file is being read, it will automatically be interpolated).
--
-- To define a declaration, you can do the following:
-- sequelie:declare TABLE books
-- sequelie:declare ROMANCE_CATEGORY 'romance'

-- sequelie:query get

-- Declarations can be used by using the format {$$key} such as {$$TABLE}. These are immediately replaced
-- upon startup to reduce waste.
SELECT * FROM {$$TABLE} WHERE id = $1

-- sequelie:query get_with_field

-- You can also enable "literal interpolations" other than declarations, these are interpolations
-- that can be done via code. You can use the format {&key} to indicate a literal interpolation,
-- but you shouldn't use it for actual values, or user-input as it bares an SQL-Injection risk.
--
-- To interpolate in the code, you can use `sequelie.GetAndTransform("books.get_with_field", sequelie.Map{"field":"id"})`
-- WARNING: You shouldn't use this in user-inputted values or in values, in general. This bares the risk of SQL-injection
-- And is intended for direct-injection such as fields and related.
--
-- Also, this can be less-performant since it may use fmt.Sprint if there are no marshal methods.
SELECT * FROM {$$TABLE} WHERE {&field} = $1

-- sequelie:query get_romance_books

SELECT * FROM {$$TABLE} WHERE category = {$$ROMANCE_CATEGORY}

-- sequelie:query test
-- NOTE: This is used directly in the test.go file to test out transformers, you shouldn't follow this pattern
-- as much as possible (e.g. the use of literal interpolation in the values).
SELECT * FROM {$$TABLE} WHERE {&field} = {&value}