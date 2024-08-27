package book

type Page struct {
	PageLimit  int `query:"pageLimit"`
	PageNumber int `query:"pageNumber"`
}
