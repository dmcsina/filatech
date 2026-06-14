package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// UserRequest represents the JSON request body for creating/updating users
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse represents the JSON response for user operations
type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ErrorResponse represents an error message response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// healthCheck handles GET /v1/healthcheck
func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":      "available",
		"environment": "dev",
		"version":     "1.0.0",
	}

	app.sendJSONResponse(w, response, http.StatusOK)
}

// createUser handles POST /users
func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req UserRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.Unmarshal(body, &req)

	// Validate required fields
	// if req.Username == "" {
	// 	app.sendErrorResponse(w, "Username is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.Email == "" {
	// 	app.sendErrorResponse(w, "Email is required", http.StatusBadRequest)
	// 	return
	// }
	// if req.Password == "" {
	// 	app.sendErrorResponse(w, "Password is required", http.StatusBadRequest)
	// 	return
	// }

	// Create user in database
	err = app.UserModel.Create(req.Username, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			app.sendErrorResponse(w, "Username or email already exists", http.StatusConflict)
			return
		}
		app.sendErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := UserResponse{
		Username: req.Username,
		Email:    req.Email,
	}
	app.sendJSONResponse(w, response, http.StatusCreated)
}

// getUpdateDeleteUser handles:
// GET /users/{username} - Get user details
// PUT /users/{username} - Update user password
// DELETE /users/{username} - Delete user
func (app *application) getUpdateDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract username from URL path
	// URL pattern: /users/username
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	username := strings.TrimSpace(path)

	if username == "" {
		app.sendErrorResponse(w, "Username is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.getUser(w, r, username)
	case http.MethodPut:
		app.updateUser(w, r, username)
	case http.MethodDelete:
		app.deleteUser(w, r, username)
	default:
		app.sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getUser handles GET /users/{username}
func (app *application) getUser(w http.ResponseWriter, r *http.Request, username string) {
	user, err := app.UserModel.GetByUsername(username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			app.sendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}
		app.sendErrorResponse(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		Username: user.Username,
		Email:    user.Email,
	}
	app.sendJSONResponse(w, response, http.StatusOK)
}

// updateUser handles PUT /users/{username}
func (app *application) updateUser(w http.ResponseWriter, r *http.Request, username string) {
	var req UserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For update, only password is typically updated
	// You can extend this to update email as well
	if req.Password == "" {
		app.sendErrorResponse(w, "New password is required", http.StatusBadRequest)
		return
	}

	// Verify user exists before updating
	_, err = app.UserModel.GetByUsername(username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			app.sendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}
		app.sendErrorResponse(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	// Update password
	err = app.UserModel.UpdatePassword(username, req.Password)
	if err != nil {
		app.sendErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":  "User updated successfully",
		"username": username,
	}
	app.sendJSONResponse(w, response, http.StatusOK)
}

// deleteUser handles DELETE /users/{username}
func (app *application) deleteUser(w http.ResponseWriter, r *http.Request, username string) {
	err := app.UserModel.Delete(username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			app.sendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}
		app.sendErrorResponse(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":  "User deleted successfully",
		"username": username,
	}
	app.sendJSONResponse(w, response, http.StatusOK)
}

// Helper function to send JSON responses
func (app *application) sendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log the error but don't try to send another response
		app.logger.Printf("Error encoding JSON response: %v", err)
	}
}

// Helper function to send error responses
func (app *application) sendErrorResponse(w http.ResponseWriter, message string, status int) {
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	}
	app.sendJSONResponse(w, response, status)
}
