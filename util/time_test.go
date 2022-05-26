package util

import (
	"fmt"
	"testing"
	"time"
)

func TestTime2String(t *testing.T) {
	fmt.Println(Time2String(time.Now()))
}
