package controller

import (
	"blog/api/service"
	"blog/models"
	"blog/util"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	token := c.GetHeader("token")
	if !util.CheckToken("2806374351z@gmail.com", token) {
		util.ErrorJSON(c, http.StatusBadRequest, "token illegal")
		return
	}
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

//DeleteUser -> calls DeleteUser service to delete user
func (u *UserController) DeleteUser(c *gin.Context) {
	token := c.GetHeader("token")
	if !util.CheckToken("2806374351z@gmail.com", token) {
		util.ErrorJSON(c, http.StatusBadRequest, "token illegal")
		return
	}

	idParam := c.Param("id")
	fmt.Println(idParam)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "id invalid")
		return
	}

	err = u.service.DeleteUser(id)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Failed to delete user")
		return
	}
	response := &util.Response{
		Success: true,
		Message: "Deleted Sucessfully",
	}
	util.SuccessJSON(c, http.StatusOK, response)
}

//LoginUser : Generates JWT Token for validated user
func (u *UserController) LoginUser(c *gin.Context) {
	var user models.UserLogin
	salt := os.Getenv("TOKEN_SALT")
	var hmacSampleSecret []byte = []byte(salt)
	if err := c.ShouldBindJSON(&user); err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Inavlid Json Provided")
		return
	}
	dbUser, err := u.service.LoginUser(user)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "Invalid Login Credentials")
		return
	}
	claims := models.MyClaims{Email: dbUser.Email, StandardClaims: jwt.StandardClaims{
		NotBefore: time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 3600*72,
		Issuer:    "verda",
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
	fmt.Println("生成的token: ", tokenString)
	c.JSON(http.StatusOK, response)
}
