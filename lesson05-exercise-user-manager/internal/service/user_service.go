package service

import (
	"strings"

	"galvin/lession05-exercise-user-management/internal/models"
	"galvin/lession05-exercise-user-management/internal/repository"
	"galvin/lession05-exercise-user-management/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers(search string, page, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternal)
	}

	var filteredUsers []models.User
	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)

			if strings.Contains(name, search) || strings.Contains(email, search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		filteredUsers = users
	}

	// --- Pagination algorithm ---
	//
	//   index:  0   1   2   3   4   5   6   7
	//          [A] [B] [C] [D] [E] [F] [G] [H]   len = 8
	//
	//   start = (page - 1) * limit
	//   end   = start + limit
	//
	//   With limit = 3:
	//     page 1 -> start 0, end 3  -> A B C
	//     page 2 -> start 3, end 6  -> D E F
	//     page 3 -> start 6, end 9  -> G H     (end is clamped to 8)
	//
	// Note: `start` is inclusive, `end` is exclusive, which is exactly how
	// a Go slice expression `s[start:end]` works.

	// Page is out of range (e.g. page 5 when we only have 8 items).
	// Return an empty page instead of panicking on the slice.
	start := (page - 1) * limit
	if start >= len(filteredUsers) {
		return []models.User{}, nil
	}

	// The last page is usually not full, so clamp `end` to the real length.
	// Without this, page 3 above would slice [6:9] and panic (out of range).
	end := start + limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	return filteredUsers[start:end], nil
}

func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}

	user.UUID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "faild to hash password", utils.ErrCodeInternal)
	}

	user.Password = string(hashedPassword)

	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "faild to create user", utils.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) GetUserByUUID(uuid string) (models.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	return user, nil
}

func (us *userService) UpdateUser(uuid string, updatedUser models.User) (models.User, error) {
	updatedUser.Email = utils.NormalizeString(updatedUser.Email)

	if u, exist := us.repo.FindByEmail(updatedUser.Email); exist && u.UUID != uuid {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}

	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}

	currentUser.Name = updatedUser.Name
	currentUser.Email = updatedUser.Email
	currentUser.Age = updatedUser.Age
	currentUser.Status = updatedUser.Status
	currentUser.Level = updatedUser.Level

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "faild to hash password", utils.ErrCodeInternal)
		}

		currentUser.Password = string(hashedPassword)
	}

	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "faild to update user", utils.ErrCodeInternal)
	}

	return currentUser, nil
}

func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.WrapError(err, "faild to hash password", utils.ErrCodeInternal)
	}

	return nil
}
