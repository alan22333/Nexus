package controller

import (
	"Nuxus/internal/res"
	"Nuxus/internal/service"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	tagService *service.TagService
}

func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

func (tc *TagController) ListTags(c *gin.Context) {
	sortedBy := c.Param("sort")
	if sortedBy == "" {
		sortedBy = "post_count"
	}

	resDto, err := tc.tagService.ListTags(sortedBy)
	if err != nil {
		c.Error(err)
		return
	}

	res.OkWithData(c, resDto)
}
