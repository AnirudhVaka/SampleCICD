// test/s3_terratest.go
package test

import (
    "testing"
    "strings"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

// TestS3BucketWebsiteConfig verifies that the S3 bucket has been set up correctly with website hosting
func TestS3BucketWebsiteConfig(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        // Path to where your Terraform code is located
        TerraformDir: "../",
    }

    // Clean up resources with terraform destroy at the end of the test
    defer terraform.Destroy(t, terraformOptions)

    // Run terraform init and apply
    terraform.InitAndApply(t, terraformOptions)

    // Verify the bucket name is as expected
    bucketID := terraform.Output(t, terraformOptions, "bucket_id")
    assert.Contains(t, bucketID, "samplewebsitebucket", "Bucket ID should contain 'samplewebsitebucket'")

    // Check that the website endpoint is correctly formatted
    websiteURL := terraform.Output(t, terraformOptions, "website_url")
    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")
    assert.True(t, strings.Contains(websiteURL, "s3-website"), "Website URL should contain 's3-website'")
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
    response, err := http_helper.HttpGet(t, "http://" + websiteURL, nil)

    // Verify we can reach the website endpoint and get a 200 status
    assert.NoError(t, err, "Should be able to reach the website URL")
    assert.Equal(t, 200, response.StatusCode, "Expected status code 200")
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
