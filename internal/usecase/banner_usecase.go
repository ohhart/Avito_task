package usecase

import (
	"errors"
	"fmt"

	"Avito_task/internal/auth"
	"Avito_task/internal/db"
	"Avito_task/internal/entity"
)

var (
	ErrUnauthorized   = errors.New("неавторизованный запрос")
	ErrInvalidParams  = errors.New("неверные параметры")
	ErrBannerNotFound = errors.New("баннер не найден")
	ErrCreateBanner   = errors.New("ошибка при создании баннера")
	ErrUpdateBanner   = errors.New("ошибка при обновлении баннера")
	ErrDeleteBanner   = errors.New("ошибка при удалении баннера")
)

// BannerUseCase представляет интерфейс для работы с баннерами
type BannerUseCase struct {
	BannerRepository db.BannerRepository
	TokenService     *auth.TokenService
}

// NewBannerUseCase создает новый экземпляр BannerUseCase
func NewBannerUseCase(bannerRepo db.BannerRepository, tokenService auth.TokenService) *BannerUseCase {
	return &BannerUseCase{
		BannerRepository: bannerRepo,
		TokenService:     &tokenService,
	}
}

// GetUserBanner получает баннер для пользователя по заданным параметрам
func (uc *BannerUseCase) GetUserBanner(id int, token string) (*entity.Banner, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", ErrUnauthorized)
	}

	// Получаем баннер из репозитория с флагом useLastRevision
	banner, err := uc.BannerRepository.GetBannerByID(id, true) // Используем флаг true для получения самой актуальной информации
	if err != nil {
		return nil, fmt.Errorf("%w", ErrBannerNotFound)
	}

	return banner, nil
}

// GetAllBanners получает все баннеры с учетом фильтров по фиче, тегу, лимиту и оффсету
func (uc *BannerUseCase) GetAllBanners(tagID, featureID, limit, offset int, token string) ([]*entity.Banner, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", ErrUnauthorized)
	}

	// Проверяем, что limit и offset не равны нулю
	if limit == 0 || offset < 0 {
		return nil, fmt.Errorf("%w", ErrInvalidParams)
	}

	// Получаем все баннеры из репозитория с учетом фильтров
	banners, err := uc.BannerRepository.GetAllBanners(tagID, featureID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrBannerNotFound)
	}

	return banners, nil
}

// CreateBanner создает новый баннер
func (uc *BannerUseCase) CreateBanner(tagIDs []int, featureID int, content map[string]interface{}, isActive bool, token string) (*entity.Banner, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", ErrUnauthorized)
	}

	// Проверяем, что tagIDs не пустой и featureID не равен нулю
	if len(tagIDs) == 0 || featureID == 0 {
		return nil, fmt.Errorf("%w", ErrInvalidParams)
	}

	// Преобразуем содержимое баннера в формат JSON
	jsonStructure, err := entity.MapToJSON(content)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования JSON: %w", err)
	}

	// Создаем новый баннер
	newBanner := &entity.Banner{
		JSONStructure: jsonStructure,
		FeatureID:     featureID,
		TagIDs:        tagIDs,
		IsActive:      isActive,
	}

	err = uc.BannerRepository.CreateBanner(newBanner)
	if err != nil {
		// Возвращаем ошибку с сообщением об ошибке при создании баннера
		return nil, fmt.Errorf("ошибка при создании баннера: %w", ErrCreateBanner)
	}

	return newBanner, nil
}

// UpdateBanner обновляет информацию о баннере
func (uc *BannerUseCase) UpdateBanner(id int, tagIDs []int, featureID int, content map[string]interface{}, isActive bool, token string) (*entity.Banner, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", ErrUnauthorized)
	}

	// Проверяем, что tagIDs не пустой и featureID не равен нулю
	if len(tagIDs) == 0 || featureID == 0 {
		return nil, fmt.Errorf("%w", ErrInvalidParams)
	}

	// Преобразуем содержимое баннера в формат JSON
	jsonStructure, err := entity.MapToJSON(content)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования JSON: %w", err)
	}

	// Обновляем баннер в репозитории
	updatedBanner := &entity.Banner{
		ID:            id,
		JSONStructure: jsonStructure,
		FeatureID:     featureID,
		TagIDs:        tagIDs,
		IsActive:      isActive,
	}

	err = uc.BannerRepository.UpdateBanner(updatedBanner)
	if err != nil {
		// Возвращаем ошибку с сообщением об ошибке при обновлении баннера
		return nil, fmt.Errorf("ошибка при обновлении баннера: %w", ErrUpdateBanner)
	}

	return updatedBanner, nil
}

// DeleteBanner удаляет баннер по его ID
func (uc *BannerUseCase) DeleteBanner(id int, token string) error {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return fmt.Errorf("ошибка авторизации: %w", ErrUnauthorized)
	}

	err := uc.BannerRepository.DeleteBannerByID(id)
	if err != nil {
		// Возвращаем ошибку с сообщением об ошибке при удалении баннера
		return fmt.Errorf("ошибка при удалении баннера: %w", ErrDeleteBanner)
	}

	return nil
}
