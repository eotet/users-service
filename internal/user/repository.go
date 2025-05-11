package user

import (
	"errors"

	errs "github.com/eotet/users-service/internal/errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user User) (User, error)
	GetUserByID(id uint32) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint32, user UpdateUserRequest) (User, error)
	DeleteUserByID(id uint32) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user User) (User, error) {
	if err := r.db.Where("email = ?", user.Email).First(&User{}).Error; err == nil {
		return User{}, errs.ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}

	if err := r.db.Create(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) GetUserByID(id uint32) (User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return User{}, errs.ErrUserNotFound
	}
	return user, nil
}

func (r *repository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) UpdateUserByID(id uint32, user UpdateUserRequest) (User, error) {
	var existingUser User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		return User{}, errs.ErrUserNotFound
	}

	if user.Email != nil {
		if err := r.db.Where("email = ? AND id != ?", *user.Email, id).First(&User{}).Error; err == nil {
			return User{}, errs.ErrUserAlreadyExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, err
		}
		existingUser.Email = *user.Email
	}
	if user.Password != nil {
		existingUser.Password = *user.Password
	}

	if err := r.db.Save(&existingUser).Error; err != nil {
		return User{}, err
	}

	return existingUser, nil
}

func (r *repository) DeleteUserByID(id uint32) error {
	var existingUser User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		return errs.ErrUserNotFound
	}

	if err := r.db.Delete(&existingUser).Error; err != nil {
		return err
	}

	return nil
}
