package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestS3Website(t *testing.T) {
	t.Parallel()

	// Specify the path to your Terraform code
	tfOptions := &terraform.Options{
		TerraformDir: "../terraform", // Update this path accordingly

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"bucket_name": "samplewebsitebucket",
		},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	}

	// Clean up resources with 'terraform destroy' at the end of the test
	defer terraform.Destroy(t, tfOptions)

	// This will run 'terraform init' and 'terraform apply'. Fail the test if there are any errors.
	initAndApply := terraform.InitAndApply(t, tfOptions)

	// Check that the bucket was created or exists
	bucketExists := aws.S3BucketExists(t, "samplewebsitebucket")
	assert.True(t, bucketExists, "Expected S3 bucket to exist")

	if initAndApply {
		// Optionally check if the website configuration is set up correctly
		bucketWebsiteURL := terraform.Output(t, tfOptions, "website_url")

		assert.NotEmpty(t, bucketWebsiteURL, "Expected website URL to be set")

		fmt.Printf("Website URL: %s\n", bucketWebsiteURL)

		// Here you could add more assertions to check if the website is accessible
	}
}
