package webroutes

import (
	"OttBot2/models"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func PostRegister(c *gin.Context) {
	newUser := &models.User{}
	c.BindJSON(&newUser)
	newUser.IsAdmin = false //Grr
	err := models.InsertUser(newUser.UserName, newUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"success": 1, "message": "Registered"})
		return
	}
}
