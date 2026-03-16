package repository

import (
  "fmt"
  "strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
  return &Repository{}, nil
}

type Location struct {
  ID          int
  Name        string
  Description string
  ImagePath   string
  VideoPath   string
  Players     string
}

func (r *Repository) GetLocations() ([]Location, error) {

  locations := []Location{
    {
      ID:          1,
      Name:        "Тирсфальские леса",
      Description: "Заброшенные леса Лордерона с богатыми залежами золота и лесными угодьями.",
      ImagePath:   "http://localhost:9000/locations-images/map1.jpg",
	    VideoPath:   "http://localhost:9000/locations-videos/map1.mp4",
	    Players: "3-6",
    },
    {
      ID:          2,
      Name:        "Пустоши Дурхота",
      Description: "Сухие красные степи, где золото встречается часто, а дерево можно найти только в редких оазисах.",
      ImagePath:   "http://localhost:9000/locations-images/map2.jpg",
	    VideoPath:   "http://localhost:9000/locations-videos/map2.mp4",
	    Players: "2-4",
    },
    {
      ID:          3,
      Name:        "Ледяная Корона",
      Description: "Суровая мерзлота с незначительными запасами древесины, но богатыми золотыми жилами глубоко во льдах.",
      ImagePath:   "http://localhost:9000/locations-images/map3.jpg",
	    VideoPath:   "http://localhost:9000/locations-videos/map3.mp4",
	    Players: "4-8",
    },
    {
      ID:          4,
      Name:        "Болота Печали",
      Description: "Топкое негостеприимное место, бедное на любые ресурсы — только выживание и контроль над ограниченными источниками.",
      ImagePath:   "http://localhost:9000/locations-images/map4.jpg",
	    VideoPath:   "http://localhost:9000/locations-videos/map4.mp4",
	    Players: "1-3",
    },
    {
      ID:          5,
      Name:        "Пылающие степи",
      Description: "Выжженная демоническая земля, где нет деревьев, но золото течет рекой (в прямом смысле — жидкое золото в лаве).",
      ImagePath:   "http://localhost:9000/locations-images/map5.jpg",
	    VideoPath:   "http://localhost:9000/locations-videos/map5.mp4",
	    Players: "2-4",
    },
  }

  if len(locations) == 0 {
    return nil, fmt.Errorf("массив пустой")
  }

  return locations, nil
}

func (r *Repository) GetLocation(id int) (Location, error) {
	locations, err := r.GetLocations()
	if err != nil {
		return Location{}, err
	}

	for _, location := range locations {
		if location.ID == id {
			return location, nil
		}
	}
	return Location{}, fmt.Errorf("локация не найдена")
}

func (r *Repository) GetLocationsByName(name string) ([]Location, error) {
	locations, err := r.GetLocations()
	if err != nil {
		return []Location{}, err
	}

	var result []Location
	for _, location := range locations {
		if strings.Contains(strings.ToLower(location.Name), strings.ToLower(name)) {
			result = append(result, location)
		}
	}

	return result, nil
}

type PlayersLocationRequest struct {
	ID   int
	Nickname string
}

type PlayersChosenLocations struct {
	RequestID int
	LocationID  int
	Priority  int
}

var playersLocationRequests = map[PlayersLocationRequest][]PlayersChosenLocations{
	{ID: 1, Nickname: "Player_Alpha"}: {
		{
			RequestID:  1,
			LocationID: 2,
			Priority:   1,
		},
		{
			RequestID:  1,
			LocationID: 3,
			Priority:   2,
		},
		{
			RequestID:  1,
			LocationID: 4,
			Priority:   3,
		},
	},
	{ID: 2, Nickname: "Player_Beta"}: {
		{
			RequestID:  2,
			LocationID: 2,
			Priority:   1,
		},
		{
			RequestID:  2,
			LocationID: 4,
			Priority:   2,
		},
	},
	{ID: 3, Nickname: "Player_Gamma"}: {
		{
			RequestID:  3,
			LocationID: 5,
			Priority:   1,
		},
	},
}


func (r *Repository) GetPlayersLocationsForRequest(requestID int) (PlayersLocationRequest, []PlayersChosenLocations, error) {
	for req, locations := range playersLocationRequests {
		if req.ID == requestID {
			return req, locations, nil
		}
	}
	return PlayersLocationRequest{}, nil, fmt.Errorf("запрос с ID %d не найден", requestID)
}
