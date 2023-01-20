package service

import (
	"github.com/aliyevazam/tages/internal/controller/storage"
	"github.com/aliyevazam/tages/internal/controller/storage/repo"
	"github.com/aliyevazam/tages/internal/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type TagesService struct {
	Storage   storage.IStorage
	Logger    *logger.Logger
	FileStore repo.DisckFileStore
}

func NewTagesService(db *sqlx.DB, l *logger.Logger, fileStore repo.DisckFileStore) *TagesService {
	return &TagesService{
		Storage:   storage.NewStoragePg(db),
		Logger:    l,
		FileStore: fileStore,
	}
}
