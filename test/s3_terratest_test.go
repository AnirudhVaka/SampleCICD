package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestS3WebsiteHttpAccess(t *testing.T) {
	t.Parallel()

	awsRegion := "ap-south-1"
	bucketName := "samplewebsitebucket"

	// Terraform options for applying configuration
	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	// Ensure Terraform destroy is run at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Apply Terraform code
	terraform.InitAndApply(t, terraformOptions)

	// Set up AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	assert.NoError(t, err)

	// Verify the bucket is configured as a website
	s3Client := s3.New(sess)
	_, err = s3Client.GetBucketWebsite(&s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucketName),
	})
	assert.NoError(t, err, "The S3 bucket should be configured as a website")

	// Construct the website URL
	websiteURL := fmt.Sprintf("http://%s.s3-website-%s.amazonaws.com", bucketName, awsRegion)

	// Make an HTTP GET request to the website URL
	resp, err := http.Get(websiteURL)
	assert.NoError(t, err, "HTTP request to the S3 website should not error")

	// Check that the response status is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status 200 OK for the S3 website")
}
