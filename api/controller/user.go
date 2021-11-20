package controller

import (
	"blog/api/service"
	"blog/models"
	"blog/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

//UserController struct
type UserController struct {
	service service.UserService
}

//NewUserController : NewUserController
func NewUserController(s service.UserService) UserController {
	return UserController{
		service: s,
	}
}

// FindAllUser -> call FindAllUser service for get user info
func (u UserController) FindAllUser(c *gin.Context) {
	token := c.Request.Header.Get("token")
	fmt.Printf(token)
	users, total, err := u.service.FindAllUser()
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Fail to find queston")
		fmt.Println(err)
		return
	}

	respArr := make([]gin.H, 0, 0)

	for _, user := range *users {
		resp := user.ResponseMap()
		respArr = append(respArr, resp)
	}

	c.JSON(http.StatusOK, &util.Response{
		Success: true,
		Message: "User result set",
		Data: gin.H{
			"rows":      respArr,
			"total_row": total,
		},
	})
}

//CreateUser ->  calls CreateUser services for validated user
func (u *UserController) CreateUser(c *gin.Context) {
	var user models.UserRegister
	if err := c.ShouldBind(&user); err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Inavlid Json Provided")
		return
	}

	hashPassword, _ := util.HashPassword(user.Password)
	user.Password = hashPassword

	err := u.service.CreateUser(user)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Failed to create user")
		return
	}

	util.SuccessJSON(c, http.StatusOK, "Successfully Created user")
}

//LoginUser : Generates JWT Token for validated user
func (u *UserController) LoginUser(c *gin.Context) {
	var user models.UserLogin
	var hmacSampleSecret []byte
	if err := c.ShouldBindJSON(&user); err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Inavlid Json Provided")
		return
	}
	dbUser, err := u.service.LoginUser(user)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Invalid Login Credentials")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": dbUser,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Failed to get token")
		return
	}
	response := &util.Response{
		Success: true,
		Message: "Token generated sucessfully",
		Data:    tokenString,
	}
	c.JSON(http.StatusOK, response)
}
