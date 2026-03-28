package repository

import (
	"errors"
	"locations-project/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) CreateUser(user ds.User) (ds.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *Repository) GetUserByEmail(email string) (ds.User, error) {
	var user ds.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.User{}, errors.New("user not found")
		}
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(id uint) (ds.User, error) {
	var user ds.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.User{}, errors.New("user not found")
		}
		return ds.User{}, err
	}
	return user, nil
}

func (r *Repository) UpdateUser(id uint, user ds.User) error {
	tx := r.db.Model(&ds.User{}).Where("id = ?", id).Updates(user)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *Repository) CheckCredentials(email, password string) (ds.User, error) {
	var user ds.User
	err := r.db.Where("email = ? AND password = ?", email, password).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.User{}, errors.New("invalid credentials")
		}
		return ds.User{}, err
	}
	return user, nil
}
