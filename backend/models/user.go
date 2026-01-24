package models

import (
	"database/sql"
	"memo-studio/backend/database"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUser 创建用户
func CreateUser(username, password, email string) (*User, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := database.DB.Exec(
		"INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		username, string(hashedPassword), email,
	)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(int(userID))
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	var password string
	err := database.DB.QueryRow(
		"SELECT id, username, password, email, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &password, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	// 不返回密码
	return user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := database.DB.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// VerifyPassword 验证密码
func VerifyPassword(username, password string) (*User, error) {
	var userID int
	var hashedPassword string
	var email string
	var createdAt time.Time

	err := database.DB.QueryRow(
		"SELECT id, password, email, created_at FROM users WHERE username = ?",
		username,
	).Scan(&userID, &hashedPassword, &email, &createdAt)

	if err != nil {
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, sql.ErrNoRows // 密码错误
	}

	return &User{
		ID:        userID,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}
