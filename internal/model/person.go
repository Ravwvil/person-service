package model

// Person представляет полную информацию о человеке.
type Person struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Surname     string  `db:"surname" json:"surname"`
	Patronymic  *string `db:"patronymic" json:"patronymic,omitempty"`
	Age         int     `db:"age" json:"age"`
	Gender      string  `db:"gender" json:"gender"`
	Nationality string  `db:"nationality" json:"nationality"`
}

// CreatePersonRequest — входная DTO для POST /people.
type CreatePersonRequest struct {
	Name       string  `json:"name" example:"Dmitriy" binding:"required"`
	Surname    string  `json:"surname" example:"Ushakov" binding:"required"`
	Patronymic *string `json:"patronymic,omitempty" example:"Vasilevich"`
}

// ToPerson преобразует DTO в модель для сохранения.
func (r CreatePersonRequest) ToPerson() Person {
	return Person{
		Name:       r.Name,
		Surname:    r.Surname,
		Patronymic: r.Patronymic,
	}
}

// UpdatePersonRequest — входная DTO для PUT /people/{id}.
type UpdatePersonRequest struct {
	Name        string  `json:"name" example:"Dmitriy"`
	Surname     string  `json:"surname" example:"Ushakov"`
	Patronymic  *string `json:"patronymic,omitempty" example:"Vasilevich"`
	Age         int     `json:"age" example:"30"`
	Gender      string  `json:"gender" example:"male"`
	Nationality string  `json:"nationality" example:"RU"`
}

// ToPerson преобразует DTO в модель, включая ID.
func (r UpdatePersonRequest) ToPerson(id int) Person {
	return Person{
		ID:          id,
		Name:        r.Name,
		Surname:     r.Surname,
		Patronymic:  r.Patronymic,
		Age:         r.Age,
		Gender:      r.Gender,
		Nationality: r.Nationality,
	}
}

// ErrorResponse описывает формат ошибок.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid JSON"`
}
