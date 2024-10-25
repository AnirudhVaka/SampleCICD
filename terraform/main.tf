# Check if the S3 bucket exists
data "aws_s3_bucket" "existing" {
  bucket = "samplewebsitebucket"
}

locals {
  bucket_exists = try(data.aws_s3_bucket.existing.id != "", false)
}

# Conditionally create the bucket if it does not exist
resource "aws_s3_bucket" "website_bucket" {
  count  = local.bucket_exists ? 0 : 1
  bucket = "samplewebsitebucket"
}

# Apply a website configuration if the bucket is created by Terraform
resource "aws_s3_bucket_website_configuration" "website" {
  count  = local.bucket_exists ? 0 : 1
  bucket = aws_s3_bucket.website_bucket[0].id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

# Disable Block Public Access to allow public bucket policies
resource "aws_s3_bucket_public_access_block" "public_access" {
  count = local.bucket_exists ? 0 : 1
  bucket = coalesce(
    try(data.aws_s3_bucket.existing.id, null),
    try(aws_s3_bucket.website_bucket[0].id, null)
  )

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# Apply a bucket policy to allow public read access to objects
resource "aws_s3_bucket_policy" "website_policy" {
  bucket = coalesce(
    try(data.aws_s3_bucket.existing.id, null),
    try(aws_s3_bucket.website_bucket[0].id, null)
  )

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Principal = "*",
        Action    = "s3:GetObject",
        Resource  = "arn:aws:s3:::samplewebsitebucket/*"
      }
    ]
  })
}

# Output the S3 website URL, depending on whether the bucket is pre-existing or newly created
output "website_url" {
  value       = coalesce(data.aws_s3_bucket.existing.website_endpoint, try(aws_s3_bucket_website_configuration.website[0].website_endpoint, null))
  description = "The URL of the static website hosted on S3"
}
