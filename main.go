package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

type Secrets struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	SecretAccessKey string `yaml:"SecretAccessKey"`
	Region          string `yaml:"Region"`
	BucketName      string `yaml:"BucketName"`
}

func main() {
	s, err := getKeys("./secret.yaml")
	if err != nil {
		panic(err)
	}
	c := s3.New(&aws.Config{}) // TODO: 型が合わない
	//Credentials:credentials.NewStaticCredentials(s.AccessKeyId, s.SecretAccessKey, ""),
	//	Region: s.Region,

	info, err := ioutil.ReadDir("./contents")
	if err != nil {
		panic(err)
	}
	for _, v := range info {
		if !v.IsDir() {
			fr, err := os.Open("./contents/" + v.Name())
			if err != nil {
				panic(err)
			}
			resp, err := c.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(s.BucketName),
				Key:    aws.String(v.Name()),
				Body:   fr,
			})
			log.Println(awsutil.StringValue(resp))
			fr.Close()
		}
	}
}

func getKeys(path string) (*Secrets, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := &Secrets{}
	r := bufio.NewReader(f)
	if err != nil {
		return nil, err
	}

	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(l, &s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}
