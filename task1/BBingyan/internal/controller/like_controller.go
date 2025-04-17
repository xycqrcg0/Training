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

var (
	HASLIKED   = "2"
	LIKE       = "1"
	UNLIKE     = "0"
	NEEDUNLIKE = "-1"
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

	//先确定email的合法性
	emailKey := fmt.Sprintf("email:%s", data.LikedUser)
	if v, err := global.RedisDB.Get(emailKey).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			if _, er := model.GetUserByEmail(data.LikedUser); er != nil {
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

	likeKey := fmt.Sprintf("userlike:%s:%s", email, data.LikedUser)
	likesKey := fmt.Sprintf("userlikes:%s", data.LikedUser)

	//再获取点赞数/创建该kv对
	if _, err := global.RedisDB.Get(likesKey).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			//读数据库
			if l, err := model.GetUserLikes(data.LikedUser); err != nil {
				log.Errorf("Fail to read postgres,error:%v", err)
				return c.JSON(http.StatusInternalServerError, &param.Response{
					Status: false,
					Msg:    "Internal Server",
				})
			} else {
				if _, er := global.RedisDB.Set(likesKey, l, time.Hour*5).Result(); er != nil {
					log.Errorf("Fail to write in redis,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			}
		}
	}

	//j用来判断之后要怎么记录kv对
	j, err := global.RedisDB.Get(likeKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			//查数据库
			ok, er := model.HasLikeUserShip(email, data.LikedUser)
			if er != nil {
				log.Errorf("Fail to read postgres,error:%v", er)
				return c.JSON(http.StatusInternalServerError, &param.Response{
					Status: false,
					Msg:    "Internal Server",
				})
			}
			if !ok {
				//数据库里没有该条记录 -> 未写入的未点赞状态
				j = UNLIKE
			} else {
				//数据库里已经有该条记录 -> 已写入的已点赞状态
				j = HASLIKED
			}
		}
	}

	var likes string
	_, er := global.RedisDB.TxPipelined(func(pipe redis.Pipeliner) error {
		switch j {
		case HASLIKED:
			pipe.Set(likeKey, NEEDUNLIKE, time.Hour*5)
			pipe.Decr(likesKey)
			break
		case NEEDUNLIKE:
			pipe.Set(likeKey, HASLIKED, time.Hour*5)
			pipe.Incr(likesKey)
			break
		case UNLIKE:
			pipe.Set(likeKey, LIKE, time.Hour*5)
			pipe.Incr(likesKey)
			break
		case LIKE:
			pipe.Set(likeKey, UNLIKE, time.Hour*5)
			pipe.Decr(likesKey)
			break
		}
		likes, _ = pipe.Get(likesKey).Result()
		return nil
	})
	if er != nil {
		log.Errorf("Fail to finish redis transaction,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal Server",
		})
	}

	return c.JSON(http.StatusOK, &param.Response{
		Status: true,
		Msg:    "",
		Data: &param.LikeResponse{
			Likes: likes,
		},
	})
}

func LikePost(c echo.Context) error {
	email := c.Get("identification").(string)
	var data param.PostLikeRequest
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, &param.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	//先确定post的合法性
	postKey := fmt.Sprintf("post:%d", data.LikedPost)
	if v, err := global.RedisDB.Get(postKey).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			if _, er := model.GetPostById(data.LikedPost); er != nil {
				if errors.Is(er, gorm.ErrRecordNotFound) {
					if _, e := global.RedisDB.Set(postKey, param.INVALID, time.Minute*5).Result(); e != nil {
						log.Errorf("Fail to write in redis,error:%v", err)
						return c.JSON(http.StatusInternalServerError, &param.Response{
							Status: false,
							Msg:    "Internal Server",
						})
					}
					return c.JSON(http.StatusBadRequest, &param.Response{
						Status: false,
						Msg:    "nonexistent post",
					})
				} else {
					log.Errorf("Fail to read postgres,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			} else {
				if _, e := global.RedisDB.Set(postKey, param.VALID, time.Minute*5).Result(); e != nil {
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

	likeKey := fmt.Sprintf("postlike:%s:%d", email, data.LikedPost)
	likesKey := fmt.Sprintf("postlikes:%d", data.LikedPost)

	//再获取点赞数/创建该kv对
	if _, err := global.RedisDB.Get(likesKey).Result(); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			//读数据库
			if l, err := model.GetPostLikes(data.LikedPost); err != nil {
				log.Errorf("Fail to read postgres,error:%v", err)
				return c.JSON(http.StatusInternalServerError, &param.Response{
					Status: false,
					Msg:    "Internal Server",
				})
			} else {
				if _, er := global.RedisDB.Set(likesKey, l, time.Hour*5).Result(); er != nil {
					log.Errorf("Fail to write in redis,error:%v", err)
					return c.JSON(http.StatusInternalServerError, &param.Response{
						Status: false,
						Msg:    "Internal Server",
					})
				}
			}
		}
	}

	//j用来判断之后要怎么记录kv对
	j, err := global.RedisDB.Get(likeKey).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Errorf("Fail to read redis,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal Server",
			})
		} else {
			//查数据库
			ok, er := model.HasLikePostShip(email, data.LikedPost)
			if er != nil {
				log.Errorf("Fail to read postgres,error:%v", er)
				return c.JSON(http.StatusInternalServerError, &param.Response{
					Status: false,
					Msg:    "Internal Server",
				})
			}
			if !ok {
				//数据库里没有该条记录 -> 未写入的未点赞状态
				j = UNLIKE
			} else {
				//数据库里已经有该条记录 -> 已写入的已点赞状态
				j = HASLIKED
			}
		}
	}

	var likes string
	_, er := global.RedisDB.TxPipelined(func(pipe redis.Pipeliner) error {
		switch j {
		case HASLIKED:
			pipe.Set(likeKey, NEEDUNLIKE, time.Hour*5)
			pipe.Decr(likesKey)
			break
		case NEEDUNLIKE:
			pipe.Set(likeKey, HASLIKED, time.Hour*5)
			pipe.Incr(likesKey)
			break
		case UNLIKE:
			pipe.Set(likeKey, LIKE, time.Hour*5)
			pipe.Incr(likesKey)
			break
		case LIKE:
			pipe.Set(likeKey, UNLIKE, time.Hour*5)
			pipe.Decr(likesKey)
			break
		}
		likes, _ = pipe.Get(likesKey).Result()
		return nil
	})
	if er != nil {
		log.Errorf("Fail to finish redis transaction,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal Server",
		})
	}

	return c.JSON(http.StatusOK, &param.Response{
		Status: true,
		Msg:    "",
		Data: &param.LikeResponse{
			Likes: likes,
		},
	})
}
