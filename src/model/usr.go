package model

import (
	"example.com/m/v2/src/helper"
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Time  int64  `json:"time"`
}

func (m *Movie) Create() error {

	sq := "INSERT INTO movie (title,time) VALUES ($1, $2) RETURNING id"
	err := helper.App.DB.QueryRow(sq, m.Title, m.Time).Scan(&m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Movie) Update(id int64) error {

	sq := "UPDATE movie SET title = $1, time = $2 WHERE id = $3"
	_, err := helper.App.DB.Exec(sq, m.Title, m.Time, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *Movie) Delete(id int64) error {

	sq := "DELETE FROM movie WHERE id = $1"
	_, err := helper.App.DB.Exec(sq, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *Movie) Get(id int64) error {

	sq := "SELECT id, title, time FROM movie WHERE id = $1"
	err := helper.App.DB.QueryRow(sq, id).Scan(&m.ID, &m.Title, &m.Time)
	if err != nil {
		return err
	}
	return nil
}

func (m *Movie) GetAll() ([]Movie, error) {

	rows, err := helper.App.DB.Query("SELECT id,title,time FROM movie")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var Movies []Movie

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var Movie Movie
		if err := rows.Scan(&Movie.ID, &Movie.Title, &Movie.Time); err != nil {
			return Movies, err
		}
		Movies = append(Movies, Movie)
	}
	if err = rows.Err(); err != nil {
		return Movies, err
	}
	return Movies, nil
}

func (m *Movie) GetAllByTime() ([]Movie, error) {

	rows, err := helper.App.DB.Query("SELECT id,title,time FROM movie ORDER BY time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var Movies []Movie

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var Movie Movie
		if err := rows.Scan(&Movie.ID, &Movie.Title, &Movie.Time); err != nil {
			return Movies, err
		}
		Movies = append(Movies, Movie)
	}
	if err = rows.Err(); err != nil {
		return Movies, err
	}
	return Movies, nil
}
