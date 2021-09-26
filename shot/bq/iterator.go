package bq

import (
	"cloud.google.com/go/bigquery"
)

type Iterator interface {
	Next(dst interface{}) error
}

type IteratorWrapper struct {
	iterator *bigquery.RowIterator
}

func newRawIterator(iterator *bigquery.RowIterator) Iterator {

	return &IteratorWrapper{iterator: iterator}
}

func (i *IteratorWrapper) Next(dst interface{}) error {

	return i.iterator.Next(dst)
}
