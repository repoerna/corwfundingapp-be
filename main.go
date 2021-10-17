package main

import (
	"bwacrowdfunding/auth"
	"bwacrowdfunding/campaign"
	"bwacrowdfunding/handler"
	"bwacrowdfunding/helper"
	"bwacrowdfunding/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/bwacrowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("db connected")

	userRepo := user.NewRepository(db)
	campaignRepo := campaign.NewRepository(db)

	authSvc := auth.NewService()
	userSvc := user.NewService(userRepo)
	campaignSvc := campaign.NewService(campaignRepo)

	userHandler := handler.NewUserHandler(authSvc, userSvc)
	campaignHandler := handler.NewCampaignHandler(campaignSvc)

	router := gin.Default()
	router.Static("/images", "./images")

	apiV1 := router.Group("api/v1")
	apiV1.POST("/users", userHandler.RegisterUser)
	apiV1.POST("/sessions", userHandler.Login)
	apiV1.POST("/email_checkers", userHandler.CheckEmailAvailability)
	apiV1.POST("/avatars", authMiddleware(authSvc, userSvc), userHandler.UploadAvatar)

	apiV1.GET("/campaigns", campaignHandler.GetCampaigns)
	apiV1.GET("/campaigns/:id", campaignHandler.GetCampaign)
	apiV1.POST("/campaigns", authMiddleware(authSvc, userSvc), campaignHandler.CreateCampaign)
	apiV1.PUT("/campaigns/:id", authMiddleware(authSvc, userSvc), campaignHandler.UpdateCampaign)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		var tokenString string
		headerSplit := strings.Split(authHeader, " ")
		if len(headerSplit) == 2 {
			tokenString = headerSplit[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}
