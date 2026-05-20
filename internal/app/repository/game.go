package repository

import (
	"errors"
	"fmt"
	"locations-project/internal/app/ds"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// GetRequests получает список заявок с фильтрацией по статусу и диапазону даты формирования
func (r *Repository) GetRequests(status *ds.GameStatus, startDate, endDate *time.Time) ([]ds.PlayersLocationGame, error) {
	var requests []ds.PlayersLocationGame
	query := r.db.Where("status != ? AND status != ?", ds.GameStatusDeleted, ds.GameStatusDraft)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if startDate != nil {
		query = query.Where("formed_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("formed_at <= ?", *endDate)
	}

	err := query.Preload("Creator").Preload("Moderator").Find(&requests).Error
	return requests, err
}

// GetRequest получает одну заявку по ID
func (r *Repository) GetRequest(id uint) (ds.PlayersLocationGame, error) {
	var request ds.PlayersLocationGame
	err := r.db.Preload("Creator").Preload("Moderator").Where("id = ? AND status != ?", id, ds.GameStatusDeleted).First(&request).Error
	return request, err
}

// GetRequestWithLocations получает заявку со списком локаций
func (r *Repository) GetRequestWithLocations(id uint) (ds.PlayersLocationGame, []ds.PlayersChosenLocation, error) {
	req, err := r.GetRequest(id)
	if err != nil {
		return ds.PlayersLocationGame{}, nil, err
	}

	var locations []ds.PlayersChosenLocation
	// ✅ Исправлено: request_id вместо players_location_game_id
	err = r.db.Preload("Location").Where("request_id = ?", id).Find(&locations).Error
	if err != nil {
		return ds.PlayersLocationGame{}, nil, err
	}

	return req, locations, nil
}

// GetPlayersLocationsForRequest получает информацию о заявке пользователя
func (r *Repository) GetPlayersLocationsForRequest(requestID int) (ds.PlayersLocationGame, []ds.PlayersChosenLocation, error) {
	return r.GetRequestWithLocations(uint(requestID))
}

// GetDraftRequestInfo получает информацию о черновике заявки пользователя
func (r *Repository) GetDraftRequestInfo(creatorID uint) (ds.PlayersLocationGame, []ds.PlayersChosenLocation, error) {
	var request ds.PlayersLocationGame
	// ✅ Исправлено: creatorID теперь приходит параметром, запрос по request_id
	err := r.db.Preload("Creator").Preload("Moderator").Where("creator_id = ? AND status = ?", creatorID, ds.GameStatusDraft).First(&request).Error
	if err != nil {
		return ds.PlayersLocationGame{}, nil, err
	}

	var chosenLocations []ds.PlayersChosenLocation
	err = r.db.Preload("Location").Where("request_id = ?", request.ID).Find(&chosenLocations).Error
	if err != nil {
		return ds.PlayersLocationGame{}, nil, err
	}

	return request, chosenLocations, nil
}

// UpdateRequest обновляет поля заявки
func (r *Repository) UpdateRequest(id uint, request ds.PlayersLocationGame) error {
	var existingRequest ds.PlayersLocationGame
	err := r.db.Where("id = ? AND status != ?", id, ds.GameStatusDeleted).First(&existingRequest).Error
	if err != nil {
		return err
	}
	return r.db.Model(&existingRequest).Updates(request).Error
}

// UpdateRequestStatus обновляет статус заявки с проверкой допустимых переходов
func (r *Repository) UpdateRequestStatus(id uint, newStatus ds.GameStatus, moderatorID *uint) error {
	var request ds.PlayersLocationGame
	err := r.db.Where("id = ? AND status != ?", id, ds.GameStatusDeleted).First(&request).Error
	if err != nil {
		return err
	}

	if !r.isValidStatusTransition(request.Status, newStatus) {
		return fmt.Errorf("недопустимый переход статуса с %s на %s", request.Status, newStatus)
	}

	updates := map[string]interface{}{
		"status": newStatus,
	}

	switch newStatus {
	case ds.GameStatusFormed:
		now := time.Now()
		updates["formed_at"] = now
	case ds.GameStatusCompleted, ds.GameStatusRejected:
		now := time.Now()
		updates["completed_at"] = now
		if moderatorID != nil {
			updates["moderator_id"] = *moderatorID
		}
	}

	return r.db.Model(&ds.PlayersLocationGame{}).Where("id = ?", id).Updates(updates).Error
}

// isValidStatusTransition проверяет допустимость перехода статуса
func (r *Repository) isValidStatusTransition(current, new ds.GameStatus) bool {
	validTransitions := map[ds.GameStatus][]ds.GameStatus{
		ds.GameStatusDraft:     {ds.GameStatusDeleted, ds.GameStatusFormed},
		ds.GameStatusFormed:    {ds.GameStatusCompleted, ds.GameStatusRejected},
		ds.GameStatusCompleted: {},
		ds.GameStatusRejected:  {},
		ds.GameStatusDeleted:   {},
	}

	allowedStatuses, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == new {
			return true
		}
	}
	return false
}

// CreateRequestWithLocation создаёт новую заявку и добавляет в неё локацию
func (r *Repository) CreateRequestWithLocation(locationID uint, creatorID uint) (ds.PlayersLocationGame, error) {
	request := ds.PlayersLocationGame{
		Nickname:  "",
		Status:    ds.GameStatusDraft,
		CreatorID: creatorID, // ✅ Из токена
	}

	err := r.db.Create(&request).Error
	if err != nil {
		return ds.PlayersLocationGame{}, err
	}

	err = r.AddLocationToRequest(request.ID, locationID)
	if err != nil {
		return ds.PlayersLocationGame{}, err
	}

	return request, nil
}

// AddLocationToRequest добавляет локацию в заявку (приоритет 1 по умолчанию)
func (r *Repository) AddLocationToRequest(requestID, locationID uint) error {
	chosenLocation := ds.PlayersChosenLocation{
		RequestID:  requestID, // ✅ Исправлено: было PlayersLocationGameID
		LocationID: locationID,
		Priority:   1,
	}
	err := r.db.Create(&chosenLocation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("локация уже добавлена в заявку")
		}
		return err
	}
	return nil
}

