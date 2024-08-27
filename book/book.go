package book

import (
	"context"
	"encore.app/db"
)

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

type ListBooksResponse struct {
	Books []Book `json:"books"`
}

//encore:api auth method=GET path=/books
func ListBooks(ctx context.Context, page Page) (*ListBooksResponse, error) {
	// page
	// limit
	// - these are going to be query parameters on the GET endpoint.
	// - Should not use request bodies on GET endpoints

	defaultPageLimit := 2
	if page.PageLimit == 0 {
		page.PageLimit = defaultPageLimit
	}

	// page limit and page number from user input needs to translate into SQL offset and limit
	offset := (page.PageNumber - 1) * page.PageLimit

	var books []Book
	rows, err := db.Bookstoredb.Query(ctx, `
        SELECT id, title, author, price 
        FROM books
        OFFSET $1
        LIMIT $2
    `, offset, page.PageLimit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return &ListBooksResponse{Books: books}, err
	}

	return &ListBooksResponse{Books: books}, nil
}
