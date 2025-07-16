package models

type User struct {
	Login    string
	Password []byte
	Salt     string
}
