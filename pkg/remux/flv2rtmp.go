// Copyright 2020, Chef.  All rights reserved.
// https://github.com/yangjing0630/go-stream
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package remux

import (
	"github.com/yangjing0630/go-stream/pkg/base"
	"github.com/yangjing0630/go-stream/pkg/httpflv"
	"github.com/yangjing0630/go-stream/pkg/rtmp"
)

func FlvTagHeader2RtmpHeader(in httpflv.TagHeader) (out base.RtmpHeader) {
	out.MsgLen = in.DataSize
	out.MsgTypeId = in.Type
	out.MsgStreamId = rtmp.Msid1
	switch in.Type {
	case httpflv.TagTypeMetadata:
		out.Csid = rtmp.CsidAmf
	case httpflv.TagTypeAudio:
		out.Csid = rtmp.CsidAudio
	case httpflv.TagTypeVideo:
		out.Csid = rtmp.CsidVideo
	}
	out.TimestampAbs = in.Timestamp
	return
}

// @return 返回的内存块引用参数输入的内存块
func FlvTag2RtmpMsg(tag httpflv.Tag) (msg base.RtmpMsg) {
	msg.Header = FlvTagHeader2RtmpHeader(tag.Header)
	msg.Payload = tag.Raw[11 : 11+msg.Header.MsgLen]
	return
}
