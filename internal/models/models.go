package models


/* health response */
type HealthResponse struct {
	Status string `json:"status"`
}

/* username and password */
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
