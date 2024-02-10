package main

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
	"log"
	"math"
	"path/filepath"
	"strings"
)

const (
	nDropBoxWorkers = 5
	nAwsWorkers     = 5
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	dropboxClient, err := NewDropBoxClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create DropBox client: %w", err)
	}

	awsClient, err := NewAwsClient()
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	if false {
		if err := awsClient.CleanReset(ctx); err != nil {
			return fmt.Errorf("failed to reset bucket: %w", err)
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	downloadCh := make(chan File, nDropBoxWorkers)
	type Upload struct {
		Path string
		Data []byte
	}
	uploadCh := make(chan Upload, nAwsWorkers)

	// producers
	g.Go(func() error {
		defer close(downloadCh)
		err := dropboxClient.IterateFiles(ctx, math.MaxUint32, func(ctx context.Context, file File) error {
			if file.Tag == "file" {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case downloadCh <- file:
					log.Println("iterate:", file.ID, file.Size, file.PathDisplay)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to iterate: %w", err)
		}
		return nil
	})

	// download from dropbox
	dbWorkers := atomic.NewInt32(nDropBoxWorkers)
	for i := 0; i < nDropBoxWorkers; i++ {
		g.Go(func() error {
			// last one out turns off the lights
			defer func() {
				if dbWorkers.Dec() == 0 {
					close(uploadCh)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case work, ok := <-downloadCh:
					if !ok {
						return nil // all done
					}

					path := work.PathDisplay
					basePath := filepath.Base(path)

					data, err := dropboxClient.Download(ctx, path)
					if err != nil {
						return fmt.Errorf("failed to download: %w", err)
					}
					text := string(data)
					count := len(strings.Fields(text))
					if count >= 2000 && count <= 5000 {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case uploadCh <- Upload{Path: basePath, Data: data}:
							log.Println("downloaded:", count, basePath)
						}
					} else {
						log.Println("skip:", count, basePath)
					}
				}
			}

		})
	}

	// upload to AWS
	uploadCount := atomic.NewInt32(0)
	for i := 0; i < nAwsWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case work, ok := <-uploadCh:
					if !ok {
						return nil // all done
					}
					if err := awsClient.Upload(ctx, work.Path, work.Data); err != nil {
						return fmt.Errorf("failed to upload: %w", err)
					}
					log.Println("uploaded:", work.Path)
					uploadCount.Inc()
				}
			}
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	log.Println(uploadCount.Load(), "uploaded") // 263
	return nil
}
