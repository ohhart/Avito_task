package db

import (
	"database/sql"

	"Avito_task/internal/entity"
)

// BannerRepository представляет репозиторий для работы с баннерами в базе данных
type BannerRepository struct {
	DB *sql.DB
}

// NewBannerRepository создает новый экземпляр BannerRepository
func NewBannerRepository(db *sql.DB) *BannerRepository {
	return &BannerRepository{DB: db}
}

// CreateBanner создает новый баннер в базе данных
func (repo *BannerRepository) CreateBanner(banner *entity.Banner) error {
	_, err := repo.DB.Exec(`
        INSERT INTO banners (json_structure, feature_id, is_active)
        VALUES ($1, $2, $3)
    `, banner.JSONStructure, banner.FeatureID, banner.IsActive)
	if err != nil {
		return err
	}

	// Получение ID созданного баннера
	err = repo.DB.QueryRow("SELECT lastval()").Scan(&banner.ID)
	if err != nil {
		return err
	}

	// Добавление связей с тегами
	for _, tagID := range banner.TagIDs {
		_, err := repo.DB.Exec(`
            INSERT INTO banner_tags (banner_id, tag_id)
            VALUES ($1, $2)
        `, banner.ID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetBannerByID получает баннер из базы данных по его ID
func (repo *BannerRepository) GetBannerByID(id int, use_last_revision bool) (*entity.Banner, error) {
	banner := &entity.Banner{}
	var err error
	if use_last_revision {
		err = repo.DB.QueryRow(`
            SELECT id, json_structure, feature_id, is_active
            FROM banners
            WHERE id = $1
            AND created_at >= NOW() - interval '5 minutes'
            ORDER BY created_at DESC
            LIMIT 1
        `, id).Scan(&banner.ID, &banner.JSONStructure, &banner.FeatureID, &banner.IsActive)
	} else {
		err = repo.DB.QueryRow(`
            SELECT id, json_structure, feature_id, is_active
            FROM banners
            WHERE id = $1
        `, id).Scan(&banner.ID, &banner.JSONStructure, &banner.FeatureID, &banner.IsActive)
	}
	if err != nil {
		return nil, err
	}

	// Получение связанных тегов
	rows, err := repo.DB.Query(`
        SELECT tag_id
        FROM banner_tags
        WHERE banner_id = $1
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tagID int
		if err := rows.Scan(&tagID); err != nil {
			return nil, err
		}
		banner.TagIDs = append(banner.TagIDs, tagID)
	}

	return banner, nil
}

// UpdateBanner обновляет информацию о баннере в базе данных
func (repo *BannerRepository) UpdateBanner(banner *entity.Banner) error {
	_, err := repo.DB.Exec(`
        UPDATE banners
        SET json_structure = $1, feature_id = $2, is_active = $3
        WHERE id = $4
    `, banner.JSONStructure, banner.FeatureID, banner.IsActive, banner.ID)
	if err != nil {
		return err
	}

	// Удаление старых связей с тегами
	_, err = repo.DB.Exec(`
        DELETE FROM banner_tags
        WHERE banner_id = $1
    `, banner.ID)
	if err != nil {
		return err
	}

	// Добавление новых связей с тегами
	for _, tagID := range banner.TagIDs {
		_, err := repo.DB.Exec(`
            INSERT INTO banner_tags (banner_id, tag_id)
            VALUES ($1, $2)
        `, banner.ID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteBannerByID удаляет баннер из базы данных по его ID
func (repo *BannerRepository) DeleteBannerByID(id int) error {
	_, err := repo.DB.Exec(`
        DELETE FROM banners
        WHERE id = $1
    `, id)
	if err != nil {
		return err
	}

	// Удаление связей с тегами
	_, err = repo.DB.Exec(`
        DELETE FROM banner_tags
        WHERE banner_id = $1
    `, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBanners получает все баннеры с учетом фильтров по фиче, тегу, лимиту и оффсету
func (repo *BannerRepository) GetAllBanners(tagID, featureID, limit, offset int) ([]*entity.Banner, error) {
	// Формируем SQL-запрос с учетом фильтров, лимита и оффсета
	query := `
        SELECT id, json_structure, feature_id, is_active
        FROM banners
        WHERE 1=1
    `
	args := []interface{}{}

	// Добавляем условия фильтров, если они указаны
	if tagID != 0 {
		query += " AND id IN (SELECT banner_id FROM banner_tags WHERE tag_id = $1)"
		args = append(args, tagID)
	}
	if featureID != 0 {
		query += " AND feature_id = $2"
		args = append(args, featureID)
	}

	query += " ORDER BY id LIMIT $3 OFFSET $4"

	// Выполняем запрос и получаем результат
	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*entity.Banner
	for rows.Next() {
		banner := &entity.Banner{}
		if err := rows.Scan(&banner.ID, &banner.JSONStructure, &banner.FeatureID, &banner.IsActive); err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	return banners, nil
}
