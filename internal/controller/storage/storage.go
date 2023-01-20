package storage

import (
	"github.com/aliyevazam/tages/internal/controller/storage/postgres"
	"github.com/aliyevazam/tages/internal/controller/storage/repo"
	"github.com/jmoiron/sqlx"
)

// IStorage ...
type IStorage interface {
	Tages() repo.FileStorageI
}
type storagePg struct {
	db        *sqlx.DB
	tagesRepo repo.FileStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:        db,
		tagesRepo: postgres.New(db),
	}
}
func (s storagePg) Tages() repo.FileStorageI {
	return s.tagesRepo
}
