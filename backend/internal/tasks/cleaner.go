package tasks

import (
	"context"
	"log"
	"os"
	"time"

	"vpj/internal/config"
	"vpj/internal/models"
)

// FileCleanerTask implements the Task interface for cleaning expired files
type FileCleanerTask struct {
	interval time.Duration
}

// NewFileCleanerTask creates a new file cleaner task
func NewFileCleanerTask() *FileCleanerTask {
	// Default interval: 1 hour, can be configured
	interval := 1 * time.Hour
	return &FileCleanerTask{
		interval: interval,
	}
}

// Name returns the task name
func (f *FileCleanerTask) Name() string {
	return "file-cleaner"
}

// Interval returns the task execution interval
func (f *FileCleanerTask) Interval() time.Duration {
	return f.interval
}

// Run executes the file cleanup task
func (f *FileCleanerTask) Run(ctx context.Context) error {
	log.Println("Starting file cleanup...")

	expiredTimePoint := time.Now().Add(-time.Duration(config.AppConfig.FileExpireTime) * time.Hour)

	var files []models.File
	if err := models.DB.Where("created_at < ?", expiredTimePoint).Find(&files).Error; err != nil {
		log.Printf("Error querying expired files: %v", err)
		return err
	}

	if len(files) == 0 {
		log.Println("No expired files to clean")
		return nil
	}

	deletedCount := 0
	errorCount := 0

	for _, file := range files {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			log.Println("File cleanup cancelled")
			return ctx.Err()
		default:
		}

		// Delete physical file
		if err := os.Remove(file.Path); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Error deleting file %s: %v", file.Path, err)
				errorCount++
			}
			// Continue even if file doesn't exist
		}

		// Delete database record
		if err := models.DB.Delete(&file).Error; err != nil {
			log.Printf("Error deleting file record %d: %v", file.ID, err)
			errorCount++
			continue
		}

		deletedCount++
	}

	log.Printf("File cleanup completed: %d files deleted, %d errors", deletedCount, errorCount)
	return nil
}
