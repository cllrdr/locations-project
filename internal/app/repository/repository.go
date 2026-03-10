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
      ImagePath:   "img/map1.jpg",
	  VideoPath:   "img/map1.mp4",
	  Players: "3-6",
    },
    {
      ID:          2,
      Name:        "Пустоши Дурхота",
      Description: "Сухие красные степи, где золото встречается часто, а дерево можно найти только в редких оазисах.",
      ImagePath:   "img/map2.jpg",
	  VideoPath:   "img/map2.mp4",
	  Players: "2-4",
    },
    {
      ID:          3,
      Name:        "Ледяная Корона",
      Description: "Суровая мерзлота с незначительными запасами древесины, но богатыми золотыми жилами глубоко во льдах.",
      ImagePath:   "img/map3.jpg",
	  VideoPath:   "img/map3.mp4",
	  Players: "4-8",
    },
    {
      ID:          4,
      Name:        "Болота Печали",
      Description: "Топкое негостеприимное место, бедное на любые ресурсы — только выживание и контроль над ограниченными источниками.",
      ImagePath:   "img/map4.jpg",
	  VideoPath:   "img/map4.mp4",
	  Players: "1-3",
    },
    {
      ID:          5,
      Name:        "Пылающие степи",
      Description: "Выжженная демоническая земля, где нет деревьев, но золото течет рекой (в прямом смысле — жидкое золото в лаве).",
      ImagePath:   "img/map5.jpg",
	  VideoPath:   "img/map5.mp4",
	  Players: "2-4",
    },
  }

  if len(locations) == 0 {
    return nil, fmt.Errorf("массив пустой")
  }

  return locations, nil
}

func (r *Repository) GetLocation(id int) (Location, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	locations, err := r.GetLocations()
	if err != nil {
		return Location{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, location := range locations {
		if location.ID == id {
			return location, nil // если нашли, то просто возвращаем найденную локацию без ошибок
		}
	}
	return Location{}, fmt.Errorf("локация не найдена") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
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