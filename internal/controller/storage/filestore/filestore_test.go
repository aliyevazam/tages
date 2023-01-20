package filestore

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	tages "github.com/aliyevazam/tages/genproto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testFileFolder := "../../../../files"
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	fileService := tages.NewTagesServiceClient(conn)
	filePath := fmt.Sprintf("%s/laptop.jpg", testFileFolder)
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := fileService.UploadFile(context.Background())
	require.NoError(t, err)

	// imageType := filepath.Ext(filePath)
	req := &tages.UploadRequest{
		Data: &tages.UploadRequest_FileName{
			FileName: "laptop.jpg",
		},
	}
	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &tages.UploadRequest{
			Data: &tages.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}
	savedfilePath := fmt.Sprintf("%s/%s", testFileFolder, "laptop.jpg") //imageType\\)
	require.FileExists(t, savedfilePath)
}

func TestClientDownloadImage(t *testing.T) {
	t.Parallel()

	testFileFolder := "../../../../test"
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	imageName := "asd.jpg"
	fileService := tages.NewTagesServiceClient(conn)
	stream, err := fileService.DownloadFile(context.Background(), &tages.DowloandRequest{
		FileName: imageName,
	})
	require.NoError(t, err)
	fileData := bytes.Buffer{}
	imageSize := 0
	for {
		err := contextError(stream.Context())
		if err != nil {
			break
		}

		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("no more data")
			break
		}
		fmt.Println(t)
		require.NoError(t, err)
		if err != nil {
			fmt.Printf("cannot receive chunk data %s\n", err.Error())
			break
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size

		_, err = fileData.Write(chunk)
		require.NoError(t, err)
		if err != nil {
			break
		}
	}

	filePath := fmt.Sprintf("%s/%s", testFileFolder, imageName)
	fmt.Println(filePath)
	file, err := os.Create(filePath)
	// defer func(file){
	// 	err := file.Close()
	// 	if err != nil {
	// 		fmt.Println("Error while close file")
	// 	}
	// }
	require.NoError(t, err)

	_, err = fileData.WriteTo(file)

	require.NoError(t, err)
	savedfilePath := fmt.Sprintf("%s/%s", testFileFolder, "laptop.jpg")
	require.FileExists(t, savedfilePath)
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
