// @Title
// @Description
// @Author Jairo 2024/5/13 16:41
// @Email jairoguo@163.com

package http

import (
	"math/rand"
	"time"
)

var random *rand.Rand

// DefaultBackoff 线性间隔
func DefaultBackoff(_ int) time.Duration {
	return 1 * time.Second
}

// ExponentialBackoff 指数间隔 returns ever-increasing backoff by a power of 2
func ExponentialBackoff(i int) time.Duration {
	return time.Duration(1<<uint(i)) * time.Second
}

// ExponentialJitterBackoff 指数间隔+随机时间 returns ever-increasing backoff by a power of 2
// with +/- 0-33% to prevent synchronized requests.
func ExponentialJitterBackoff(i int) time.Duration {
	return jitter(int(1 << uint(i)))
}

// LinearBackoff 线性间隔  returns increasing durations, each a second longer than the last
func LinearBackoff(i int) time.Duration {
	return time.Duration(i) * time.Second
}

// LinearJitterBackoff 线性间隔+随机时间 returns increasing durations, each a second longer than the last
// with +/- 0-33% to prevent synchronized requests.
func LinearJitterBackoff(i int) time.Duration {
	return jitter(i)
}

// jitter keeps the +/- 0-33% logic in one place
func jitter(i int) time.Duration {
	ms := i * 1000

	maxJitter := ms / 3

	// ms ± rand

	ms += random.Intn(2*maxJitter) - maxJitter

	// a jitter of 0 messes up the time.Tick chan
	if ms <= 0 {
		ms = 1
	}

	return time.Duration(ms) * time.Millisecond
}
