package postgres

import (
	pb "github.com/aliyevazam/tages/genproto"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type tagesRepo struct {
	db *sqlx.DB
}

func NewTagesRepo(db *sqlx.DB) *tagesRepo {
	return &tagesRepo{
		db: db,
	}

}

func (t *tagesRepo) CreateOrUpdateFileInfo(FileName string) error {
	inserter := sq.Insert("files").Columns("file_name").
		Values(FileName).Suffix("ON CONFLICT (file_name) DO UPDATE SET updated_at=current_timestamp").
		PlaceholderFormat(sq.Dollar).
		RunWith(t.db)
	_, err := inserter.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (t *tagesRepo) GetFileInfo(*pb.Empty) (*pb.GetFile, error) {
	getter := sq.Select("file_name", "created_at", "updated_at").RunWith(t.db)
	rows, err := getter.Query()
	if err != nil {
		return &pb.GetFile{}, err
	}
	response := &pb.GetFile{}

	for rows.Next() {
		file := &pb.FileInfo{}
		err := rows.Scan(
			&file.FileName, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return &pb.GetFile{}, err
		}
		response.File = append(response.File, file)
	}
	return response, nil

}
