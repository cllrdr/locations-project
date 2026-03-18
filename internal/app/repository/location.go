package repository

import (
	"fmt"
	"locations-project/internal/app/ds"
)

func (r *Repository) GetLocations() ([]ds.Location, error) {
	var locations []ds.Location
	err := r.db.Find(&locations).Error
	if err != nil {
		return nil, err
	}
	if len(locations) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return locations, nil
}

func (r *Repository) GetLocation(id int) (ds.Location, error) {
	location := ds.Location{}
	err := r.db.Where("id = ?", id).First(&location).Error
	if err != nil {
		return ds.Location{}, err
	}
	return location, nil
}

func (r *Repository) GetLocationsByName(name string) ([]ds.Location, error) {
	var locations []ds.Location
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, nil
}
