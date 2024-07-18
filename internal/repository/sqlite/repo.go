package sqlite

import (
	"taskapi/internal/repository"

	"gorm.io/gorm"
)

type repo struct {
	conn *gorm.DB
}

// New 依賴注入
func New(conn *gorm.DB) repository.Repositorier {
	return &repo{
		conn: conn,
	}
}
