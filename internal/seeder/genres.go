package seeder

import (
	"context"
	"github.com/google/uuid"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/pkg/helper"
	"log"
)

// Seed data
var defaultGenres = []entity.Genre{
	{
		ID:   uuid.New().String(),
		Name: "Action",
	},
	{
		ID:   uuid.New().String(),
		Name: "Drama",
	},
	{
		ID:   uuid.New().String(),
		Name: "Horror",
	},
	{
		ID:   uuid.New().String(),
		Name: "Thriller",
	},
	{
		ID:   uuid.New().String(),
		Name: "Romance",
	},
	{
		ID:   uuid.New().String(),
		Name: "Animation",
	},
	{
		ID:   uuid.New().String(),
		Name: "Sci Fi",
	},
}

func (s *seeder) GenresData(ctx context.Context) {
	log.Println("genres seed store data starting...")

	for _, v := range defaultGenres {
		slug := helper.ToSlugFormat(v.Name)

		genre, err := s.repoGenre.FindOne(ctx, entity.Genre{
			Slug: slug,
		}, []string{"id", "name"})
		if err != nil {
			log.Printf("find one exist genre: %s seed data err: %v", v.Name, err)
			continue
		}

		if genre != nil {
			log.Printf("genre: %s already exist, skipped", v.Name)
			continue
		}

		err = s.repoGenre.Store(ctx, entity.Genre{
			ID:         v.ID,
			Name:       v.Name,
			Slug:       slug,
			ViewNumber: 0,
		})
		if err != nil {
			log.Printf("store genre: %s seed data err: %v", v.Name, err)
			continue
		}

		log.Printf("succes store seed data genre: %s", v.Name)

	}
}
