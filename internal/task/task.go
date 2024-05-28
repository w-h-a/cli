package task

type Task interface {
	Validate() error
	Plan() error
	Apply() error
	Finalize() error
	Destroy() error
}
