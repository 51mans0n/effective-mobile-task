// Package model defines core domain models used across the service.
package model

// Enriched contains predicted age, gender, and nationality.
type Enriched struct {
	Age         *int
	Gender      *string
	CountryCode *string
	Probability *float32
}
