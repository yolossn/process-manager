package process

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/yolossn/process-manager/pkg/backoff"
	"github.com/yolossn/process-manager/pkg/config"
)

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
}

// New creates a new process
func New(conf config.Command) *Process {

	bo := backoff.NewStaticBackoff(time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	return &Process{config: conf, maxRetries: conf.MaxRetries, backoff: bo, ctx: ctx, cancelFunc: cancel}
}

// Run the process until successful or until max tries is reached
func (p *Process) Run(complete chan *Process) {
begin:
	// Recreate command on every run
	// because once the command is Run it cannot be reused
	var stdout, stderr bytes.Buffer
	cmd := NewCommand(p.ctx, p.config.Command, p.config.Args, p.config.EnvStrings(), &stdout, &stderr)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	p.output = &stdout
	p.err = &stderr
	p.command = cmd
	p.tryCount++

	// Run command
	err := p.command.Run()
	fmt.Println(err)
	if err != nil {

		// If maxRetries is not reached
		if p.tryCount < p.maxRetries {

			// Backoff
			backoffDuration := p.backoff.Duration()
			time.Sleep(backoffDuration)

			// retry only if the process is not stopped
			if !p.stopped {
				goto begin
			}
		}

		// If the process is not retried set the process as completed and signal the manager
		p.completed = true
		complete <- p
		return
	}

	// On success set the process as successful and completed.
	p.completed = true
	p.successful = true
	// signal the manager
	complete <- p
}

// Stop cancels the context of the process if it is not completed already.
// Incase the process didn't start yet the cancelled context will raise an error once it is started
func (p *Process) Stop() {

	// stop process
	p.stopped = true

	// If the process is already completed return
	if p.completed {
		return
	}

	// Refer: https://stackoverflow.com/questions/52346262/how-to-call-cancel-when-using-exec-commandcontext-in-a-goroutine
	// Use the cancel function to kill the process
	p.cancelFunc()

}

func (p *Process) IsSuccessful() bool {
	return p.successful
}

func (p *Process) Output() string {
	return p.output.String()
}

func (p *Process) Error() string {
	return p.err.String()
}

func (p *Process) MaxRetries() int {
	return p.maxRetries
}

func (p *Process) Retries() int {
	return p.tryCount - 1
}

func (p Process) String() string {
	return fmt.Sprintln(p.command.Path, p.command.Args)
}

func NewCommand(ctx context.Context, command string, args []string, env []string, stdout *bytes.Buffer, stderr *bytes.Buffer) *exec.Cmd {

	// TODO: Not sure if the args must be expanded
	// expandedArgs := []string{}
	// for _, arg := range args {
	// 	expandedArgs = append(expandedArgs, os.ExpandEnv(arg))
	// }

	cmd := exec.CommandContext(ctx, command, args...)
	// Add os environment for the process

	// TODO: Not sure if the process must have access to current os env
	// cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, env...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
