package repo

import (
	"bytes"

	pb "github.com/aliyevazam/tages/genproto"
)

type FileStorageI interface {
	Save(FileName string, fileData bytes.Buffer) error
	GetImage(FileName string, stream pb.TagesService_DownloadFileServer) error
	GetFileInfo(*pb.Empty) (*pb.GetFile, error)
}
