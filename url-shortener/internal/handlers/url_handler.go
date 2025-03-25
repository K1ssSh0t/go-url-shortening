package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"url-shortener/internal/service"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid URL",
			"details": err.Error(),
		})
		return
	}

	url, err := h.urlService.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create short URL",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, url)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.urlService.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve URL",
			"details": err.Error(),
		})
		return
	}

	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url.OriginalURL)
}

func (h *URLHandler) UpdateShortURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var req struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid URL",
			"details": err.Error(),
		})
		return
	}

	updatedURL, err := h.urlService.UpdateShortURL(c.Request.Context(), shortCode, req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update URL",
			"details": err.Error(),
		})
		return
	}

	if updatedURL == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	c.JSON(http.StatusOK, updatedURL)
}

func (h *URLHandler) DeleteShortURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	err := h.urlService.DeleteShortURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete URL",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *URLHandler) GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.urlService.GetURLStats(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve URL stats",
			"details": err.Error(),
		})
		return
	}

	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	c.JSON(http.StatusOK, url)
}
