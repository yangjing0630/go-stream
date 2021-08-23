package httpflv

import "os"

// TODO chef: 结构体重命名为FileWriter，文件名重命名为file_writer.go。所有写流文件的（flv,hls,ts）统一重构

type FlvFileWriter struct {
	fp *os.File
}

func (ffw *FlvFileWriter) Open(filename string) (err error) {
	ffw.fp, err = os.Create(filename)
	return
}

func (ffw *FlvFileWriter) WriteRaw(b []byte) (err error) {
	if ffw.fp == nil {
		return ErrHttpflv
	}
	_, err = ffw.fp.Write(b)
	return
}

func (ffw *FlvFileWriter) WriteFlvHeader() (err error) {
	if ffw.fp == nil {
		return ErrHttpflv
	}
	_, err = ffw.fp.Write(FlvHeader)
	return
}

func (ffw *FlvFileWriter) WriteTag(tag Tag) (err error) {
	if ffw.fp == nil {
		return ErrHttpflv
	}
	_, err = ffw.fp.Write(tag.Raw)
	return
}

func (ffw *FlvFileWriter) Dispose() error {
	if ffw.fp == nil {
		return ErrHttpflv
	}
	return ffw.fp.Close()
}

func (ffw *FlvFileWriter) Name() string {
	if ffw.fp == nil {
		return ""
	}
	return ffw.fp.Name()
}
