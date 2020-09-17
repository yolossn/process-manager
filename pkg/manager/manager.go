package manager

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/yolossn/process-manager/pkg/config"
	"github.com/yolossn/process-manager/pkg/process"
)

// Manager runs all the processes and allows to stop them.
type Manager interface {
	Run() chan struct{} // Runs all the process and returns a chan to signal completion of all processes
	Stop()              // Stops all the process and returns
	SuccessCount() int
	FailCount() int
	Status() map[string]status
}

type manager struct {
	total        int
	completed    int
	successful   int
	failed       int
	processes    []*process.Process
	completeChan chan *process.Process
	mu           sync.RWMutex
}

// New creates a new manager to manage the processes.
func New(commands []config.Command) Manager {

	var processes []*process.Process
	for _, command := range commands {
		proc := process.New(command)
		processes = append(processes, proc)
	}

	return &manager{
		total:     len(processes),
		processes: processes,
	}
}

// Run all the process and returns a chan to signal completion of all processes.
func (m *manager) Run() chan struct{} {

	// channel for processes to signal completion.
	completeChan := make(chan *process.Process, 1)

	// Handle SIGKILL and SIGINT
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, os.Kill)
	go func(chan os.Signal) {
		<-kill
		// debug
		fmt.Println("Killing the processes")
		m.Stop()
		os.Exit(1)
	}(kill)

	// Run all the process as goroutines
	for _, process := range m.processes {
		go process.Run(completeChan)
	}

	stop := make(chan struct{})
	go func() {
		for {
			select {
			case proc := <-completeChan:

				m.mu.Lock()
				m.completed++
				if proc.IsSuccessful() {
					m.successful++
				} else {
					m.failed++
				}
				m.mu.Unlock()

				m.mu.RLock()
				completed := m.completed
				total := m.total
				m.mu.RUnlock()

				if completed == total {
					stop <- struct{}{}
					break
				}
			}
		}
	}()
	return stop
}

// Stop all the processes and return.
func (m *manager) Stop() {
	// Stop all processes
	for _, process := range m.processes {
		go process.Stop()
	}

	return
}

// SuccessCount returns the number of process which completed succesfully.
func (m *manager) SuccessCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.successful
}

// FailCount returns the number of process which failed.
func (m *manager) FailCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.failed
}

type status struct {
	Output       string
	Err          string
	IsSuccessful bool
}

// Status returns the status of each process as a map
func (m *manager) Status() map[string]status {
	outputs := make(map[string]status)

	for _, process := range m.processes {
		processStatus := status{Output: process.Output(), Err: process.Error(), IsSuccessful: process.IsSuccessful()}
		outputs[process.String()] = processStatus
	}
	return outputs
}
