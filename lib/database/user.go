package database

import (
	"alta/be4/mvc/config"
	"alta/be4/mvc/middlewares"
	"alta/be4/mvc/models"
)

func GetUsers() (interface{}, error) {
	var users []models.User
	tx := config.DB.Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

func GetUser(id int) (interface{}, int, error) {
	user := models.User{}

	tx := config.DB.Find(&user, id)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, 0, nil
	}
	return user, 1, nil
}

func CreateUser(user *models.User) (interface{}, error) {
	tx := config.DB.Create(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func EditUser(newData *models.User, id int) (interface{}, int, error) {
	user := models.User{}
	tx := config.DB.Find(&user, id)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, 0, nil
	}

	txUpdate := config.DB.Model(&user).Updates(newData)
	if txUpdate.Error != nil {
		return models.User{}, 0, txUpdate.Error
	}
	if txUpdate.RowsAffected == 0 {
		return models.User{}, 0, nil
	}
	return user, 1, nil

}

func DeleteUser(id int) (interface{}, int, error) {
	user := models.User{}

	tx := config.DB.Delete(&user, id)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, 0, nil
	}
	return user, 1, nil
}

func LoginUsers(user *models.User) (interface{}, error) {
	var err error
	tx := config.DB.Where("email=? AND password=?", user.Email, user.Password).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	user.Token, err = middlewares.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}
	if er := config.DB.Save(user).Error; er != nil {
		return nil, err
	}
	return user, nil
}

func GetDetailUsers(id int) (interface{}, error) {
	var user models.User

	if e := config.DB.Find(&user, id).Error; e != nil {
		return nil, e
	}
	return user, nil
}
