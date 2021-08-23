package core

import (
	"fmt"
	"Course-Net/final_project2/config"
	model "Course-Net/final_project2/models"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func Register(c *gin.Context) {
	var (
		user model.User
		//err error
	)

	if (c.PostForm("username") == "" ||
		c.PostForm("fullname") == "" ||
		c.PostForm("email") == "" ||
		c.PostForm("address") == "" ||
		c.PostForm("password") == "") {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Status Bad Request",
			"message" : "Username, fullname, email, adress, password, and role is required",
		})
		c.Abort()
		return
	}
	user.Username = c.PostForm("username")
	user.Fullname = c.PostForm("fullname")
	user.Email = c.PostForm("email")
	user.Address = c.PostForm("address")
	user.Password = encryptPassword(c.PostForm("password"))
	user.Role, _ = strconv.ParseInt(c.PostForm("role"), 10, 64)

	result := config.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"nessage" : fmt.Sprintf("%s", result.Error),
		})
		c.Abort()
		return
	}

	token := createToken(&user)
	var userToken model.User
	result = config.Db.First(&userToken, "id =?", user.ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"message" : "Erro when getting user data",
		})
		c.Abort()
		return
	}
	userToken.Token = token
	result = config.Db.Save(&userToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"message" : "Error when creating and saving token",
		})
		c.Abort()
		return
	}
	c.Set("user_token", token)
	c.JSON(http.StatusOK, gin.H{
		"message" : "Register success",
		"user" : user,
		"token" : token,
	})
}

func Login (c *gin.Context) {

	if (c.PostForm("username") == "" ||
		c.PostForm("password") == "") {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Status Bad Request",
			"message" : "Username and password is required",
		})
		c.Abort()
		return
	}

	username := c.PostForm("username");
	password := encryptPassword(c.PostForm("password"));
	
	var (
		user model.User
	)

	result := config.Db.First(&user, "username = ? and password = ?", username, password)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status" : "Not Founf",
			"nessage" : fmt.Sprintf("%s", "User Not Found"),
		})
		c.Abort()
		return
	} else {
		token := createToken(&user)
		user.Token = token
		saveToken := config.Db.Save(&user)
		if saveToken.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status" : "Internal Server Error",
				"messages" : fmt.Sprintf("%s", saveToken.Error),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status" : "Success",
			"messages" : "Found User",
			"token" : token,
		})
	}
}

func Logout (c *gin.Context) {

	authHeader := c.Request.Header.Get("Authorization")
	auth := strings.Split(authHeader, " ")
	token, _ := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	claim, _ := token.Claims.(jwt.MapClaims)
	userID := uint(claim["user_id"].(float64))
	var user model.User
	err := config.Db.First(&user, "id = ?", userID)
	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status" : "Not Found",
			"messages" : "User Not Found",
		})
		c.Abort()
	} else {
		saveUser := config.Db.Model(&user).Update("token", nil)
		if saveUser.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status" : "Internal Server Error",
				"messages" : fmt.Sprintf("%s", saveUser.Error),
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status" : "success",
			"messages" : "Logout Success",
		})
	}
}

func encryptPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(h[:])
}

func createToken(user *model.User) string {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.Role,
		"exp" : time.Now().AddDate(0, 0, 1).Unix(),
		"issued_at": time.Now().Unix(),
	})

	tokenString, _ := token.SignedString(key)
	return tokenString
}

func IsAuthenticated(token string) bool {
	var user model.User
	result := config.Db.Find(&user, "token = ?", token)
	if result.Error != nil {
		panic("Error when checing token")
	} else if user.ID == 0 {
		return false
	}

	return true
}