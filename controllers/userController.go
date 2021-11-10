package controllers

import (
	"alta/be4/mvc/lib/database"
	"alta/be4/mvc/middlewares"
	"alta/be4/mvc/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetUsersController(c echo.Context) error {
	//memanggil fungsi yang ada di folder lib/database
	users, err := database.GetUsers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "failed to load data",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "success load data users",
		"data":    users,
	})
}

func CreateUserController(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	result, err := database.CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "failed to create data",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "success create user",
		"data":    result,
	})
}

func GetOneUserController(c echo.Context) error {
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"message": "false param",
		})
	}
	user, rowAffected, err := database.GetUser(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if rowAffected == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "user id not found",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    user,
	})
}

func UpdateUserController(c echo.Context) error {
	// responses := map[string]interface{}{
	// 	"message": "failed to update",
	// }

	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"message": "false param",
		})
	}

	user := models.User{}
	c.Bind(&user)
	resultUser, rowsAffected, err := database.EditUser(&user, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed",
		})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong user id",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success update user",
		"data":    resultUser,
	})

}

func DeleteUserController(c echo.Context) error {
	id, e := strconv.Atoi(c.Param("id"))
	if e != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"message": "false param",
		})
	}
	_, rowAffected, err := database.DeleteUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if rowAffected == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "user id not found",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success delete user",
		"iduser":  id,
	})

}

func LoginUsersController(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	users, err := database.LoginUsers(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "failed login",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "success login",
		"data":    users,
	})
}

func GetUserDetailController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	loggedInUserId := middlewares.ExtractTokenUserId(c)

	if loggedInUserId != id {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
			"status":  "failed",
			"message": "access forbidden",
		})
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "failed fetch data",
		})
	}

	users, err := database.GetDetailUsers(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"status":  "failed",
			"message": "failed fetch data",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "success get detail user",
		"data":    users,
	})
}
