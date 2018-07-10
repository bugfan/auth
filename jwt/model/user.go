package model

type User struct {
	Id       int64    `json:"id"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Profile  *Profile `json:"profile"`
}
type Profile struct {
}
