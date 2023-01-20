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
	fmt.Printf("Start upload %d gorutine\n", i)
	stream, err := c.UploadFile(context.Background())
	if err != nil {
		fmt.Println("Error while streaming", err)
		return
	}
	req := &pb.UploadRequest{
		Data: &pb.UploadRequest_FileName{
			FileName: fileName,
		},
	}

	err = stream.Send(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	testFileFolder := "./uploadfile"
	filePath := fmt.Sprintf("%s/%s", testFileFolder, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error while open file", err)
		return
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		i, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while reading file to upload", err)
			return

		}
		req := &pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:i],
			},
		}
		err = stream.Send(req)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DowloandFile(c pb.TagesServiceClient, i int, fileName string) {
	fmt.Printf("Start dowloand %d gorutine\n", i)
	testFileFolder := "./dowloandfile"
	stream, err := c.DownloadFile(context.Background(), &pb.DowloandRequest{
		FileName: fileName,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fileData := bytes.Buffer{}
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

		_, err = fileData.Write(chunk)
		if err != nil {
			break
		}
	}
	err = stream.CloseSend()
	if err != nil {
		fmt.Println("Error while close send", err)
		return
	}

	filePath := fmt.Sprintf("%s/%s", testFileFolder, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = fileData.WriteTo(file)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func GetFile(c pb.TagesServiceClient, i int) {
	fmt.Printf("Get files %d gorutine", i)
	response, err := c.GetFileInfo(context.Background(), &pb.Empty{})
	if err != nil {
		return
	}
	fmt.Println("get images", i, response)
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
		return
	}
	defer conn.Close()
	c := pb.NewTagesServiceClient(conn)

	files := []string{"test1.jpg", "test2.jpg",
		"test3.jpg", "test4.jpg", "test5.png", "test6.jpg",
		"test7.jpg", "test8.jpg", "test9.png", "test10.jpg",
		"test11.jpg", "test12.jpg"}

	fmt.Println("Start uploading files")
	time.Sleep(time.Second * 2)
	for i, file := range files {
		go UploadFile(c, i, file)
	}
	fmt.Println("Stop uploading files")

	fmt.Println("Start dowloand files")
	time.Sleep(time.Second * 10)
	for i, file := range files {
		DowloandFile(c, i, file)
	}
	fmt.Println("Stop dowloand files")
	time.Sleep(time.Second * 2)
	for i := 0; i <= 101; i++ {
		go GetFile(c, i)
	}
	time.Sleep(time.Second * 4)
}
