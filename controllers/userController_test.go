package controllers

import (
	"alta/be4/mvc/config"
	"alta/be4/mvc/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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
		size       int
	}{
		{
			name:       "success",
			path:       "/users",
			expectCode: 200,
			size:       1,
		},
	}

	e := InitEchoTestAPI()
	InsertMockDataUserToDB()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for index, testCase := range testCases {
		c.SetPath(testCase.path)

		// Assertions
		if assert.NoError(t, GetUsersController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			body := rec.Body.String()
			var responses UsersResponseSuccess
			err := json.Unmarshal([]byte(body), &responses)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, testCases[index].size, len(responses.Data))
			// assert.Equal(t, userJSON, rec.Body.String())
			// assert.True(t, strings.HasPrefix(body, testCase.expectBodyStartsWith))

		}
	}

}

func TestGetUsersControllerTableNotFound(t *testing.T) {
	var testCases = struct {
		name       string
		path       string
		expectCode int
	}{

		name:       "failed",
		path:       "/users",
		expectCode: 400,
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

		name:       "success",
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

		name:       "failed",
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
