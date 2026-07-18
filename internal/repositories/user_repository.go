package repositories

import (
	"context"
	"errors"
	"fmt"
	"task-management-api/internal/apperrors"
	"task-management-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(
	ctx context.Context,
	user *models.User,
) error {
	err := r.db.WithContext(ctx).Create(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrConflict
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *userRepository) FindByID(
	ctx context.Context,
	id uint,
) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("find user by email: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) Update(
	ctx context.Context,
	user *models.User,
) error {
	err := r.db.WithContext(ctx).Save(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrConflict
		}

		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		return fmt.Errorf("delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}