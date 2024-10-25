package test

import (
    "encoding/json"
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/gruntwork-io/terratest/modules/http-helper"
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

    // Get the website URL from Terraform output in JSON format
    outputJson, err := terraform.OutputJsonE(t, terraformOptions, "website_url")
    if err != nil {
        t.Fatalf("Failed to get website_url output: %v", err)
    }

    // Unmarshal JSON to extract URL as a string
    var websiteURL string
    if err := json.Unmarshal([]byte(outputJson), &websiteURL); err != nil {
        t.Fatalf("Failed to parse website_url output: %v", err)
    }

    assert.NotEmpty(t, websiteURL, "Website URL should not be empty")

    // Test the website endpoint using the HTTP helper
    maxRetries := 10
    timeBetweenRetries := 10 * time.Second

    // Check if the website endpoint is reachable with an HTTP 200 response
    http_helper.HttpGetWithRetry(t, "http://"+websiteURL, nil, 200, "", maxRetries, timeBetweenRetries)
}
