package controller

import (
	"BBingyan/internal/config"
	"BBingyan/internal/controller/param"
	"BBingyan/internal/log"
	"BBingyan/internal/model"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func AddPost(c echo.Context) error {
	user := c.Get("identification").(string)
	var data param.PostRequest
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, &param.Response{
			Status: false,
			Msg:    "Invalid request",
		})
	}

	j := false
	for _, tag := range config.Config.Curd.Tags {
		if tag == data.Tag {
			j = true
			break
		}
	}
	if !j {
		return c.JSON(http.StatusBadRequest, &param.Response{
			Status: false,
			Msg:    "Nonexistent tag",
		})
	}

	newPost := &model.Post{
		Author:  user,
		Title:   data.Title,
		Tag:     data.Tag,
		Content: data.Content,
	}

	err := model.AddPost(newPost)
	if err != nil {
		log.Errorf("Fail to write in postgres,error:%v", err)
		return c.JSON(http.StatusInternalServerError, &param.Response{
			Status: false,
			Msg:    "Internal server",
		})
	}

	return c.JSON(http.StatusCreated, &param.Response{
		Status: true,
		Msg:    "Create post Successfully",
	})
}

func DeletePost(c echo.Context) error {
	user := c.Get("identification").(string)
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)

	err := model.DeletePostById(user, id)
	if err != nil {
		if errors.Is(err, errors.New("none")) {
			return c.JSON(http.StatusBadRequest, &param.Response{
				Status: false,
				Msg:    "Invalid request",
			})
		} else {
			log.Errorf("Fail to write in postgres,error:%v", err)
			return c.JSON(http.StatusInternalServerError, &param.Response{
				Status: false,
				Msg:    "Internal server",
			})
		}
	}

	return c.JSON(http.StatusOK, &param.Response{
		Status: false,
		Msg:    "Delete post successfully",
	})
}
