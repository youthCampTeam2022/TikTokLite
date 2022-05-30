package util

import (
	"github.com/syyongx/go-wordsfilter"
	"log"
	"os"
	"strings"
)

const SensitiveTextPathMain = "./util/sensitive_txt/sensitive_word.txt"
//const SensitiveTextPathTest = "./sensitive_txt/sensitive_word.txt"

// Wf 暂时只加了暴力相关的敏感词用于测试
var Wf *wordsfilter.WordsFilter
var WfRoot map[string]*wordsfilter.Node

// FilterInit 初始化过滤器
func FilterInit()  {
	f, err := os.Open(SensitiveTextPathMain)
	if err != nil {
		log.Println("err:", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("err in close")
		}
	}(f)
	var b = make([]byte, 4096)
	_, err = f.Read(b)
	if err != nil {
		log.Println("err:", err)
		return
	}
	l := strings.Split(string(b), "\r\n")
	Wf = wordsfilter.New()
	WfRoot = Wf.Generate(l)
	if err != nil{
		log.Println(err)
	}
}

// Filtration 主要过滤接口，将敏感词替换为**，因过滤后删除了空格，需要做判断
func Filtration(text string)(newText string,ok bool) {
	newText = Wf.Replace(text, WfRoot)
	buffer := strings.ReplaceAll(text," ","")
	if newText == buffer{
		return text,true
	}
	return newText,false
}