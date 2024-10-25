resource "aws_s3_bucket" "website_bucket" {
  bucket = "samplewebvaultbucket"  # Replace with a unique name
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

resource "aws_s3_bucket_object" "index" {
  bucket = aws_s3_bucket.website_bucket.bucket
  key    = "index.html"
  source = "${path.module}/../website/index.html"
  acl    = "public-read"
  content_type = "text/html"
}

output "website_url" {
  value = aws_s3_bucket.website_bucket.website_endpoint
}
