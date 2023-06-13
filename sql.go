package sequelie

type Query string

func (q *Query) String() string {
	return string(*q)
}

func (q *Query) Bytes() []byte {
	return []byte(*q)
}

func (q *Query) Interpolate(transformers Map) string {
	return transform(q.String(), transformers, &Settings)
}
