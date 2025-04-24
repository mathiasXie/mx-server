package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/dto"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/model"
	"github.com/mathiasXie/gin-web/applications/xiaozhi-server/internal/service"
	"github.com/mathiasXie/gin-web/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserCenterHandler struct {
	UserService *service.UserService
}

func (h *UserCenterHandler) SignIn(ctx *gin.Context) {
	var signInRequest dto.SignInRequest
	if err := ctx.ShouldBindJSON(&signInRequest); err != nil {
		dto.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.UserService.GetUserByPhoneEmailOrUsername(signInRequest.Username)
	if err != nil {
		dto.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	// 验证密码，使用bcrypt.CompareHashAndPassword
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInRequest.Password)); err != nil {
		dto.Fail(ctx, http.StatusUnauthorized, "密码错误")
		return
	}

	// 生成token
	token, err := utils.MakeToken(int(user.Id), user.UserName)
	if err != nil {
		dto.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(ctx, gin.H{
		"token":  token,
		"userId": user.Id,
	})
}

func (h *UserCenterHandler) SignUp(ctx *gin.Context) {
	var signUpRequest dto.SignUpRequest
	if err := ctx.ShouldBindJSON(&signUpRequest); err != nil {
		dto.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if filed := h.UserService.CheckUserExists(signUpRequest); filed != "" {
		dto.Fail(ctx, http.StatusBadRequest, fmt.Sprintf("%s已存在", filed))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		dto.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	user := model.AiUser{
		UserName: signUpRequest.Username,
		Password: string(hashedPassword),
		Email:    signUpRequest.Email,
		Phone:    signUpRequest.Phone,
	}

	if err := h.UserService.CreateUser(user); err != nil {
		dto.Fail(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	dto.Success(ctx, gin.H{
		"message": "注册成功",
	})

}
