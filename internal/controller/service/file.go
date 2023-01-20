package service

import (
	"bytes"
	"context"

	// "encoding/json"
	"fmt"
	"io"
	"log"
	// "time"

	pb "github.com/aliyevazam/tages/genproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// const maxImageSize = 1 << 20

func (b *TagesService) UploadFile(stream pb.TagesService_UploadFileServer) error {
	req, err := stream.Recv()
	fmt.Println(req.GetFileName(), err)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}
	FileName := req.GetFileName()
	fmt.Printf("image-name %s ", FileName)

	fileData := bytes.Buffer{}
	// imageSize := 0
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size: %d", size)

		// imageSize += size
		// if imageSize > maxImageSize {
		// 	return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
		// }

		// write slowly
		// time.Sleep(time.Second)

		_, err = fileData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}
	err = b.fileStore.Save(FileName, fileData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image in folder: %v", err))
	}
	// res := &pb.FileInfo{
	// 	FileName:  FileName,
	// 	CreatedAt: time.Now().String(),
	// 	UpdatedAt: time.Now().String(),
	// }
	
	if err != nil {
		return err
	}
	var m interface{}
	err = stream.RecvMsg(m)
	if err != nil {
		return err
	}
	return nil
}

func (b *TagesService) DownloadFile(req *pb.DowloandRequest, stream pb.TagesService_DownloadFileServer) error {
	err := b.fileStore.GetImage(req.FileName, stream)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot get image from folder: %v", err))
	}
	return nil
}

func (b *TagesService) GetFileInfo(ctx context.Context, req *pb.Empty) (*pb.GetFile, error) {
	res, err := b.fileStore.GetFileInfo(req)
	if err != nil {
		return &pb.GetFile{}, err
	}
	return res, nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
