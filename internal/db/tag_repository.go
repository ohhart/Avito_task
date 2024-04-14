package db

import (
	"database/sql"

	"Avito_task/internal/entity"
)

// TagRepository представляет репозиторий для работы с тегами в базе данных
type TagRepository struct {
	DB *sql.DB
}

// NewTagRepository создает новый экземпляр TagRepository
func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{DB: db}
}

// CreateTag создает новый тег в базе данных
func (repo *TagRepository) CreateTag(tag *entity.Tag) error {
	_, err := repo.DB.Exec(`
        INSERT INTO tags (name)
        VALUES ($1)
    `, tag.Name)
	if err != nil {
		return err
	}

	// Получение ID созданного тега
	err = repo.DB.QueryRow("SELECT lastval()").Scan(&tag.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetTagByID получает тег из базы данных по его ID
func (repo *TagRepository) GetTagByID(id int) (*entity.Tag, error) {
	tag := &entity.Tag{}
	err := repo.DB.QueryRow(`
		SELECT id, name
		FROM tags
		WHERE id = $1	
	`, id).Scan(&tag.ID, &tag.Name)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag обновляет информацию о теге в базе данных
func (repo *TagRepository) UpdateTag(id int, newName string) error {
	_, err := repo.DB.Exec(`
		UPDATE tags
		SET name = $1
		WHERE id = $2
	`, newName, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTagByID удаляет тег из базы данных по его ID
func (repo *TagRepository) DeleteTagByID(id int) error {
	_, err := repo.DB.Exec(`
		DELETE FROM tags
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	return nil
}
