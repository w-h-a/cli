package state

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/task"
)

type stateChecker struct {
	options task.TaskOptions
}

func (s *stateChecker) Validate() error {
	stateStore := viper.GetString("state-store")

	if err := s.validateConfig(); err != nil {
		fmt.Fprintf(os.Stdout, "remote state backend in %s is invalid\n", stateStore)

		return err
	}

	fmt.Fprintf(os.Stdout, "remote state backend in %s is valid\n", stateStore)

	return nil
}

func (s *stateChecker) Plan() error {
	return nil
}

func (s *stateChecker) Apply() error {
	return nil
}

func (s *stateChecker) Destroy() error {
	return nil
}

func (s *stateChecker) Finalize() error {
	return nil
}

func (s *stateChecker) validateConfig() error {
	stateStore := viper.GetString("state-store")

	switch stateStore {
	case "aws":
		return s.validateAWS()
	default:
		return fmt.Errorf("remote state backend in %s is not supported", stateStore)
	}
}

func (s *stateChecker) validateAWS() error {
	config := &aws.Config{
		Region: aws.String(viper.GetString("aws-region")),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return fmt.Errorf("failed to generate an aws session: %v", err)
	}

	s3Client := s3.New(sess)

	bucket := viper.GetString("aws-s3-bucket")

	if _, err := s3Client.PutObject(
		&s3.PutObjectInput{
			Key:    aws.String(s.options.Name),
			Bucket: aws.String(bucket),
			Body:   strings.NewReader(s.options.Name),
		},
	); err != nil {
		return fmt.Errorf("failed to put an object into the remote state backend: %v", err)
	}

	read, err := s3Client.GetObject(
		&s3.GetObjectInput{
			Key:    aws.String(s.options.Name),
			Bucket: aws.String(bucket),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to read back an object from the remote state backend: %v", err)
	}

	defer read.Body.Close()

	body, err := io.ReadAll(read.Body)
	if err != nil {
		return fmt.Errorf("failed to read the body of an object from the remote state backend: %v", err)
	}

	if string(body) != s.options.Name {
		return fmt.Errorf("read back an incorrect value from the remote state backend: want %s, got %s", s.options.Name, string(body))
	}

	if _, err := s3Client.DeleteObject(
		&s3.DeleteObjectInput{
			Key:    aws.String(s.options.Name),
			Bucket: aws.String(bucket),
		},
	); err != nil {
		return fmt.Errorf("failed to delete object from the remote state backend: %v", err)
	}

	return nil
}

func NewTask(opts ...task.TaskOption) task.Task {
	options := task.NewTaskOptions(opts...)

	s := &stateChecker{
		options: options,
	}

	return s
}
