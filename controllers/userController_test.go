package controllers

import (
	"alta/be4/mvc/config"
	"alta/be4/mvc/constants"
	"alta/be4/mvc/middlewares"
	"alta/be4/mvc/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

//with slice
type UsersResponseSuccess struct {
	Status  string
	Message string
	Data    []models.User
}

//without slice
type SingleUserResponseSuccess struct {
	Status  string
	Message string
	Data    models.User
}

type ResponseFailed struct {
	Status  string
	Message string
}

func InitEchoTestAPI() *echo.Echo {
	config.InitDBTest()
	e := echo.New()
	return e
}

var (
	mock_data_user = models.User{
		Name:     "alta",
		Email:    "alta@gmail.com",
		Password: "12345",
	}
)

func InsertMockDataUserToDB() error {
	var err error
	if err = config.DB.Save(&mock_data_user).Error; err != nil {
		return err
	}
	return nil
}

func TestGetUsersControllerSuccess(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		expectSize int
	}{
		{
			name:       "success to get all data users",
			path:       "/users",
			expectCode: http.StatusOK,
			expectSize: 1,
		},
	}

	e := InitEchoTestAPI()
	//add data user to table users
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)

	for index, testCase := range testCases {
		context.SetPath(testCase.path)

		// Assertions
		if assert.NoError(t, GetUsersController(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			body := rec.Body.String()
			var responses UsersResponseSuccess
			err := json.Unmarshal([]byte(body), &responses)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, testCases[index].expectSize, len(responses.Data))
			assert.Equal(t, "alta", responses.Data[0].Name)

		}
	}

}

func TestGetUsersControllerTableNotFound(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "failed to get all data users",
		path:       "/users",
		expectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	// Drop table user in DB test to create failed condition
	config.DB.Migrator().DropTable(&models.User{})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath(testCases.path)

	// Call function on controller
	GetUsersController(c)

	body := rec.Body.String()

	var responses ResponseFailed
	err := json.Unmarshal([]byte(body), &responses)
	assert.Equal(t, testCases.expectCode, rec.Code)
	if err != nil {
		assert.Error(t, err, "error")
	}
	assert.Equal(t, "failed", responses.Status)
	assert.Equal(t, "failed to load data", responses.Message)
}

func TestCreateUserControllerSuccess(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "success to create user",
		path:       "/users",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()

	body, err := json.Marshal(mock_data_user)
	if err != nil {
		t.Error(t, err, "error")
	}

	//send data using request body with HTTP Method POST
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateUserController(c)) {
		bodyResponses := rec.Body.String()
		var user SingleUserResponseSuccess

		err := json.Unmarshal([]byte(bodyResponses), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, rec.Code)
		assert.Equal(t, "alta", user.Data.Name)
		assert.Equal(t, "alta@gmail.com", user.Data.Email)
		assert.Equal(t, "success", user.Status)
	}

}

func TestCreateUserControllerFailed(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "failed to create user",
		path:       "/users",
		expectCode: 400,
	}

	e := InitEchoTestAPI()
	// Drop table user in DB test to create failed condition
	config.DB.Migrator().DropTable(&models.User{})

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath(testCases.path)

	// Call function on controller
	if assert.NoError(t, CreateUserController(c)) {
		bodyResponses := rec.Body.String()
		var user ResponseFailed

		err := json.Unmarshal([]byte(bodyResponses), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, rec.Code)
		assert.Equal(t, "failed", user.Status)
	}
}

func TestGetOneUserControllerSuccess(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "success to get one data user",
		path:       "/users/:id",
		expectCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")

	if assert.NoError(t, GetOneUserController(context)) {

		var response SingleUserResponseSuccess
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "alta", response.Data.Name)

	}
}

func TestGetOneUserControllerFailedChar(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "failed to get one data user",
		path:       "/users/:id",
		expectCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("#")

	// type Response struct {
	// 	Message string      `json:"message"`
	// 	Data    models.User `json:"data"`
	// }

	if assert.NoError(t, GetOneUserController(context)) {

		var response SingleUserResponseSuccess
		res_body := res.Body.String()
		err := json.Unmarshal([]byte(res_body), &response)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "failed", response.Status)

	}
}

//get user detail controller using jwt
func TestGetUserDetailControllerSuccess(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "get user detail using jwt",
		path:       "/jwt/users/:id",
		expectCode: http.StatusOK,
	}

	mock_user2 := models.User{
		Name:     "alta",
		Password: "12345",
		Email:    "alta@gmail.com",
	}

	e := InitEchoTestAPI()
	InsertMockDataUserToDB()

	//create token
	var user models.User
	tx := config.DB.Where("email = ? AND password = ?", mock_user2.Email, mock_user2.Password).First(&user)
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath(testCases.path)
	context.SetParamNames("id")
	context.SetParamValues("1")
	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllerTesting())(context)

	// type Response struct {
	// 	Message string      `json:"message"`
	// 	Data    models.User `json:"data"`
	// }

	var response SingleUserResponseSuccess

	res_body := res.Body.String()
	json.Unmarshal([]byte(res_body), &response)

	t.Run("GET /jwt/users/:id", func(t *testing.T) {
		assert.Equal(t, testCases.expectCode, res.Code)
		assert.Equal(t, "alta", response.Data.Name)
		assert.Equal(t, "alta@gmail.com", response.Data.Email)
	})

	// if assert.NoError(t, GetUserDetailControllers(context)) {

	// 	var response Response
	// 	res_body := res.Body.String()
	// 	err := json.Unmarshal([]byte(res_body), &response)
	// 	if err != nil {
	// 		assert.Error(t, err, "error")
	// 	}

	// 	assert.Equal(t, testCases.expectCode, res.Code)
	// 	assert.Equal(t, "alta", response.Data.Name)

	// }
}
