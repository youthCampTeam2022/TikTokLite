package service

import (
	"TikTokLite/model"
)

type IFollowService interface {
	//RelationAction(req *RelationActionRequest) (*RelationActionResponse, error)
	RedisRelationAction(req *RelationActionRequest) (resp *RelationActionResponse, err error)
	RedisFollowList(uid int64) (resp *RelationFollowListResponse, err error)
	RedisFollowerList(uid int64) (resp *RelationFollowListResponse, err error)
}
type FollowService struct {
	FollowRepository model.IFollowRepository
}

func NewFollowService() *FollowService {
	return &FollowService{
		FollowRepository: model.NewFollowManagerRepository(),
	}
}

type RelationActionRequest struct {
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
	ToUserID   int64  `json:"to_user_id"`
	ActionType int32  `json:"action_type"`
}
type RelationActionResponse struct {
	Response
}

type RelationFollowListRequest struct {
	UserID int64
	Token  string
}
type RelationFollowListResponse struct {
	Response
	Users []User `json:"user_list"`
}

//func (s *FollowService) RelationAction(req *RelationActionRequest) (resp *RelationActionResponse, err error) {
//	if req.ActionType == 1 {
//		f := &model.Follow{
//			FollowerID: req.UserID,
//			UserID:     req.ToUserID,
//		}
//		err = s.FollowRepository.Insert(f)
//	} else {
//		err = s.FollowRepository.Delete(req.ToUserID, req.UserID)
//	}
//	resp = new(RelationActionResponse)
//	resp.Response = BuildResponse(err)
//	return
//}

func (s *FollowService) RedisRelationAction(req *RelationActionRequest) (resp *RelationActionResponse, err error) {
	if req.ActionType == 1 {
		err = s.FollowRepository.RedisInsert(req.ToUserID, req.UserID)
	} else {
		err = s.FollowRepository.RedisDelete(req.ToUserID, req.UserID)
	}
	resp = new(RelationActionResponse)
	resp.Response = BuildResponse(err)
	return
}

func (s *FollowService) RedisFollowList(uid int64) (resp *RelationFollowListResponse, err error) {
	var follows []int64
	follows, err = s.FollowRepository.RedisGetFollowList(uid)
	if err != nil {
		return nil, err
	}
	resp = new(RelationFollowListResponse)
	resp.Users = BuildUserList(uid, follows, s.FollowRepository)
	return
}

func (s *FollowService) RedisFollowerList(uid int64) (resp *RelationFollowListResponse, err error) {
	var followers []int64
	followers, err = s.FollowRepository.RedisGetFollowerList(uid)
	if err != nil {
		return nil, err
	}
	resp = new(RelationFollowListResponse)
	resp.Users = BuildUserList(uid, followers, s.FollowRepository)
	return
}
