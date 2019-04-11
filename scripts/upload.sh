#!/usr/bin/env bash
set -e
set -u
set -o pipefail

echo 'Uploading builds to AWS S3...'
aws s3 sync ${PWD}/build s3://lvnacapital/algorand --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --region us-west-2

# echo "Creating invalidation for AWS Cloudfront"
# aws configure set preview.cloudfront true
# aws cloudfront create-invalidation --distribution-id E3B5Z3LYG19QSL --paths /algorand
