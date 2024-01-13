package repository

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrDataNotFound

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(d *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: d,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	err := repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	return err
}

func (repo *UserRepository) FindByEmail(ctx context.Context,
	email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx context.Context,
	id int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, id)
	return repo.toDomain(u), err
}

func (repo *UserRepository) UpdateById(ctx context.Context, u domain.User) error {
	return repo.dao.UpdateById(ctx, dao.User{Id: u.Id, NickName: u.NickName, BirthDay: u.BirthDay.UnixMilli(), AboutMe: u.AboutMe})
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		NickName: u.NickName,
		BirthDay: time.UnixMilli(u.BirthDay),
		AboutMe:  u.AboutMe,
	}
}
