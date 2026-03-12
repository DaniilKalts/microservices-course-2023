package diagnostic

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type PingChecker struct {
	name   string
	pinger Pinger
}

func NewPingChecker(name string, pinger Pinger) *PingChecker {
	return &PingChecker{name: name, pinger: pinger}
}

func (c *PingChecker) Name() string {
	return c.name
}

func (c *PingChecker) Check(ctx context.Context) error {
	return c.pinger.Ping(ctx)
}
