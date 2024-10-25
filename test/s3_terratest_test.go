package test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestS3BucketExistence(t *testing.T) {
	t.Parallel()

	awsRegion := "ap-south-1"

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Apply Terraform code
	terraform.InitAndApply(t, terraformOptions)

	// Check if the S3 bucket exists using AWS SDK
	bucketName := "samplewebsitebucket"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	assert.NoError(t, err)

	s3Client := s3.New(sess)
	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	bucketExists := err == nil
	assert.True(t, bucketExists, "The S3 bucket should exist")
}
