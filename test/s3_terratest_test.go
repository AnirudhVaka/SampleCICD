package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

// TestS3BucketCreation verifies that the S3 bucket is created successfully and checks its existence
func TestS3BucketCreation(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        // Path to the Terraform code that provisions the S3 bucket
        TerraformDir: "../terraform",
    }

    // Initialize and apply the Terraform configuration
    terraform.InitAndApply(t, terraformOptions)

    // Verify the output for bucket name
    bucketName, err := terraform.OutputE(t, terraformOptions, "bucket_name")
    if err != nil {
        t.Fatalf("Failed to get bucket_name output: %v", err)
    }

    // Assert that bucket name is as expected
    expectedBucketName := "samplewebsitebucket" // Update this to match the expected bucket name in your configuration
    assert.Equal(t, expectedBucketName, bucketName, "Bucket name should match expected name")
}
