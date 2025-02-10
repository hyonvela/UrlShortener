package handlers

import (
	"net/http"

	"example.com/m/internal/entity"
	"example.com/m/internal/usecase"
	"example.com/m/pkg/logging"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc  *usecase.Usecase
	log *logging.Logger
}

func NewHandler(uc *usecase.Usecase, logger *logging.Logger) *Handler {
	return &Handler{uc, logger}
}

func (h *Handler) ShortenUrl(c *gin.Context) {
	req := entity.LongUrl{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	str, err := h.uc.ShortenUrl(req.LongUrl, c.Request.Context())
	if err != nil {
		h.log.Errorf("an error occurred: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	answer := entity.ShortUrl{}
	answer.ShortUrl = str
	c.JSONP(http.StatusOK, answer)
}

func (h *Handler) GetLongUrl(c *gin.Context) {
	req := entity.ShortUrl{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var answer entity.LongUrl
	answer.LongUrl, err = h.uc.GetLongUrl(req.ShortUrl, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, answer)
}
