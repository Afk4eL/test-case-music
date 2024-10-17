package repos

import (
	"errors"
	"strconv"
	"test-case/internal/models"
	"test-case/internal/utils/logger"
	"test-case/internal/utils/paginates"

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
	result := query.Order("id asc").Scopes(paginates.SongPaginate(page, limit)).Find(&songs)
	if result.Error != nil {
		logger.Logger.Info().Interface("Error occured: ", result.Error).Msg(op)
		return nil, result.Error
	}

	return songs, nil
}

func (r *songRepo) GetSongText(id string) (string, error) {
	const op = "storage.repos.GetSongText"

	songId, err := strconv.Atoi(id)
	if err != nil {
		logger.Logger.Info().Interface("Error occured: ", err).Msg(op)
		return "", err
	}

	var text string
	result := r.database.Model(&models.Song{}).
		Select("text").Where("id = ?", songId).Scan(&text)
	if result.Error != nil {
		logger.Logger.Info().Interface("Error occured: ", result.Error).Msg(op)
		return "", result.Error
	}

	return text, nil
}

func (r *songRepo) DeleteSong(id string) error {
	const op = "storage.repos.DeleteSong"

	songId, err := strconv.Atoi(id)
	if err != nil {
		logger.Logger.Info().Interface("Error occured: ", err).Msg(op)
		return err
	}

	result := r.database.Where(&models.Song{Id: uint(songId)}).Delete(&models.Song{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		logger.Logger.Info().Interface("Error occured: ", result.Error).Msg(op)
		return result.Error
	}

	return nil
}

func (r *songRepo) UpdateSong(updatedSong models.Song) error {
	const op = "storage.repos.UpdateSong"

	var oldSong models.Song
	result := r.database.Where("id = ?", updatedSong.Id).First(&oldSong)
	if result.Error != nil {
		logger.Logger.Info().Interface("Error occured: ", result.Error).Msg(op)
		return result.Error
	}

	oldSong.Song = updatedSong.Song
	oldSong.Group = updatedSong.Group
	oldSong.Text = updatedSong.Text
	oldSong.ReleaseDate = updatedSong.ReleaseDate
	oldSong.Link = updatedSong.Link

	if result := r.database.Save(&oldSong); result.Error != nil {
		logger.Logger.Info().Interface("Error occured: ", result.Error).Msg(op)
		return result.Error
	}

	return nil
}

func (r *songRepo) AddSong(newSong models.Song) (uint, error) {
	const op = "storage.repos.AddSong"

	var group models.Group
	result := r.database.Where("name = ?", newSong.Band).First(&group)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Logger.Info().Str("Group not found", "Creating...").Msg(op)

			group = models.Group{Name: newSong.Band}
			result := r.database.Create(&group)
			if result.Error != nil {
				logger.Logger.Error().Interface("Error: ", result.Error).Msg(op)
				return 0, result.Error
			}
		} else {
			logger.Logger.Error().Interface("Error: ", result.Error).Msg(op)
			return 0, result.Error
		}
	}

	var maxId uint
	r.database.Model(&models.Song{}).
		Select("MAX(id)").Scan(&maxId)

	newSong.Id = maxId + 1
	newSong.GroupId = group.Id

	result = r.database.Create(&newSong)
	if result.Error != nil {
		logger.Logger.Info().Interface("Error occured: ", result.Error).Msg(op)
		return 0, result.Error
	}

	return newSong.Id, nil
}
