package data

import "database/sql"

type Gun struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Model        string `json:"model"`
	Caliber      string `json:"caliber"`
	Price        int    `json:"price"`
	Availability bool   `json:"availability"`
}

type GunModel struct {
	DB *sql.DB
}

func (m GunModel) Insert(gun *Gun) error {
	query := `
	INSERT INTO guns (name, model, caliber, price, availability)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, name`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{gun.Name, gun.Model, gun.Caliber, gun.Price, gun.Availability}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&gun.ID, &gun.Name)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m GunModel) Get(id int64) (*Gun, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m GunModel) Update(gun *Gun) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m GunModel) Delete(id int64) error {
	return nil
}
