package repo

import (

	pb "github.com/aliyevazam/tages/genproto"
)

type FileStorageI interface {
	GetFileInfo(*pb.Empty) (*pb.GetFile, error)
	CreateOrUpdateFileInfo(FileName string) error
}
