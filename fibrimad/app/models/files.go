package models

import (
	"database/sql"
	"time"
)

// WorkOrderFile defines a work order document
type WorkOrderFile struct {
	ID          int64
	WorkOrderID int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Path        string
	Filename    string
}

// CreateWorkOrderDocument creates a work order document
func (w *WorkOrderFile) CreateWorkOrderDocument(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO work_order_files (work_order_id, name, path, created_at, updated_at, filename) VALUES (?, ?, ?, ?, ?, ?)",
		w.WorkOrderID,
		w.Name,
		w.Path,
		time.Now(),
		time.Now(),
		w.Filename,
	)
}

// GetFileByID returns the given work order file
func GetFileByID(db *sql.DB, id int64) (*WorkOrderFile, error) {
	file := &WorkOrderFile{}

	row := db.QueryRow(
		"SELECT filename, id, work_order_id, name, path FROM work_order_files WHERE id = ?",
		id,
	)

	if err := row.Scan(&file.Filename, &file.ID, &file.WorkOrderID, &file.Name, &file.Path); err != nil {
		return nil, err
	}

	return file, nil
}

// EditFile edits the given file
func (w *WorkOrderFile) EditFile(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"UPDATE work_order_files SET name = ?, filename = ?, path = ?, updated_at = ? WHERE id = ?",
		w.Name,
		w.Filename,
		w.Path,
		time.Now(),
		w.ID,
	)
}

// DeleteWorkOrderFile deletes the given work order file
func (w *WorkOrderFile) DeleteWorkOrderFile(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_files WHERE id = ?",
		w.ID,
	)
}

// GetFiles retrieves all the work order files
func (w *WorkOrder) GetFiles(db *sql.DB) ([]*WorkOrderFile, error) {
	fileList := []*WorkOrderFile{}

	rows, err := db.Query(
		"SELECT filename, id, created_at, updated_at, name, path FROM work_order_files WHERE work_order_id = ?",
		w.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		file := &WorkOrderFile{}
		if err := rows.Scan(&file.Filename, &file.ID, &file.CreatedAt, &file.UpdatedAt, &file.Name, &file.Path); err != nil {
			return nil, err
		}
		fileList = append(fileList, file)
	}

	return fileList, nil
}

// DeleteWorkOrderFilesTransaction deletes all work order files
func (w *WorkOrder) DeleteWorkOrderFilesTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_files WHERE work_order_id = ?",
		w.ID,
	)
}
