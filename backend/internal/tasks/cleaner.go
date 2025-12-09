package tasks

import (
	"log"
	"os"
	"time"

	"vpj/internal/config"
	"vpj/internal/models"
)

type Cleaner struct {
	ticker *time.Ticker
	done   chan bool
}

func NewCleaner() *Cleaner {
	return &Cleaner{
		done: make(chan bool),
	}
}

func (c *Cleaner) Start() {
	// Run immediately on start
	c.clean()

	// Then run every hour
	c.ticker = time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.clean()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *Cleaner) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.done <- true
}

func (c *Cleaner) clean() {
	log.Println("Starting file cleanup...")

	expiredTimePoint := time.Now().Add(-time.Duration(config.AppConfig.FileExpireTime) * time.Hour)

	var files []models.File
	if err := models.DB.Where("created_at < ?", expiredTimePoint).Find(&files).Error; err != nil {
		log.Printf("Error querying expired files: %v", err)
		return
	}

	if len(files) == 0 {
		log.Println("No expired files to clean")
		return
	}

	deletedCount := 0
	for _, file := range files {
		// Delete physical file
		if err := os.Remove(file.Path); err != nil {
			log.Printf("Error deleting file %s: %v", file.Path, err)
			continue
		}

		// Delete database record
		if err := models.DB.Delete(&file).Error; err != nil {
			log.Printf("Error deleting file record %d: %v", file.ID, err)
			continue
		}

		deletedCount++
	}

	log.Printf("Cleaned %d expired files", deletedCount)
}

