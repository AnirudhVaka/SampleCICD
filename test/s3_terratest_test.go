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

	// Verify if the S3 bucket exists
	bucketName := "samplewebsitebucket" // replace with your bucket name
	bucketExists, err := aws.DoesS3BucketExistE(t, awsRegion, bucketName)
	assert.NoError(t, err)
	assert.True(t, bucketExists, "The S3 bucket should exist")
}
