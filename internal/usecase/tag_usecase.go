// tag_usecase.go
package usecase

import (
	"fmt"

	"Avito_task/internal/auth"
	"Avito_task/internal/db"
	"Avito_task/internal/entity"
)

// TagUseCase представляет интерфейс для работы с тегами
type TagUseCase struct {
	TagRepository db.TagRepository
	TokenService  *auth.TokenService
}

// NewTagUseCase создает новый экземпляр TagUseCase
func NewTagUseCase(tagRepo db.TagRepository, tokenService auth.TokenService) *TagUseCase {
	return &TagUseCase{
		TagRepository: tagRepo,
		TokenService:  &tokenService,
	}
}

// CreateTag создает новый тег
func (uc *TagUseCase) CreateTag(name string, token string) (*entity.Tag, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", err)
	}

	newTag := &entity.Tag{Name: name}
	if err := uc.TagRepository.CreateTag(newTag); err != nil {
		return nil, fmt.Errorf("ошибка при создании нового тега: %w", err)
	}

	return newTag, nil
}

// UpdateTag обновляет информацию о теге
func (uc *TagUseCase) UpdateTag(id int, newName string, token string) (*entity.Tag, error) {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %w", err)
	}

	tag, err := uc.TagRepository.GetTagByID(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении тега: %w", err)
	}

	tag.Name = newName
	if err := uc.TagRepository.UpdateTag(id, newName); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении информации о теге: %w", err)
	}

	return tag, nil
}

// DeleteTag удаляет тег по ID
func (uc *TagUseCase) DeleteTag(id int, token string) error {
	// Проверка токена администратора
	if err := uc.TokenService.VerifyAdminToken(token); err != nil {
		return fmt.Errorf("ошибка авторизации: %w", err)
	}

	if err := uc.TagRepository.DeleteTagByID(id); err != nil {
		return fmt.Errorf("ошибка при удалении тега: %w", err)
	}

	return nil
}
