package test

import (
	"TikTokLite/model"
	"fmt"
	"testing"
)

func TestHotFeed(t *testing.T) {
	model.Init()
	model.BuildHotFeed()
	fmt.Println(model.PullHotFeed(20))
}
