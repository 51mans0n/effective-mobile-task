package model

import "time"

type Person struct {
	ID             string    `db:"id"`
	Name           string    `db:"name"`
	Surname        string    `db:"surname"`
	Patronymic     *string   `db:"patronymic"`
	Age            *int      `db:"age"`
	Gender         *string   `db:"gender"`
	CountryCode    *string   `db:"country_code"`
	NatProbability *float32  `db:"nat_probability"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
