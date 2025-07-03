package s3

import (
	"os"

	"git.sxxfuture.net/taojiayi/super-cp/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/linolabx/lino_s3"
)

var _ core.Target = &S3{}

type S3 struct {
	*lino_s3.LinoS3Path
}

func (s *S3) Init(dsn string) error {
	s.LinoS3Path = lino_s3.MustLoadS3Path(dsn)
	return nil
}

func (s *S3) Upload(spFile core.SPFile) error {
	metadata := spFile.Metadata

	file, err := os.Open(spFile.LocalPath)
	if err != nil {
		return err
	}
	defer file.Close()

	objectInput := s3.PutObjectInput{
		Body:     file,
		Metadata: metadata,
	}
	objectInput.ContentType = aws.String(metadata["Content-Type"])
	objectInput.ContentDisposition = aws.String(metadata["Content-Disposition"])
	objectInput.CacheControl = aws.String(metadata["Cache-Control"])
	objectInput.ContentLanguage = aws.String(metadata["Content-Language"])
	objectInput.ContentEncoding = aws.String(metadata["Content-Encoding"])

	if _, err := s.LinoS3Path.Object(spFile.RemotePath).Put(objectInput); err != nil {
		return err
	}

	return nil
}

func init() {
	core.Targets["@s3"] = &S3{}
}
