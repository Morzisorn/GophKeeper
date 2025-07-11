package models

import "time"

type ItemType string

func (it *ItemType) String() string {
	return string(*it)
}

const (
	ItemTypeCREDENTIALS ItemType = "CREDENTIALS"
	ItemTypeTEXT        ItemType = "TEXT"
	ItemTypeBINARY      ItemType = "BINARY"
	ItemTypeCARD        ItemType = "CARD"
)

var Types []ItemType = []ItemType{ItemTypeCREDENTIALS, ItemTypeTEXT, ItemTypeBINARY, ItemTypeCARD}

type Item struct {
	ID        string
	UserLogin string
	Name      string
	Type      ItemType
	Data      Data
	Meta      Meta
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Meta struct {
	Map map[string]string
}

type Data interface {
	GetType() ItemType
	//Edit() error
	//Encrypt() interface{}
	//Decrypt() interface{}
}

type Credentials struct {
	ItemID   string
	Login    string
	Password string
}

func (c Credentials) GetType() ItemType {
	return ItemTypeCREDENTIALS
}

type Text struct {
	ItemID  string
	Content string
}

func (t Text) GetType() ItemType {
	return ItemTypeTEXT
}

type Binary struct {
	ItemID  string
	Content []byte
}

func (b Binary) GetType() ItemType {
	return ItemTypeBINARY
}

type Card struct {
	ItemID         string
	Number         string
	ExpiryDate     string
	SecurityCode   string
	CardholderName string
}

func (c Card) GetType() ItemType {
	return ItemTypeCARD
}
