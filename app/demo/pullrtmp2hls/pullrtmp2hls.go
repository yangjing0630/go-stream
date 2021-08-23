package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/yangjing0630/go-stream/pkg/base"
	"github.com/yangjing0630/go-stream/pkg/hls"
	"github.com/yangjing0630/go-stream/pkg/rtmp"
)

var TaskMap sync.Map

type TaskInfo struct {
	ps *rtmp.PullSession
}

func startTask(w http.ResponseWriter, r *http.Request) {
	rtmpUrl := r.FormValue("rtmpUrl")
	taskId := r.FormValue("taskId")
	if rtmpUrl == "" && taskId == "" {
		fmt.Fprintf(w, "参数为空\n")
		return
	}
	go pullRtmp2Hls(rtmpUrl, "", 3000, 6, taskId)
	fmt.Fprintf(w, "success!!\n")
}

func pullRtmp2Hls(url string, hlsOutPath string, fragmentDurationMs int, fragmentNum int, taskId string) {
	nazalog.Infof("parse flag succ. url=%s, hlsOutPath=%s, fragmentDurationMs=%d, fragmentNum=%d",
		url, hlsOutPath, fragmentDurationMs, fragmentNum)

	hlsMuxerConfig := hls.MuxerConfig{
		OutPath:            hlsOutPath,
		FragmentDurationMs: fragmentDurationMs,
		FragmentNum:        fragmentNum,
	}
	ctx, err := base.ParseRtmpUrl(url)
	if err != nil {
		nazalog.Fatalf("parse rtmp url failed. url=%s, err=%+v", url, err)
	}
	streamName := ctx.LastItemOfPath

	hlsMuexer := hls.NewMuxer(streamName, true, &hlsMuxerConfig, nil)
	hlsMuexer.Start()

	pullSession := rtmp.NewPullSession(func(option *rtmp.PullSessionOption) {
		option.PullTimeoutMs = 10000
		option.ReadAvTimeoutMs = 10000
	})
	err = pullSession.Pull(url, func(msg base.RtmpMsg) {
		hlsMuexer.FeedRtmpMessage(msg)
	})
	TaskMap.Store(taskId, TaskInfo{
		ps: pullSession,
	})
	if err != nil {
		nazalog.Fatalf("pull rtmp failed. err=%+v", err)
	}
	//err = <-pullSession.WaitChan()
	select {
	case err = <-pullSession.WaitChan():
		TaskMap.Delete(taskId)
		nazalog.Errorf("< session.Wait [%s] err=%+v", pullSession.UniqueKey(), err)
		return
	}
}

func stopTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("taskId")
	value, ok := TaskMap.Load(taskId)
	if ok && value != nil {
		taskInfo := value.(TaskInfo)
		if err := taskInfo.ps.Dispose(); err != nil {
			fmt.Fprintf(w, "卧槽 结束流失败")
			return
		}
	}
	TaskMap.Range(func(key, value interface{}) bool {
		fmt.Println("Key =", key, "Value =", value)
		return true
	})
	fmt.Fprintf(w, "%v is %v", value, ok)
}

func main() {
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})
	defer nazalog.Sync()

	http.HandleFunc("/startTask", startTask)
	http.HandleFunc("/stopTask", stopTask)

	http.ListenAndServe(":9898", nil)

	return
}

func parseFlag() (url string, hlsOutPath string, fragmentDurationMs int, fragmentNum int) {
	i := flag.String("i", "", "specify pull rtmp url")
	o := flag.String("o", "", "specify ouput hls file")
	d := flag.Int("d", 3000, "specify duration of each ts file in millisecond")
	n := flag.Int("n", 6, "specify num of ts file in live m3u8 list")
	flag.Parse()
	if *i == "" {
		flag.Usage()
		eo := filepath.FromSlash("./pullrtmp2hls/")
		_, _ = fmt.Fprintf(os.Stderr, `Example:
  %s -i rtmp://127.0.0.1:19350/live/test110 -o %s
  %s -i rtmp://127.0.0.1:19350/live/test110 -o %s -d 5000 -n 5
`, os.Args[0], eo, os.Args[0], eo)
		base.OsExitAndWaitPressIfWindows(1)
	}
	return *i, *o, *d, *n
}
