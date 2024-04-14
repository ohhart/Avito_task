package usecase

import (
	"fmt"

	"Avito_task/internal/auth"
	"Avito_task/internal/db"
	"Avito_task/internal/entity"
)

// FeatureUseCase представляет интерфейс для работы с фичами
type FeatureUseCase struct {
	FeatureRepository *db.FeatureRepository
	TokenService      *auth.TokenService
}

// NewFeatureUseCase создает новый экземпляр FeatureUseCase
func NewFeatureUseCase(featureRepo *db.FeatureRepository, tokenService *auth.TokenService) *FeatureUseCase {
	return &FeatureUseCase{
		FeatureRepository: featureRepo,
		TokenService:      tokenService,
	}
}

// CreateFeature создает новую фичу
func (uc *FeatureUseCase) CreateFeature(name string, token string) (*entity.Feature, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", err)
	}

	newFeature := &entity.Feature{Name: name}
	if err := uc.FeatureRepository.CreateFeature(newFeature); err != nil {
		return nil, fmt.Errorf("ошибка при создании новой фичи: %w", err)
	}

	return newFeature, nil
}

// UpdateFeature обновляет информацию о фиче
func (uc *FeatureUseCase) UpdateFeature(id int, newName string, token string) (*entity.Feature, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", err)
	}

	if err := uc.FeatureRepository.UpdateFeature(id, newName); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении информации о фичи: %w", err)
	}

	feature, err := uc.FeatureRepository.GetFeatureByID(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении фичи: %w", err)
	}

	return feature, nil
}

// DeleteFeature удаляет фичу по ID
func (uc *FeatureUseCase) DeleteFeature(id int, token string) error {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	if err := uc.FeatureRepository.DeleteFeatureByID(id); err != nil {
		return fmt.Errorf("ошибка при удалении фичи: %w", err)
	}

	return nil
}
