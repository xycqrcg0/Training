package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"ums/internal/controller/params"
	"ums/internal/models"
	"ums/internal/utils"
)

// Login 要求email不重复，用email登录
func Login(c echo.Context) error {
	userInfo := &params.UserLoginRequest{}
	if err := c.Bind(&userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	user, err := models.GetUserByEmail(userInfo.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusForbidden, &params.Response{
				Status: false,
				Msg:    "Wrong email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "Wrong email or password",
		})
	}

	//身份认证成功，生成jwt
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "Login successfully",
		Data: &params.TokenResponse{
			Token: token,
		},
	})
}

func Register(c echo.Context) error {
	userInfo := &params.UserRegisterRequest{}
	if err := c.Bind(&userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	if userInfo.Password == "" {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}
	userInfo.Password = string(hashed)

	if userInfo.Email == "" || userInfo.Name == "" {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}
	//检查邮箱格式

	//检查邮箱是否存在
	_, err = models.GetUserByEmail(userInfo.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, &params.Response{
				Status: false,
				Msg:    err.Error(),
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest,&params.Response{
			Status: false,
			Msg: "Email has existed",
		})
	}

	//填入数据库
	newUser := &models.User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Password: userInfo.Password,
	}
	if err := models.AddUser(newUser); err != nil {
		return c.JSON(http.StatusInternalServerError,&params.Response{
			Status: false,
			Msg: err.Error(),
		})
	}

	//返回一个token
	token, err := utils.GenerateJWT(newUser.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,&params.Response{
			Status: false,
			Msg: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated,&params.Response{
		Status: true,
		Msg: "Register successfully",
		Data: &params.TokenResponse{
			Token: token,
		},
	})
}
