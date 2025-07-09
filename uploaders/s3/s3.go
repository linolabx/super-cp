package s3

import (
	"fmt"
	"os"
	"slices"
	"sync"

	"git.sxxfuture.net/taojiayi/super-cp/core"
	"git.sxxfuture.net/taojiayi/super-cp/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/linolabx/lino_s3"
	"github.com/spf13/viper"
)

type S3Uploader struct {
	*lino_s3.LinoS3Path
}

func (s *S3Uploader) String() string {
	return fmt.Sprintf("s3[%s]", s.LinoS3Path.Object("").Key())
}

var ExcludedHeaders = []string{
	"Content-Type",
	"Content-Disposition",
	"Cache-Control",
	"Content-Language",
	"Content-Encoding",
}

func (s *S3Uploader) UploadSingle(spFile *core.SourceFile) error {
	file, err := os.Open(spFile.LocalPath)
	if err != nil {
		return err
	}
	defer file.Close()

	objectInput := s3.PutObjectInput{
		Body:     file,
		Metadata: map[string]string{},
	}

	objectInput.ContentType = aws.String(spFile.Metadata["Content-Type"])
	objectInput.ContentDisposition = aws.String(spFile.Metadata["Content-Disposition"])
	objectInput.CacheControl = aws.String(spFile.Metadata["Cache-Control"])
	objectInput.ContentLanguage = aws.String(spFile.Metadata["Content-Language"])
	objectInput.ContentEncoding = aws.String(spFile.Metadata["Content-Encoding"])

	for k, v := range spFile.Metadata {
		if slices.Contains(ExcludedHeaders, k) {
			continue
		}
		objectInput.Metadata[k] = v
	}

	if viper.GetBool("dry-run") {
		return nil
	}

	if _, err := s.LinoS3Path.Object(spFile.RemotePath).Put(objectInput); err != nil {
		return err
	}

	return nil
}

func (s *S3Uploader) Upload(files []*core.SourceFile) error {
	semaphore := make(chan struct{}, viper.GetInt("concurrency"))

	var wg sync.WaitGroup
	wg.Add(len(files))

	// Channel to signal all goroutines to stop
	stopChan := make(chan struct{})

	// Channel to collect errors
	errChan := make(chan error, len(files))

	for _, file := range files {
		go func(f *core.SourceFile) {
			defer wg.Done()

			select {
			case <-stopChan:
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }() // Release
			}

			if err := s.UploadSingle(f); err != nil {
				errChan <- fmt.Errorf("failed to upload %s: %v", f.LocalPath, err)
				close(stopChan) // Signal all goroutines to stop
				return
			}

			select {
			case <-stopChan:
				return
			default:
				utils.Verbosef("Uploaded: %s -> %s", f.LocalPath, f.RemotePath)
			}
		}(file)
	}

	wg.Wait()

	// Check if any error occurred
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func init() {
	core.RegisterUploader("s3", func(dsn string) (core.Uploader, error) {
		s3Path, err := lino_s3.LoadS3Path(dsn)
		if err != nil {
			return nil, err
		}
		return &S3Uploader{s3Path}, nil
	})
}
