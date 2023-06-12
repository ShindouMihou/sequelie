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
SELECT * FROM {$$TABLE} AND {$INSERT:articles.get}
