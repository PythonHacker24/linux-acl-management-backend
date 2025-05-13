package models

type Config struct {

}

type HealthResponse struct {
    Status string `json:"status"`
}

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
