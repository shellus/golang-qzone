package qzone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type responsePack struct {
	Code    int
	Message string
	Data    *json.RawMessage
}

type PhotoPagePack struct {
	limit        int
	PhotoList    []*PhotoPack
	TotalInAlbum int
	TotalInPage  int
}
type PhotoPack struct {
	Raw          string
	Lloc         string
	Rawshoottime string // 拍摄时间
	Uploadtime   string
	Name         string
	Is_video     bool
}

type VideoPagePack struct {
	Photos []*VideoPack
}
type VideoPack struct {
	Lloc       string
	PicKey     string `json:"picKey"`
	Is_video   bool
	Video_info *VideoInfo
	Uploadtime string
	ShootTime  int64 `json:"shootTime"`
	Name       string
}

type VideoInfo struct {
	Vid       string
	Video_url string
}

func DownloadFactory(pCh chan *PhotoPack) {
	for p := range pCh {
		if p.Is_video {
			DownloadVideo(p)
		} else {
			//downloadPhoto(p)
		}
	}
	wg.Done()
}

func genFileName(p *PhotoPack) string {
	//originName := strings.ReplaceAll(p.Rawshoottime, ":", "-")
	// \/:*?"<>|
	// ! => @
	// * => -
	name := p.Lloc
	name = strings.ReplaceAll(name, "!", "@")
	name = strings.ReplaceAll(name, "*", "-")

	var ext string
	if p.Is_video {
		ext = "mp4"
	} else {
		ext = "jpg"
	}
	path := downloadDir + `\` + name + "." + ext
	// 自动重命名
	//i := 1
	//for ;fileExists(path); {
	//	path = downloadDir + `\` + name + "("+strconv.Itoa(i)+")" + "." + ext
	//	i = i + 1
	//	fmt.Printf("文件存在，新文件名：%s\n", path)
	//}

	return path
}

func downloadPhoto(p *PhotoPack) error {
	var err error
	path := genFileName(p)
	if fileExists(path) {
		return fmt.Errorf("文件存在，跳过下载：%s\n", path)
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", p.Rawshoottime, time.Local)
	if err != nil {
		return fmt.Errorf("parse time [%s] err: %s\n", p.Rawshoottime, err.Error())
	}

	return downloadFile(path, p.Raw, t)
}

func DownloadVideo(p *PhotoPack) (err error) {
	videoUrl, err := queryVideoUrl(p)
	if err != nil {
		return err
	}

	path := genFileName(p)
	if fileExists(path) {
		return fmt.Errorf("文件存在，跳过下载：%s\n", path)
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", p.Rawshoottime, time.Local)
	if err != nil {
		return fmt.Errorf("parse time [%s] err: %s\n", p.Rawshoottime, err.Error())
	}

	err = downloadFile(path, videoUrl, t)
	if err != nil {
		fmt.Printf("downloadFile err: %s\n", err.Error())
	}

	return nil
}

func queryVideoUrl(p *PhotoPack) (string, error) {
	if p.Is_video != true {
		return "", fmt.Errorf("not video so fail: %v", p)
	}
	picKey := p.Lloc
	url := "https://h5.qzone.qq.com/proxy/domain/photo.qzone.qq.com/fcgi-bin/cgi_floatview_photo_list_v2?" +
		"g_tk=" + g_tk + "&" +
		"callback=viewer_Callback&" +
		"t=414765976&" +
		"topicId=" + topicId + "&" +
		"picKey=" + picKey + "&" +
		"shootTime=&" +
		"cmtOrder=1&" +
		"fupdate=1&" +
		"plat=qzone&" +
		"source=qzone&" +
		"cmtNum=10&" +
		"likeNum=5&" +
		"inCharset=utf-8&" +
		"outCharset=utf-8&" +
		"callbackFun=viewer&" +
		"offset=0&" +
		"number=15&" +
		"uin=" + qq + "&" +
		"hostUin=" + qq + "&" +
		"appid=4&" +
		"isFirst=1&" +
		"sortOrder=1&" +
		"showMode=1&" +
		"need_private_comment=1&" +
		"prevNum=9&" +
		"postNum=18&" +
		"_=1584434348557"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	r.Header.Set("Cookie", cookieStr)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("resp.StatusCode != 200 : %d", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	l := len(bytes)

	var repack responsePack
	err = json.Unmarshal(bytes[16:l-2], &repack)
	if err != nil {
		return "", fmt.Errorf("Unmarshal responsePack fail: %s", string(bytes[16:l-2]))
	}

	var photoPagePack = VideoPagePack{}
	err = json.Unmarshal(*repack.Data, &photoPagePack)
	if err != nil {
		return "", fmt.Errorf("Unmarshal VideoPagePack fail: %s", string(*repack.Data))
	}
	videoUrl := ""
	for _, photo := range photoPagePack.Photos {
		if photo.Lloc == p.Lloc {
			videoUrl = photo.Video_info.Video_url
		}
	}
	if videoUrl == "" {
		return "", fmt.Errorf("lloc: %s\n not fonud videoUrl: %s", p.Lloc, string(bytes[16:l-2]))
	}
	return videoUrl, nil
}

func requestPhotoList(page int, pageSize int) *PhotoPagePack {

	pageStart := pageSize * (page - 1)
	url := "https://h5.qzone.qq.com/proxy/domain/photo.qzone.qq.com/fcgi-bin/cgi_list_photo?" +
		"g_tk=" + g_tk + "&" +
		"callback=shine9_Callback&" +
		"t=138268231&" +
		"mode=0&" +
		"idcNum=4&" +
		"hostUin=" + qq + "&" +
		"topicId=" + topicId + "&" +
		"noTopic=0&" +
		"uin=" + qq + "&" +
		"pageStart=" + strconv.Itoa(pageStart) + "&" +
		"pageNum=" + strconv.Itoa(pageSize) + "&" +
		"skipCmtCount=0&" +
		"singleurl=1&" +
		"batchId=&" +
		"notice=0&" +
		"appid=4&" +
		"inCharset=utf-8&" +
		"outCharset=utf-8&" +
		"source=qzone&" +
		"plat=qzone&" +
		"outstyle=json&" +
		"format=jsonp&" +
		"json_esc=1&" +
		"callbackFun=shine9&" +
		"_=1584434364077"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	r.Header.Set("Cookie", cookieStr)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	l := len(bytes)

	var repack responsePack
	err = json.Unmarshal(bytes[16:l-2], &repack)
	if err != nil {
		panic(err)
	}
	var photoPagePack = PhotoPagePack{}
	err = json.Unmarshal(*repack.Data, &photoPagePack)
	if err != nil {
		panic(err)
	}
	return &photoPagePack
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return !info.IsDir()
}
