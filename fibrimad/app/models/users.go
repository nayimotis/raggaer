package models

import (
	"database/sql"
	"time"
)

// User defines an application user
type User struct {
	ID         int64
	Username   string
	Password   string
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsAssigned bool
}

// DeleteFromWorkOrderTransaction deletes a user from a work order using a MySQL transaction
func (u *User) DeleteFromWorkOrderTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM work_order_users WHERE user_id = ?",
		u.ID,
	)
}

// DeleteUserTransaction deletes the given user using a MySQL transaction
func (u *User) DeleteUserTransaction(db *sql.Tx) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM users WHERE id = ?",
		u.ID,
	)
}

// DeleteUser deletes the given user
func (u *User) DeleteUser(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"DELETE FROM users WHERE id = ?",
		u.ID,
	)
}

// EditUser edits the given user
func (u *User) EditUser(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"UPDATE users SET username = ?, password = ?, role = ?, updated_at = ? WHERE id = ?",
		u.Username,
		u.Password,
		u.Role,
		time.Now(),
		u.ID,
	)
}

// UserList retrieves all the users of the database
func UserList(db *sql.DB) ([]*User, error) {
	userList := []*User{}

	rows, err := db.Query("SELECT id, username, role, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}

	return userList, nil
}

// CreateUser creates the given user
func (u *User) CreateUser(db *sql.DB) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO users (username, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		u.Username,
		u.Password,
		u.Role,
		time.Now(),
		time.Now(),
	)
}

// GetUserByID retrieves a user by identifier
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	row := db.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE id = ?",
		id,
	)

	user := User{}

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	)

	user := User{}

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
