package dao

import (
	"context"

	"github.com/mathiasXie/curd"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
	"gorm.io/gorm"
)

type UserDAO struct {
	userModel *curd.Model[model.User]
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		userModel: curd.NewModel[model.User](db),
	}
}

func (u *UserDAO) AddUser(ctx context.Context, user *model.User) (int64, error) {
	insert, err := u.userModel.Insert(ctx, map[string]interface{}{
		"user_name": user.UserName,
	})
	if err != nil {
		return 0, err
	}
	return insert.Id, nil
}

func (u *UserDAO) GetAllUser(ctx context.Context) ([]*model.User, int64, error) {
	users, total, err := u.userModel.SelectAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u *UserDAO) GetUser(ctx context.Context, Id int64) (*model.User, error) {
	user, err := u.userModel.FindOne(ctx, Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDAO) FindUser(ctx context.Context, userName string) (*model.User, error) {
	user, total, err := u.userModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"user_name": userName},
		},
	})
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return nil, nil
	}
	return user[0], nil
}
