package terraform

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
	"github.com/w-h-a/cli/internal/task"
)

const (
	tfS3BackendTemplate = `terraform {
		backend "s3" {
		  bucket         = "{{.StateBucket}}"
		  dynamodb_table = "{{.LockTable}}"
		  key            = "{{.Key}}"
		  region         = "{{.Region}}"
		}
	  }
	  `

	tfS3RemoteStateTemplate = `data "terraform_remote_state" "{{.RemoteStateName}}" {
		backend = "s3"

		config = {
		  bucket         = "{{.StateBucket}}"
		  dynamodb_table = "{{.LockTable}}"
		  key            = "{{.Key}}"
		  region         = "{{.Region}}"
		}
	  }
	  `
)

type terraformExecutor struct {
	options task.TaskOptions
}

// Validate attempts to fetch terraform and run `terraform init` and `terraform validate`
func (t *terraformExecutor) Validate() error {
	if err := os.MkdirAll(t.options.Path, 0o777); err != nil {
		return err
	}

	u, err := url.Parse(t.options.Source)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http":
		fallthrough
	case "https":
		if err := t.executeGitClone(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("scheme %s is not supported", u.Scheme)
	}

	if err := t.writeStateFiles(); err != nil {
		return err
	}

	if err := t.executeTerraform(context.Background(), "init"); err != nil {
		return err
	}

	if err := t.executeTerraform(context.Background(), "validate"); err != nil {
		return err
	}

	return nil
}

func (t *terraformExecutor) Plan() error {
	return t.executeTerraform(context.Background(), "plan")
}

func (t *terraformExecutor) Apply() error {
	return nil
}

func (t *terraformExecutor) Destroy() error {
	return nil
}

func (t *terraformExecutor) Finalize() error {
	// return os.RemoveAll(t.options.Path)
	return nil
}

func (t *terraformExecutor) executeTerraform(ctx context.Context, args ...string) error {
	// set up terraform command
	tf := exec.CommandContext(ctx, "terraform", args...)
	tf.Dir = t.options.Path
	tf.Env = os.Environ()

	for k, v := range t.options.EnvVars {
		tf.Env = append(tf.Env, fmt.Sprintf("%s=%s", k, v))
	}

	tfVars := map[string]string{}
	if v, ok := t.options.Context.Value("tf_vars_key").(map[string]string); ok {
		tfVars = v
	}
	for k, v := range tfVars {
		tf.Env = append(tf.Env, fmt.Sprintf("TF_VAR_%s=%s", k, v))
	}

	stdout, err := tf.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdoutpip failed: %v", err)
	}

	stderr, err := tf.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderrpipe failed: %v", err)
	}

	// wait so we don't truncate output from terraform
	ioWait := make(chan struct{})
	defer func() {
		// wait for both routines (see below) to finish
		// so we capture everything
		<-ioWait
		<-ioWait
	}()

	for _, ioPair := range []struct {
		in  io.ReadCloser
		out *os.File
	}{
		{in: stdout, out: os.Stdout},
		{in: stderr, out: os.Stderr},
	} {
		go func(name string, in io.ReadCloser, out *os.File, done chan<- struct{}) {
			defer func() {
				done <- struct{}{}
			}()

			defer in.Close()

			reader := bufio.NewReader(in)

			for {
				s, err := reader.ReadString('\n')
				if err == nil || err == io.EOF {
					if len(strings.TrimSpace(s)) != 0 {
						fmt.Fprintf(out, "[%s] %s", name, s)
					}
					if err == io.EOF {
						return
					}
				} else {
					fmt.Fprintf(out, "[%s] error: %s\n", name, err.Error())
					return
				}
			}

		}(t.options.Name, ioPair.in, ioPair.out, ioWait)
	}

	if err := tf.Start(); err != nil {
		return fmt.Errorf("failed to execute terraform: %v", err)
	}

	return tf.Wait()
}

func (t *terraformExecutor) executeGitClone() error {
	fmt.Fprintf(os.Stdout, "cloning repo %s\n", t.options.Source)

	if _, err := git.PlainClone(
		t.options.Path,
		false,
		&git.CloneOptions{
			URL:      t.options.Source,
			Progress: os.Stdout,
		},
	); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "successfully cloned repo %s\n", t.options.Source)

	return nil
}

func (t *terraformExecutor) writeStateFiles() error {
	stateStore := viper.GetString("state-store")

	switch stateStore {
	case "aws":
		if err := t.writeBackendFileAWS(); err != nil {
			return err
		}

		if err := t.writeRemoteStatesFileAWS(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s is not a supported remote state backend", stateStore)
	}

	fmt.Fprintf(os.Stdout, "successfully wrote state files to %s\n", t.options.Path)

	return nil
}

func (t *terraformExecutor) writeBackendFileAWS() error {
	// create the file
	f, err := os.OpenFile(filepath.Join(t.options.Path, "backend-config.tf"), os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	// defer its closing
	defer f.Close()

	// create the write template
	backend := template.Must(template.New(t.options.Name + "backend").Parse(tfS3BackendTemplate))

	// write to the file
	if err := backend.Execute(f, struct {
		StateBucket string
		LockTable   string
		Key         string
		Region      string
	}{
		StateBucket: viper.GetString("aws-s3-bucket"),
		LockTable:   viper.GetString("aws-dynamodb-table"),
		Key:         t.options.Name,
		Region:      viper.GetString("aws-region"),
	}); err != nil {
		return err
	}

	return nil
}

func (t *terraformExecutor) writeRemoteStatesFileAWS() error {
	// create the file
	f, err := os.OpenFile(filepath.Join(t.options.Path, "remote-state-data-sources.tf"), os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}

	// defer its closing
	defer f.Close()

	// create the write template
	remote := template.Must(template.New(t.options.Name + "remote").Parse(tfS3RemoteStateTemplate))

	// get the desired remote states
	states := map[string]string{}
	if rs, ok := t.options.Context.Value("tf_remote_states_key").(map[string]string); ok {
		states = rs
	}

	// write to the file for each remote state requested
	for k, v := range states {
		if err := remote.Execute(f, struct {
			RemoteStateName string
			StateBucket     string
			LockTable       string
			Key             string
			Region          string
		}{
			RemoteStateName: k,
			StateBucket:     viper.GetString("aws-s3-bucket"),
			LockTable:       viper.GetString("aws-dynamodb-table"),
			Key:             v,
			Region:          viper.GetString("aws-region"),
		}); err != nil {
			return err
		}
	}

	return nil
}

func NewTask(opts ...task.TaskOption) task.Task {
	options := task.NewTaskOptions(opts...)

	t := &terraformExecutor{
		options: options,
	}

	return t
}
