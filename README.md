This is a tool we use to tag our instances in DataDog with their application name and environment. It currently only works for Elastic Beanstalk applications, and assumes that the environment name in Elastic Beanstalk follows this format: `<application name>-<environment>`. the `<application name>` can have hyphens, but the last hyphen is assumed to separate the application name from the environment.

For example:

1. `dispatch-api-dev` -> `application:"dispatch-api", environment:"dev"`
2. `email-sender-prod` -> `application:"email-sender", environment:"prod"`


# Installation
## Dependencies
1. Go 1.5
2. Environment variable: `GO15VENDOREXPERIMENT="1"`

## Instructions

```bash
$ cd $GOPATH/src/github.com/DispatchMe
$ git clone git@github.com:DispatchMe/datadog-tagger
$ cd datadog-tagger
$ go build
$ go install
```

# Usage

Make sure you have `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` exported in your terminal. You can also automatically the below flags with `DATADOG_API_KEY`, `DATADOG_APP_KEY`, and `AWS_REGION` environment variables, respectively.

```bash
$ datadog-tagger --apiKey="<datadog API key>" --appKey="<datadog app key>" --awsRegion="<aws region for ec2 instances>"
```
