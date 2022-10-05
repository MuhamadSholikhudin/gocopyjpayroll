package fileprocessmodel

import (
	"database/sql"

	"gocopyjpayroll/config"
	"gocopyjpayroll/entities"
)

type FileprocessModel struct {
	db *sql.DB
}

func New() *FileprocessModel {
	db, err := config.DBConnection()
	if err != nil {
		panic(err)
	}
	return &FileprocessModel{db: db}
}

func (m *FileprocessModel) FindAll(fileprocess *[]entities.Fileprocess) error {
	rows, err := m.db.Query("select * from fileprocess")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var data entities.Fileprocess
		rows.Scan(
			&data.Id,
			&data.Periode,
			&data.File,
			&data.Category,
			&data.CreatedAt,
			&data.UpdatedAt)
		*fileprocess = append(*fileprocess, data)
	}

	return nil
}

func (m *FileprocessModel) Create(fileprocess *entities.Fileprocess) error {
	result, err := m.db.Exec("insert into fileprocess (periode, file, category) values(?,?,?)",
		fileprocess.Periode, fileprocess.File, fileprocess.Category)

	if err != nil {
		return err
	}

	lastInsertId, _ := result.LastInsertId()
	fileprocess.Id = lastInsertId
	return nil
}

func (m *FileprocessModel) Find(id int64, fileprocess *entities.Fileprocess) error {
	return m.db.QueryRow("select * from fileprocess where id = ?", id).Scan(
		&fileprocess.Id,
		&fileprocess.Periode,
		&fileprocess.File,
		&fileprocess.Category)
}

func (m *FileprocessModel) Update(fileprocess entities.Fileprocess) error {

	_, err := m.db.Exec("update fileprocess set periode = ?, text = ?, category = ? where id = ?",
		fileprocess.Periode, fileprocess.File, fileprocess.Category, fileprocess.Id)

	if err != nil {
		return err
	}

	return nil
}

func (m *FileprocessModel) Delete(id int64) error {
	_, err := m.db.Exec("delete from fileprocess where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
