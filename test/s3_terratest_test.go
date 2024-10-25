package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestS3BucketVersioning(t *testing.T) {
	// Specify the Terraform options
	terraformOptions := &terraform.Options{
		// Path to the Terraform code
		TerraformDir: "../terraform",
	}

	// Apply and defer cleanup
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Check if the S3 bucket exists and has versioning enabled
	bucketID := terraform.Output(t, terraformOptions, "bucket_id")
	region := "ap-south-1"
	isVersioningEnabled := aws.IsS3BucketVersioningEnabled(t, region, bucketID)
	assert.True(t, isVersioningEnabled, "Bucket versioning should be enabled.")
}
