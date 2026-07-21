package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"galvin/lession05-exercise-user-management/internal/models"
	"galvin/lession05-exercise-user-management/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetAllUsers(search string, page, limit int) ([]models.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	if search != "" {
		keyword := strings.ToLower(search)
		filtered := []models.User{}
		for _, user := range users {
			if strings.Contains(strings.ToLower(user.Name), keyword) ||
				strings.Contains(strings.ToLower(user.Email), keyword) {
				filtered = append(filtered, user)
			}
		}
		users = filtered
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	if start >= len(users) {
		return []models.User{}, nil
	}
	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

func (s *userService) CreateUser(user models.User) (models.User, error) {
	if _, found := s.repo.FindByEmail(user.Email); found {
		return models.User{}, ErrEmailExists
	}

	if user.UUID == "" {
		uuid, err := newUUID()
		if err != nil {
			return models.User{}, err
		}
		user.UUID = uuid
	}

	if err := s.repo.Create(user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *userService) GetUserByUUID(uuid string) (models.User, error) {
	user, found := s.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUser(uuid string, user models.User) (models.User, error) {
	existing, found := s.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, ErrUserNotFound
	}

	if user.Email != "" && user.Email != existing.Email {
		if _, taken := s.repo.FindByEmail(user.Email); taken {
			return models.User{}, ErrEmailExists
		}
		existing.Email = user.Email
	}

	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Age > 0 {
		existing.Age = user.Age
	}
	if user.Password != "" {
		existing.Password = user.Password
	}
	if user.Status > 0 {
		existing.Status = user.Status
	}
	if user.Level > 0 {
		existing.Level = user.Level
	}

	if err := s.repo.Update(uuid, existing); err != nil {
		return models.User{}, err
	}

	return existing, nil
}

func (s *userService) DeleteUser(uuid string) error {
	if _, found := s.repo.FindByUUID(uuid); !found {
		return ErrUserNotFound
	}
	return s.repo.Delete(uuid)
}

// newUUID tạo UUID v4 bằng crypto/rand để khỏi thêm dependency mới.
func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10

	return strings.Join([]string{
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	}, "-"), nil
}
