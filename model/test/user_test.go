/*
对数据库users表操作的测试
*/
package test

import (
	"TikTokLite/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

var (
	mock sqlmock.Sqlmock
	err  error
	db   *sql.DB
)

// TestMain是在当前package下，最先运行的一个函数，常用于初始化
func TestMain(m *testing.M) {
	//把匹配器设置成相等匹配器，不设置默认使用正则匹配
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {

		panic(err)
	}
	model.DB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// m.Run 是调用包下面各个Test函数的入口
	os.Exit(m.Run())
}

//插入user实例测试
func TestInsert(t *testing.T) {
	user := &model.User{
		Name:     "haogee",
		Password: "123456",
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`password`) VALUES (?,?,?,?,?)").
		WithArgs(user.CreatedAt, user.UpdatedAt, nil, user.Name, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := model.NewUserManagerRepository().Insert(user)
	assert.Nil(t, err)
}

//根据id查询测试
func TestGetById(t *testing.T) {
	user := &model.User{
		Name:     "haogee",
		Password: "123456",
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = 10
	mock.ExpectQuery("SELECT * FROM `users` WHERE id=? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1").
		WithArgs(user.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "password"}).
				AddRow(user.ID, user.CreatedAt, user.UpdatedAt, nil, user.Name, user.Password))
	u := new(model.User)
	err := model.NewUserManagerRepository().GetById(u, 10)
	assert.Equal(t, user, u)
	assert.Nil(t, err)
}

//根据Name查询测试
func TestGetByName(t *testing.T) {
	user := &model.User{
		Name:     "haogee",
		Password: "123456",
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = 10
	mock.ExpectQuery("SELECT * FROM `users` WHERE name=? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1").
		WithArgs(user.Name).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "password"}).
				AddRow(user.ID, user.CreatedAt, user.UpdatedAt, nil, user.Name, user.Password))
	u := new(model.User)
	err := model.NewUserManagerRepository().GetByName(u, "haogee")
	assert.Equal(t, user, u)
	assert.Nil(t, err)
}

//测试用户名是否重复
func TestIsExists(t *testing.T) {
	user := &model.User{
		Name:     "haogee",
		Password: "123456",
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = 10
	mock.ExpectQuery("SELECT * FROM `users` WHERE name=? AND `users`.`deleted_at` IS NULL LIMIT 1").
		WithArgs(user.Name).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "password"}).
				AddRow(user.ID, user.CreatedAt, user.UpdatedAt, nil, user.Name, user.Password))
	err := model.NewUserManagerRepository().IsExists(user.Name)
	diffErr := errors.New(fmt.Sprintf("The user already exists with UserName:%s", user.Name))
	assert.Equal(t, err, diffErr)
	err = model.NewUserManagerRepository().IsExists("laowang")
	assert.Nil(t, err)
}
