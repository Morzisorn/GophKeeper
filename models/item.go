package models

import (
	"fmt"
	"time"
)

type ItemType string

func (it *ItemType) String() string {
	return string(*it)
}

const (
	ItemTypeUNSPECIFIED ItemType = "UNSPECIFIED"
	ItemTypeCREDENTIALS ItemType = "CREDENTIALS"
	ItemTypeTEXT        ItemType = "TEXT"
	ItemTypeBINARY      ItemType = "BINARY"
	ItemTypeCARD        ItemType = "CARD"
)

var ItemTypes []ItemType = []ItemType{ItemTypeCREDENTIALS, ItemTypeTEXT, ItemTypeBINARY, ItemTypeCARD}

type Item struct {
	ID        [16]byte
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
	Login    string
	Password string
}

func (c Credentials) GetType() ItemType {
	return ItemTypeCREDENTIALS
}

type Text struct {
	Content string
}

func (t Text) GetType() ItemType {
	return ItemTypeTEXT
}

type Binary struct {
	Content []byte
}

func (b Binary) GetType() ItemType {
	return ItemTypeBINARY
}

type Card struct {
	Number         string
	ExpiryDate     string
	SecurityCode   string
	CardholderName string
}

func (c Card) GetType() ItemType {
	return ItemTypeCARD
}

func (t ItemType) CreateDataByType() (Data, error) {
	switch t {
	case ItemTypeCREDENTIALS:
		return &Credentials{}, nil
	case ItemTypeTEXT:
		return &Text{}, nil
	case ItemTypeBINARY:
		return &Binary{}, nil
	case ItemTypeCARD:
		return &Card{}, nil
	default:
		return nil, fmt.Errorf("unknown item type: %s", t)
	}
}
