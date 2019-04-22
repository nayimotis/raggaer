package models

import (
	"database/sql"
	"time"
)

// UserLog defines a log action of an user
type UserLog struct {
	ID        int64
	UserID    int64
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AddLogMessage adds a log message for the given user
func (u *User) AddLogMessage(db *sql.DB, msg string) (sql.Result, error) {
	return db.Exec(
		"INSERT INTO user_logs (user_id, message, created_at, updated_at) VALUES (?, ?, ?, ?)",
		u.ID,
		msg,
		time.Now(),
		time.Now(),
	)
}

// UserLogs retrieves all the user logs
func UserLogs(db *sql.DB, u *User) ([]*UserLog, error) {
	logList := []*UserLog{}

	rows, err := db.Query(
		"SELECT id, message, created_at, updated_at FROM user_logs WHERE user_id = ? ORDER BY id DESC",
		u.ID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		l := &UserLog{}
		if err := rows.Scan(&l.ID, &l.Message, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		logList = append(logList, l)
	}

	return logList, nil
}
