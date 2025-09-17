package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CutAndCombineOptions struct {
	InputFile  string
	OutputFile string
	Ranges     []string
	NoAudio    bool
	ProgressCB func(percent string)
	ErrorCB    func(msg string)
	FinishedCB func()
}

func CutAndCombine(opt CutAndCombineOptions) error {
	tempDir := os.TempDir()
	var tempFiles []string

	defer func() {
		for _, tf := range tempFiles {
			os.Remove(tf)
		}
	}()

	for i, r := range opt.Ranges {
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid range format: %s", r)
		}
		start, end := parts[0], parts[1]

		tempFile := filepath.Join(tempDir, fmt.Sprintf("part_%d.mp4", i))
		args := []string{"-i", opt.InputFile, "-ss", start, "-to", end, "-c", "copy"}
		if opt.NoAudio {
			args = append(args, "-an")
		}
		args = append(args, tempFile)

		cmd := exec.Command("ffmpeg", args...)
		stderr, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			return err
		}

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if opt.ErrorCB != nil && (strings.Contains(line, "Error") || strings.Contains(line, "Invalid")) {
				opt.ErrorCB(line)
			}
			if strings.Contains(line, "time=") && opt.ProgressCB != nil {
				percent := float64(i+1) / float64(len(opt.Ranges)) * 100
				opt.ProgressCB(fmt.Sprintf("%.2f%%", percent))
			}
		}

		if err := cmd.Wait(); err != nil {
			return err
		}
		tempFiles = append(tempFiles, tempFile)
	}

	listFile := filepath.Join(tempDir, "list.txt")
	f, err := os.Create(listFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, tf := range tempFiles {
		_, _ = f.WriteString(fmt.Sprintf("file '%s'\n", tf))
	}

	concatArgs := []string{"-f", "concat", "-safe", "0", "-i", listFile, "-c", "copy", opt.OutputFile}
	cmd := exec.Command("ffmpeg", concatArgs...)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if opt.ErrorCB != nil && (strings.Contains(line, "Error") || strings.Contains(line, "Invalid")) {
			opt.ErrorCB(line)
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if opt.FinishedCB != nil {
		opt.FinishedCB()
	}

	os.Remove(listFile)

	return nil
}
