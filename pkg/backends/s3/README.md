# S3 Backend

The S3 backend enables `confd` to pull an YAML/JSON file from AWS S3 containing all keys.

## Configuration

The S3 backend utilizes the AWS SDK which utilizes the same options required by
the AWS CLI. The backend minimally requires setting the following:

-   `AWS_ACCESS_KEY_ID`
-   `AWS_SECRET_ACCESS_KEY`
-   `AWS_DEFAULT_REGION` and/or `AWS_REGION`

### Environment Variables

Environment variables can be used to provide the required configurations to
`confd`. They will override configurations set in the config and credentials
files.

```
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=us-east-2
```

### Config and Credentials Files

AWS credentials and configuration can be stored in the standard AWS CLI config
files. These may be set up manually or via `aws configure`

\~/.aws/credentials

```
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

\~/.aws/config

```
[default]
region=us-east-2
```
