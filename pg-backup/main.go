package main

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DumpAll       bool   `envconfig:"DUMP_ALL"`
	BucketName    string `envconfig:"BUCKET_NAME"`
	ObjectKey     string `envconfig:"OBJECT_KEY"`
	UploadTempDir string `envconfig:"UPLOAD_TEMP_DIR"`
}

type PgBackup struct {
	client *s3.Client
	config Config
}

func (p *PgBackup) CompressAndSave(ctx context.Context, name string, data io.Reader) error {
	if c, ok := data.(io.Closer); ok {
		defer c.Close()
	}

	tempFile, err := os.CreateTemp(p.config.UploadTempDir, "pg-backup-*")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	_ = os.Remove(tempFile.Name())

	gzWriter := gzip.NewWriter(tempFile)

	if _, err := io.Copy(gzWriter, data); err != nil {
		_ = gzWriter.Close()
		return err
	}
	_ = gzWriter.Close()

	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	_, err = p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &p.config.BucketName,
		Key:    &name,
		Body:   tempFile,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *PgBackup) GetDumpCommand() string {
	if p.config.DumpAll {
		return "pg_dumpall"
	} else {
		return "pg_dump"
	}
}

func (p *PgBackup) BackupDB(ctx context.Context) error {
	dumpCommand := exec.Command(p.GetDumpCommand())
	dumpCommand.Stderr = os.Stderr
	dumpCommandStdout, err := dumpCommand.StdoutPipe()
	if err != nil {
		return err
	}
	defer dumpCommandStdout.Close()

	if err := dumpCommand.Start(); err != nil {
		return err
	}

	if err := p.CompressAndSave(ctx, p.config.ObjectKey, dumpCommandStdout); err != nil {
		return err
	}

	if err := dumpCommand.Wait(); err != nil {
		return err
	}

	return nil
}

func MustLoadAWSConfig(ctx context.Context) aws.Config {
	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	ctx := context.Background()

	client := s3.NewFromConfig(MustLoadAWSConfig(ctx))

	var config Config
	envconfig.MustProcess("", &config)

	pgBackup := PgBackup{
		client: client,
		config: config,
	}

	if err := pgBackup.BackupDB(ctx); err != nil {
		log.Fatal(err)
	}
}
