package usecase

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"Avito_task/internal/auth"
	"Avito_task/internal/db"
	"Avito_task/internal/entity"
)

// UserUseCase представляет интерфейс для работы с пользователями
type UserUseCase struct {
	UserRepository db.UserRepository
	TokenService   *auth.TokenService
}

// NewUserUseCase создает новый экземпляр UserUseCase
func NewUserUseCase(userRepo db.UserRepository, tokenService *auth.TokenService) *UserUseCase {
	return &UserUseCase{
		UserRepository: userRepo,
		TokenService:   tokenService,
	}
}

// HashPassword хеширует пароль пользователя
func (uc *UserUseCase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash проверяет соответствие хеша пароля
func (uc *UserUseCase) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// RegisterUser регистрирует нового пользователя
func (uc *UserUseCase) RegisterUser(username, password, role string) (*entity.User, error) {
	// Hash the password before storing it
	hashedPassword, err := uc.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	// Save the user to the database
	err = uc.UserRepository.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// AuthenticateUser аутентифицирует пользователя (логин)
func (uc *UserUseCase) AuthenticateUser(username, password string) (*entity.User, error) {
	user, err := uc.UserRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Check if the password matches the stored hash
	if !uc.CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("неправильный пароль")
	}

	return user, nil
}

// GetUserByID получает информацию о пользователе по его ID
func (uc *UserUseCase) GetUserByID(id int) (*entity.User, error) {
	user, err := uc.UserRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (uc *UserUseCase) UpdateUser(id int, username, password, role string) (*entity.User, error) {
	// Получите пользователя из базы данных
	user, err := uc.UserRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Обновите данные пользователя
	user.Username = username
	user.Role = role
	// Хешируйте новый пароль, если он был предоставлен
	if password != "" {
		hashedPassword, err := uc.HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashedPassword
	}

	// Вызовите метод UpdateUser из UserRepository
	err = uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUserByID удаляет пользователя по его ID
func (uc *UserUseCase) DeleteUserByID(id int) error {
	err := uc.UserRepository.DeleteUserByID(id)
	if err != nil {
		return err
	}

	return nil
}
