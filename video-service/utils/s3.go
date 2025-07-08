// package utils

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// )

// // func GeneratePresignedURL(fileName string) (string, error) {
// // 	fmt.Println("REGION:", os.Getenv("AWS_REGION"))
// // 	fmt.Println("BUCKET:", os.Getenv("S3_BUCKET"))

// // 	sess, err := session.NewSession(&aws.Config{
// // 		Region: aws.String(os.Getenv("AWS_REGION")),
// // 		Credentials: credentials.NewStaticCredentials(
// // 			os.Getenv("AWS_ACCESS_KEY_ID"),
// // 			os.Getenv("AWS_SECRET_ACCESS_KEY"),
// // 			"",
// // 		),
// // 		// Optional but safe for local dev
// // 		S3ForcePathStyle: aws.Bool(true),
// // 	})
// // 	if err != nil {
// // 		return "", fmt.Errorf("Session creation failed: %w", err)
// // 	}

// // 	svc := s3.New(sess)

// // 	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
// // 		Bucket: aws.String(os.Getenv("S3_BUCKET")),
// // 		Key:    aws.String(fileName),
// // 		// ACL:    aws.String("public-read"),
// // 	})

// // 	urlStr, err := req.Presign(15 * time.Minute)
// // 	if err != nil {
// // 		return "", fmt.Errorf("Failed to sign request: %w", err)
// // 	}

// //		return urlStr, nil
// //	}
// func GeneratePresignedPutURL(userID, fileName string) (string, string, error) {
// 	timestamp := time.Now().Unix()
// 	uniqueKey := fmt.Sprintf("videos/%s/%d_%s", userID, timestamp, fileName)

// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(os.Getenv("AWS_REGION")),
// 		Credentials: credentials.NewStaticCredentials(
// 			os.Getenv("AWS_ACCESS_KEY_ID"),
// 			os.Getenv("AWS_SECRET_ACCESS_KEY"),
// 			"",
// 		),
// 		S3ForcePathStyle: aws.Bool(true),
// 	})
// 	if err != nil {
// 		return "", "", fmt.Errorf("Session creation failed: %w", err)
// 	}

// 	svc := s3.New(sess)

// 	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
// 		Bucket: aws.String(os.Getenv("S3_BUCKET")),
// 		Key:    aws.String(uniqueKey),
// 	})

// 	urlStr, err := req.Presign(15 * time.Minute)
// 	if err != nil {
// 		return "", "", fmt.Errorf("Failed to sign request: %w", err)
// 	}

// 	return urlStr, uniqueKey, nil
// }

// func GeneratePresignedGetURL(fileName string) (string, error) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(os.Getenv("AWS_REGION")),
// 		Credentials: credentials.NewStaticCredentials(
// 			os.Getenv("AWS_ACCESS_KEY_ID"),
// 			os.Getenv("AWS_SECRET_ACCESS_KEY"),
// 			"",
// 		),
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	svc := s3.New(sess)

// 	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
// 		Bucket: aws.String(os.Getenv("S3_BUCKET")),
// 		Key:    aws.String(fileName),
// 	})

// 	return req.Presign(15 * time.Minute)
// }

// video-service/utils/s3.go
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GeneratePresignedPutURL(userID, fileName string) (string, string, error) {
	timestamp := time.Now().Unix()
	uniqueKey := fmt.Sprintf("videos/%s/%d_%s", userID, timestamp, fileName)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return "", "", fmt.Errorf("Session creation failed: %w", err)
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(uniqueKey),
		ContentType: aws.String("video/mp4"),
		// Remove ACL to avoid signing complications
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("Failed to sign request: %w", err)
	}

	return urlStr, uniqueKey, nil
}

// Option 2: With ACL (if you need public access) - Alternative implementation
func GeneratePresignedPutURLWithACL(userID, fileName string) (string, string, error) {
	timestamp := time.Now().Unix()
	uniqueKey := fmt.Sprintf("videos/%s/%d_%s", userID, timestamp, fileName)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return "", "", fmt.Errorf("Session creation failed: %w", err)
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(uniqueKey),
		ContentType: aws.String("video/mp4"),
		ACL:         aws.String("public-read"),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("Failed to sign request: %w", err)
	}

	return urlStr, uniqueKey, nil
}

func GeneratePresignedGetURL(fileName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(fileName),
	})

	return req.Presign(15 * time.Minute)
}

func DeleteFromS3(objectKey string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return fmt.Errorf("Session creation failed: %w", err)
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return fmt.Errorf("Failed to delete object from S3: %w", err)
	}

	return nil
}
