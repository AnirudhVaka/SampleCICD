# Provider configuration for AWS
provider "aws" {
  region = "ap-south-1"
}

# Data source to check if the bucket exists
data "aws_s3_bucket" "existing_bucket" {
  bucket = "samplewebsitebucket"
  # Ignore errors if the bucket doesn't exist, which will allow conditional creation
  count  = length(try(aws_s3_bucket.website_bucket.id, [])) == 0 ? 0 : 1
}

# Conditionally create the bucket if it does not already exist
resource "aws_s3_bucket" "website_bucket" {
  bucket = "samplewebsitebucket"
  
  # Only create the bucket if the data source didn't find an existing one
  count = length(data.aws_s3_bucket.existing_bucket) == 0 ? 1 : 0

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

# Bucket policy to allow public read access
resource "aws_s3_bucket_policy" "website_policy" {
  bucket = aws_s3_bucket.website_bucket[0].id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = "*",
        Action    = "s3:GetObject",
        Resource  = "${aws_s3_bucket.website_bucket[0].arn}/*"
      }
    ]
  })
  
  # Only apply the policy if the bucket was created
  count = length(data.aws_s3_bucket.existing_bucket) == 0 ? 1 : 0
}

output "website_url" {
  value = aws_s3_bucket.website_bucket[0].website_endpoint
  description = "The URL of the static website hosted on S3"
  # Only output if the bucket was created
  depends_on = [aws_s3_bucket.website_bucket]
}
