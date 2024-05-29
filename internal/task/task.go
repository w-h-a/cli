package task

type Task interface {
	Validate() error
	Plan() error
	Apply() error
	Destroy() error
	Finalize() error
}
