package controllers

import (
	"fmt"
	"net/http"

	"github.com/ankitanwar/GoAPIUtils/errors"
	oauth "github.com/ankitanwar/e-Commerce/Middleware/oAuth"
	"github.com/ankitanwar/e-Commerce/User/domain/users"
	"github.com/ankitanwar/e-Commerce/User/services"
	"github.com/gin-gonic/gin"
)

func getUserid(request *http.Request) (string, *errors.RestError) {
	userID := request.Header.Get("X-Caller-Id")
	if userID == "" {
		return "", errors.NewBadRequest("Invalid User ID")
	}
	return userID, nil
}

//CreateUser : To create the user
func CreateUser(c *gin.Context) {
	var newUser users.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		err := errors.NewBadRequest("Invalid Request")
		c.JSON(err.Status, err)
		return
	}

	result, saverr := services.UserServices.CreateUser(newUser)
	if saverr != nil {
		c.JSON(saverr.Status, saverr)
		return
	}
	c.JSON(http.StatusCreated, result.MarshallUser(oauth.IsPublic(c.Request)))
}

//GetUser : To get the user from the database
func GetUser(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	userid, userErr := getUserid(c.Request)
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	user, err := services.UserServices.GetUser(userid)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.MarshallUser(oauth.IsPublic(c.Request)))

}

//UpdateUser :To Update the value of particaular user
func UpdateUser(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	var user = users.User{}
	userid, userErr := getUserid(c.Request)
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	user.ID = userid
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequest("Invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	updatedUser, err := services.UserServices.UpdateUser(user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, updatedUser.MarshallUser(oauth.IsPublic(c.Request)))
}

//DeleteUser :To Delete the user with given id
func DeleteUser(c *gin.Context) {
	userid, userErr := getUserid(c.Request)
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	if err := services.UserServices.DeleteUser(userid); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"Status": "User Deleted"})
}

//FindByStatus : To find all the users by given status
func FindByStatus(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UserServices.FindByStatus(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.MarshallUser(oauth.IsPublic(c.Request)))

}

//Login : to verify user email and password
func Login(c *gin.Context) {
	verifyUser := users.LoginRequest{}
	if err := c.ShouldBindJSON(&verifyUser); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	user, err := services.UserServices.LoginUser(verifyUser)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.MarshallUser(oauth.IsPublic(c.Request)))

}

//GetAddress : To Get the address of the given user
func GetAddress(c *gin.Context) {
	err := oauth.AuthenticateRequest(c.Request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	userID, err := getUserid(c.Request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	fmt.Println("The value of userID is ", userID)
	address, err := services.UserServices.GetAddress(userID)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusAccepted, address)
}

//AddAddress : To Get the address of the given user
func AddAddress(c *gin.Context) {
	err := oauth.AuthenticateRequest(c.Request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	address := &users.UserAddress{}
	bindErr := c.ShouldBindJSON(address)
	fmt.Println("The value of bindErr is", bindErr)
	if bindErr != nil {
		c.JSON(http.StatusBadRequest, "Error while binding to the json")
		return
	}
	userID, err := getUserid(c.Request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	err = services.UserServices.AddAddress(userID, address)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusAccepted, "Address has been added successfully")
}
