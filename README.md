This loads all of our instances from EC2 via the AWS API and updates their tags in DataDog to reflect the **environment** and **application**. For Elastic Beanstalk apps only.

# Installation
```bash
$ cd $GOPATH/src/github.com/DispatchMe
$ git clone git@github.com:DispatchMe/datadog-tagger
$ cd datadog-tagger
$ go-getter Gofile
$ go build
$ go install
```

# Usage

Make sure you have `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` exported in your terminal. Ask Jason for the below keys.

```bash
$ datadog-tagger --apiKey="<datadog API key>" --appKey="<datadog app key>"
```
