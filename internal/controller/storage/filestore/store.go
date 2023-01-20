package filestore

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	tages "github.com/aliyevazam/tages/genproto"
)

type databaseSchema struct {
	FileInfo map[string]*tages.FileInfo `json:"users"`
}

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

func CreateFileInfo(req databaseSchema) error {
	dat, err := json.Marshal(req)
	fmt.Println("keldi create file infoda ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile() (databaseSchema, error) {
	fmt.Println("keldi readfilega")
	dat, err := json.Marshal(path)
	if err != nil {
		fmt.Println("1 error")
		return databaseSchema{}, err
	}
	db := databaseSchema{}
	err = json.Unmarshal(dat, &db)
	if err != nil {
		fmt.Println("2 error")
		fmt.Println(err)
		return databaseSchema{}, err
	}
	return db, nil
}

func (store *DiskFileStore) GetFileInfo(*tages.Empty) (*tages.FileInfo, error) {
	dat, err := json.Marshal(path)
	if err != nil {
		return &tages.FileInfo{}, err
	}
	response := &tages.FileInfo{}
	err = json.Unmarshal(dat, response)
	if err != nil {
		return &tages.FileInfo{}, err
	}
	return response, nil
}
