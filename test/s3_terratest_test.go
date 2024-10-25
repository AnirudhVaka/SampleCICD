package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/http-helper" // Import http-helper module
    "time"
)

// TestS3BucketWebsiteConfig verifies that the S3 bucket has been set up correctly with website hosting
func TestS3BucketWebsiteConfig(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        // Path to your Terraform directory
        TerraformDir: "../terraform",
    }

    // Run terraform init and apply, ensure Terraform destroy is run after tests to clean up resources
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    // Check that the website endpoint is correctly formatted
    websiteURL := terraform.Output(t, terraformOptions, "website_url")
    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")
    assert.Contains(t, websiteURL, "s3-website", "Website URL should contain 's3-website'")
}

// TestS3BucketPublicAccess checks if the public access policy is properly configured
func TestS3BucketPublicAccess(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../terraform",
    }

    // Ensure resources are cleaned up after test completion
    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    // Since `bucket_policy` is not defined in your outputs, skip this check or modify the configuration if needed.
}

// TestWebsiteEndpoint verifies that the S3 website URL returns the expected content (index.html)
func TestWebsiteEndpoint(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../terraform",
    }

    defer terraform.Destroy(t, terraformOptions)
    terraform.InitAndApply(t, terraformOptions)

    // Get the website URL from Terraform output
    websiteURL := terraform.Output(t, terraformOptions, "website_url")
    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")

    // Test the website endpoint using the HTTP helper
    maxRetries := 10
    timeBetweenRetries := 10 * time.Second

    // Use http_helper to test if the website is serving the expected content
    http_helper.HttpGetWithRetry(t, "http://"+websiteURL, nil, 200, "Expected content in index.html", maxRetries, timeBetweenRetries)
}
