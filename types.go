type Inspiration struct {
	Links Links `json:"links"`
	Meta Meta `json:"meta"`
	PageProps PageProps `json:"pageProps"`
	Data []any `json:"data"`
}
type Links struct {
	First string `json:"first"`
	Last any `json:"last"`
	Prev any `json:"prev"`
	Next any `json:"next"`
}
type Meta struct {
	To any `json:"to"`
	CurrentPage int `json:"current_page"`
	CurrentPageURL string `json:"current_page_url"`
	From any `json:"from"`
	Path string `json:"path"`
	PerPage int `json:"per_page"`
}
type PageProps struct {
	ShortTitle string `json:"shortTitle"`
	Description string `json:"description"`
	Byline string `json:"byline"`
	Image string `json:"image"`
	WishListUuid string `json:"wishListUuid"`
	Title string `json:"title"`
}
