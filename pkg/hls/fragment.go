package hls

import (
	"github.com/q191201771/naza/pkg/filesystemlayer"
)

type Fragment struct {
	fp       filesystemlayer.IFile
	filename string
}

func (f *Fragment) OpenFile(filename string) (err error) {
	f.fp, err = fslCtx.Create(filename)
	if err != nil {
		return
	}
	f.filename = filename
	return
}

func (f *Fragment) WriteFile(b []byte) (err error) {
	_, err = f.fp.Write(b)
	return
}

func (f *Fragment) CloseFile() error {
	return f.fp.Close()
}

func (f *Fragment) FileName() string {
	return f.filename
}
