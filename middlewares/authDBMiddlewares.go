package middlewares

import (
	"alta/be4/mvc/config"
	"alta/be4/mvc/models"

	"github.com/labstack/echo/v4"
)

func BasicAuthDB(email, password string, c echo.Context) (bool, error) {
	var db = config.DB
	var user models.User
	tx := db.Where("email=? AND password=?", email, password).First(&user)
	if tx.Error != nil {
		return false, nil
	}
	return true, nil
}
