package amap

import (
	"fmt"
	"testing"
)

var (
	amap *Amap
)

func setup() {
	amapKey := ""
	amap = NewAmap(amapKey)
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func TestAmapGetGeo(t *testing.T) {
	res,err := amap.GetGeo("广州市番禺区万达广场", "广州")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println(fmt.Sprintf("%+v", res))
		fmt.Println(res.Geocodes[0])
		fmt.Println(res.Geocodes[0].Location)
	}

}
