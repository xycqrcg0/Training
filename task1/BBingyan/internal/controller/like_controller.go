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
	"strconv"
	"time"
)

func LikeUser(c echo.Context) error {
	email := c.Get("identification").(string)
	var data param.UserLikeRequest
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, &param.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}
	likesKey := fmt.Sprintf("userlikes:%s", data.LikedUser)
	likeKey := fmt.Sprintf("userlikeship:%s:%s", email, data.LikedUser)
	unlikeKey := fmt.Sprintf("userunlikeship:%s:%s", email, data.LikedUser)
	//区分一下,likeKey的值：0->数据库里无该条记录；1->数据库里有该条记录

	//先判断一下有没有点过赞
	ok := false
	if j, err := global.RedisDB.Get(likeKey).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			var ee error
			ok, ee = model.HasLikeUserShip(email, data.LikedUser)
			if ee != nil {
				log.Errorf("Fail to read likeship from postgres,error:%v", err)
				return c.JSON(http.StatusInternalServerError, &param.Response{
					Status: false,
					Msg:    "Internal server error",
				})
			}
			//说明redis里没有，数据库里有信息，可能是重复点赞
			if ok {
				//如果没有取消过点赞，那就是重复点赞
				if _, err := global.RedisDB.Get(unlikeKey).Result(); err != nil {
					if !errors.Is(err, redis.Nil) {
						log.Errorf("Fail to read likeship from redis.error:%v", err)
						return c.JSON(http.StatusInternalServerError, &param.Response{
							Status: false,
							Msg:    "Internal server error",
						})
					}
				} else {
					return c.JSON(http.StatusForbidden, &param.Response{
						Status: false,
						Msg:    "Has liked",
					})
				}
			}
		} else {
			log.Errorf("Fail to read likeship from redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal server error",
			})
		}
	} else {
		ok = j == "1"
	}

	var likesString string
	//先看redis里有没有点赞信息
	_, err := global.RedisDB.Get(likesKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//没有就在数据库里查查
			l, err := model.GetLikes(data.LikedUser)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					//没查到，那么这个email就是不合法的
					return c.JSON(http.StatusBadRequest, &param.Response{
						Status: false,
						Msg:    "nonexistent email",
					})
				} else {
					//500
					log.Errorf("Fail to read likes from postgres,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal server error",
					})
				}
			} else {
				//数据库里查到了，就把它写进redis
				if _, err := global.RedisDB.Set(likesKey, l, time.Hour*5).Result(); err != nil {
					log.Errorf("Fail to write in redis,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal server error",
					})
				}
			}
		} else {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal server error",
			})
		}
	}

	_, err = global.RedisDB.TxPipelined(func(pipe redis.Pipeliner) error {
		if _, e := pipe.Get(unlikeKey).Result(); e != nil {
			if !errors.Is(e, redis.Nil) {
				return e
			}
		} else {
			pipe.Del(unlikeKey)
		}
		pipe.Incr(likesKey)

		if ok {
			pipe.Set(likeKey, 1, time.Hour*5)
		} else {
			pipe.Set(likeKey, 0, time.Hour*5)
		}

		var e error
		likesString, e = pipe.Get(likesKey).Result()
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		log.Errorf("Fail to finish redis transaction,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal server error",
		})
	}

	likes, _ := strconv.Atoi(likesString)
	return c.JSON(http.StatusOK, &param.Response{
		Status: true,
		Msg:    "",
		Data: &param.UserLikeResponse{
			Likes: likes,
		},
	})
}
