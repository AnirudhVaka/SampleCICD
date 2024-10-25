package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/http-helper"
    "time"
)

// TestS3BucketWebsite verifies that the S3 bucket is created or exists, and that the website URL is accessible
func TestS3BucketWebsite(t *testing.T) {
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
    assert.NotEmpty(t, bucketName, "Bucket name should not be empty")

    // Verify the website URL output
    websiteURL, err := terraform.OutputE(t, terraformOptions, "website_url")
    if err != nil {
        t.Fatalf("Failed to get website_url output: %v", err)
    }
    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")

    // Test the website endpoint using the HTTP helper if a website URL is provided
    if websiteURL != "" {
        maxRetries := 10
        timeBetweenRetries := 10 * time.Second
        http_helper.HttpGetWithRetry(t, "http://"+websiteURL, nil, 200, "", maxRetries, timeBetweenRetries)
    }
}
