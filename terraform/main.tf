# Local variable to control bucket creation
locals {
  create_bucket = true  # Set to false if you don't want to create the bucket
}

# Conditionally create the S3 bucket
resource "aws_s3_bucket" "website_bucket" {
  count  = local.create_bucket ? 1 : 0
  bucket = "samplewebsitebucket"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

# Apply a bucket policy for public access
resource "aws_s3_bucket_policy" "website_policy" {
  count  = local.create_bucket ? 1 : 0
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
}

output "website_url" {
  value       = aws_s3_bucket.website_bucket[0].website_endpoint
  description = "The URL of the static website hosted on S3"
  depends_on  = [aws_s3_bucket.website_bucket]
}