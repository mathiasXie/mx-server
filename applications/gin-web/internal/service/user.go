package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/gin-web/dto"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/dao"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
	"github.com/mathiasXie/gin-web/applications/gin-web/loader"
	"github.com/mathiasXie/gin-web/pkg/logger"
)

type UserService struct {
	userDAO *dao.UserDAO
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{
		userDAO: dao.NewUserDAO(loader.GetDB(ctx)),
	}
}

func (u *UserService) AddUser(ctx context.Context, req *dto.AddUserRequest) error {
	user := &model.User{UserName: req.UserName}
	_, err := u.userDAO.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUsers(ctx *gin.Context) ([]*dto.UserResponse, int64, error) {
	users, total, err := u.userDAO.GetAllUser(ctx)
	if err != nil {
		logger.CtxError(ctx, "[UserService]GetUsers Error:", err.Error())
		return nil, 0, err
	}
	return u.transUsers2Resp(users), total, nil
}

//func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
//	user, err := u.userDAO.FindUser(ctx, req.UserName)
//	if err != nil {
//		return nil, err
//	}
//	if user == nil {
//		return nil, nil
//	}
//	token, _ := utils.MakeToken(uint64(user.Id), user.UserName)
//	return &dto.LoginResponse{
//		Id:    user.Id,
//		Token: token,
//	}, nil
//}

func (u *UserService) transUsers2Resp(users []*model.User) []*dto.UserResponse {

	resp := make([]*dto.UserResponse, 0)
	for _, user := range users {
		resp = append(resp, u.transUser2Resp(user))
	}
	return resp
}

func (u *UserService) transUser2Resp(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       user.Id,
		UserName: user.UserName,
	}
}
