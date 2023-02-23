package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Gun struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Type         string `json:"type"`
	Caliber      string `json:"caliber"`
	Price        int    `json:"price"`
	Availability bool   `json:"availability"`
}

type GunModel struct {
	DB *sql.DB
}

func (m GunModel) Insert(gun *Gun) error {
	query := `
	INSERT INTO guns (name, model, caliber, price, availability, type)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, name`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{gun.Name, gun.Model, gun.Caliber, gun.Price, gun.Availability, gun.Type}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&gun.ID, &gun.Name)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m GunModel) Get(id int64) (*Gun, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, "name", model, caliber, price, availability, type
	FROM guns where id = $1`
	var gun Gun
	err := m.DB.QueryRow(query, id).Scan(
		&gun.ID,
		&gun.Name,
		&gun.Model,
		&gun.Caliber,
		&gun.Price,
		&gun.Availability,
		&gun.Type,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &gun, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m GunModel) Update(gun *Gun) error {
	query := `UPDATE public.guns
	SET "name"=$1, model=$2, caliber=$3, price=$4, availability=$5
	WHERE id=$6
	RETURNING id;`

	args := []any{
		gun.Name,
		gun.Model,
		gun.Caliber,
		gun.Price,
		gun.Availability,
		gun.ID,
	}
	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&gun.ID)
}

func (m GunModel) List(title string, genres string, filters Filters) ([]*Gun, error) {
	query := `SELECT id, "name", model, caliber, price, availability, type
	FROM guns WHERE (LOWER(name) = LOWER($1) OR $1 = '') and  (LOWER(type) = LOWER($2) OR $2 = '') ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, title, genres)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guns := []*Gun{}

	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var gun Gun

		err := rows.Scan(
			&gun.ID,
			&gun.Name,
			&gun.Model,
			&gun.Caliber,
			&gun.Price,
			&gun.Availability,
			&gun.Type,
		)
		if err != nil {
			return nil, err
		}

		guns = append(guns, &gun)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return guns, nil
}
