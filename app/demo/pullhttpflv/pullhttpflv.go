package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/yangjing0630/go-stream/pkg/base"
	"github.com/yangjing0630/go-stream/pkg/httpflv"
)

// 拉取HTTP/HTTPS-FLV的流
//
// TODO
// - 存储成flv文件
// - 拉取HTTP-FLV流进行分析参见另外一个demo：analyseflvts

var _dir = "temp/"
var _m = make(chan bool)

func isExistDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func autoDir(dir string) {
	if exist, _ := isExistDir(dir); exist {
		return
	}
	err := os.MkdirAll(dir, 0777)
	fmt.Println(err)
}

func main() {

	f, url, isStore, gNum := parseFlag()

	fps := strings.Split(_dir+f, "/")
	var path string
	for i := 0; i < len(fps)-1; i++ {
		path += fps[i] + "/"
	}
	autoDir(path)
	for i := 0; i < gNum; i++ {
		go func(currentI int) {
			pullFlv(fmt.Sprintf("%s%d_%s", _dir, currentI, f), url, isStore)
		}(i)
	}

	http.ListenAndServe("0.0.0.0:6060", nil)
}

func pullFlv(filename, url string, isStore bool) {
	var err error
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	defer nazalog.Sync()

	session := httpflv.NewPullSession()
	var httpFlvWriter httpflv.FlvFileWriter

	if isStore {
		err = httpFlvWriter.Open(filename)
		nazalog.Assert(nil, err)
		defer httpFlvWriter.Dispose()
		err = httpFlvWriter.WriteRaw(httpflv.FlvHeader)
		nazalog.Assert(nil, err)
	}
	err = session.Pull(url, func(tag httpflv.Tag) {
		if isStore {
			err := httpFlvWriter.WriteTag(tag)
			nazalog.Assert(nil, err)
		}
		switch tag.Header.Type {
		case httpflv.TagTypeMetadata:
			nazalog.Info(hex.Dump(tag.Payload()))
		case httpflv.TagTypeAudio:
		case httpflv.TagTypeVideo:
		}
	})
	nazalog.Assert(nil, err)
	if err = <-session.WaitChan(); err != nil {
		_m <- true
	}
	nazalog.Assert(nil, err)
}

func parseFlag() (f, utl string, iStore bool, num int) {
	url := flag.String("i", "", "specify http-flv url")
	filename := flag.String("f", "", "flv filename")
	isStore := flag.Bool("s", false, "is store")
	n := flag.Int("n", 1, "goroutine number")
	flag.Parse()
	f = *filename
	if *url == "" {
		flag.Usage()
		base.OsExitAndWaitPressIfWindows(1)
	}
	if *filename == "" {
		fs := strings.Split(*url, "/")
		f = fmt.Sprintf("%d_%s", time.Now().Unix(), fs[len(fs)-1])
	}
	return strings.TrimLeft(f, "/"), *url, *isStore, *n
}
