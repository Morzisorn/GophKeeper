package repositories

import (
	"gophkeeper/config"
	"gophkeeper/internal/server/repositories/database"
)


type Storage interface{
	database.Database
	//Close() error
}

func NewStorage(cfg *config.Config) (Storage, error){
	return database.NewPGDB(cfg)
}