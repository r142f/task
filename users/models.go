package users

import "backend-trainee-assignment-2023/config"

type User struct {
	Id int
}

func InsertUser(user *User) error {
	_, err := config.DB.Exec("INSERT INTO Users (UserId) VALUES ($1);", user.Id)
	return err
}

func DeleteUserById(userId int) error {
	_, err := config.DB.Exec("DELETE FROM Users WHERE UserId=$1;", userId)
	return err
}
