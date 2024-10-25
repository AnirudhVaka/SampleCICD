// test/s3_terratest.go
package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestS3Bucket(t *testing.T) {
    t.Parallel()

    // Define Terraform options
    terraformOptions := &terraform.Options{
        TerraformDir: "../", // Path to your main Terraform directory
    }

    // Ensure Terraform destroy is run after tests to clean up resources
    defer terraform.Destroy(t, terraformOptions)

    // Run Terraform init and apply
    terraform.InitAndApply(t, terraformOptions)

    // Fetch and verify the S3 bucket name
    bucketID := terraform.Output(t, terraformOptions, "bucket_id")
    assert.Contains(t, bucketID, "samplewebsitebucket")

    // Fetch and verify the website endpoint
    websiteURL := terraform.Output(t, terraformOptions, "website_url")
    assert.Contains(t, websiteURL, "s3-website")
}
