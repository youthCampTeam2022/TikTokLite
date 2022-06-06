package service

import (
	"TikTokLite/model"
	"errors"
	"fmt"
	"strconv"
	"time"
)

//IUserService service层用户信息操作接口
type IUserService interface {
	UpdateUser(u *model.User) (Response, error)
	UserLogin(req *UserLoginOrRegisterRequest) (*UserLoginOrRegisterResponse, error)
	UserRegister(req *UserLoginOrRegisterRequest) (*UserLoginOrRegisterResponse, error)
	UserInfo(req *UserInfoRequest) (*UserInfoResponse, error)
}

//UserService 实现service层用户信息操作的接口方法
type UserService struct {
	//数据库操作接口
	UserRepository model.IUserRepository
}

func NewUserService() *UserService {
	return &UserService{
		UserRepository: model.NewUserManagerRepository(),
	}
}

//UserLoginOrRegisterRequest 登录和注册所用的req
type UserLoginOrRegisterRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

//UserLoginOrRegisterResponse 登录和注册所用的resp
type UserLoginOrRegisterResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

//UserInfoRequest 用户信息的req
type UserInfoRequest struct {
	UserId int64
	Token  string
}

//UserInfoResponse 用户信息的resp
type UserInfoResponse struct {
	Response
	User
}

//创建新用户
func (s *UserService) createUser(u *model.User) error {
	//insert一个实例
	return s.UserRepository.Insert(u)
}

//根据id获取用户信息
func (s *UserService) getUserById(u *model.User, id uint) (rep Response, err error) {
	err = s.UserRepository.GetById(u, id)
	rep = BuildResponse(err)
	return
}

//UpdateUser 更新用户基本信息
func (s *UserService) UpdateUser(u *model.User) (rep Response, err error) {
	err = s.UserRepository.Update(u)
	rep = BuildResponse(err)
	return
}

//检查用户名是否存在
func (s *UserService) userIsExists(username string) error {
	return s.UserRepository.IsExists(username)
}

//检查用户名密码
func (s *UserService) checkUser(req *UserLoginOrRegisterRequest) (u *model.User, err error) {
	u = new(model.User)
	err = s.UserRepository.GetByName(u, req.Name)
	if err != nil {
		return
	}
	if u.Password != req.Password {
		err = errors.New("password error")
		return
	}
	return
}

//UserLogin 用户登录，先检查用户名密码的正确性，再获取对应token
func (s *UserService) UserLogin(req *UserLoginOrRegisterRequest) (resp *UserLoginOrRegisterResponse, err error) {
	resp = new(UserLoginOrRegisterResponse)
	var u *model.User
	u, err = s.checkUser(req)

	if err != nil {
		resp.Response = BuildResponse(err)
		return
	}
	resp.UserId = int64(u.ID)
	resp.Token, resp.Response, err = GetToken(u)

	//conn := model.RedisCache.Conn()
	conn := model.RedisCache.AsynConn()
	defer conn.Close()
	//loginTimeKey := fmt.Sprintf("%s:%s", strconv.FormatInt(resp.UserId, 10), "loginTime")
	//_, err = conn.Do("SET", loginTimeKey, time.Now().UnixMilli())
	loginTimeKey := fmt.Sprintf("%s", strconv.FormatInt(resp.UserId, 10))
	//_, err = conn.Do("HSET", "aliveUser", loginTimeKey, time.Now().UnixMilli())
	_, err = conn.AsyncDo("HSET", "aliveUser", loginTimeKey, time.Now().UnixMilli())
	return
}

//UserRegister 用户注册，先判断username是否已存在，再插入user实例，获取对应token
func (s *UserService) UserRegister(req *UserLoginOrRegisterRequest) (resp *UserLoginOrRegisterResponse, err error) {
	resp = new(UserLoginOrRegisterResponse)
	if err = s.userIsExists(req.Name); err != nil {
		resp.Response = BuildResponse(err)
		return
	}
	var u model.User
	u.Name = req.Name
	u.Password = req.Password
	if err = s.createUser(&u); err != nil {
		resp.Response = BuildResponse(err)
		return
	}
	resp.UserId = int64(u.ID)
	//GetToken生成token
	resp.Token, resp.Response, err = GetToken(&u)
	return
}

//UserInfo 获取用户基本信息
func (s *UserService) UserInfo(req *UserInfoRequest) (resp *UserInfoResponse, err error) {
	//var u model.User
	resp = new(UserInfoResponse)
	//resp.Response, err = s.getUserById(&u, uint(req.UserId))
	resp.Response = BuildResponse(nil)
	resp.User = BuildUser(req.UserId, req.UserId, model.NewFollowManagerRepository())
	//resp.FollowCount = 0
	//resp.FollowerCount = 0
	//resp.Id = int64(u.ID)
	//resp.IsFollow = false
	//resp.Name = u.Name
	return
}
