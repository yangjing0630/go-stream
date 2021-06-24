// Copyright 2021, Chef.  All rights reserved.
// https://github.com/yangjing0630/go-stream
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"context"
	"time"

	"github.com/yangjing0630/go-stream/pkg/base"
	"github.com/yangjing0630/go-stream/pkg/rtmp"
)

// 检查远端rtmp流是否能正常拉取
func StreamExist(url string) error {
	const (
		timeoutMs = 10000
	)

	errChan := make(chan error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutMs*time.Millisecond)
	defer cancel()

	s := rtmp.NewPullSession()
	defer s.Dispose()

	go func() {
		// 有的场景只有音频没有视频，所以我们不检查视频
		var hasNotify bool
		var readMetadata bool
		var readAudio bool
		err := s.Pull(url, func(msg base.RtmpMsg) {
			if hasNotify {
				return
			}

			switch msg.Header.MsgTypeId {
			case base.RtmpTypeIdMetadata:
				readMetadata = true
			case base.RtmpTypeIdAudio:
				readAudio = true
			}
			if readMetadata && readAudio {
				hasNotify = true
				errChan <- nil
			}
		})
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
