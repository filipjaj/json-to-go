type Test struct {
	CreatedAt time.Time `json:"created_at"`
	Profile Profile `json:"profile"`
	Posts []Posts `json:"posts"`
	UserID int `json:"user_id"`
	Username string `json:"username"`
	Email string `json:"email"`
}
type Profile struct {
	Age int `json:"age"`
	Active bool `json:"active"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}
type Posts struct {
	PostID int `json:"post_id"`
	Title string `json:"title"`
	Views int `json:"views"`
}
