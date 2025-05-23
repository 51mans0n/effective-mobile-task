// Package service implements business logic for people management.
package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/51mans0n/effective-mobile-task/internal/client"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"github.com/51mans0n/effective-mobile-task/internal/repository"
)

// Service provides business logic for person enrichment and storage.
type Service struct {
	repo      repository.PeopleRepo
	agify     client.Enricher
	genderize client.Enricher
	nat       client.Enricher
}

// New creates a new Service instance with repository and API clients.
func New(repo repository.PeopleRepo) *Service {
	return &Service{
		repo:      repo,
		agify:     client.NewAgify(),
		genderize: client.NewGenderize(),
		nat:       client.NewNationalize(),
	}
}

// Create → fan-out к внешним API, мерж, сохранение
func (s *Service) Create(ctx context.Context, p *model.Person) error {
	var enr [3]*model.Enriched

	// этот контекст только для fan-out (внешние API)
	ctxFan, cancelFan := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFan()

	g, ctxFan := errgroup.WithContext(ctxFan)

	g.Go(func() error {
		var err error
		enr[0], err = s.agify.Enrich(ctxFan, p.Name)
		if err != nil {
			fmt.Println("agify error:", err)
		}
		return err
	})

	g.Go(func() error {
		var err error
		enr[1], err = s.genderize.Enrich(ctxFan, p.Name)
		if err != nil {
			fmt.Println("genderize error:", err)
		}
		return err
	})

	g.Go(func() error {
		var err error
		enr[2], err = s.nat.Enrich(ctxFan, p.Name)
		if err != nil {
			fmt.Println("nationalize error:", err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// merge into Person
	if enr[0].Age != nil {
		p.Age = enr[0].Age
	}
	if enr[1].Gender != nil {
		p.Gender = enr[1].Gender
	}
	if enr[2].CountryCode != nil {
		p.CountryCode = enr[2].CountryCode
		p.NatProbability = enr[2].Probability
	}

	// используем исходный ctx без таймаута (важно!)
	return s.repo.Create(ctx, p)
}

// List returns a filtered and paginated list of people.
func (s *Service) List(ctx context.Context, f repository.ListFilter) (*repository.PaginatedPeople, error) {
	return s.repo.List(ctx, f)
}

// Get returns a person by UUID.
func (s *Service) Get(ctx context.Context, id string) (*model.Person, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateName updates name/surname/patronymic fields of a person.
func (s *Service) UpdateName(ctx context.Context, id string, n, sname, patr string) (bool, error) {
	return s.repo.UpdateName(ctx, id, n, sname, patr)
}

// Delete removes a person from storage by UUID.
func (s *Service) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}
