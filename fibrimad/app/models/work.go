package models

import (
	"database/sql"
	"time"
)

// WorkOrder defines a work order
type WorkOrder struct {
	ID          int64
	Code        string
	Description string
	State       string
	UserList    []*User
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WorkOrderUser defines a user for a work order
type WorkOrderUser struct {
	ID          int64
	UserID      int64
	WorkOrderID int64
}

// GetWorkOrderByID retrieves a work order by identifier
func GetWorkOrderByID(db *sql.DB, id int64) (*WorkOrder, error) {
	order := &WorkOrder{}
	row := db.QueryRow(
		"SELECT id, code, description, state, start_date, end_date FROM work_orders WHERE id = ?",
		id,
	)
	if err := row.Scan(&order.ID, &order.Code, &order.Description, &order.State, &order.StartDate, &order.EndDate); err != nil {
		return nil, err
	}
	return order, nil
}

// DeleteWorkOrder deletes the given work order
func (w *WorkOrder) DeleteWorkOrder(db *sql.DB) (sql.Result, error) {
	return db.Exec("DELETE FROM work_orders WHERE id = ?", w.ID)
}

// DeleteWorkOrderTransaction deletes the given work order using a MySQL transaction
func (w *WorkOrder) DeleteWorkOrderTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec("DELETE FROM work_orders WHERE id = ?", w.ID)
}

// LoadWorkOrders loads all the application work orders
func LoadWorkOrders(db *sql.DB) ([]*WorkOrder, error) {
	workOrderList := []*WorkOrder{}

	rows, err := db.Query("SELECT id, code, description, state, start_date, end_date FROM work_orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		order := &WorkOrder{}
		if err := rows.Scan(&order.ID, &order.Code, &order.Description, &order.State, &order.StartDate, &order.EndDate); err != nil {
			return nil, err
		}
		workOrderList = append(workOrderList, order)
	}

	return workOrderList, nil
}

// DeleteAssignedUser deletes the given assigned user
func (w *WorkOrder) DeleteAssignedUser(db *sql.DB, u *User) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_users WHERE work_order_id = ? AND user_id = ?",
		w.ID,
		u.ID,
	)
}

// DeleteAssignedUsersTransaction deletes all users of the given work order using a MySQL transaction
func (w *WorkOrder) DeleteAssignedUsersTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_users WHERE work_order_id = ?",
		w.ID,
	)
}

// CreateWorkOrder crestes the given work order
func (w *WorkOrder) CreateWorkOrder(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO work_orders (code, description, state, start_date, end_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		w.Code,
		w.Description,
		w.State,
		w.StartDate,
		w.EndDate,
		time.Now(),
		time.Now(),
	)
}

// AddUserToWorkOrder adds a user to the given work order
func (w *WorkOrder) AddUserToWorkOrder(db *sql.DB, user *User) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO work_order_users (user_id, work_order_id) VALUES (?, ?)",
		user.ID,
		w.ID,
	)
}

// IsUserAssigned checks if the user is assigned to the given order
func (w *WorkOrder) IsUserAssigned(db *sql.DB, user *User) (bool, error) {
	assigned := false

	row := db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM work_order_users WHERE user_id = ? AND work_order_id = ?)",
		user.ID,
		w.ID,
	)

	if err := row.Scan(&assigned); err != nil {
		return false, err
	}

	return assigned, nil
}

// GetUsers retrieves all users assigned to a work order
func (w *WorkOrder) GetUsers(db *sql.DB) ([]*User, error) {
	userList := []*User{}

	rows, err := db.Query(
		"SELECT a.id, a.username, a.password, a.role, a.created_at, a.updated_at FROM users a, work_order_users b WHERE b.work_order_id = ? AND b.user_id = a.id",
		w.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}

	return userList, nil
}

// EditWorkOrder edits the given work order
func (w *WorkOrder) EditWorkOrder(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"UPDATE work_orders SET code = ?, description = ?, start_date = ?, end_date = ?, state = ?, updated_at = ?",
		w.Code,
		w.Description,
		w.StartDate,
		w.EndDate,
		w.State,
		time.Now(),
	)
}

// GetUserWorkOrders retrieves all work orders of the given user
func GetUserWorkOrders(db *sql.DB, u *User) ([]*WorkOrder, error) {
	orderList := []*WorkOrder{}

	rows, err := db.Query(
		"SELECT a.id, a.code, a.description, a.start_date, a.end_date, a.state FROM work_orders a, work_order_users b WHERE b.work_order_id = a.id AND b.user_id = ?",
		u.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		order := &WorkOrder{}
		if err := rows.Scan(&order.ID, &order.Code, &order.Description, &order.StartDate, &order.EndDate, &order.State); err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}

	return orderList, nil
}
