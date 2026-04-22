package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"locations-project/internal/app/ds"
	"strings"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func (r *Repository) GetLocations(locationName string) ([]ds.Location, error) {
	var locations []ds.Location
	query := r.db.Where("is_deleted = ?", false)

	if locationName != "" {
		query = query.Where("name ILIKE ?", "%"+locationName+"%")
	}

	err := query.Find(&locations).Error
	if locations == nil {
		locations = []ds.Location{}
	}
	return locations, err
}

func (r *Repository) GetLocation(id uint) (ds.Location, error) {
	var location ds.Location
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&location).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Location{}, errors.New("location not found")
		}
		return ds.Location{}, err
	}
	return location, nil
}

func (r *Repository) GetLocationsByName(name string) ([]ds.Location, error) {
	return r.GetLocations(name)
}

func (r *Repository) CreateLocation(location ds.Location) (ds.Location, error) {
	err := r.db.Create(&location).Error
	return location, err
}

func (r *Repository) UpdateLocation(id uint, location ds.Location) error {
	tx := r.db.Model(&ds.Location{}).Where("id = ? AND is_deleted = ?", id, false).Updates(location)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("location not found")
	}
	return nil
}

func (r *Repository) DeleteLocation(id uint) error {
	tx := r.db.Model(&ds.Location{}).Where("id = ?", id).Update("is_deleted", true)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("location not found")
	}
	return nil
}

func (r *Repository) UpdateLocationImage(id uint, imagePath string) error {
	tx := r.db.Model(&ds.Location{}).Where("id = ? AND is_deleted = ?", id, false).Update("image_path", imagePath)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("location not found")
	}
	return nil
}

func (r *Repository) UploadFileToMinIO(ctx context.Context, fileName string, fileReader io.Reader, fileSize int64, contentType string) error {
	// Проверяем доступность MinIO
	_, err := r.minio.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("MinIO is not accessible: %v", err)
	}

	// Загружаем файл
	_, err = r.minio.PutObject(ctx, r.bucket, fileName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (r *Repository) DeleteFileFromMinIO(ctx context.Context, fileName string) error {
	// Извлекаем имя файла из пути
	parts := strings.Split(fileName, "/")
	fileName = parts[len(parts)-1]

	err := r.minio.RemoveObject(ctx, r.bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}
	return nil
}
