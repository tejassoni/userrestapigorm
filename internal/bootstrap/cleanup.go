package bootstrap

import "log/slog"

type cleanupTask struct {
	name string
	fn   func() error
}

type Cleanup struct {
	tasks []cleanupTask
}

func NewCleanup() *Cleanup {
	return &Cleanup{}
}

func (c *Cleanup) Add(name string, fn func() error) {
	c.tasks = append(c.tasks, cleanupTask{
		name: name,
		fn:   fn,
	})
}

func (c *Cleanup) Run(logger *slog.Logger) {

	for i := len(c.tasks) - 1; i >= 0; i-- {

		task := c.tasks[i]

		logger.Info("cleanup", "resource", task.name)

		if err := task.fn(); err != nil {

			logger.Error(
				"cleanup failed",
				"resource", task.name,
				"error", err,
			)
		}
	}
}
