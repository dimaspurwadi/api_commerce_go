package middleware

import "Course-Net/final_project2/core"
import "net/http"
import "os"
import "strings"
import "github.com/dgrijalva/jwt-go"
import "github.com/gin-gonic/gin"

func isAdmin() gin.HandlerFunc {
	return CheckJWT(1)
}

func CheckJWT(role uint) gin.HandlerFunc {
	var (
		token *jwt.Token
		err error
	)

	return func(c * gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		auth := strings.Split(authHeader, " ")

		if len(auth) == 2 {
			token, err = jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status" : "Internal Server Error",
					"message" : "Error when parsing JWT",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status" : "Unauthorized",
				"message" : "Token is required",
			})
			c.Abort()
			return
		}

		if token.Valid {
			isAuth := core.IsAuthenticated(auth[1])

			if isAuth {
				claim, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status" : "Internal Server Error",
						"message" : "Failed to claim this token",
					})
					c.Abort()
					return
				}

				user_id := uint(claim["user_id"].(float64))
				role_id := uint(claim["role_id"].(float64))

				if role == role_id {
					c.Set("user_id", user_id)
					c.Set("role_id", role_id)
				} else {
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"status" : "Unprocessable Entity",
						"message" : " You are not allowed to access thin endpoint",
					})
					c.Abort()
					return
				}
 			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "AUthorization required, please login/register first",
				})
				c.Abort()
				return
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "Invalid token format",
				})
				c.Abort()
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "Token has expired",
				})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "Invalid token",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status" : "Unauthorized",
				"message" : "Can't proccess this token",
			})
			c.Abort()
			return
		}
	}
}

func CheckJWTWithoutRole() gin.HandlerFunc {
	var (
		token *jwt.Token
		err error
	)

	return func(c * gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		auth := strings.Split(authHeader, " ")

		if len(auth) == 2 {
			token, err = jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status" : "Internal Server Error",
					"message" : "Error when parsing JWT",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status" : "Unauthorized",
				"message" : "Token is required",
			})
			c.Abort()
			return
		}

		if token.Valid {
			isAuth := core.IsAuthenticated(auth[1])

			if isAuth {
				_, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status" : "Internal Server Error",
						"message" : "Failed to claim this token",
					})
					c.Abort()
					return
				}
 			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "AUthorization required, please login/register first",
				})
				c.Abort()
				return
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "Invalid token format",
				})
				c.Abort()
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "Token has expired",
				})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status" : "Unauthorized",
					"message" : "Invalid token",
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status" : "Unauthorized",
				"message" : "Can't proccess this token",
			})
			c.Abort()
			return
		}
	}
}