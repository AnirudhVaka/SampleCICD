package test

import (
    "testing"

    "github.com/gruntwork-io/terratest/modules/aws"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestS3Website(t *testing.T) {
    t.Parallel()

    // Specify the path to your Terraform code
    tfOptions := &terraform.Options{
        TerraformDir: "../path/to/your/terraform/code", // Update this path accordingly

        Vars: map[string]interface{}{
            "bucket_name": "samplewebsitebucket",
        },

        NoColor: true,
    }

    // Clean up resources with 'terraform destroy' at the end of the test
    defer terraform.Destroy(t, tfOptions)

    // Run 'terraform init' and 'terraform apply'. Fail the test if there are any errors.
    terraform.InitAndApply(t, tfOptions)

    // Check if the S3 bucket exists
    bucketName := "samplewebsitebucket"
    bucketExists := aws.GetS3Bucket(t, bucketName)
    
    assert.NotNil(t, bucketExists, "Expected S3 bucket to exist")

    // Check if the website configuration is set up correctly
    bucketWebsiteURL := terraform.Output(t, tfOptions, "website_url")
    
    assert.NotEmpty(t, bucketWebsiteURL, "Expected website URL to be set")
    
    t.Logf("Website URL: %s\n", bucketWebsiteURL)
}
