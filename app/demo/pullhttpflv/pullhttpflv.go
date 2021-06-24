package main

import (
	"encoding/hex"
	"flag"
	"fmt"
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

var _dir="temp/"

func isExistDir(dir string)(bool,error)  {
	_,err:=os.Stat(dir)
	if err==nil{
		return true,nil
	}
	if os.IsNotExist(err){
		return false, err
	}
	return false, err
}

func autoDir(dir string)  {
	if exist,_:=isExistDir(dir);exist{
		return
	}
	err:=os.MkdirAll(dir,0777)
	fmt.Println(err)
}

func main() {
	f,url,isStore := parseFlag()
	fps:=strings.Split(_dir+f,"/")
	var path string
	for i:=0;i<len(fps)-1;i++{
		path+=fps[i]+"/"
	}
	autoDir(path)
	pullFlv(_dir+f,url,isStore)
}

func pullFlv(filename,url string,isStore bool){
	var err error
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	defer nazalog.Sync()

	session := httpflv.NewPullSession()
	var httpFlvWriter httpflv.FlvFileWriter

	if isStore{
		err = httpFlvWriter.Open(filename)
		nazalog.Assert(nil, err)
		defer httpFlvWriter.Dispose()
		err = httpFlvWriter.WriteRaw(httpflv.FlvHeader)
		nazalog.Assert(nil, err)
	}
	err = session.Pull(url, func(tag httpflv.Tag) {
		if isStore{
			err:=httpFlvWriter.WriteTag(tag)
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
	err = <-session.WaitChan()
	nazalog.Assert(nil, err)
}

func parseFlag() (f,utl string,iStore bool) {
	url := flag.String("i", "", "specify http-flv url")
	filename:=flag.String("f","","flv filename")
	isStore:=flag.Bool("s",false,"is store")
	flag.Parse()
	f=*filename
	if *url == "" {
		flag.Usage()
		base.OsExitAndWaitPressIfWindows(1)
	}
	if *filename==""{
		fs:=strings.Split(*url,"/")
		f=fmt.Sprintf("%d_%s",time.Now().Unix(),fs[len(fs)-1])
	}
	return strings.TrimLeft(f,"/"),*url,*isStore
}
