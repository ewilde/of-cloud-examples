package function

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func saveToS3(result string, hostname string) (string, error) {

	s := createSession()

	bucket := os.Getenv("aws_bucket")
	if bucket == "" {
		log.Fatalf("could not retrieve aws_bucket from env var")
	}

	fileName := fmt.Sprintf("%s.txt", hostname)
	_, err := save(s, bucket, []byte(result), fileName)
	if err != nil {
		return "", fmt.Errorf("error saving file %s to bucket. %v", fileName, err)
	}

	return fileName, nil
}

func save(s *session.Session, bucket string, data []byte, key string) (*s3.PutObjectOutput, error) {

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	result, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(data),
		ContentLength:        aws.Int64(int64(len(data))),
		ContentType:          aws.String(http.DetectContentType(data)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return result, err
}

func createSession() *session.Session {
	accessKey, err := getAPISecret("aws-access-key-id")
	if err != nil {
		log.Printf("could not retrieve secret aws-access-key-id, using default session")
		return createSessionDefaultSession()
	}

	secret, err := getAPISecret("aws-secret-access-key")
	if err != nil {
		log.Printf("could not retrieve secret aws-secret-access-key, using default session")
		return createSessionDefaultSession()
	}

	region := os.Getenv("aws_region")
	if region == "" {
		log.Fatalf("could not retrieve aws_region from env var")
	}

	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(string(accessKey), string(secret), ""),
	})

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func createSessionDefaultSession() *session.Session {
	return session.Must(session.NewSession())
}
