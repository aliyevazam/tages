package repo

import (
	"bytes"

	pb "github.com/aliyevazam/tages/genproto"
)

type DisckFileStore interface {
	Save(FileName string, fileData bytes.Buffer) error
	DownloadFile(FileName string, stream pb.TagesService_DownloadFileServer) error
}
