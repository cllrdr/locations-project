package repository

import (
	"fmt"
	"locations-project/internal/app/ds"
	"math/rand"
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
	creatorID := uint(1)

	// Получаем пользователя для заполнения Nickname
	var user ds.User
	err := r.db.First(&user, creatorID).Error
	nickname := ""
	if err == nil {
		nickname = user.Name
	}

	request := ds.PlayersLocationRequest{
		Nickname:  nickname,
		Status:    ds.RequestStatusDraft,
		CreatorID: creatorID,
	}

	err = r.db.Create(&request).Error
	if err != nil {
		return ds.PlayersLocationRequest{}, err
	}

	err = r.AddLocationToRequest(request.ID, locationID)
	if err != nil {
		return ds.PlayersLocationRequest{}, err
	}

	return request, nil
}

// AddLocationToRequest добавляет локацию в заявку (приоритет 1 по умолчанию)
func (r *Repository) AddLocationToRequest(requestID, locationID uint) error {
	chosenLocation := ds.PlayersChosenLocation{
		RequestID:  requestID,
		LocationID: locationID,
		Priority:   1,
	}

	return r.db.Create(&chosenLocation).Error
}

// DeleteRequest меняет статус заявки на "удалён"
func (r *Repository) DeleteRequest(requestID uint) error {
	return r.db.Exec("UPDATE players_location_requests SET status = ? WHERE id = ?", ds.RequestStatusDeleted, requestID).Error
}

// ChooseRandomLocation выбирает случайную локацию с учётом весов приоритетов
// Формула: P(locationᵢ) = priorityᵢ / Σ(priorityⱼ)
func (r *Repository) ChooseRandomLocation(locs []ds.PlayersChosenLocation) ds.PlayersChosenLocation {
	if len(locs) == 0 {
		return ds.PlayersChosenLocation{}
	}
	total := 0
	for _, l := range locs {
		total += l.Priority
	}
	randVal, cum := rand.Intn(total), 0
	for _, l := range locs {
		if cum += l.Priority; randVal < cum {
			return l
		}
	}
	return locs[0]
}

// ChooseRandomLocationForUser выбирает случайную локацию из корзины пользователя
// (заявка остаётся в статусе "черновик")
func (r *Repository) ChooseRandomLocationForUser(creatorID uint) (ds.PlayersChosenLocation, error) {
	// Получаем черновик заявки пользователя
	draftRequest, chosenLocations, err := r.GetDraftRequestInfo()
	if err != nil {
		return ds.PlayersChosenLocation{}, err
	}

	if len(chosenLocations) == 0 {
		return ds.PlayersChosenLocation{}, fmt.Errorf("в корзине нет локаций")
	}

	// Выбираем случайную локацию с учётом приоритетов
	chosen := r.ChooseRandomLocation(chosenLocations)

	// Сохраняем ID выбранной локации в заявке
	err = r.db.Model(&draftRequest).Update("randomed_location", chosen.LocationID).Error
	if err != nil {
		return ds.PlayersChosenLocation{}, err
	}

	return chosen, nil
}
