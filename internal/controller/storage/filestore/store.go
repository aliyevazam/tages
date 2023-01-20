package filestore

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	tages "github.com/aliyevazam/tages/genproto"
)

type DiskFileStore struct {
	mutex      sync.RWMutex
	fileFolder string
	Files      map[string]*FileInfo
}
type FileInfo struct {
	Name string
	Path string
}

const path = "./files/info.json"

func NewDiskFileStore(fileFolder string) *DiskFileStore {
	return &DiskFileStore{
		fileFolder: fileFolder,
		Files:      make(map[string]*FileInfo),
	}
}

func (store *DiskFileStore) Save(
	fileName string,
	fileData bytes.Buffer,
) error {

	filePath := fmt.Sprintf("%s/%s", store.fileFolder, fileName)
	fmt.Println(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	_, err = fileData.WriteTo(file)
	if err != nil {
		return fmt.Errorf("cannot write image to file: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.Files[fileName] = &FileInfo{
		Name: fileName,
		Path: filePath,
	}
	return nil
}

func (store *DiskFileStore) GetImage(
	FileName string,
	stream tages.TagesService_DownloadFileServer,
) error {

	filePath := fmt.Sprintf("%s/%s", store.fileFolder, FileName)

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("cannot open image file: %w", err)
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		size += n

		req := &tages.DowloandResponse{
			ChunkData: buffer[:n],
		}

		err = stream.Send(req)
		if err != nil {
			return err
		}
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()

	return nil
}

