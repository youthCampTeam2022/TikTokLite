package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
)

func Videos(c *gin.Context) {
	name := c.Query("name")
	filePath := "./static/videos/" + name
	fileTmp, errByOpenFile := os.Open(filePath)
	defer fileTmp.Close()
	//获取文件的名称
	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	if errByOpenFile != nil {
		c.Redirect(http.StatusFound, "/404")
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
	return
}
func Covers(c *gin.Context) {
	name := c.Query("name")
	filePath := "./static/covers/" + name
	fileTmp, errByOpenFile := os.Open(filePath)
	defer fileTmp.Close()
	//获取文件的名称
	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	if errByOpenFile != nil {
		log.Println("获取文件失败")
		c.Redirect(http.StatusFound, "/404")
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
	return
}
