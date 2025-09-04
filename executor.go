package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExecutorService handles execution of Base CLI commands
type ExecutorService struct {
	basePath string
	cmdPath  string
}

// NewExecutorService creates a new executor service
func NewExecutorService() *ExecutorService {
	// Try to find the base command in various locations
	basePath := findBasePath()
	cmdPath := findCmdPath()

	return &ExecutorService{
		basePath: basePath,
		cmdPath:  cmdPath,
	}
}

// findBasePath attempts to locate the base executable
func findBasePath() string {
	// Check if base is in PATH
	if path, err := exec.LookPath("base"); err == nil {
		return path
	}

	// Check common installation locations
	candidates := []string{
		"/usr/local/bin/base",
		"/usr/bin/base",
		"./cmd/base",
		"../cmd/base",
		"./base",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Default to assuming base is in PATH
	return "base"
}

// findCmdPath attempts to locate the cmd directory for the Base CLI
func findCmdPath() string {
	candidates := []string{
		"../cmd",
		"./cmd",
		"../../cmd",
	}

	for _, candidate := range candidates {
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			return candidate
		}
	}

	return ""
}

// ExecuteGenerate executes the base generate command
func (e *ExecutorService) ExecuteGenerate(name string, fields []string) (string, error) {
	args := []string{"generate", name}
	args = append(args, fields...)

	return e.executeBaseCommand(args...)
}

// ExecuteStart executes the base start command
func (e *ExecutorService) ExecuteStart(reload, docs bool) (string, error) {
	args := []string{"start"}

	if reload {
		args = append(args, "-r")
	}
	if docs {
		args = append(args, "-d")
	}

	return e.executeBaseCommand(args...)
}

// ExecuteNew executes the base new command
func (e *ExecutorService) ExecuteNew(name, path string) (string, error) {
	args := []string{"new", name}

	if path != "" {
		args = append(args, "--path", path)
	}

	return e.executeBaseCommand(args...)
}

// ExecuteDestroy executes the base destroy command
func (e *ExecutorService) ExecuteDestroy(name string) (string, error) {
	return e.executeBaseCommand("destroy", name)
}

// ExecuteDocs executes the base docs command
func (e *ExecutorService) ExecuteDocs() (string, error) {
	return e.executeBaseCommand("docs")
}

// executeBaseCommand executes a base command with the given arguments
func (e *ExecutorService) executeBaseCommand(args ...string) (string, error) {
	// Try using the base CLI if available
	if e.basePath != "" {
		cmd := exec.Command(e.basePath, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("base command failed: %v\nOutput: %s", err, string(output))
		}
		return string(output), nil
	}

	// Fallback to direct Go execution if cmd path is available
	if e.cmdPath != "" {
		return e.executeGoDirect(args...)
	}

	return "", fmt.Errorf("base CLI not found - please install Base CLI or run from Base project directory")
}

// executeGoDirect executes Base CLI commands directly using go run
func (e *ExecutorService) executeGoDirect(args ...string) (string, error) {
	mainGo := filepath.Join(e.cmdPath, "main.go")

	// Check if main.go exists
	if _, err := os.Stat(mainGo); os.IsNotExist(err) {
		return "", fmt.Errorf("base CLI main.go not found at %s", mainGo)
	}

	// Prepare go run command
	goArgs := []string{"run", mainGo}
	goArgs = append(goArgs, args...)

	cmd := exec.Command("go", goArgs...)
	cmd.Dir = e.cmdPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go run failed: %v\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// IsBaseAvailable checks if Base CLI is available
func (e *ExecutorService) IsBaseAvailable() bool {
	return e.basePath != "" || e.cmdPath != ""
}

// GetStatus returns the status of the executor service
func (e *ExecutorService) GetStatus() string {
	var status []string

	if e.basePath != "" {
		status = append(status, fmt.Sprintf("Base CLI: %s", e.basePath))
	}

	if e.cmdPath != "" {
		status = append(status, fmt.Sprintf("Cmd Path: %s", e.cmdPath))
	}

	if len(status) == 0 {
		return "Base CLI not found"
	}

	return strings.Join(status, "\n")
}
