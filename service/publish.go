package service

import (
	"TikTokLite/model"
	"TikTokLite/setting"
	"TikTokLite/util"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

func PublishVideo(data *multipart.FileHeader, userId int64, c *gin.Context) (string, string, error) {
	filename := filepath.Base(data.Filename)
	videoData, err := data.Open()
	if err != nil {
		log.Println("获取data数据失败：", err)
		return "", "", err
	}
	defer videoData.Close()
	//获取uuid，拼接视频名称，方便调试就先加上user_id和视频名称
	uuid := util.GetUUID()
	//先把视频保存本地，再制作封面，再一起上传到七牛云，完成后删除本地视频和封面
	videoName := fmt.Sprintf("%s_%d_%s", uuid, userId, filename)
	videoPath := setting.Conf.VideoPathPrefix + videoName
	//先保存本地然后取出第一帧之,(后可选一起上传至七牛云)
	if err := c.SaveUploadedFile(data, videoPath); err != nil {
		log.Println("本地存储video失败", err)
		return "", "", err
	}
	//截取第一帧做封面
	coverName, err := getCoverName(videoName)
	if err != nil {
		log.Println("获取coverName失败：", err)
		return "", "", err
	}
	coverPath := setting.Conf.CoverPathPrefix + coverName
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-y", "-f", "mjpeg", "-ss", "0.1", "-t", "0.001", coverPath)
	if err := cmd.Run(); err != nil {
		log.Println("执行ffmpeg失败：", err)
		return "", "", err
	}
	var playUrl, coverUrl string
	//上传至七牛云
	if setting.Conf.PublishConfig.Mode {
		cover, err := os.Open(coverPath)
		if err != nil {
			log.Println("创建cover失败：", err)
			return "", "", err
		}
		defer os.Remove(videoPath)
		//因为先进后出，所以得先关闭链接之后再删除
		defer os.Remove(coverPath)
		defer cover.Close()
		//最后上传至七牛云
		co, err := cover.Stat()
		if err != nil {
			log.Println("打开cover.Stat失败：", err)
			return "", "", err
		}
		err = uploadVideoToQiNiuCloud(videoData, cover, videoName, coverName, data.Size, co.Size())
		if err != nil {
			log.Println("七牛云上传失败：", err)
			return "", "", err
		}
		playUrl = setting.Conf.QiNiuCloudPlayUrlPrefix + videoName
		coverUrl = setting.Conf.QiNiuCloudCoverUrlPrefix + coverName
		return playUrl, coverUrl, nil
	}
	playUrl = fmt.Sprintf("http://%s:%d/static/videos/%s", setting.Conf.LocalIP, setting.Conf.Port, videoName)
	coverUrl = fmt.Sprintf("http://%s:%d/static/covers/%s", setting.Conf.LocalIP, setting.Conf.Port, coverName)
	return playUrl, coverUrl, nil
}

func uploadVideoToQiNiuCloud(video, cover multipart.File, videoName, coverName string, videoSize, coverSize int64) error {
	//上传的路径+文件名
	videoKey := fmt.Sprintf("videos/%s", videoName)
	coverKey := fmt.Sprintf("covers/%s", coverName)
	//上传凭证
	mac := qbox.NewMac(setting.Conf.AccessKey, setting.Conf.SecretKey)
	putPolicy := storage.PutPolicy{
		Scope: setting.Conf.BucketName,
	}
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		// 空间对应的机房
		Zone: &storage.ZoneHuanan,
		// 是否使用https域名
		UseHTTPS: true,
		// 上传是否使用CDN上传加速
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	//额外参数
	putExtra := storage.PutExtra{
		//Params: map[string]string{
		//	"x:name": "github logo",
		//},
	}
	err := formUploader.Put(context.Background(), &ret, upToken, videoKey, video, videoSize, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(ret.Key, ret.Hash) //打印此次上传的一些信息
	err = formUploader.Put(context.Background(), &ret, upToken, coverKey, cover, coverSize, &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(ret.Key, ret.Hash)
	return nil
}

func getCoverName(s string) (string, error) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			return fmt.Sprintf("%s.jpeg", s[:i]), nil
		}
	}
	return "", errors.New("文件名格式不合法")
}

func CreateVideo(v *model.Video) error {
	return v.Create()
}

func GetVideoList(userId, toUserId int64) ([]Video, error) {
	followService := NewFollowService()
	//数据库表格式的videos
	videos, err := model.GetVideosByUserId(userId)
	if err != nil {
		log.Println("getVideosByUserId failed:", err)
		return nil, err
	}
	videoList := make([]Video, len(videos))
	//user复用一个就行
	author := BuildUser(userId, toUserId, followService.FollowRepository)
	//不能直接用video
	for i := range videoList {
		videoId := int64(videos[i].ID)
		isFavorite, err := model.IsFavorite(userId, videoId)
		if err != nil {
			return nil, err
		}
		videoList[i].Id = videoId
		videoList[i].Author = author
		videoList[i].Title = videos[i].Title
		videoList[i].PlayUrl = videos[i].PlayUrl
		videoList[i].CoverUrl = videos[i].CoverUrl
		videoList[i].FavoriteCount = model.GetFavoriteNum(videoId)
		videoList[i].CommentCount = model.GetCommentNum(videoId)
		videoList[i].IsFavorite = isFavorite
	}
	return videoList, nil
}
