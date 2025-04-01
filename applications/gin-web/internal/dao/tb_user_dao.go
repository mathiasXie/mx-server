package dao

import (
	"context"

	"github.com/mathiasXie/curd"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
)

type TbUserDAO struct {
	userModel *curd.Model[model.TbUser]
}

func (u *TbUserDAO) AddUser(ctx context.Context, user *model.TbUser) (int64, error) {
	insert, err := u.userModel.Insert(ctx, map[string]interface{}{
		"user_name": user.UserName,
		"email":     user.Email,
		"password":  user.Password,
	})
	if err != nil {
		return 0, err
	}
	return insert.UserId, nil
}

func (u *TbUserDAO) GetAllUser(ctx context.Context) ([]*model.TbUser, int64, error) {
	users, total, err := u.userModel.SelectAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
