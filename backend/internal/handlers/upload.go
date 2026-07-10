package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/ledger/backend/pkg/response"
)

type UploadHandler struct {
	db *gorm.DB
}

func NewUploadHandler(db *gorm.DB) *UploadHandler {
	return &UploadHandler{db: db}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, 401, "Unauthorized")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "No file uploaded")
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		response.Error(c, 400, "Invalid file type. Only images are allowed")
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10<<20 {
		response.Error(c, 400, "File too large. Max 10MB")
		return
	}

	// Create upload directory
	userDir := filepath.Join("uploads", userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		response.Error(c, 500, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(userDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		response.Error(c, 500, "Failed to save file")
		return
	}

	// Return URL (in production, this would be a CDN URL)
	url := fmt.Sprintf("/api/v1/uploads/%s", filename)
	response.Success(c, UploadResponse{
		URL:      url,
		Filename: file.Filename,
		Size:     file.Size,
	})
}

// CashFlowHandler handles cash flow projections
