package handlers

import (
	"encoding/json" // Import encoding/json
	"file-vault/backend/internal/database"
	"file-vault/backend/internal/email"
	"file-vault/backend/internal/models"
	"fmt" // Import fmt for error handling
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/postgrest-go" // Import postgrest-go
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the registration of a new user
func RegisterUser(db *postgrest.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser struct {
			Username    string `json:"username" binding:"required"`
			Email       string `json:"email" binding:"required"`
			Password    string `json:"password" binding:"required"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			DateOfBirth string `json:"date_of_birth"` // Assuming YYYY-MM-DD format
			PhoneNumber string `json:"phone_number"`
		}

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user already exists
		var existingUsers []struct {
			FirstName     string `json:"first_name"`
			EmailVerified bool   `json:"email_verified"`
		}
		respBody, _, err := db.From("users").Select("first_name,email_verified", "exact", false).Filter("email", "eq", newUser.Email).Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check for existing user: %v", err)})
			return
		}
		if err := json.Unmarshal(respBody, &existingUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal existing user data: %v", err)})
			return
		}

		if len(existingUsers) > 0 {
			if existingUsers[0].EmailVerified {
				c.JSON(http.StatusConflict, gin.H{"error": "A verified user with this email already exists."})
				return
			}

			// User exists but is not verified, resend OTP
			otp := fmt.Sprintf("%06d", rand.Intn(1000000))
			otpExpiresAt := time.Now().Add(15 * time.Minute)
			updateData := map[string]interface{}{
				"otp":            otp,
				"otp_expires_at": otpExpiresAt,
			}
			_, _, updateErr := db.From("users").Update(updateData, "", "").Filter("email", "eq", newUser.Email).Execute()
			if updateErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update OTP for existing user: %v", updateErr)})
				return
			}
			if err := email.SendOTP(existingUsers[0].FirstName, newUser.Email, otp); err != nil {
				log.Printf("Failed to resend OTP to %s: %v", newUser.Email, err)
			}
			c.JSON(http.StatusOK, gin.H{"message": "User already exists. A new verification OTP has been sent to your email."})
			return
		}

		// User does not exist, proceed with creation
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Insert user into the 'users' table using postgrest-go
		user := map[string]interface{}{
			"username":      newUser.Username,
			"email":         newUser.Email,
			"password_hash": string(hashedPassword),
			"first_name":    newUser.FirstName,
			"last_name":     newUser.LastName,
			"phone_number":  newUser.PhoneNumber,
		}

		// Add date_of_birth only if it's provided
		if newUser.DateOfBirth != "" {
			user["date_of_birth"] = newUser.DateOfBirth
		}

		// Generate OTP
		otp := fmt.Sprintf("%06d", rand.Intn(1000000))
		otpExpiresAt := time.Now().Add(15 * time.Minute)

		user["otp"] = otp
		user["otp_expires_at"] = otpExpiresAt
		user["email_verified"] = false

		// The first return value is []byte (response body), second is count (int64), third is error
		_, _, insertDBErr := db.From("users").Insert(user, false, "", "", "").Execute() // Use a new variable for the DB error
		if insertDBErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %v", insertDBErr)})
			return
		}

		// Send OTP via email
		if err := email.SendOTP(newUser.FirstName, newUser.Email, otp); err != nil {
			log.Printf("Failed to send OTP to %s: %v", newUser.Email, err)
			// Note: In a real app, you might want to handle this more gracefully
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully. Please verify your email with the OTP sent."})
	}
}

// LoginUser handles user login and token generation
func LoginUser(db *postgrest.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Query user by email using postgrest-go
		var users []struct {
			UserID        string `json:"user_id"`
			Username      string `json:"username"`
			Email         string `json:"email"`
			PasswordHash  string `json:"password_hash"`
			EmailVerified bool   `json:"email_verified"`
			FirstName     string `json:"first_name"`
			LastName      string `json:"last_name"`
			IsAdmin       bool   `json:"is_admin"`
		}

		// The first return value is []byte (response body), second is count (int64), third is error
		respBody, _, err := db.From("users").Select("user_id,username,email,password_hash,email_verified,first_name,last_name,is_admin", "exact", false).Filter("email", "eq", credentials.Email).Execute() // Correctly assign all three return values
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve user: %v", err)})
			return
		}

		// Unmarshal the response body into the users slice
		if len(respBody) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if err := json.Unmarshal(respBody, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal user data: %v", err)})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		user := users[0]
		storedPassword := user.PasswordHash

		bcryptErr := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password))
		if bcryptErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if !user.EmailVerified {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before logging in."})
			return
		}

		// At this point, the password is correct. Generate a JWT token.
		// (Implementation of JWT generation will be added in the next step)
		c.JSON(http.StatusOK, gin.H{
			"message":    "Login successful",
			"user_id":    user.UserID,
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"is_admin":   user.IsAdmin,
		})
	}
}

// VerifyOTP handles the verification of a user's email using an OTP
func VerifyOTP(db *postgrest.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Email string `json:"email" binding:"required"`
			OTP   string `json:"otp" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var users []struct {
			OTP           string `json:"otp"`
			OTPExpiresAt  string `json:"otp_expires_at"`
			EmailVerified bool   `json:"email_verified"`
		}

		respBody, _, err := db.From("users").Select("otp,otp_expires_at,email_verified", "exact", false).Filter("email", "eq", payload.Email).Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve user: %v", err)})
			return
		}

		if err := json.Unmarshal(respBody, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal user data: %v", err)})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		user := users[0]

		if user.EmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified."})
			return
		}

		if user.OTP != payload.OTP {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
			return
		}

		// The time from Supabase might not have timezone info, so we parse it manually.
		// The layout must match the format returned by the database.
		const layout = "2006-01-02T15:04:05.999999"
		otpExpiresAt, parseErr := time.Parse(layout, user.OTPExpiresAt)
		if parseErr != nil {
			// Try parsing without fractional seconds as a fallback
			otpExpiresAt, parseErr = time.Parse("2006-01-02T15:04:05", user.OTPExpiresAt)
			if parseErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse expiration time: %v", parseErr)})
				return
			}
		}

		if time.Now().After(otpExpiresAt) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
			return
		}

		// Update user to mark as verified and clear OTP fields
		updateData := map[string]interface{}{
			"email_verified": true,
			"otp":            nil,
			"otp_expires_at": nil,
		}

		_, _, updateErr := db.From("users").Update(updateData, "", "").Filter("email", "eq", payload.Email).Execute()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update user: %v", updateErr)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
	}
}

