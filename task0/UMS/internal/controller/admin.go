package controller

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
	"ums/internal/config"
	"ums/internal/controller/params"
	"ums/internal/models"
	"ums/internal/utils"
)

func AdminLogin(c echo.Context) error {
	admin := &params.AdminRequest{}
	if err := c.Bind(admin); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	adminUser, err := models.GetAdminByName(admin.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	if err := utils.ComparePassword(adminUser.Password, admin.Password); err != nil {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "Wrong name or password",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"identification": admin.Name,
		"role":           "admin",
		"exp":            time.Now().Add(time.Second * time.Duration(config.Config.JWT.Exp)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(config.Config.JWT.Key))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    "Fail to generate token",
		})
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "admin login successfully",
		Data: &params.TokenResponse{
			Token: "Bearer " + signedToken,
		},
	})

}

func AddNewAdmin(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not an admin ",
		})
	}
	adminReq := &params.AdminRequest{}
	if err := c.Bind(adminReq); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}
	if adminReq.Name == "" || adminReq.Password == "" {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	_, err := models.GetAdminByName(adminReq.Name)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusInternalServerError, &params.Response{
				Status: false,
				Msg:    err.Error(),
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    "Name has existed",
		})
	}

	hashed, err := utils.HashPassword(adminReq.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	newAdmin := &models.Admin{
		Name:     adminReq.Name,
		Password: hashed,
	}

	err = models.AddAdmin(newAdmin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, &params.Response{
		Status: true,
		Msg:    "add new admin successfully",
	})
}

// DeleteAdmin 这个就删自己吧
func DeleteAdmin (c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not an admin ",
		})
	}
	name := c.Get("identification").(string)
	err := models.DeleteAdmin(name)
	if err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return c.JSON(http.StatusBadRequest, &params.Response{
				Status: false,
				Msg:    "nonexistent name",
			})
		}else {
			return c.JSON(http.StatusInternalServerError, &params.Response{
				Status: false,
				Msg:    err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, &params.Response{
		Status: false,
		Msg:    "delete admin successfully",
	})
}

func GetUserInfo(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not a admin",
		})
	}

	data := &params.HandleUserRequest{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	user, err := models.GetUserByEmail(data.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, &params.Response{
				Status: false,
				Msg:    "Nonexistent user",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, &params.Response{
				Status: false,
				Msg:    err.Error(),
			})
		}
	}

	res := &params.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "",
		Data:   res,
	})
}

func GetAllUsers(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not an admin",
		})
	}

	users, err := models.GetAllUsersInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	usersInfo := make([]params.UserResponse, 0)
	for _, user := range users {
		usersInfo = append(usersInfo, params.UserResponse{
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "",
		Data:   usersInfo,
	})
}

func DropUserByEmail(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, &params.Response{
			Status: false,
			Msg:    "not a admin",
		})
	}

	data := &params.HandleUserRequest{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	_, err := models.GetUserByEmail(data.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, &params.Response{
				Status: false,
				Msg:    "Nonexistent user",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, &params.Response{
				Status: false,
				Msg:    err.Error(),
			})
		}
	}

	err = models.DeleteUser(data.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &params.Response{
			Status: false,
			Msg:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &params.Response{
		Status: true,
		Msg:    "drop user successfully",
	})
}
