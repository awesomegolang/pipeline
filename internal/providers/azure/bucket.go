package azure

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/Azure/azure-storage-blob-go/2016-05-31/azblob"
	"github.com/goph/emperror"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const containerUrlTemplate = "https://%s.blob.core.windows.net/%s"

// CreateBucket creates a bucket.
func CreateBucket(
	ctx context.Context,
	client storage.AccountsClient,
	resourceGroup string,
	storageAccount string,
	bucketName string,
) error {
	key, err := GetStorageAccountKey(ctx, client, resourceGroup, storageAccount)
	if err != nil {
		return emperror.With(
			errors.Wrap(err, "failed to create bucket"),
			"resource-group", resourceGroup,
			"bucket", bucketName,
		)
	}

	p := azblob.NewPipeline(azblob.NewSharedKeyCredential(storageAccount, key), azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf(containerUrlTemplate, storageAccount, bucketName)) // TODO: handle error
	containerURL := azblob.NewContainerURL(*URL, p)

	_, err = containerURL.GetPropertiesAndMetadata(ctx, azblob.LeaseAccessConditions{})
	if err != nil && err.(azblob.StorageError).ServiceCode() == azblob.ServiceCodeContainerNotFound { // Bucket not found, so create it
		_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
		if err != nil {
			return emperror.With(
				errors.Wrap(err, "failed to create bucket"),
				"resource-group", resourceGroup,
				"bucket", bucketName,
			)
		}
	}

	return nil
}

// DeleteBucketState removes a bucket from the database.
func DeleteBucketState(db *gorm.DB, bucketID uint) error {
	return db.Delete(&ObjectStoreBucketModel{ID: bucketID}).Error
}
