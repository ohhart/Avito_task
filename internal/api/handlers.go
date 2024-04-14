package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Avito_task/internal/entity"
	"Avito_task/internal/usecase"
)

// BannerHandlers представляет обработчики запросов для баннеров
type BannerHandlers struct {
	BannerUseCase *usecase.BannerUseCase
}

// NewBannerHandlers создает новый экземпляр BannerHandlers
func NewBannerHandlers(bannerUseCase *usecase.BannerUseCase) *BannerHandlers {
	return &BannerHandlers{
		BannerUseCase: bannerUseCase,
	}
}

// GetUserBannerHandler обработчик для получения баннера пользователя по ID
func (h *BannerHandlers) GetUserBannerHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	token := c.GetHeader("Authorization")
	banner, err := h.BannerUseCase.GetUserBanner(id, token)
	if err != nil {
		switch err {
		case usecase.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		case usecase.ErrBannerNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Баннер для пользователя не найден"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, banner)
}

// GetAllBannersHandler обработчик для получения всех баннеров с учетом фильтров
func (h *BannerHandlers) GetAllBannersHandler(c *gin.Context) {
	tagIDStr := c.Query("tag_id")
	featureIDStr := c.Query("feature_id")
	limitStr := c.DefaultQuery("limit", "0")
	offsetStr := c.DefaultQuery("offset", "0")

	tagID, _ := strconv.Atoi(tagIDStr)
	featureID, _ := strconv.Atoi(featureIDStr)
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	token := c.GetHeader("Authorization")
	banners, err := h.BannerUseCase.GetAllBanners(tagID, featureID, limit, offset, token)
	if err != nil {
		switch err {
		case usecase.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, banners)
}

// CreateBanner обработчик для создания нового баннера
func (h *BannerHandlers) CreateBanner(c *gin.Context) {
	// Получаем параметры из контекста Gin
	var req entity.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	token := c.GetHeader("Authorization")

	// Вызываем метод usecase для создания нового баннера
	newBanner, err := h.BannerUseCase.CreateBanner(req.TagIDs, req.FeatureID, req.Content, req.IsActive, token)
	if err != nil {
		// Обработка ошибок и отправка соответствующих HTTP-ответов
		switch err {
		case usecase.ErrInvalidParams:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		case usecase.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	// Отправляем созданный баннер в качестве ответа
	c.JSON(http.StatusCreated, newBanner)
}

// DeleteBannerHandler обработчик для удаления баннера по его ID
func (h *BannerHandlers) DeleteBannerHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	token := c.GetHeader("Authorization")
	err = h.BannerUseCase.DeleteBanner(id, token)
	if err != nil {
		switch err {
		case usecase.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		case usecase.ErrBannerNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Баннер не найден"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
