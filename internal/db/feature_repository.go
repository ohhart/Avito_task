package db

import (
	"database/sql"

	"Avito_task/internal/entity"
)

// FeatureRepository представляет репозиторий для работы с фичами в базе данных
type FeatureRepository struct {
	DB *sql.DB
}

// NewFeatureRepository создает новый экземпляр FeatureRepository
func NewFeatureRepository(db *sql.DB) *FeatureRepository {
	return &FeatureRepository{DB: db}
}

// Метод для создания новой фичи
func (fr *FeatureRepository) CreateFeature(feature *entity.Feature) error {
	_, err := fr.DB.Exec(`
        INSERT INTO features (name)
        VALUES ($1)
    `, feature.Name)
	if err != nil {
		return err
	}

	// Получение ID созданной фичи
	err = fr.DB.QueryRow("SELECT lastval()").Scan(&feature.ID)
	if err != nil {
		return err
	}

	return nil
}

// Метод для получения фичи по ID
func (fr *FeatureRepository) GetFeatureByID(id int) (*entity.Feature, error) {
	feature := &entity.Feature{}
	err := fr.DB.QueryRow(`
		SELECT id, name
		FROM features
		WHERE id = $1
	`, id).Scan(&feature.ID, &feature.Name)
	if err != nil {
		return nil, err
	}

	return feature, nil
}

// Метод для обновления фичи
func (fr *FeatureRepository) UpdateFeature(id int, newName string) error {
	_, err := fr.DB.Exec(`
		UPDATE features
		SET name = $1
		WHERE id = $2
	`, newName, id)
	if err != nil {
		return err
	}

	return nil
}

// Метод для удаления фичи по ID
func (fr *FeatureRepository) DeleteFeatureByID(id int) error {
	_, err := fr.DB.Exec(`
		DELETE FROM features
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	return nil
}
