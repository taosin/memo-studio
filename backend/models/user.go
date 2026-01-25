package models

import (
	"database/sql"
	"fmt"
	"memo-studio/backend/database"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	MustChangePassword bool `json:"must_change_password"`
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
		"INSERT INTO users (username, password, email, is_admin) VALUES (?, ?, ?, 0)",
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
		"SELECT id, username, password, email, is_admin, must_change_password, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &password, &user.Email, &user.IsAdmin, &user.MustChangePassword, &user.CreatedAt)

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
		"SELECT id, username, email, is_admin, must_change_password, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.MustChangePassword, &user.CreatedAt)

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
	var isAdmin bool
	var mustChange bool
	var createdAt time.Time

	err := database.DB.QueryRow(
		"SELECT id, password, email, is_admin, must_change_password, created_at FROM users WHERE username = ?",
		username,
	).Scan(&userID, &hashedPassword, &email, &isAdmin, &mustChange, &createdAt)

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
		IsAdmin:   isAdmin,
		MustChangePassword: mustChange,
		CreatedAt: createdAt,
	}, nil
}

func UpdateMe(userID int, newUsername, newEmail string) (*User, error) {
	newUsername = strings.TrimSpace(newUsername)
	newEmail = strings.TrimSpace(newEmail)
	if newUsername == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}
	_, err := database.DB.Exec(
		"UPDATE users SET username = ?, email = ? WHERE id = ?",
		newUsername, newEmail, userID,
	)
	if err != nil {
		return nil, err
	}
	return GetUserByID(userID)
}

func ChangePassword(userID int, oldPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return fmt.Errorf("新密码长度至少为6位")
	}
	var hashed string
	err := database.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&hashed)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(oldPassword)); err != nil {
		return sql.ErrNoRows
	}
	newHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHashed), userID)
	if err == nil {
		_, _ = database.DB.Exec("UPDATE users SET must_change_password = 0 WHERE id = ?", userID)
	}
	return err
}

type CreateUserInput struct {
	Username string
	Password string
	Email    string
	IsAdmin  bool
}

func AdminCreateUser(in CreateUserInput) (*User, error) {
	in.Username = strings.TrimSpace(in.Username)
	in.Email = strings.TrimSpace(in.Email)
	if len(in.Username) < 3 {
		return nil, fmt.Errorf("用户名长度至少为3位")
	}
	if len(in.Password) < 6 {
		return nil, fmt.Errorf("密码长度至少为6位")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	adminVal := 0
	if in.IsAdmin {
		adminVal = 1
	}
	res, err := database.DB.Exec(
		"INSERT INTO users (username, password, email, is_admin) VALUES (?, ?, ?, ?)",
		in.Username, string(hashedPassword), in.Email, adminVal,
	)
	if err != nil {
		return nil, err
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetUserByID(int(id64))
}

func AdminListUsers() ([]User, error) {
	rows, err := database.DB.Query("SELECT id, username, email, is_admin, must_change_password, created_at FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.IsAdmin, &u.MustChangePassword, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

func AdminUpdateUser(id int, username, email string, isAdmin bool) (*User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if username == "" {
		return nil, fmt.Errorf("用户名不能为空")
	}
	adminVal := 0
	if isAdmin {
		adminVal = 1
	}
	_, err := database.DB.Exec("UPDATE users SET username = ?, email = ?, is_admin = ? WHERE id = ?", username, email, adminVal, id)
	if err != nil {
		return nil, err
	}
	return GetUserByID(id)
}

func AdminDeleteUser(id int) error {
	// 防止误删默认管理员（可按需调整）
	var username string
	_ = database.DB.QueryRow("SELECT username FROM users WHERE id = ?", id).Scan(&username)
	if username == "admin" {
		return fmt.Errorf("不能删除默认管理员")
	}
	_, err := database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
