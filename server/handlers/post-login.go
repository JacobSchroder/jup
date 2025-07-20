package handlers

import (
	"log/slog"
	"net/http"
)

func HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse login form", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "on"

	// Basic validation
	if email == "" || password == "" {
		slog.Warn("Login attempt with missing credentials")
		// In a real app, you'd render the form with error messages
		http.Redirect(w, r, "/login?error=missing_credentials", http.StatusSeeOther)
		return
	}

	// TODO: Implement actual authentication logic
	// For now, we'll just log the attempt and redirect
	slog.Info("Login attempt", "email", email, "remember", remember)

	// Simulate authentication - in real app, verify against database
	if email == "admin@example.com" && password == "password" {
		// TODO: Set session/JWT token
		slog.Info("Successful login", "email", email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Authentication failed
	slog.Warn("Failed login attempt", "email", email)
	http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
}
