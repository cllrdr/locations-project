package repository

import (
	"locations-project/internal/app/ds"
)

// GetPlayersLocationsForRequest получает информацию о заявке пользователя
func (r *Repository) GetPlayersLocationsForRequest(requestID int) (ds.PlayersLocationRequest, []ds.PlayersChosenLocation, error) {
	var request ds.PlayersLocationRequest
	var chosenLocations []ds.PlayersChosenLocation

	err := r.db.Joins("Creator").First(&request, requestID).Error
	if err != nil {
		return request, nil, err
	}

	err = r.db.Where("request_id = ?", requestID).Find(&chosenLocations).Error
	return request, chosenLocations, err
}

// GetDraftRequestInfo получает информацию о черновике заявки пользователя
func (r *Repository) GetDraftRequestInfo() (ds.PlayersLocationRequest, []ds.PlayersChosenLocation, error) {
	creatorID := uint(1)

	var request ds.PlayersLocationRequest
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.RequestStatusDraft).First(&request).Error
	if err != nil {
		return ds.PlayersLocationRequest{}, nil, err
	}

	var chosenLocations []ds.PlayersChosenLocation
	err = r.db.Where("request_id = ?", request.ID).Find(&chosenLocations).Error
	if err != nil {
		return ds.PlayersLocationRequest{}, nil, err
	}

	return request, chosenLocations, nil
}

// CreateRequestWithLocation создаёт новую заявку и добавляет в неё локацию
func (r *Repository) CreateRequestWithLocation(locationID uint) (ds.PlayersLocationRequest, error) {
	request := ds.PlayersLocationRequest{
		Nickname:  "",
		Status:    ds.RequestStatusDraft,
		CreatorID: 1,
	}

	err := r.db.Create(&request).Error
	if err != nil {
		return ds.PlayersLocationRequest{}, err
	}

	err = r.AddLocationToRequest(request.ID, locationID)
	if err != nil {
		return ds.PlayersLocationRequest{}, err
	}

	return request, nil
}

// AddLocationToRequest добавляет локацию в заявку (приоритет 0 по умолчанию)
func (r *Repository) AddLocationToRequest(requestID, locationID uint) error {
	chosenLocation := ds.PlayersChosenLocation{
		RequestID:  requestID,
		LocationID: locationID,
		Priority:   0,
	}

	return r.db.Create(&chosenLocation).Error
}

// DeleteRequest меняет статус заявки на "удалён"
func (r *Repository) DeleteRequest(requestID uint) error {
	return r.db.Exec("UPDATE players_location_requests SET status = ? WHERE id = ?", ds.RequestStatusDeleted, requestID).Error
}
