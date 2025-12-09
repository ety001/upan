package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"vpj/internal/config"
	"vpj/internal/models"

	"github.com/gin-gonic/gin"
)

type MainHandler struct{}

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}

// GetConfig returns the configuration for the frontend
func (h *MainHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"file_max_size":   config.AppConfig.FileMaxSize / (1024 * 1024), // Convert to MB
		"file_expire_time": config.AppConfig.FileExpireTime,
	})
}

// Upload handles file upload
func (h *MainHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("o")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "请选择要上传的文件!",
		})
		return
	}

	// Check file size
	if file.Size > config.AppConfig.FileMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  fmt.Sprintf("文件大小限制: %dMB", config.AppConfig.FileMaxSize/(1024*1024)),
		})
		return
	}

	// Generate unique code
	code := h.generateUniqueCode()
	if code == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "生成提取码失败",
		})
		return
	}

	// Save file
	filename := file.Filename
	ext := filepath.Ext(filename)
	savePath := filepath.Join(config.AppConfig.StoragePath, code+ext)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "保存文件失败",
		})
		return
	}

	// Get MIME type
	mime := file.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}

	// Save to database
	fileRecord := models.File{
		Code:     code,
		Path:     savePath,
		Filename: filename,
		Mime:     mime,
		Size:     file.Size,
	}

	if err := models.DB.Create(&fileRecord).Error; err != nil {
		// Clean up file if database save fails
		os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "保存文件记录失败",
		})
		return
	}

	// Calculate expired_at
	expiredAt := time.Now().Add(time.Duration(config.AppConfig.FileExpireTime) * time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"code":      code,
		"expired_at": expiredAt.Unix(),
	})
}

// GetFile handles file download
func (h *MainHandler) GetFile(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "提取码不能为空",
		})
		return
	}

	var fileRecord models.File
	if err := models.DB.Where("code = ?", code).First(&fileRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "提取码不存在",
		})
		return
	}

	// Check if file exists
	if _, err := os.Stat(fileRecord.Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文件不存在",
		})
		return
	}

	// Open file
	file, err := os.Open(fileRecord.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "打开文件失败",
		})
		return
	}
	defer file.Close()

	// Set headers
	c.Header("Pragma", "public")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Cache-Control", "private")
	c.Header("Content-Type", "application/force-download")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileRecord.Filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Connection", "close")

	// Stream file
	c.DataFromReader(http.StatusOK, fileRecord.Size, fileRecord.Mime, file, nil)
}

// generateUniqueCode generates a unique 6-character code
func (h *MainHandler) generateUniqueCode() string {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		bytes := make([]byte, 3)
		if _, err := rand.Read(bytes); err != nil {
			continue
		}
		code := hex.EncodeToString(bytes)

		// Check if code exists
		var count int64
		models.DB.Model(&models.File{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			return code
		}
	}
	return ""
}

