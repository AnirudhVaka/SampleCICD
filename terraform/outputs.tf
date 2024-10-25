# Output the S3 website URL, depending on whether the bucket is pre-existing or newly created
output "website_url" {
  value       = coalesce(data.aws_s3_bucket.existing.website_endpoint, try(aws_s3_bucket_website_configuration.website[0].website_endpoint, null))
  description = "The URL of the static website hosted on S3"
}
