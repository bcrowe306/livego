package api

import (
	"livego/backends"
	"livego/backends/models"
	"livego/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Start : This function start the GIN http server.
func Start(port string) {
	if port == "" {
		port = "8080"
	}
	log.Println("Starting API")
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.LoadHTMLFiles("./www/index.html")
	r.Static("/assets", "./www./assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"user": "bcrowe"})
	})
	v := r.Group("/api")
	{
		// /api/routes ---------
		v.GET("/routes", func(c *gin.Context) {
			r := backends.DB.SelectAll()
			c.JSON(http.StatusOK, r)
		})
		v.GET("/routes/:id", func(c *gin.Context) {
			id := c.Param("id")
			r, exist := backends.DB.Select(id)
			if !exist {
				c.JSON(http.StatusBadRequest, r)
			}
			c.JSON(http.StatusOK, r)
		})
		v.POST("/routes", func(c *gin.Context) {
			var rt models.Route
			err := c.BindJSON(&rt)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Error in request body.",
				})
			}
			newRoute, err := backends.DB.Insert(rt)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Error Insterting and saving route.",
				})
			}
			c.JSON(http.StatusOK, newRoute)
		})
		v.DELETE("/routes/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := backends.DB.Delete(id); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Error Insterting and saving route.",
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Route deleted successfully",
			})
		})
		v.PUT("/routes/:id", func(c *gin.Context) {
			id := c.Param("id")
			var rt models.Route
			c.BindJSON(&rt)
			updatedRoute, err := backends.DB.Update(id, rt)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Error Updating route and saving route.",
				})
			}
			c.JSON(http.StatusOK, updatedRoute)
		})

		// /api/config ---------------------------------
		v.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, config.Data)
		})
	}

	r.Run(":" + port)
}
