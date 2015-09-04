package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var instanceTags = make(map[string][]string)

func fail(err error) {
	println(err.Error())
	os.Exit(1)
}
func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "apiKey",
			Usage:  "API key",
			EnvVar: "DATADOG_API_KEY",
		},
		cli.StringFlag{
			Name:   "appKey",
			Usage:  "App key",
			EnvVar: "DATADOG_APP_KEY",
		},
		cli.StringFlag{
			Name:   "awsRegion",
			Usage:  "Region for EC2 instances",
			EnvVar: "AWS_REGION",
		},
	}

	app.Action = func(c *cli.Context) {
		client := ec2.New(&aws.Config{
			Region: aws.String(c.String("awsRegion")),
		})

		output, err := client.DescribeInstances(nil)

		if err != nil {
			fail(err)
		}

		for _, r := range output.Reservations {
			for _, instance := range r.Instances {
				tags := make([]string, 0)

				isBeanstalk := false
				for _, tag := range instance.Tags {
					switch *tag.Key {
					// Only use this one for now
					case "elasticbeanstalk:environment-name":
						isBeanstalk = true

						// By convention, the stage is the last segment when split by "-"
						spl := strings.Split(*tag.Value, "-")
						environment := spl[len(spl)-1]
						app := strings.Join(spl[0:len(spl)-1], "-")
						tags = append(tags, "environment:"+environment)
						tags = append(tags, "application:"+app)
					}
				}

				if isBeanstalk {
					instanceTags[*instance.InstanceId] = tags
				}
			}
		}

		for hostId, tags := range instanceTags {
			data := map[string]interface{}{
				"tags": tags,
			}

			fmt.Printf("Updating host %s with tags %s\n", hostId, tags)

			buf := &bytes.Buffer{}
			encoder := json.NewEncoder(buf)
			err := encoder.Encode(data)
			if err != nil {
				fail(err)
			}

			request, err := http.NewRequest("PUT", "https://app.datadoghq.com/api/v1/tags/hosts/"+hostId+"?api_key="+c.String("apiKey")+"&application_key="+c.String("appKey"), buf)
			if err != nil {
				fail(err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				fail(err)
			}

			fmt.Printf("Status code %d\n", response.StatusCode)

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fail(err)
			}

			fmt.Println(string(body))

		}
	}

	app.Run(os.Args)

}
