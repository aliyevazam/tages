package service

import (
	"github.com/aliyevazam/tages/internal/controller/storage/repo"
	"github.com/aliyevazam/tages/internal/pkg/logger"
)

type TagesService struct {
	Logger    *logger.Logger
	fileStore repo.FileStorageI
}

func NewTagesService(l *logger.Logger, fileStore repo.FileStorageI) *TagesService {
	return &TagesService{
		Logger:    l,
		fileStore: fileStore,
	}
}
