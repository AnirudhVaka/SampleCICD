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
        TerraformDir: "../",
    }

    // Ensure Terraform destroy is run after tests to clean up resources
    defer terraform.Destroy(t, terraformOptions)

    // Run terraform init and apply
    terraform.InitAndApply(t, terraformOptions)

    // Verify the bucket name is as expected
    bucketID := terraform.Output(t, terraformOptions, "bucket_id")
    assert.Contains(t, bucketID, "samplewebsitebucket", "Bucket ID should contain 'samplewebsitebucket'")

    // Check that the website endpoint is correctly formatted
    websiteURL := terraform.Output(t, terraformOptions, "website_url")
    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")
    assert.True(t, websiteURL != "", "Website URL should not be empty")
    assert.Contains(t, websiteURL, "s3-website", "Website URL should contain 's3-website'")
}

// TestS3BucketPublicAccess checks if the public access policy is properly configured
func TestS3BucketPublicAccess(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../",
    }

    // Ensure resources are cleaned up after test completion
    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    // Retrieve and verify the bucket policy
    bucketPolicy := terraform.Output(t, terraformOptions, "bucket_policy")
    assert.Contains(t, bucketPolicy, `"Effect": "Allow"`, "Bucket policy should allow public access")
    assert.Contains(t, bucketPolicy, `"Principal": "*"`, "Bucket policy should allow public access for everyone")
}

// TestWebsiteEndpoint verifies that the S3 website URL returns the expected content (index.html)
func TestWebsiteEndpoint(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../",
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    // Get the website URL from Terraform output
    websiteURL := terraform.Output(t, terraformOptions, "website_url")

    // Test the website endpoint using the HTTP helper
    maxRetries := 10
    timeBetweenRetries := 10 * time.Second

    // Use http_helper to test if the website is serving the expected content
    http_helper.HttpGetWithRetry(t, "http://"+websiteURL, nil, 200, "index.html", maxRetries, timeBetweenRetries)
}

// TestCleanup verifies the cleanup process works correctly
func TestCleanup(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "../",
    }

    // Apply the configuration
    terraform.InitAndApply(t, terraformOptions)

    // Destroy the resources and verify cleanup
    terraform.Destroy(t, terraformOptions)
}
