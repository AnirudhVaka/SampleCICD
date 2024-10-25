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

	// AWS Region
	awsRegion := "ap-south-1" // replace with your region

	// Terraform options to configure the Terraform apply
	terraformOptions := &terraform.Options{
		// Set the path to the Terraform code
		TerraformDir: "../terraform", // Update this path if necessary
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	// Ensure Terraform destroy is run at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Apply Terraform code
	terraform.InitAndApply(t, terraformOptions)

	// Check if the S3 bucket exists using AWS SDK
	bucketName := "samplewebsitebucket" // replace with your bucket name
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	assert.NoError(t, err)

	s3Client := s3.New(sess)
	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	// If err is nil, the bucket exists; otherwise, it does not
	bucketExists := err == nil
	assert.True(t, bucketExists, "The S3 bucket should exist")
}
