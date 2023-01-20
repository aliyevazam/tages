package postgres

import (
	pb "github.com/aliyevazam/tages/genproto"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type tagesRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *tagesRepo {
	return &tagesRepo{
		db: db,
	}

}

func (t *tagesRepo) CreateFileInfo(req *pb.FileInfo) error {
	inserter := sq.Insert("files").Columns("filename", "created_at", "updated_at").
		Values(req.FileName, req.CreatedAt, req.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		RunWith(t.db)
	_, err := inserter.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (t *tagesRepo) GetFileInfo (*pb.Empty)(*pb.GetFile,error) {
	getter := sq.Select("filename","created_at","updated_at").RunWith(t.db)
	result,err := getter.Query()
	if err != nil {
		return &pb.GetFile{},err
	}
	response := &pb.GetFile{}

	for result.Next() {
		file := &pb.FileInfo{}
		err := result.Scan(&file.FileName,&file.CreatedAt,&file.UpdatedAt)
		if err != nil {
			return &pb.GetFile{},err
		}
		response = append(response, file)
	}
	return response,nil


}
