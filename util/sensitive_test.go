package util

import (
	"fmt"
	"testing"
)

func TestReplace(t *testing.T) {
	FilterInit()
	a := Wf.Replace("卖炸   药", WfRoot)
	fmt.Println(a)
}

func TestFiltration(t *testing.T) {
	FilterInit()
	filtration, ok := Filtration("卖炸   药")
	fmt.Println(filtration, ok)
}
