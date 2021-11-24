package resources

import "pulley.com/shakesearch/src/book"

type QueryResponse struct {
	Books  []book.Book
	Quotes []book.Book
}
