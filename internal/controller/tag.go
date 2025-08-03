package controller

import (
	"Nuxus/internal/res"
	"Nuxus/internal/service"

	"github.com/gin-gonic/gin"
)

func ListTags(c *gin.Context) {
	sortedBy := c.Param("sort")
	if sortedBy == "" {
		sortedBy = "post_count"
	}

	resDto, err := service.ListTags(sortedBy)
	if err != nil {
		c.Error(err)
		return
	}

	res.OkWithData(c, resDto)
}
