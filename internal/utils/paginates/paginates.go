package paginates

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func SongTextPaginate(src string, page int, limit int) []string {
	parts := strings.Split(src, "\n\n")

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex > len(parts) {
		startIndex = len(parts)
	}
	if endIndex > len(parts) {
		endIndex = len(parts)
	}

	return parts[startIndex:endIndex]
}

func SongPaginate(page string, limit string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(page)
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(limit)
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
