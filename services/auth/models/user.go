package models

type User struct {
	ID          uint64      `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	Email       string      `json:"email"`
	AccountType AccountType `json:"accountType"`
	Password    string      `json:"password"`
	isActivated bool        `json:"-"`
}

type AccountType int

const (
	_ = iota
	Basic
	Google
	Twitter
)
