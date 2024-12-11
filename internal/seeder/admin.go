package seeder

import (
	"context"
	"github.com/google/uuid"
	"github.com/kholidss/movie-fest-skilltest/internal/consts"
	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/pkg/cipher"
	"github.com/kholidss/movie-fest-skilltest/pkg/masker"
	"log"
)

var (
	defaultAdminEmail    = `admin001@yopmail.com`
	defaultAdminFullName = `Admin Root`
	defaultAdminPassword = `admin@1234`
)

func (s *seeder) AdminData(ctx context.Context) {
	log.Println("admin seed store data starting...")

	adminExist, err := s.repoUser.FindOne(ctx, entity.User{
		Email: defaultAdminEmail,
	}, []string{"id", "full_name", "entity"})

	if err != nil {
		log.Printf("find one exist admin seed data err: %v", err)
		return
	}

	if adminExist != nil {
		log.Printf("admin seed data already exist, id: %s, fullname: %s", adminExist.ID, adminExist.FullName)
		return
	}

	var (
		userID = uuid.New().String()
	)

	password, err := cipher.EncryptAES256(defaultAdminPassword, s.cfg.AppConfig.AppPasswordSecret)
	if err != nil {
		log.Printf("encrypt admin seed password err: %v", err)
		return
	}

	err = s.repoUser.Store(ctx, entity.User{
		ID:       userID,
		FullName: defaultAdminFullName,
		Email:    defaultAdminEmail,
		Password: password,
		Entity:   consts.RoleEntityAdmin,
	})

	if err != nil {
		log.Printf("store admin seed data err: %v", err)
		return
	}

	log.Printf("admin seed data store success\nuser_id: %s\nemail: %s\npassword: %s\n",
		userID,
		defaultAdminEmail,
		masker.Censored(defaultAdminPassword, "*"))

}
