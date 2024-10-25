# Conditional creation of S3 bucket as described earlier

# Data source to check if the bucket exists
data "aws_s3_bucket" "existing_bucket" {
  bucket = "samplewebvaultbucket"
  count  = length(try(aws_s3_bucket.website_bucket.id, [])) == 0 ? 0 : 1
}

# Conditionally create the bucket if it does not already exist
resource "aws_s3_bucket" "website_bucket" {
  bucket = "samplewebvaultbucket"
  count = length(data.aws_s3_bucket.existing_bucket) == 0 ? 1 : 0

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

# Bucket policy for public access
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
  count = length(data.aws_s3_bucket.existing_bucket) == 0 ? 1 : 0
}

output "website_url" {
  value = aws_s3_bucket.website_bucket[0].website_endpoint
  depends_on = [aws_s3_bucket.website_bucket]
}
