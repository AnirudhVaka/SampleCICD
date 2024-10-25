resource "aws_s3_bucket_website_configuration" "website" {
  bucket = "samplewebvaultbucket"  # Existing bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

resource "aws_s3_bucket_policy" "website_policy" {
  bucket = "samplewebvaultbucket"  # Existing bucket

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "arn:aws:s3:::samplewebvaultbucket/*"
      }
    ]
  })
}
