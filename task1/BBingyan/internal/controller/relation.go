package controller

import (
	"BBingyan/internal/controller/param"
	"BBingyan/internal/global"
	"BBingyan/internal/log"
	"BBingyan/internal/model"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func FollowUser(c echo.Context) error {
	user := c.Get("identification").(string)
	followed := c.Param("email")

	//先确定email的合法性
	emailKey := fmt.Sprintf("email:%s", followed)
	if v, err := global.RedisDB.Get(emailKey).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			if _, er := model.GetUserByEmail(followed); er != nil {
				if errors.Is(er, gorm.ErrRecordNotFound) {
					if _, e := global.RedisDB.Set(emailKey, param.INVALID, time.Minute*5).Result(); e != nil {
						log.Errorf("Fail to write in redis,error:%v", err)
						return c.JSON(http.StatusInternalServerError, &param.Response{
							Status: false,
							Msg:    "Internal Server",
						})
					}
					return c.JSON(http.StatusBadRequest, &param.Response{
						Status: false,
						Msg:    "nonexistent email",
					})
				} else {
					log.Errorf("Fail to read postgres,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			} else {
				if _, e := global.RedisDB.Set(emailKey, param.VALID, time.Minute*5).Result(); e != nil {
					log.Errorf("Fail to write in redis,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			}
		}
	} else if v == param.INVALID {
		return c.JSON(http.StatusBadRequest, &param.Response{
			Status: false,
			Msg:    "nonexistent email",
		})
	}

	k := fmt.Sprintf("followship:%s:%s", user, followed)
	if _, err := global.RedisDB.Incr(k).Result(); err != nil {
		log.Errorf("Fail to write in redis,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal Server",
		})
	}

	err := model.FollowUser(user, followed)
	if err != nil {
		log.Errorf("Fail to write in postgres,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal server",
		})
	}

	return c.JSON(http.StatusOK, &param.Response{
		Status: true,
		Msg:    "",
	})
}

func UnFollowUser(c echo.Context) error {
	user := c.Get("identification").(string)
	unfollowed := c.Param("email")

	//先确定email的合法性
	emailKey := fmt.Sprintf("email:%s", unfollowed)
	if v, err := global.RedisDB.Get(emailKey).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			if _, er := model.GetUserByEmail(unfollowed); er != nil {
				if errors.Is(er, gorm.ErrRecordNotFound) {
					if _, e := global.RedisDB.Set(emailKey, param.INVALID, time.Minute*5).Result(); e != nil {
						log.Errorf("Fail to write in redis,error:%v", err)
						return c.JSON(http.StatusInternalServerError, &param.Response{
							Status: false,
							Msg:    "Internal Server",
						})
					}
					return c.JSON(http.StatusBadRequest, &param.Response{
						Status: false,
						Msg:    "nonexistent email",
					})
				} else {
					log.Errorf("Fail to read postgres,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			} else {
				if _, e := global.RedisDB.Set(emailKey, param.VALID, time.Minute*5).Result(); e != nil {
					log.Errorf("Fail to write in redis,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			}
		}
	} else if v == param.INVALID {
		return c.JSON(http.StatusBadRequest, &param.Response{
			Status: false,
			Msg:    "nonexistent email",
		})
	}

	err := model.UnfollowUser(user, unfollowed)
	if err != nil {
		log.Errorf("Fail to write in postgres,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal server",
		})
	}

	return c.JSON(http.StatusOK, &param.Response{
		Status: true,
		Msg:    "",
	})
}

func GetFollows(c echo.Context) error {
	user := c.Get("identification").(string)

}
