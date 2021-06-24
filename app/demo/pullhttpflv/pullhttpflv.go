// Copyright 2019, Chef.  All rights reserved.
// https://github.com/yangjing0630/go-stream
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"flag"
	"github.com/yangjing0630/go-stream/pkg/base"
	"github.com/yangjing0630/go-stream/pkg/httpflv"
	"github.com/q191201771/naza/pkg/nazalog"
)

// 拉取HTTP-FLV的流
//
// TODO
// - 存储成flv文件
// - 拉取HTTP-FLV流进行分析参见另外一个demo：analyseflvts。 这个demo可能可以删除掉了。

func main() {
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	defer nazalog.Sync()

	url := parseFlag()
	session := httpflv.NewPullSession()
	var httpFlvWriter httpflv.FlvFileWriter

	err := httpFlvWriter.Open("test/j443eekj.flv")
	nazalog.Assert(nil, err)
	defer httpFlvWriter.Dispose()
	err = httpFlvWriter.WriteRaw(httpflv.FlvHeader)
	nazalog.Assert(nil, err)
	err = session.Pull(url, func(tag httpflv.Tag) {
		err:=httpFlvWriter.WriteTag(tag)
		nazalog.Assert(nil, err)
		//fmt.Println(tag)
		//switch tag.Header.Type {
		//case httpflv.TagTypeMetadata:
		//	nazalog.Info(hex.Dump(tag.Payload()))
		//case httpflv.TagTypeAudio:
		//case httpflv.TagTypeVideo:
		//}
	})
	nazalog.Assert(nil, err)
	err = <-session.WaitChan()
	nazalog.Assert(nil, err)
}

func parseFlag() string {
	url := flag.String("i", "", "specify http-flv url")
	flag.Parse()
	if *url == "" {
		flag.Usage()
		base.OsExitAndWaitPressIfWindows(1)
	}
	return *url
}
