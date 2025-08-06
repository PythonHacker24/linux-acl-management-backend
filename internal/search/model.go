package search

/* struct for returning common name, mail, and username */
type User struct {
    CN       string `json:"cn"`
    Mail     string `json:"mail"`
    Username string `json:"username"`
}
