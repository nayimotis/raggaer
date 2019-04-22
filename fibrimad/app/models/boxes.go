package models

import (
	"database/sql"
	"time"
)

// WorkOrderBox defines a work order box
type WorkOrderBox struct {
	ID          int64
	WorkOrderID int64
	Code        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateBox creates a work order box
func (b *WorkOrderBox) CreateBox(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO work_order_boxes (code, work_order_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		b.Code,
		b.WorkOrderID,
		time.Now(),
		time.Now(),
	)
}

// GetBoxes returns the list of work order boxes
func (w *WorkOrder) GetBoxes(db *sql.DB) ([]*WorkOrderBox, error) {
	boxList := []*WorkOrderBox{}

	rows, err := db.Query(
		"SELECT id, code, created_at, updated_at FROM work_order_boxes WHERE work_order_id = ? ORDER BY id DESC",
		w.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		b := &WorkOrderBox{}
		if err := rows.Scan(&b.ID, &b.Code, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		boxList = append(boxList, b)
	}

	return boxList, nil
}

// GetBoxWorkOrderByID retrieves a work order box by its identifier
func GetBoxWorkOrderByID(db *sql.DB, id int64) (*WorkOrderBox, error) {
	box := &WorkOrderBox{}

	row := db.QueryRow("SELECT id, work_order_id, code, created_at, updated_at FROM work_order_boxes WHERE id = ?", id)

	if err := row.Scan(&box.ID, &box.WorkOrderID, &box.Code, &box.CreatedAt, &box.UpdatedAt); err != nil {
		return nil, err
	}

	return box, nil
}

// DeleteBoxTransaction deletes a box using a MySQL transaction
func (b *WorkOrderBox) DeleteBoxTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_boxes WHERE id = ?",
		b.ID,
	)
}

// DeleteBoxPhotosTransaction deletes all the photos of the given box
func (b *WorkOrderBox) DeleteBoxPhotosTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_box_photos WHERE box_id = ?",
		b.ID,
	)
}

// EditBoxTransaction edits the given box
func (b *WorkOrderBox) EditBoxTransaction(db *sql.Tx, name string) (sql.Result, error) {
	return db.Exec(
		"UPDATE work_order_boxes SET code = ? WHERE id = ?",
		name,
		b.ID,
	)
}
