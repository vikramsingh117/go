package utils

import (
	"bytes"
	"fmt"
	// "os"
	"os/exec"
)

// ConvertWAVToFLAC converts WAV data to FLAC format using an external encoder (e.g., flac CLI tool)
func ConvertWAVToFLAC(inputFilePath string, outputFilePath string) error {
	// Use the flac CLI tool to convert WAV to FLAC
	cmd := exec.Command("flac", inputFilePath, "-o", outputFilePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("conversion failed: %v, %v", err, stderr.String())
	}

	fmt.Println("Conversion successful:", out.String())
	return nil
}
