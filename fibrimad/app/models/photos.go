package models

import (
	"database/sql"
	"time"
)

// BoxPhoto defines a photo of a work order box
type BoxPhoto struct {
	ID        int64
	BoxID     int64
	Name      string
	Filename  string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreatePhoto adds the given photo to a box
func (b *WorkOrderBox) CreatePhoto(db *sql.DB, name string) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO work_order_box_photos (box_id, name, filename, path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		b.ID,
		name,
		"",
		"",
		time.Now(),
		time.Now(),
	)
}

// CreatePhotoTransaction adds the given photo to a box using a MySQL transaction
func (b *WorkOrderBox) CreatePhotoTransaction(db *sql.Tx, name string) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO work_order_box_photos (box_id, name, filename, path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		b.ID,
		name,
		"",
		"",
		time.Now(),
		time.Now(),
	)
}

// GetPhotos retrieves all photos from a work order
func (b *WorkOrderBox) GetPhotos(db *sql.DB) ([]*BoxPhoto, error) {
	photoList := []*BoxPhoto{}

	rows, err := db.Query(
		"SELECT id, box_id, name, filename, path, created_at, updated_at FROM work_order_box_photos WHERE box_id = ?",
		b.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		pic := &BoxPhoto{}
		if err := rows.Scan(&pic.ID, &pic.BoxID, &pic.Name, &pic.Filename, &pic.Path, &pic.CreatedAt, &pic.UpdatedAt); err != nil {
			return nil, err
		}
		photoList = append(photoList, pic)
	}

	return photoList, nil
}

// GetBoxPhotoByID retrieves a box photo by its identifier
func GetBoxPhotoByID(db *sql.DB, id int64) (*BoxPhoto, error) {
	photo := &BoxPhoto{}

	row := db.QueryRow(
		"SELECT id, box_id, name, filename, path, created_at, updated_at FROM work_order_box_photos WHERE id = ?",
		id,
	)

	if err := row.Scan(&photo.ID, &photo.BoxID, &photo.Name, &photo.Filename, &photo.Path, &photo.CreatedAt, &photo.UpdatedAt); err != nil {
		return nil, err
	}

	return photo, nil
}

// UpdateBoxPhoto updates the given box photo
func (b *BoxPhoto) UpdateBoxPhoto(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"UPDATE work_order_box_photos SET filename = ?, path = ?, updated_at = ? WHERE id = ?",
		b.Filename,
		b.Path,
		time.Now(),
		b.ID,
	)
}

// DeleteBoxPhotoTransaction deletes the given photo using a MySQL transaction
func (b *BoxPhoto) DeleteBoxPhotoTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_box_photos WHERE id = ?",
		b.ID,
	)
}

// EditBoxPhotoTransaction edits the given photo using a MySQL transaction
func (b *BoxPhoto) EditBoxPhotoTransaction(db *sql.Tx, name string) (sql.Result, error) {
	return db.Exec(
		"UPDATE work_order_box_photos SET name = ? WHERE id = ?",
		name,
		b.ID,
	)
}
