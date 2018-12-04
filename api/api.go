package api

import (
	"livego/backends"
	"livego/backends/models"
	"livego/config"
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

var identityKey = "id"

// Start : This function start the GIN http server.
func Start(port string) {
	if port == "" {
		port = "8080"
	}
	log.Println("Starting API")
	r := gin.Default()
	r.Delims("{[{", "}]}")

	// Setup JWT Authentication Middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if userID == config.Data.Username && password == config.Data.Password {
				return &User{
					UserName:  userID,
					LastName:  "User",
					FirstName: "Admin",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.LoadHTMLFiles("./www/index.html")
	// r.Static("/assets", "./www./assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"user": "bcrowe"})
	})
	r.Use(static.Serve("/assets", static.LocalFile("./www/assets", true)))
	v := r.Group("/api")
	v.Use(authMiddleware.MiddlewareFunc())
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
