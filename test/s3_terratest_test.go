package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/http-helper"
    "time"
)

// TestS3BucketWebsite verifies that the S3 website URL exists and is accessible
func TestS3BucketWebsite(t *testing.T) {
    t.Parallel()

    // Define Terraform options
    terraformOptions := &terraform.Options{
        // Path to the Terraform code that provisions the S3 bucket
        TerraformDir: "../terraform",
    }

    // Clean up resources after tests
    defer terraform.Destroy(t, terraformOptions)

    // Initialize and apply the Terraform configuration
    terraform.InitAndApply(t, terraformOptions)

    // Verify the website URL output
    websiteURL, err := terraform.OutputE(t, terraformOptions, "website_url")
    if err != nil {
        t.Fatalf("Failed to get website_url output: %v", err)
    }

    // Check if the website URL is not empty
    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")

    // Ensure the website URL is a valid format
    assert.Regexp(t, `^https?://`, websiteURL, "Website URL should start with http:// or https://")

    // Test the website endpoint using HTTP helper if a website URL is provided
    maxRetries := 10
    timeBetweenRetries := 10 * time.Second

    // Perform an HTTP GET request with retries
    http_helper.HttpGetWithRetry(t, websiteURL, nil, 200, "", maxRetries, timeBetweenRetries)
}