// ResendOTP handles resending a new OTP to the user's email
func ResendOTP(db *postgrest.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Email string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if user is already verified
		var existingUsers []struct {
			FirstName     string `json:"first_name"`
			EmailVerified bool   `json:"email_verified"`
		}
		respBody, _, err := db.From("users").Select("first_name,email_verified", "exact", false).Filter("email", "eq", payload.Email).Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check user verification status: %v", err)})
			return
		}
		if err := json.Unmarshal(respBody, &existingUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal user data: %v", err)})
			return
		}

		if len(existingUsers) > 0 && existingUsers[0].EmailVerified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified."})
			return
		}

		// Generate new OTP
		otp := fmt.Sprintf("%06d", rand.Intn(1000000))
		otpExpiresAt := time.Now().Add(15 * time.Minute)

		updateData := map[string]interface{}{
			"otp":            otp,
			"otp_expires_at": otpExpiresAt,
		}

		_, _, updateErr := db.From("users").Update(updateData, "", "").Filter("email", "eq", payload.Email).Execute()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update OTP: %v", updateErr)})
			return
		}

		// Send new OTP via email
		if err := email.SendOTP(existingUsers[0].FirstName, payload.Email, otp); err != nil {
			log.Printf("Failed to resend OTP to %s: %v", payload.Email, err)
			// Note: In a real app, you might want to handle this more gracefully
		}

		c.JSON(http.StatusOK, gin.H{"message": "A new OTP has been sent to your email."})
	}
}

func GetUserQuota(clients *database.AppClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		userModel, ok := user.(models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user model in context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"rate_limit":     userModel.RateLimit,
			"storage_quota":  userModel.StorageQuota,
		})
	}
}

// UpdatePassword handles changing a user's password
func UpdatePassword(clients *database.AppClients) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			CurrentPassword string `json:"current_password" binding:"required"`
			NewPassword     string `json:"new_password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		userModel, ok := user.(models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user model in context"})
			return
		}

		// Retrieve the user's current password hash from the database
		var users []struct {
			PasswordHash string `json:"password_hash"`
		}
		respBody, _, err := clients.Postgrest.From("users").Select("password_hash", "exact", false).Filter("user_id", "eq", userModel.UserID).Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve user password hash: %v", err)})
			return
		}
		if err := json.Unmarshal(respBody, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal user password hash data: %v", err)})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		storedPasswordHash := users[0].PasswordHash

		// Verify the current password
		bcryptErr := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(payload.CurrentPassword))
		if bcryptErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
			return
		}

		// Hash the new password
		hashedNewPassword, hashErr := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}

		// Update the password hash in the database
		updateData := map[string]interface{}{
			"password_hash": string(hashedNewPassword),
		}

		_, _, updateErr := clients.Postgrest.From("users").Update(updateData, "", "").Filter("user_id", "eq", userModel.UserID).Execute()
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update password: %v", updateErr)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}
