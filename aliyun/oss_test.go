package aliyun

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

var (
	s        *OssStore
	savePath = "test/"
)

func setup() error {

	conf := &OssConfig{
		AccessKeyID:     "",
		AccessKeySecret: "",
		Endpoint:        "",
		Bucket:          "",
	}

	var newOssErr error
	s, newOssErr = NewOssStore(conf)
	if newOssErr != nil {
		return newOssErr
	}

	return nil
}

func teardown() {
	fmt.Println("store done")
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Println(err)
		return
	}
	m.Run()
	teardown()
}

func TestUploadText(t *testing.T) {
	if err := s.UploadText("test store text", savePath+"tst.txt"); err != nil {
		t.Error(err)
	}
}

func TestUploadBytes(t *testing.T) {
	if err := s.UploadBytes([]byte("test store text"), savePath+"tsb.txt"); err != nil {
		t.Error(err)
	}
}

func TestUploadFile(t *testing.T) {
	if err := s.UploadFile("./text.txt", savePath+"tsf.txt"); err != nil {
		t.Error(err)
	}
}

func TestUploadStream(t *testing.T) {
	fd, openFileErr := os.Open("./text.txt")
	defer fd.Close()
	if openFileErr != nil {
		t.Error(openFileErr)
	}

	if err := s.UploadStream(fd, savePath+"tss.txt"); err != nil {
		t.Error(err)
	}
}

func TestGetContent(t *testing.T) {
	res, err := s.GetContent(savePath + "tss.txt")
	if err != nil {
		t.Error(err)
	}

	t.Log(string(res))
}

func TestGetContentToBuf(t *testing.T) {
	buf := new(bytes.Buffer)
	err := s.GetContentToBuf(savePath+"tss.txt", buf)
	if err != nil {
		t.Error(err)
	}

	t.Log(buf)
	t.Log(buf.String())
}

func TestGetContentToStream(t *testing.T) {
	fd, openFileErr := os.OpenFile("./text.txt", os.O_RDWR|os.O_CREATE, 0660)
	defer fd.Close()
	if openFileErr != nil {
		t.Error(openFileErr)
	}

	if err := s.GetContentToStream(savePath+"tss.txt", fd); err != nil {
		t.Error(err)
	}

	t.Log(fd)

	chunk := make([]byte, 1024)
	fd.Read(chunk)
	t.Log(string(chunk))
}

func TestDownload(t *testing.T) {
	if err := s.Download(savePath+"tss.txt", "./tss.txt"); err != nil {
		t.Error(err)
	}
}
