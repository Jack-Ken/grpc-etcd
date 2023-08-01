package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bakape/thumbnailer/v2"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("1111")
	b, err := os.ReadFile("J:\\zj_qingxunying\\code\\mini_tiktok\\static\\videos\\1676733081-149168431801831424.mp4")
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(b)

	println("1111")
	thumbData, err := getThumbnail(reader)
	if err != nil {
		log.Fatal(err)
	}
	code, url := UploadimageToQiNiu(thumbData)
	if code != 0 {
		fmt.Println(code)
		//log.Fatalln(url)
	} else {
		fmt.Println(code)
		fmt.Println(url)
	}

	//code, url := UploadimageToQiNiu(b)
	//if code != 0 {
	//	fmt.Println(code)
	//	//log.Fatalln(url)
	//} else {
	//	fmt.Println(code)
	//	fmt.Println(url)
	//}

}

// 上传图片到七牛云，然后返回状态和图片的url
func UploadimageToQiNiu(data []byte) (int, string) {

	var AccessKey = "UGBFJdhQOWUe72cPA-NWjBpJpQ3jD9SzJIJMj5fn" // 秘钥对
	var SerectKey = "otu3im7-VCo6QOnoyfkG5ryTewnNxDzUBCtLAIs_"
	var Bucket = "minitiktok"                   // 空间名称
	var ImgUrl = "rym83vvtw.hb-bkt.clouddn.com" // 自定义域名或测试域名

	// 根据userid创建视频名称
	name := "testimage"
	// 检查文件类型 MIME 的功能，例如image/jpeg , video/mp4
	contentType := http.DetectContentType(data)

	t := strings.Split(contentType, "/")
	root, suffix := t[0], t[1]

	src := bytes.NewReader(data)

	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)

	// 获取上传凭证
	upToken := putPlicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuabei, // 华北区
		UseCdnDomains: false,
		UseHTTPS:      false, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	key := root + "/" + name + "." + suffix // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err := formUploader.Put(context.Background(), &ret, upToken, key, src, int64(len(data)), &putExtra)

	if err != nil {
		code := 501
		return code, err.Error()
	}

	url := ImgUrl + "/" + ret.Key // 返回上传后的文件访问路径
	return 0, url
}

// getThumbnail Generate JPEG thumbnail from video
func getThumbnail(input io.ReadSeeker) ([]byte, error) {
	_, thumb, err := thumbnailer.Process(input, thumbnailer.Options{})
	if err != nil {
		return nil, errors.New("failed to create thumbnail")
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		return nil, errors.New("failed to create buffer")
	}
	return buf.Bytes(), nil
}
