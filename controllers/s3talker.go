package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Bucket - information per bucket
type Bucket struct {
	BucketName         string  `json:"bucketName"`
	BucketSize         float64 `json:"bucketSize"`
	BucketObjectNumber float64 `json:"bucketObjectNumber"`
}

// Buckets - list of Bucket objects
type Buckets []Bucket

// S3Summary - one JSON struct to rule them all
type S3Summary struct {
	S3Name         string  `json:"s3name"`
	S3Status       bool    `json:"s3Status"`
	S3Size         float64 `json:"s3Size"`
	S3ObjectNumber float64 `json:"s3ObjectNumber"`
	S3Buckets      Buckets `json:"s3Bucket"`
}

// S3Conn struct - keeps information about remote S3
type S3Conn struct {
	S3ConnName                      string `json:"s3_conn_name"`
	S3ConnQuota                     int64  `json:"s3_conn_quota" required:"false"`
	S3ConnAccessKey                 string `json:"s3_conn_access_key" required:"true"`
	S3ConnSecretKey                 string `json:"s3_conn_secret_key" required:"true"`
	S3ConnEndpoint                  string `json:"s3_conn_endpoint" default:"false"`
	S3ConnRegion                    string `json:"s3_conn_region" default:"default"`
	S3ConnDisableSsl                bool   `json:"s3_conn_disable_ssl" required:"true"`
	S3ConnForcePathStyle            bool   `json:"s3_conn_force_path_style" default:"true"`
	S3ConnDisableEdnpointHostPrefix bool   `json:"s3_conn_disable_endpoint_host_prefix" default:"true"`
}

// S3UsageInfo - gets s3 connection details return s3Summary
func S3UsageInfo(s3Conn S3Conn) (S3Summary, error) {
	s3All := S3Summary{}
	s3All.S3Name = s3Conn.S3ConnName

	s3Config := &aws.Config{
		Credentials:               credentials.NewStaticCredentials(s3Conn.S3ConnAccessKey, s3Conn.S3ConnSecretKey, ""),
		Endpoint:                  aws.String(s3Conn.S3ConnEndpoint),
		DisableSSL:                &s3Conn.S3ConnDisableSsl,
		DisableEndpointHostPrefix: &s3Conn.S3ConnDisableEdnpointHostPrefix,
		S3ForcePathStyle:          &s3Conn.S3ConnForcePathStyle,
		// Region aws.String("us-east-1"), // This is counter intuitive but it will fail with a non-AWS region name
		//Region:                    aws.String("default"),
		//Region: aws.String(&s3Conn.S3ConnRegion),
		Region: aws.String(s3Conn.S3ConnRegion),
	}
	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		//log.Fatal("Connection to S3 endpoint failed (log.Fatal):", err)
		fmt.Println("Connection to S3 endpoint failed :", err)
		s3All.S3Status = false
		return s3All, errors.New("s3 endpoint: unable to connect")
	}
	s3All.S3Status = true

	// Processing calculation per Bucket
	for _, b := range result.Buckets {
		size, number := countBucketSize(aws.StringValue(b.Name), s3Client)
		bucketX := Bucket{BucketName: *b.Name, BucketObjectNumber: number, BucketSize: size}
		s3All.S3Buckets = append(s3All.S3Buckets, bucketX)
		s3All.S3Size += size
		s3All.S3ObjectNumber += number
	}

	// Save s3All
	byteArray, err := json.MarshalIndent(s3All, "", "    ")
	err = ioutil.WriteFile("s3Information.json", byteArray, 0777)
	if err != nil {
		fmt.Println("Didn't manage to dump S3 details as a json file, error :", err)
	}
	return s3All, nil
}

func countBucketSize(bucketName string, s3Client *s3.S3) (float64, float64) {
	i := 0
	var bucketUsage float64
	var bucketObjects float64
	bucketUsage, bucketObjects = 0, 0

	err := s3Client.ListObjectsPages(&s3.ListObjectsInput{Bucket: aws.String(bucketName)},
		func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
			i++
			for _, obj := range p.Contents {
				bucketUsage += float64(*obj.Size)
				bucketObjects++
			}
			return true
		})

	if err != nil {
		log.Fatal(err)
		return 0, 0
	}
	return bucketUsage, bucketObjects
}
