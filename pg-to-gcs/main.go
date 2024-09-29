package main

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/api/option"
)

type Config struct {
	DumpAll         bool   `envconfig:"DUMP_ALL"`
	CredentialsJSON string `envconfig:"CREDENTIALS_JSON"`
	BucketName      string `envconfig:"BUCKET_NAME"`
	ObjectKey       string `envconfig:"OBJECT_KEY"`
}

type PgToGCS struct {
	client *storage.Client
	config Config
}

func (p *PgToGCS) CompressAndSave(ctx context.Context, name string, data io.Reader) error {
	if c, ok := data.(io.Closer); ok {
		defer c.Close()
	}

	objectWriter := p.client.Bucket(p.config.BucketName).Object(name).NewWriter(ctx)
	defer objectWriter.Close()

	gzWriter := gzip.NewWriter(objectWriter)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, data); err != nil {
		return err
	}

	return nil
}

func (p *PgToGCS) GetDumpCommand() string {
	if p.config.DumpAll {
		return "pg_dumpall"
	} else {
		return "pg_dump"
	}
}

func (p *PgToGCS) BackupDB(ctx context.Context) error {
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

func main() {
	var config Config
	envconfig.MustProcess("", &config)

	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(config.CredentialsJSON)))
	if err != nil {
		log.Fatal(err)
	}

	pgToGCS := PgToGCS{
		client: client,
		config: config,
	}

	if err := pgToGCS.BackupDB(ctx); err != nil {
		log.Fatal(err)
	}
}
