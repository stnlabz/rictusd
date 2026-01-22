package vision

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type Watcher struct {
	Process *exec.Cmd
	Active  bool
}

// Start opens the pipe to the Hailo-8 inference sidecar
func Start(modelPath string) (*Watcher, error) {
	// We call the python script that utilizes the Hailo SDK
	cmd := exec.Command("python3", "vision/hailo_inference.py", "--model", modelPath)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Watcher{Process: cmd, Active: true}, nil
}

// Monitor triggers the Go engine whenever the NPU sees a verified threat
func (w *Watcher) Monitor(onDetection func(string)) {
	scanner := bufio.NewScanner(w.Process.Stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "DETECTED:") {
			onDetection(line)
		}
	}
}
