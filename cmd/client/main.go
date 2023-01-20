package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/aliyevazam/tages/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func UploadFile(c pb.TagesServiceClient, i int, fileName string) {
	testFileFolder := "./uploadfile"
	filePath := fmt.Sprintf("%s/33.jpg", testFileFolder)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error while open file", err)
	}
	stream, err := c.UploadFile(context.Background())
	if err != nil {
		fmt.Println("Error while stream", err)
	}
	req := &pb.UploadRequest{
		Data: &pb.UploadRequest_FileName{
			FileName: "33.jpg",
		},
	}

	err = stream.Send(req)
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		i, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		req := &pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:i],
			},
		}
		err = stream.Send(req)
		if err != nil {
			fmt.Println(err)
		}
	}
	res, err := stream.CloseAndRecv()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
	}
}

func DowloandFile(c pb.TagesServiceClient, i int, fileName string) {
	testFileFolder := "./dowloandfile"
	imageName := "test1.jpg"
	stream, err := c.DownloadFile(context.Background(), &pb.DowloandRequest{
		FileName: imageName,
	})
	if err != nil {
		fmt.Println(err)
	}
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
		if err != nil {
			fmt.Printf("cannot receive chunk data %s\n", err.Error())
			break
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size

		_, err = fileData.Write(chunk)
		if err != nil {
			break
		}
	}

	filePath := fmt.Sprintf("%s/%s", testFileFolder, imageName)
	fmt.Println(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, err = fileData.WriteTo(file)
	if err != nil {
		fmt.Println(err)
	}

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

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTagesServiceClient(conn)

	files := []string{"test1.jpg", "test2.jpg",
		"test3.jpg", "test4.jpg", "test5.png", "test6.jpg",
		"test7.jpg", "test8.jpg", "test9.png", "test10.jpg",
		"test11.jpg", "test12.jpg", "test13.jpg"}

	fmt.Println("Start uploading files")
	time.Sleep(time.Second * 2)
	for i, file := range files {
		go UploadFile(c, i, file)
	}
	fmt.Println("Stop uploading files")

	fmt.Println("Start dowloand files")
	time.Sleep(time.Second * 10)
	for i, file := range files {
		go DowloandFile(c, i, file)
	}
	fmt.Println("Stop dowloand files")
	// time.Sleep(time.Second * 2)
	// for i := 0; i <= 110; i++ {

	// }
	time.Sleep(time.Second * 4)
}
