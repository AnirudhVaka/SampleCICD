package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
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

	// Verify if the S3 bucket exists by listing all S3 buckets
	bucketName := "samplewebsitebucket" // replace with your bucket name
	buckets := aws.ListS3Buckets(t, awsRegion)

	bucketExists := false
	for _, bucket := range buckets {
		if bucket == bucketName {
			bucketExists = true
			break
		}
	}

	assert.True(t, bucketExists, "The S3 bucket should exist")
}
