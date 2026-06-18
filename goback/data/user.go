package data

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // This will store the hashed password
}

type UserModel struct {
	DB         *sql.DB
	hmacSecret []byte
}

// NewUserModel creates a new user model with HMAC secret key
func NewUserModel(db *sql.DB, secretKey string) *UserModel {
	return &UserModel{
		DB:         db,
		hmacSecret: []byte(secretKey),
	}
}

// hashPassword creates HMAC-SHA256 hash of the password
func (m *UserModel) hashPassword(password string) string {
	h := hmac.New(sha256.New, m.hmacSecret)
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// verifyPassword compares a plain password with the stored hash
func (m *UserModel) verifyPassword(plainPassword, storedHash string) bool {
	hashedPassword := m.hashPassword(plainPassword)
	return hmac.Equal([]byte(hashedPassword), []byte(storedHash))
}

// Create inserts a new user into the database
func (m *UserModel) Create(username, email, password string) error {
	// Check if user already exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)"
	err := m.DB.QueryRow(checkQuery, username, email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing user: %w", err)
	}
	if exists {
		return errors.New("username or email already exists")
	}

	// Hash the password
	hashedPassword := m.hashPassword(password)

	// Insert the new user
	query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)"
	_, err = m.DB.Exec(query, username, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// GetByUsername retrieves a user by username
func (m *UserModel) GetByUsername(username string) (*User, error) {
	query := "SELECT username, email, password_hash FROM users WHERE username = $1"
	user := &User{}

	err := m.DB.QueryRow(query, username).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := "SELECT username, email, password_hash FROM users WHERE email = $1"
	user := &User{}

	err := m.DB.QueryRow(query, email).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// Authenticate verifies user credentials
func (m *UserModel) Authenticate(usernameOrEmail, password string) (*User, error) {
	var user User
	var query string

	// Check if input is email or username
	if usernameOrEmail == "" {
		return nil, errors.New("username or email is required")
	}

	// Try to find user by username or email
	query = "SELECT username, email, password_hash FROM users WHERE username = $1 OR email = $2"
	err := m.DB.QueryRow(query, usernameOrEmail, usernameOrEmail).Scan(
		&user.Username, &user.Email, &user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("error authenticating user: %w", err)
	}

	// Verify password
	if !m.verifyPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// UpdatePassword updates a user's password
func (m *UserModel) UpdatePassword(username, newPassword string) error {
	hashedPassword := m.hashPassword(newPassword)

	query := "UPDATE users SET password_hash = $1 WHERE username = $2"
	result, err := m.DB.Exec(query, hashedPassword, username)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete removes a user from the database
func (m *UserModel) Delete(username string) error {
	query := "DELETE FROM users WHERE username = $1"
	result, err := m.DB.Exec(query, username)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking delete result: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Exists checks if a user exists by username or email
func (m *UserModel) Exists(username, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)"
	err := m.DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}
	return exists, nil
}
