package terraform

import (
	"fmt"
	"net/url"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/w-h-a/cli/internal/task"
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

	return nil
}

func (t *terraformExecutor) Plan() error {
	return nil
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

func NewTask(opts ...task.TaskOption) task.Task {
	options := task.NewTaskOptions(opts...)

	t := &terraformExecutor{
		options: options,
	}

	return t
}