// DeleteRequest меняет статус заявки на "удалён"
func (r *Repository) DeleteRequest(requestID uint) error {
	var existingRequest ds.PlayersLocationGame
	err := r.db.Where("id = ? AND status != ?", requestID, ds.GameStatusDeleted).First(&existingRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return r.db.Model(&existingRequest).Update("status", ds.GameStatusDeleted).Error
}

// FormRequest формирует черновик заявки (переводит в статус "сформирован")
func (r *Repository) FormRequest(id uint) error {
	// Проверка пустоты и статуса делается в handler, здесь только смена статуса
	newStatus := ds.GameStatusFormed
	return r.UpdateRequestStatus(id, newStatus, nil)
}

// CompleteRequest завершает или отклоняет заявку
func (r *Repository) CompleteRequest(id uint, approve bool) error {
	status := ds.GameStatusRejected
	if approve {
		status = ds.GameStatusCompleted
	}

	err := r.UpdateRequestStatus(id, status, nil)
	if err != nil {
		return err
	}

	if approve && status == ds.GameStatusCompleted {
		err = r.chooseRandomLocationForRequest(id)
		if err != nil {
			return err
		}
	}

	return nil
}

// chooseRandomLocationForRequest выбирает случайную локацию и устанавливает флаг IsRandomed
func (r *Repository) chooseRandomLocationForRequest(requestID uint) error {
	var chosenLocations []ds.PlayersChosenLocation
	// ✅ Исправлено: request_id
	err := r.db.Where("request_id = ?", requestID).Find(&chosenLocations).Error
	if err != nil {
		return err
	}

	if len(chosenLocations) == 0 {
		return errors.New("в заявке нет локаций")
	}

	chosen := r.ChooseRandomLocation(chosenLocations)

	// Сбрасываем флаг IsRandomed для всех локаций этого запроса
	err = r.db.Model(&ds.PlayersChosenLocation{}).Where("request_id = ?", requestID).Update("is_randomed", false).Error
	if err != nil {
		return err
	}

	// Устанавливаем флаг IsRandomed = true для выбранной локации
	err = r.db.Model(&ds.PlayersChosenLocation{}).Where("id = ?", chosen.ID).Update("is_randomed", true).Error
	return err
}

// ChooseRandomLocation выбирает случайную локацию с учётом весов приоритетов
func (r *Repository) ChooseRandomLocation(chosenLocations []ds.PlayersChosenLocation) ds.PlayersChosenLocation {
	if len(chosenLocations) == 0 {
		return ds.PlayersChosenLocation{}
	}

	var totalPriority int
	for _, loc := range chosenLocations {
		totalPriority += loc.Priority
	}

	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(totalPriority) + 1

	cumulative := 0
	for _, loc := range chosenLocations {
		cumulative += loc.Priority
		if random <= cumulative {
			return loc
		}
	}

	return chosenLocations[len(chosenLocations)-1]
}

// ChooseRandomLocationForUser выбирает случайную локацию из корзины пользователя
func (r *Repository) ChooseRandomLocationForUser(creatorID uint) (ds.PlayersChosenLocation, error) {
	draftRequest, chosenLocations, err := r.GetDraftRequestInfo(creatorID)
	if err != nil {
		return ds.PlayersChosenLocation{}, err
	}

	if len(chosenLocations) == 0 {
		return ds.PlayersChosenLocation{}, fmt.Errorf("в корзине нет локаций")
	}

	chosen := r.ChooseRandomLocation(chosenLocations)

	err = r.db.Model(&draftRequest).Update("randomed_location", chosen.LocationID).Error
	if err != nil {
		return ds.PlayersChosenLocation{}, err
	}

	return chosen, nil
}

// RemoveLocationFromRequest удаляет локацию из заявки
func (r *Repository) RemoveLocationFromRequest(requestID, locationID uint) error {
	// ✅ Исправлено: request_id
	result := r.db.Where("request_id = ? AND location_id = ?", requestID, locationID).Delete(&ds.PlayersChosenLocation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("локация не найдена в заявке")
	}
	return nil
}

// UpdateLocationPriority обновляет приоритет локации в заявке
func (r *Repository) UpdateLocationPriority(requestID, locationID uint, priority int) error {
	// ✅ Исправлено: request_id
	result := r.db.Model(&ds.PlayersChosenLocation{}).
		Where("request_id = ? AND location_id = ?", requestID, locationID).
		Update("priority", priority)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("локация не найдена в заявке")
	}
	return nil
}

// GetRandomedLocationsCount получает количество выбранных локаций для заявки
func (r *Repository) GetRandomedLocationsCount(requestID uint) (int, error) {
	var count int64
	// ✅ Исправлено: request_id
	err := r.db.Where("request_id = ?", requestID).Model(&ds.PlayersChosenLocation{}).Count(&count).Error
	return int(count), err
}

// GetRequestsByCreator получает заявки конкретного пользователя с фильтрацией
func (r *Repository) GetRequestsByCreator(creatorID uint, status *ds.GameStatus, startDate, endDate *time.Time) ([]ds.PlayersLocationGame, error) {
	var requests []ds.PlayersLocationGame
	query := r.db.Where("creator_id = ? AND status != ? AND status != ?", creatorID, ds.GameStatusDeleted, ds.GameStatusDraft)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if startDate != nil {
		query = query.Where("formed_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("formed_at <= ?", *endDate)
	}

	err := query.Preload("Creator").Preload("Moderator").Find(&requests).Error
	return requests, err
}