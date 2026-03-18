package repository

import (
	"locations-project/internal/app/ds"
)

func (r *Repository) GetPlayersLocationsForRequest(requestID int) (ds.PlayersLocationRequest, []ds.PlayersChosenLocation, error) {
	var request ds.PlayersLocationRequest
	var chosenLocations []ds.PlayersChosenLocation

	err := r.db.First(&request, requestID).Error
	if err != nil {
		return request, nil, err
	}

	err = r.db.Where("request_id = ?", requestID).Find(&chosenLocations).Error
	return request, chosenLocations, err
}
