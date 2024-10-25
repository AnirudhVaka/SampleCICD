package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	foundRegion := aws.FindS3BucketRegion(t, bucketName, region)
	assert.NotEmpty(t, foundRegion, "Expected S3 bucket to exist")

	// Set up AWS session
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess, &aws.Config{Region: &region})

	// Test 3: Check website configuration - index document
	websiteConfig, err := s3Client.GetBucketWebsite(&s3.GetBucketWebsiteInput{
		Bucket: &bucketName,
	})
	assert.NoError(t, err, "Expected to retrieve website configuration")
	assert.Equal(t, "index.html", *websiteConfig.IndexDocument.Suffix)
	assert.Equal(t, "error.html", *websiteConfig.ErrorDocument.Key)

	// Test 4: Check bucket policy - public access allowed
	bucketPolicy, err := s3Client.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: &bucketName,
	})
	assert.NoError(t, err, "Expected to retrieve bucket policy")
	assert.Contains(t, *bucketPolicy.Policy, "s3:GetObject")
	assert.Contains(t, *bucketPolicy.Policy, "arn:aws:s3:::"+bucketName+"/*")

	// Test 5: Check that the public access block is set correctly
	publicAccessBlock, err := s3Client.GetBucketPolicyStatus(&s3.GetBucketPolicyStatusInput{
		Bucket: &bucketName,
	})
	assert.NoError(t, err, "Expected to retrieve public access block")
	assert.NotNil(t, publicAccessBlock.PolicyStatus.IsPublic, "Expected bucket to have a public access setting")

	// Test 6: Confirm the website endpoint output
	expectedWebsiteURL := "http://" + bucketName + ".s3-website-" + region + ".amazonaws.com"
	websiteURL := terraform.Output(t, terraformOptions, "website_url")
	assert.Equal(t, expectedWebsiteURL, websiteURL)
}
