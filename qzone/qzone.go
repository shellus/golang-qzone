package qzone

import (
	"fmt"
	"math"
	"sync"
)

// 含义参考 README.md
var cookieStr = ""
var g_tk = ""
var qq = ""
var topicId = ""
var downloadDir = ``

var wg sync.WaitGroup

func Login(c string, g string, q string, t string, d string) {
	cookieStr = c
	g_tk = g
	qq = q
	topicId = t
	downloadDir = d
}
func Run() {

	total := 6120.00
	pageSize := 500.00
	totalPage := int(math.Ceil(total / pageSize))

	pCh := make(chan *PhotoPack)
	go func() {
		for page := 1; page <= totalPage; page++ {
			data := requestPhotoList(page, int(pageSize))
			for _, p := range data.PhotoList {
				pCh <- p
			}
		}
		close(pCh)
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go DownloadFactory(pCh)
	}
	wg.Wait()

	fmt.Println("done !")
}
