package repos

import (
	"strconv"
	"test-case/internal/models"

	"gorm.io/gorm"
)

type SongRepository interface {
	GetSongs(filterParams map[string]string, page string, limit string) ([]models.Song, error)
	GetSongText(id string) (string, error)
	DeleteSong(id string) error
	UpdateSong(updatedSong models.Song) error
	AddSong(newSong models.Song) (uint, error)
}

type songRepo struct {
	database *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepo{database: db}
}

func paginate(page string, limit string) func(db *gorm.DB) *gorm.DB {
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

func (r *songRepo) GetSongs(filterParams map[string]string, page string, limit string) ([]models.Song, error) {
	const op = "storage.repos.GetSongs"

	query := r.database

	if filterParams["band"] != "" {
		query = query.Where("band = ?", filterParams["band"])
	}
	if filterParams["song"] != "" {
		query = query.Where("song = ?", filterParams["song"])
	}
	if filterParams["releaseDate"] != "" {
		query = query.Where("release_date = ?", filterParams["releaseDate"])
	}
	if filterParams["text"] != "" {
		query = query.Where("text LIKE ?", "%"+filterParams["text"]+"%")
	}
	if filterParams["link"] != "" {
		query = query.Where("link = ?", filterParams["link"])
	}

	var songs []models.Song
	result := query.Order("id asc").Scopes(paginate(page, limit)).Find(&songs)
	if result.Error != nil {
		return nil, result.Error
	}

	return songs, nil
}

func (r *songRepo) GetSongText(id string) (string, error) {
	const op = "storage.repos.GetSongText"

	songId, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}

	var text string
	result := r.database.Model(&models.Song{}).
		Select("text").Where("id = ?", songId).Scan(&text)
	if result.Error != nil {
		return "", result.Error
	}

	return text, nil
}

func (r *songRepo) DeleteSong(id string) error {
	const op = "storage.repos.DeleteSong"

	songId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	result := r.database.Where(&models.Song{Id: uint(songId)}).Delete(&models.Song{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *songRepo) UpdateSong(updatedSong models.Song) error {
	const op = "storage.repos.UpdateSong"

	var oldSong models.Song
	result := r.database.Where("id = ?", updatedSong.Id).First(&oldSong)
	if result.Error != nil {
		return result.Error
	}

	oldSong.Song = updatedSong.Song
	oldSong.Group = updatedSong.Group
	oldSong.Text = updatedSong.Text
	oldSong.ReleaseDate = updatedSong.ReleaseDate
	oldSong.Link = updatedSong.Link

	r.database.Save(&oldSong)
	return nil
}

func (r *songRepo) AddSong(newSong models.Song) (uint, error) {
	const op = "storage.repos.AddSong"

	var maxId uint
	r.database.Model(&models.Song{}).
		Select("MAX(id)").Scan(&maxId)

	newSong.Id = maxId + 1

	result := r.database.Create(&newSong)
	if result.Error != nil {
		return 0, result.Error
	}

	return newSong.Id, nil
}
