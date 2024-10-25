package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/http-helper" // Import http-helper module
    "time"
)

// TestWebsiteEndpoint verifies that the S3 website URL returns a successful status
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

    // Check if the website endpoint is reachable with an HTTP 200 response
    http_helper.HttpGetWithRetry(t, "http://"+websiteURL, nil, 200, "", maxRetries, timeBetweenRetries)

    // Ensure website URL output is empty post-destroy
    terraform.Destroy(t, terraformOptions)
    websiteURLAfterDestroy := terraform.Output(t, terraformOptions, "website_url")
    assert.Empty(t, websiteURLAfterDestroy, "Website URL should be empty after destroy")
}
