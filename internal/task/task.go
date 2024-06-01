package task

type Task interface {
	Options() TaskOptions
	Validate() error
	Plan() error
	Apply() error
	Destroy() error
	Finalize() error
	String() string
}
