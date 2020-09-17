package process

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/yolossn/process-manager/pkg/backoff"
	"github.com/yolossn/process-manager/pkg/config"
)

// Process is a wrapper of exec.cmd which takes care of retries based on the backoff strategy.
type Process struct {
	config  config.Command
	command *exec.Cmd

	// max and current retries
	maxRetries int
	tryCount   int

	// output and error buffers
	output *bytes.Buffer
	err    *bytes.Buffer

	// states
	successful bool
	completed  bool
	stopped    bool

	// backoff
	backoff backoff.Backoff

	ctx        context.Context
	cancelFunc func()

	mu sync.RWMutex
}

// New creates a new process.
func New(conf config.Command) *Process {

	bo := backoff.NewStaticBackoff(time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	return &Process{config: conf, maxRetries: conf.MaxRetries, backoff: bo, ctx: ctx, cancelFunc: cancel}
}

// Run the process until successful or until max tries is reached.
func (p *Process) Run(complete chan *Process) {
	for {
		// Recreate command on every run
		// because once the command is Run it cannot be reused.
		var stdout, stderr bytes.Buffer
		cmd := newCommand(p.ctx, p.config.Command, p.config.Args, p.config.EnvStrings(), &stdout, &stderr)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		p.mu.Lock()
		p.output = &stdout
		p.err = &stderr
		p.command = cmd
		p.tryCount++
		p.mu.Unlock()

		// Run command.
		err := p.command.Run()
		// debug
		// fmt.Println(p.String(), "error:", err, "tryCount:", p.tryCount, "maxRetry:", p.maxRetries)
		if err != nil {

			// If maxRetries is not reached.
			if p.tryCount < p.maxRetries {

				// Backoff.
				backoffDuration := p.backoff.Duration()
				time.Sleep(backoffDuration)

				// retry only if the process is not stopped.
				p.mu.RLock()
				stop := p.stopped
				p.mu.RUnlock()
				if !stop {
					continue
				}
			}

			// If the process is not retried set the process as completed and signal the manager.
			p.mu.Lock()
			p.completed = true
			p.mu.Unlock()
			complete <- p
			break
		}

		// On success set the process as successful and completed.
		p.mu.Lock()
		p.completed = true
		p.successful = true
		p.mu.Unlock()
		// signal the manager.
		complete <- p
		break
	}

	return
}

// Stop cancels the context of the process if it is not completed already.
// Incase the process didn't start yet the cancelled context will raise an error once it is started.
func (p *Process) Stop() {

	p.mu.Lock()
	defer p.mu.Unlock()
	// stop process.
	p.stopped = true

	// If the process is already completed return.
	if p.completed {
		return
	}

	// Refer: https://stackoverflow.com/questions/52346262/how-to-call-cancel-when-using-exec-commandcontext-in-a-goroutine
	// Use the cancel function to kill the process.
	p.cancelFunc()
}

// IsSuccessful returns the success state.
func (p *Process) IsSuccessful() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.successful
}

// Output returns the stdout of the command.
func (p *Process) Output() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.output.String()
}

// Error returns the stderr of the command.
func (p *Process) Error() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.err.String()
}

// MaxRetries returns the maximum retries of the process.
func (p *Process) MaxRetries() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.maxRetries
}

// Retries returns the number of times the process retried to run the command.
func (p *Process) Retries() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.tryCount - 1
}

// String returns the string representation of process.
func (p Process) String() string {
	return fmt.Sprint(p.command.Args)
}

func newCommand(ctx context.Context, command string, args []string, env []string, stdout *bytes.Buffer, stderr *bytes.Buffer) *exec.Cmd {

	// TODO: Not sure if expansion of args is a requirement
	// expandedArgs := []string{}
	// for _, arg := range args {
	// 	expandedArgs = append(expandedArgs, os.ExpandEnv(arg))
	// }

	cmd := exec.CommandContext(ctx, command, args...)

	// TODO: Not sure if the process must have access to current os env
	// Add os environment for the process
	// cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, env...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd
}
