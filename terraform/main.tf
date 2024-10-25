# Step 1: Check if the S3 bucket exists
data "aws_s3_bucket" "existing" {
  bucket = "samplewebsitebucket"  # Replace with your bucket name
}

# Step 2: Conditionally create the bucket if it does not exist
resource "aws_s3_bucket" "website_bucket" {
  count  = data.aws_s3_bucket.existing.id != "" ? 0 : 1
  bucket = "samplewebsitebucket"
}

# Step 3: Configure the bucket as a website if it's created by Terraform
resource "aws_s3_bucket_website_configuration" "website" {
  count  = length(aws_s3_bucket.website_bucket) > 0 ? 1 : 0
  bucket = aws_s3_bucket.website_bucket[0].id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

# Step 4: Apply a bucket policy to allow public read access
resource "aws_s3_bucket_policy" "website_policy" {
  bucket = coalesce(try(data.aws_s3_bucket.existing.id, null), try(aws_s3_bucket.website_bucket[0].id, null))

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

output "website_url" {
  value       = coalesce(data.aws_s3_bucket.existing.website_endpoint, try(aws_s3_bucket.website_bucket[0].website_endpoint, null))
  description = "The URL of the static website hosted on S3"
}
