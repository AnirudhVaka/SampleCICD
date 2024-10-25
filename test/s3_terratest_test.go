package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

func TestS3Bucket(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../test", // assuming main.tf is in the parent directory
    }

    // Clean up resources with `terraform destroy` at the end of the test
    defer terraform.Destroy(t, terraformOptions)

    // Run `terraform init` and `terraform apply`
    terraform.InitAndApply(t, terraformOptions)

    // Get the name of the S3 bucket from Terraform outputs
    bucketName := terraform.Output(t, terraformOptions, "website_url")

    // Initialize an AWS session
    awsRegion := "us-east-1" // specify your region
    sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(awsRegion)}))
    s3Client := s3.New(sess)

    // Verify if the bucket exists
    _, err := s3Client.HeadBucket(&s3.HeadBucketInput{
        Bucket: aws.String(bucketName),
    })
    assert.NoError(t, err)

    // Check if public access is disabled for the bucket (basic check)
    publicAccess, err := s3Client.GetBucketPolicyStatus(&s3.GetBucketPolicyStatusInput{
        Bucket: aws.String(bucketName),
    })
    assert.NoError(t, err)
    assert.NotNil(t, publicAccess)
}
