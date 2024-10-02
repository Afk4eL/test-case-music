package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"test-case/internal/models"
	"test-case/storage/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SongHandler struct {
	repo repos.SongRepository
}

func NewSongHandler(repos repos.SongRepository) SongHandler {
	return SongHandler{repo: repos}
}

// GetSongs godoc
//
// @Summary Get songs
// @Description Retrieve the list of all songs with pagination and filtering
// @Tags songs
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit of songs per page"
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song name"
// @Success 200 {array} models.Song
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H "Song doesn't exist"
// @Router /get-songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	filterParams := make(map[string]string)

	for key, values := range queryParams {
		if len(values) > 0 {
			filterParams[key] = values[0]
		}
	}

	result, err := h.repo.GetSongs(filterParams, c.Query("page"), c.Query("limit"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"Error": "Song doesnt exist"})
			return
		}
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// GetSongText godoc
//
// @Summary Get song text
// @Description Retrieve the text of a song by its ID with pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param songId query string true "Song ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit of couplets per page"
// @Success 200 {object} map[string]string "Paginated song couplets"
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H "Song doesn't exist"
// @Router /get-song-text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	result, err := h.repo.GetSongText(c.Query("songId"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"Error": "Song doesnt exist"})
			return
		}
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	parts := strings.Split(result, "\n\n")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex > len(parts) {
		startIndex = len(parts)
	}
	if endIndex > len(parts) {
		endIndex = len(parts)
	}

	paginatedParts := parts[startIndex:endIndex]

	type Response struct {
		Text []map[string]string `json:"song_text"`
	}

	var response Response

	for _, part := range paginatedParts {
		response.Text = append(response.Text, map[string]string{
			"couplet": part,
		})
	}

	c.JSON(200, response)
}

// DeleteSong godoc
//
// @Summary Delete a song
// @Description Delete a song by its ID
// @Tags songs
// @Accept json
// @Produce json
// @Param songId query string true "Song ID"
// @Success 200 {object} gin.H "OK: Song deleted"
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H "Song doesn't exist"
// @Router /delete-song [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	if err := h.repo.DeleteSong(c.Query("songId")); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"Error": "Song doesnt exist"})
			return
		}
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"OK": "Song deleted"})
}

// UpdateSong godoc
//
// @Summary Update a song
// @Description Update an existing song by providing updated details
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Updated song object"
// @Success 200 {object} gin.H "OK: Song updated"
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H "Song doesn't exist"
// @Router /update-song [post]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	var updatedSong models.Song
	if err := c.BindJSON(&updatedSong); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := h.repo.UpdateSong(updatedSong); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"Error": "Song doesnt exist"})
			return
		}
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"OK": "Song updated"})
}

// AddSong godoc
//
// @Summary Add a new song
// @Description Add a new song with details fetched from an external API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "New song object"
// @Success 200 {object} gin.H "OK: Song created, New song ID"
// @Failure 400 {object} gin.H
// @Router /add-song [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var newSong models.Song
	if err := c.BindJSON(&newSong); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	baseURL, _ = url.JoinPath(baseURL, "/info")

	params := url.Values{}
	params.Add("group", newSong.Group)
	params.Add("song", newSong.Song)

	apiURL, err := url.Parse(baseURL)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	apiURL.RawQuery = params.Encode()

	response, err := http.Get(apiURL.String())
	if err != nil {
		c.JSON(response.StatusCode, gin.H{"Error": err.Error()})
		return
	}

	type SongDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	var songDetail SongDetail
	if err := json.NewDecoder(response.Body).Decode(&songDetail); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	newSong.ReleaseDate = songDetail.ReleaseDate
	newSong.Text = songDetail.Text
	newSong.Link = songDetail.Link

	id, err := h.repo.AddSong(newSong)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"OK": "Song created", "New song Id": id})
}
