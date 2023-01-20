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
type StoragePg struct {
	db        *sqlx.DB
	tagesRepo repo.FileStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *StoragePg {
	return &StoragePg{
		db:        db,
		tagesRepo: postgres.NewTagesRepo(db),
	}
}
func (s StoragePg) Tages() repo.FileStorageI {
	return s.tagesRepo
}
