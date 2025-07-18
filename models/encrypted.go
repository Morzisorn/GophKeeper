package models

import "time"

type EncryptedItem struct {
	ID        [16]byte
	UserLogin string
	Name      string
	Type      ItemType
	EncryptedData      EncryptedData
	Meta      Meta
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EncryptedData struct {
	EncryptedContent  string `json:"encrypted_content"`
	Nonce string `json:"nonce"`
}

