package dashboard

import (
	"context"
	"net/http"

	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AdminUi struct {
	r      *gin.Engine
	client client.Client
}

func NewAdminUi(c client.Client) *AdminUi {
	r := gin.Default()
	r.StaticFS("/adminui", http.Dir("./adminui"))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	return &AdminUi{r: r, client: c}
}

func (this AdminUi) Start(context.Context) error {
	this.r.GET("/configs", func(c *gin.Context) {
		list := v1.DbConfigList{}
		if err := this.client.List(c, &list); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(200, list.Items)
	})
	return this.r.Run(":9003")
}
