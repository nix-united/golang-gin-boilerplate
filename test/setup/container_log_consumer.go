package setup

import (
	"github.com/testcontainers/testcontainers-go"
)

var _ testcontainers.LogConsumer = (*containerLogsConsumer)(nil)

// containerLogsConsumer collects logs from a container.
type containerLogsConsumer struct {
	logs []byte
}

func newContainerLogsConsumer() *containerLogsConsumer {
	return &containerLogsConsumer{}
}

// Accept records a log message from the container.
// It implements [testcontainers.LogConsumer] interface.
func (c *containerLogsConsumer) Accept(log testcontainers.Log) {
	c.logs = append(c.logs, log.Content...)
}

// Collect returns collected all logs.
func (c *containerLogsConsumer) Collect() []byte {
	return c.logs
}
