package dashboard

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AdminUi struct {
	r *gin.Engine
}

func NewAdminUi() *AdminUi {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
	return &AdminUi{r: r}
}

func (this AdminUi) Start(context.Context) error {
	return this.r.Run(":9003")
}
