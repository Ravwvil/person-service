package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"person-service/internal/model"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewPersonRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(p *model.Person) (int, error) {
	query := `
		INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES (:name, :surname, :patronymic, :age, :gender, :nationality)
		RETURNING id`
	rows, err := r.db.NamedQuery(query, p)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, nil
}

func (r *PersonRepository) GetAll(filters map[string]interface{}, limit, offset int) ([]model.Person, error) {
	query := `SELECT * FROM people WHERE 1=1`
	args := []interface{}{}
	argID := 1

	for k, v := range filters {
		query += fmt.Sprintf(" AND %s = $%d", k, argID)
		args = append(args, v)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	var people []model.Person
	err := r.db.Select(&people, query, args...)
	return people, err
}

func (r *PersonRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM people WHERE id=$1", id)
	return err
}

func (r *PersonRepository) Update(id int, p *model.Person) error {
	query := `
		UPDATE people SET name=:name, surname=:surname, patronymic=:patronymic,
			age=:age, gender=:gender, nationality=:nationality
		WHERE id=:id`
	p.ID = id
	_, err := r.db.NamedExec(query, p)
	return err
}
