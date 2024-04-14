package db

import (
	"Avito_task/internal/entity"
	"database/sql"
)

// UserRepository представляет репозиторий для работы с сущностью пользователя в базе данных.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

// CreateUser создает нового пользователя в базе данных.
func (ur *UserRepository) CreateUser(user *entity.User) error {
	_, err := ur.db.Exec(`
        INSERT INTO users (username, password, token, is_admin)
        VALUES ($1, $2, $3, $4)
    `, user.Username, user.PasswordHash, user.Token, user.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID возвращает пользователя из базы данных по его ID.
func (ur *UserRepository) GetUserByID(id int) (*entity.User, error) {
	user := &entity.User{}
	err := ur.db.QueryRow(`
        SELECT id, username, password, token, is_admin
        FROM users
        WHERE id = $1
    `, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Token, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername получает пользователя из базы данных по его имени пользователя (Username)
func (ur *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	err := ur.db.QueryRow(`
        SELECT id, username, password, token, is_admin
        FROM users
        WHERE username = $1
    `, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Token, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе в базе данных.
func (ur *UserRepository) UpdateUser(user *entity.User) error {
	_, err := ur.db.Exec(`
        UPDATE users
        SET username = $1, password = $2, token = $3, is_admin = $4
        WHERE id = $5
    `, user.Username, user.PasswordHash, user.Token, user.IsAdmin, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserByID удаляет пользователя из базы данных по его ID.
func (ur *UserRepository) DeleteUserByID(id int) error {
	_, err := ur.db.Exec(`
        DELETE FROM users
        WHERE id = $1
    `, id)
	if err != nil {
		return err
	}
	return nil
}
