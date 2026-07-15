package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	socketPath := os.Getenv("HERDR_SOCKET_PATH")
	binPath := os.Getenv("HERDR_BIN_PATH")
	paneID := os.Getenv("HERDR_PANE_ID")

	if socketPath == "" {
		fail("HERDR_SOCKET_PATH is not set")
	}
	if binPath == "" {
		binPath = "herdr"
	}

	tr := newTransport(socketPath)
	defer tr.Close()

	if err := splitPane(binPath, paneID); err != nil {
		fail("split pane: " + err.Error())
	}

	root, err := exportLayout(tr, paneID)
	if err != nil {
		fail("export layout: " + err.Error())
	}

	if root.Type == "pane" {
		fmt.Println("single pane, nothing to equalize")
		return
	}

	targets := EqualizeColumns(root)

	var applied int
	for _, t := range targets {
		if t.Direction != "right" {
			continue
		}
		if err := setSplitRatio(tr, paneID, t.Path, t.Ratio); err != nil {
			fmt.Fprintf(os.Stderr, "warning: set_split_ratio failed: %v\n", err)
			continue
		}
		applied++
	}

	fmt.Printf("split and equalized %d column(s)\n", applied+1)
}

func splitPane(binPath, paneID string) error {
	if paneID == "" {
		return fmt.Errorf("HERDR_PANE_ID is empty, cannot determine target pane")
	}
	cmd := exec.Command(binPath, "pane", "split", "--pane", paneID, "--direction", "right", "--focus")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(out))
	}
	return nil
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "even-columns:", msg)
	os.Exit(1)
}
