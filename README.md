# AWS SES Prometheus Exporter

A Prometheus metrics exporter for AWS SES.

## Metrics

| Metric  | Labels | Description |
| ------  | ------ | ----------- |
| ses\_max24hoursend | aws_region | The maximum number of emails allowed to be sent in a rolling 24 hours. |
| ses\_maxsendrate | aws_region | The maximum rate of emails allowed to be sent per second. |
| ses\_sentlast24hours | aws_region | The number of emails sent in the last 24 hours. |
| ses\_Bounces | aws_region | The number of emails of emails that have bounced. |
| ses\_Complaints | aws_region | Number of unwanted emails that were rejected by recipients. |
| ses\_DeliveryAttempts | aws_region | Number of emails that have been sent. |
| ses\_Rejects | aws_region | Number of emails rejected by Amazon SES. |

For more information see the [AWS SES Documentation](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SESDeveloperGuide/sqs-message-attributes.html)

## Configuration

Credentials to AWS are provided in the following order:

- Environment variables (AWS\_ACCESS\_KEY\_ID and AWS\_SECRET\_ACCESS\_KEY)
- Shared credentials file (~/.aws/credentials)
- IAM role for Amazon EC2

For more information see the [AWS SDK Documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

## Environment Variables
| Variable      | Default Value | Description                                                  |
|---------------|:---------|:-------------------------------------------------------------|
| PORT          | 9435     | The port for metrics server                                  |
| ENDPOINT      | metrics  | The metrics endpoint                                         |



## Running

```docker run -d -p 9435:9435 bruceleo1969/ses-exporter```

You can provide the AWS credentials as environment variables depending upon your security rules configured in AWS;

```docker run -d -p 9435:9435 -e AWS_ACCESS_KEY_ID=<access_key> -e AWS_SECRET_ACCESS_KEY=<secret_key> -e AWS_REGION=<region>  bruceleo1969/ses-exporter```

