// @Title
// @Description
// @Author Jairo 2024/5/13 15:14
// @Email jairoguo@163.com

package http

import "time"

type Config struct {
	EnableRetry     bool
	Retries         int
	Timeout         time.Duration
	backoffStrategy BackoffStrategy
}

type Option func(*Config)

func WithEnableRetry(enableRetry bool) Option {
	return func(o *Config) {
		o.EnableRetry = enableRetry
	}
}

func WithRetries(retries int) Option {
	return func(o *Config) {
		o.Retries = retries
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Config) {
		o.Timeout = timeout
	}
}

func WithBackoffStrategy(backoffStrategy BackoffStrategy) Option {
	return func(o *Config) {
		o.backoffStrategy = backoffStrategy
	}
}

func (c *Config) BuildHttpOption(options ...Option) {
	o := &Config{}
	for _, option := range options {
		option(o)
	}
	c.EnableRetry = o.EnableRetry
	c.Retries = o.Retries
	c.Timeout = o.Timeout
	if o.backoffStrategy != nil {
		c.backoffStrategy = o.backoffStrategy
	} else {
		c.backoffStrategy = DefaultBackoff

	}

}
