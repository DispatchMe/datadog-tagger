package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var instanceTags = make(map[string][]string)

func main() {

	var apiKey string
	var appKey string

	flag.StringVar(&apiKey, "apiKey", "", "DataDog API key")
	flag.StringVar(&appKey, "appKey", "", "DataDog Application Key")

	flag.Parse()

	client := ec2.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	output, err := client.DescribeInstances(nil)

	if err != nil {
		log.Fatal(err)
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
				instanceTags[*instance.InstanceID] = tags
			}
		}
	}

	for hostId, tags := range instanceTags {
		data := map[string]interface{}{
			"tags": tags,
		}

		log.Printf("Updating host %s with tags %s\n", hostId, tags)

		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		err := encoder.Encode(data)
		if err != nil {
			log.Fatal(err)
		}

		request, err := http.NewRequest("PUT", "https://app.datadoghq.com/api/v1/tags/hosts/"+hostId+"?api_key="+apiKey+"&application_key="+appKey, buf)
		if err != nil {
			log.Fatal(err)
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Status code %d\n", response.StatusCode)

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))

	}

}
