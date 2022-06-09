package service

import (
	"TikTokLite/model"
	"TikTokLite/setting"
	"TikTokLite/util"
	bytes2 "bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func PublishVideo(data *multipart.FileHeader, userId int64, video model.Video, c *gin.Context) error {
	filename := filepath.Base(data.Filename)
	//获取uuid，拼接视频名称，方便调试就先加上user_id和视频名称
	uuid := util.GetUUID()
	//先把视频保存本地，再制作封面，再一起上传到七牛云，完成后删除本地视频和封面
	//增加本地视频压缩（改变码率），再上传
	oldVideoName := fmt.Sprintf("old_%s_%d_%s", uuid, userId, filename)
	oldVideoPath := setting.Conf.VideoPathPrefix + oldVideoName
	videoName := fmt.Sprintf("%s_%d_%s", uuid, userId, filename)
	videoPath := setting.Conf.VideoPathPrefix + videoName
	//先保存本地然后压缩后再取出第一帧,(后可选一起上传至七牛云)
	if err := c.SaveUploadedFile(data, oldVideoPath); err != nil {
		log.Println("本地存储video失败", err)
		return err
	}
	//先截取第一帧做封面，再进行压缩
	coverName, err := getCoverName(videoName)
	if err != nil {
		log.Println("获取coverName失败：", err)
		return err
	}
	coverPath := setting.Conf.CoverPathPrefix + coverName
	cmd := exec.Command("ffmpeg", "-i", oldVideoPath, "-y", "-f", "mjpeg", "-ss", "0.1", "-t", "0.001", coverPath)
	if err := cmd.Run(); err != nil {
		log.Println("执行ffmpeg截取封面失败：", err)
		return err
	}
	var playUrl, coverUrl string
	//上传至七牛云
	if setting.Conf.PublishConfig.Mode {
		playUrl = setting.Conf.QiNiuCloudPlayUrlPrefix + videoName
		coverUrl = setting.Conf.QiNiuCloudCoverUrlPrefix + coverName
		video.PlayUrl = playUrl
		video.CoverUrl = coverUrl
		go func() {
			//压缩视频
			compressedVideo(oldVideoPath, videoPath)
			//上传
			err = uploadVideoToQiNiuCloud(videoName, coverName, videoPath, coverPath, video)
			if err != nil {
				log.Println("七牛云上传失败：", err)
			}
		}()
		return nil
	}
	playUrl = fmt.Sprintf("http://%s:%d/static/videos/?name=%s", setting.Conf.LocalIP, setting.Conf.Port, videoName)
	coverUrl = fmt.Sprintf("http://%s:%d/static/covers/?name=%s", setting.Conf.LocalIP, setting.Conf.Port, coverName)
	video.PlayUrl = playUrl
	video.CoverUrl = coverUrl
	go func() {
		compressedVideo(oldVideoPath, videoPath)
		CreateVideo(&video)
	}()
	return nil
}

func compressedVideo(oldVideoPath, videoPath string) {
	defer os.Remove(oldVideoPath)
	//压缩视频（减小码率）
	cmd := exec.Command("ffmpeg", "-i", oldVideoPath, "-b:v", "1.5M", videoPath)
	if err := cmd.Run(); err != nil {
		log.Println("执行ffmpeg压缩视频失败：", err)
		return
	}
}

func uploadVideoToQiNiuCloud(videoName, coverName, videoPath, coverPath string, video model.Video) error {
	videoData, err := os.Open(videoPath)
	if err != nil {
		log.Println("创建cover失败：", err)
		return err
	}
	cover, err := os.Open(coverPath)
	if err != nil {
		log.Println("创建cover失败：", err)
		return err
	}
	defer os.Remove(videoPath)
	//因为先进后出，所以得先关闭链接之后再删除
	defer os.Remove(coverPath)
	defer cover.Close()
	defer videoData.Close()
	//最后上传至七牛云
	videoDataStat, err := videoData.Stat()
	if err != nil {
		log.Println("打开videoData.Stat失败：", err)
		return err
	}
	coverStat, err := cover.Stat()
	if err != nil {
		log.Println("打开cover.Stat失败：", err)
		return err
	}
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
	err = formUploader.Put(context.Background(), &ret, upToken, videoKey, videoData, videoDataStat.Size(), &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(ret.Key, ret.Hash) //打印此次上传的一些信息
	err = formUploader.Put(context.Background(), &ret, upToken, coverKey, cover, coverStat.Size(), &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	CreateVideo(&video)
	return nil
}

func uploadVideoToCloud(videoPath, videoName string) error {
	buf := bytes2.Buffer{}
	bodyWriter := multipart.NewWriter(&buf)
	fileWriter, _ := bodyWriter.CreateFormFile("video", videoPath)
	f, _ := os.Open(videoPath)
	defer f.Close()
	io.Copy(fileWriter, f)
	contenType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url := fmt.Sprintf("http://0.0.0.0:0000/upload_video?video_name=%s", videoName)
	http.Post(url, contenType, &buf)
	return nil
}
func uploadCoverToCloud(coverPath, coverName string) error {
	buf := bytes2.Buffer{}
	bodyWriter := multipart.NewWriter(&buf)
	fileWriter, _ := bodyWriter.CreateFormFile("cover", coverPath)
	f, _ := os.Open(coverPath)
	defer f.Close()
	io.Copy(fileWriter, f)
	contenType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	url := fmt.Sprintf("http://0.0.0.0:0000/upload_cover?cover_name=%s", coverName)
	http.Post(url, contenType, &buf)
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
	videos, err := model.GetVideosByUserId(toUserId)
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
