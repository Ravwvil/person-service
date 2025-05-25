package models

type Person struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Surname     string  `db:"surname" json:"surname"`
	Patronymic  *string `db:"patronymic" json:"patronymic,omitempty"`
	Age         int     `db:"age" json:"age"`
	Gender      string  `db:"gender" json:"gender"`
	Nationality string  `db:"nationality" json:"nationality"`
}
