package test

import (
	"TikTokLite/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpsert(t *testing.T) {
	f := &model.Follow{
		UserID:     1,
		FollowerID: 2,
	}
	f.UpdatedAt = time.Now()
	f.CreatedAt = time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `follows` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`follower_id`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`),`follower_id`=VALUES(`follower_id`)").
		WithArgs(f.CreatedAt, f.UpdatedAt, nil, f.UserID, f.FollowerID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := model.NewFollowManagerRepository().Insert(f)
	assert.Nil(t, err)
}
