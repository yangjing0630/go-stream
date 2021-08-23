package mpegts

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/q191201771/naza/pkg/nazalog"
	"io/ioutil"
)

//使用minio分布式存储扩展ts存储

type MinioWriter struct {
	EndPoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	FileName        string
	MinioClient     *minio.Client
}

func (mw *MinioWriter) NewMinio() error {
	minioClient, err := minio.New(mw.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(mw.AccessKeyID, mw.SecretAccessKey, ""),
		Secure: mw.UseSSL,
	})
	if err != nil {
		return err
	}
	if minioClient == nil {
		nazalog.Debugf("minioClient is nil aiaiaiai...")
	}
	mw.MinioClient = minioClient
	mw.BucketName = "candice"
	return nil
}

func (mw *MinioWriter) Create(filename string) (err error) {
	mw.FileName = filename
	return
}

func (mw *MinioWriter) Write(filename string, b []byte) error {
	if mw.MinioClient == nil {
		return ErrMpegts
	}
	reader := bytes.NewReader(b)
	_, err := mw.MinioClient.PutObject(context.Background(), mw.BucketName, filename, reader, int64(len(b)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func (mw *MinioWriter) FPutObject(path, filename string) error {
	_, MinioErr := mw.MinioClient.FPutObject(context.Background(), "candice", filename, path, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if MinioErr != nil {
		nazalog.Errorf("[minio FPutObject error]")
	}
	return MinioErr
}

func (mw *MinioWriter) Dispose() (err error) {
	return
}

func (mw *MinioWriter) Name() string {
	return mw.FileName
}

func (mw *MinioWriter) ReadFile(filename string) (b []byte, err error) {
	if mw.MinioClient == nil {
		return
	}
	reader, err := mw.MinioClient.GetObject(context.Background(), mw.BucketName, filename, minio.GetObjectOptions{})

	b, err = ioutil.ReadAll(reader)
	return
}
