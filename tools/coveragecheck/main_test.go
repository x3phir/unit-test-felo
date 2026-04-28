package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateCoverage_ReturnsCoveragePercentage(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "coverage.out")
	content := "mode: atomic\nexample.go:1.1,2.2 2 1\nexample.go:3.1,4.2 2 0\n"
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	coverage, err := calculateCoverage(file)
	if err != nil {
		t.Fatalf("calculateCoverage() error = %v", err)
	}
	if coverage != 50 {
		t.Fatalf("coverage = %v, want 50", coverage)
	}
}

func TestCalculateCoverage_InvalidHeaderReturnsError(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "coverage.out")
	if err := os.WriteFile(file, []byte("bad header\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := calculateCoverage(file); err == nil {
		t.Fatal("calculateCoverage() error = nil, want error")
	}
}
