package book

import "pulley.com/shakesearch/src/book"

var books []book.Book

func GetBooks() []book.Book {
	if len(books) > 0 {
		return books
	} else {
		books, _, _, _ = book.Parse("../../completeworks.txt")
		return books
	}
}
