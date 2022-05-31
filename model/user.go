package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:char(10);"`
	Password string `gorm:"type:char(32);"`
}

//IUserRepository 对用户表操作接口
type IUserRepository interface {
	Insert(u *User) error
	GetById(u *User, id uint) error
	GetByName(u *User, username string) error
	Update(u *User) error
	IsExists(username string) error
}

//UserManagerRepository 实现了IUserRepository接口
type UserManagerRepository struct {
	db *gorm.DB
}

func NewUserManagerRepository() *UserManagerRepository {
	return &UserManagerRepository{DB}
}

//Insert 插入一个user实例
func (r *UserManagerRepository) Insert(u *User) error {
	return r.db.Create(u).Error
}

//GetById 根据id获取实例
func (r *UserManagerRepository) GetById(u *User, id uint) error {
	return r.db.Where("id=?", id).First(u).Error
}

//GetByName 根据username获取实例
func (r *UserManagerRepository) GetByName(u *User, username string) error {
	return r.db.Where("name=?", username).First(u).Error
}

//Update 更新
func (r *UserManagerRepository) Update(u *User) error {
	return r.db.Save(u).Error
}

//IsExists 判断username是否已存在表中，存在的话返回error
func (r *UserManagerRepository) IsExists(username string) error {
	err := r.db.Where("name=?", username).Take(&User{}).Error
	if err == nil {
		return errors.New(fmt.Sprintf("The user already exists with UserName:%s", username))
	} else {
		return nil
	}
}
