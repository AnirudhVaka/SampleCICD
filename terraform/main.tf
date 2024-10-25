# Step 1: Data source to check if the bucket exists
data "aws_s3_bucket" "existing" {
  bucket = "samplewebsitebucket"  # Replace with your bucket name
}

# Step 2: Conditionally create the bucket if it doesn’t exist
resource "aws_s3_bucket" "website_bucket" {
  count  = data.aws_s3_bucket.existing.id != "" ? 0 : 1
  bucket = "samplewebsitebucket"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

# Step 3: Apply a bucket policy to allow public read access (skip if it already exists)
resource "aws_s3_bucket_policy" "website_policy" {
  bucket = coalesce(try(data.aws_s3_bucket.existing.id, null), aws_s3_bucket.website_bucket[0].id)

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
  
  # Only apply the policy if the bucket was created by Terraform or if it lacks this policy
  count = data.aws_s3_bucket.existing.id != "" ? 1 : 0
}

output "website_url" {
  value       = coalesce(data.aws_s3_bucket.existing.website_endpoint, aws_s3_bucket.website_bucket[0].website_endpoint)
  description = "The URL of the static website hosted on S3"
}
