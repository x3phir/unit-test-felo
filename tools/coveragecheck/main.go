package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file := flag.String("file", "coverage.out", "coverage profile file")
	threshold := flag.Float64("threshold", 70, "minimum coverage percentage")
	flag.Parse()

	coverage, err := calculateCoverage(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "coverage check failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("coverage: %.2f%%\n", coverage)
	if coverage < *threshold {
		fmt.Fprintf(os.Stderr, "coverage threshold not met: got %.2f%% want %.2f%%\n", coverage, *threshold)
		os.Exit(1)
	}
}

func calculateCoverage(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var total, covered int64
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if lineNumber == 1 {
			if !strings.HasPrefix(line, "mode:") {
				return 0, errors.New("invalid coverage header")
			}
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 3 {
			return 0, fmt.Errorf("invalid profile line %d", lineNumber)
		}

		statements, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return 0, err
		}
		count, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return 0, err
		}

		total += statements
		if count > 0 {
			covered += statements
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, errors.New("no statements found in coverage profile")
	}

	return float64(covered) / float64(total) * 100, nil
}
