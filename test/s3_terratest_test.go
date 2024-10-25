package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func cleanup(t *testing.T, terraformOptions *terraform.Options) {
	terraform.Destroy(t, terraformOptions)
}

// Test if the S3 bucket and its website configuration are created and accessible
func TestS3WebsiteBucket(t *testing.T) {
	t.Parallel()
	region := "ap-south-1"              // Update to your region if different
	bucketName := "samplewebsitebucket" // Matches bucket name in main.tf

	// Configure Terraform options
	terraformOptions := &terraform.Options{
		TerraformDir: "../..", // Path to the directory where main.tf is located
		Vars: map[string]interface{}{
			"bucket": bucketName,
		},
	}

	// Cleanup resources with Terraform destroy at the end
	defer cleanup(t, terraformOptions)

	// Run `terraform init` and `terraform apply` and fail if any errors occur
	terraform.InitAndApply(t, terraformOptions)

	// Test 1: Check if the bucket was created
	actualBucketID := terraform.Output(t, terraformOptions, "bucket_id")
	assert.Equal(t, bucketName, actualBucketID)

	// Test 2: Verify bucket exists in AWS
	bucketExists := aws.DoesS3BucketExist(t, region, bucketName)
	assert.True(t, bucketExists, "Expected S3 bucket to exist")

	// Test 3: Check website configuration - index document
	websiteConfig := aws.GetS3BucketWebsiteConfiguration(t, region, bucketName)
	assert.Equal(t, "index.html", websiteConfig.IndexDocumentSuffix)
	assert.Equal(t, "error.html", websiteConfig.ErrorDocumentKey)

	// Test 4: Check bucket policy - public access allowed
	bucketPolicy := aws.GetS3BucketPolicy(t, region, bucketName)
	assert.Contains(t, bucketPolicy, "s3:GetObject")
	assert.Contains(t, bucketPolicy, "arn:aws:s3:::"+bucketName+"/*")

	// Test 5: Check that the public access block is set correctly
	publicAccessBlock := aws.GetS3BucketPublicAccessBlock(t, region, bucketName)
	assert.False(t, publicAccessBlock.BlockPublicAcls)
	assert.False(t, publicAccessBlock.BlockPublicPolicy)
	assert.False(t, publicAccessBlock.IgnorePublicAcls)
	assert.False(t, publicAccessBlock.RestrictPublicBuckets)

	// Test 6: Confirm the website endpoint output
	expectedWebsiteURL := "http://" + bucketName + ".s3-website-" + region + ".amazonaws.com"
	websiteURL := terraform.Output(t, terraformOptions, "website_url")
	assert.Equal(t, expectedWebsiteURL, websiteURL)
}
