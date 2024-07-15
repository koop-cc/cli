package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/supabase/cli/pkg/config"
)

func (s *StorageAPI) UpsertBuckets(ctx context.Context, bucketConfig config.BucketConfig, filter ...func(string) bool) error {
	buckets, err := s.ListBuckets(ctx)
	if err != nil {
		return err
	}
	exists := make(map[string]string, len(buckets))
	for _, b := range buckets {
		exists[b.Name] = b.Id
	}
	for name, bucket := range bucketConfig {
		props := BucketProps{
			Public:           bucket.Public,
			FileSizeLimit:    int(bucket.FileSizeLimit),
			AllowedMimeTypes: bucket.AllowedMimeTypes,
		}
		// Update bucket properties if already exists
		if bucketId, ok := exists[name]; ok {
			for _, keep := range filter {
				if !keep(bucketId) {
					continue
				}
			}
			fmt.Fprintln(os.Stderr, "Updating storage bucket:", bucketId)
			body := UpdateBucketRequest{
				Id:          bucketId,
				BucketProps: &props,
			}
			if _, err := s.UpdateBucket(ctx, body); err != nil {
				return err
			}
		} else {
			fmt.Fprintln(os.Stderr, "Creating storage bucket:", name)
			body := CreateBucketRequest{
				Name:        name,
				BucketProps: &props,
			}
			if _, err := s.CreateBucket(ctx, body); err != nil {
				return err
			}
		}
	}
	return nil
}