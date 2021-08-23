package mpegts

import "os"

type FileWriter struct {
	fp *os.File
}

func (fw *FileWriter) Create(filename string) (err error) {
	fw.fp, err = os.Create(filename)
	return
}

func (fw *FileWriter) Write(b []byte) (err error) {
	if fw.fp == nil {
		return ErrMpegts
	}
	_, err = fw.fp.Write(b)
	return
}

func (fw *FileWriter) Dispose() error {
	if fw.fp == nil {
		return ErrMpegts
	}
	return fw.fp.Close()
}

func (fw *FileWriter) Name() string {
	if fw.fp == nil {
		return ""
	}
	return fw.fp.Name()
}
