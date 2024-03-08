package tui

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/james-do2024/ghi/client"
	"github.com/james-do2024/ghi/config"
)

type TuiState struct {
	DirMap         map[int]string
	FileContent    *string
	LessRSupported bool
}

func Init() *TuiState {
	return &TuiState{
		DirMap:         make(map[int]string),
		LessRSupported: checkLessRSupport(),
	}
}

func (t *TuiState) Interact(req *client.RestRequest) error {
	for {
		// If we're displaying a file, try to highlight and page it
		if t.FileContent != nil {
			pageCode(t)
		} else {
			// Golang's map output isn't guaranteed to be in order,
			// so we address that here
			keysToSort := make([]int, 0, len(t.DirMap))
			for i := range t.DirMap {
				keysToSort = append(keysToSort, i)
			}
			sort.Ints(keysToSort)
			for _, name := range keysToSort {
				fmt.Printf("[%d] %s\n", name, t.DirMap[name])
			}
		}

		// Get user input, treat failure as fatal
		input, err := displayPrompt()
		if err != nil {
			return fmt.Errorf("error reading input: %v", err)
		}

		// Handle user input
		switch input {
		case "q", "Q": // Be a little forgiving
			fmt.Println("Exiting...")
			os.Exit(config.ExitOk)
		case "..": // Move up one directory
			req.NavigateUp()
		case "^": // Move to the root directory
			req.NavigateRoot()
		default: // Check if input is a number for navigation
			if idx, err := strconv.Atoi(input); err == nil {
				req.NavigateIndex(idx, t.DirMap)
			} else {
				fmt.Println("invalid command:", input)
			}
		}

		// Update content after input handling
		if err := t.UpdateContent(req); err != nil {
			return fmt.Errorf("error updating content: %v", err)
		}
	}
}

func (t *TuiState) UpdateContent(req *client.RestRequest) error {
	// Clear DirMap of stale data
	t.DirMap = make(map[int]string)

	fileContent, directoryContent, err := req.GetContent()
	if err != nil {
		return err
	}

	// if we have fileContent, no need to worry about directory contents
	if fileContent != nil {
		t.FileContent = fileContent
		return nil
	}

	// Updating DirMap contents
	for i, content := range directoryContent {
		t.DirMap[i+1] = *content.Name
	}
	t.FileContent = nil // Use this as a marker for what kind of content we have

	return nil
}

func (t *TuiState) Page() error {
	page := exec.Command("less", "-R")
	page.Stdin = bytes.NewReader([]byte(*t.FileContent))
	page.Stdout = os.Stdout
	page.Stderr = os.Stderr

	if err := page.Run(); err != nil {
		return fmt.Errorf("less invocation failed: %v", err)
	}
	return nil
}

func displayPrompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ") // This serves as our prompt

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}
	return strings.TrimSpace(input), nil // Return trimmed input
}

// Very naive way of checking whether `less -R` works, but doable
func checkLessRSupport() bool {
	var out bytes.Buffer

	cmd := exec.Command("less", "--help")

	cmd.Stdout = &out // Capture output
	err := cmd.Run()
	if err != nil {
		// The command might fail if 'less' is not installed,
		// but output of --help is common to all modern versions of less
		log.Printf("error executing 'less --help': %v\n", err)
		return false
	}

	// Search for '-R' in the help text
	helpText := out.String()
	return bytes.Contains([]byte(helpText), []byte("-R"))
}

func pageCode(t *TuiState) {
	// Attempt to colorize the file content
	var output *string
	colorized, err := t.Colorize()
	if err != nil {
		// If colorization fails, use the original file content for paging
		log.Printf("Failed to colorize file content: %v\n", err)
		output = t.FileContent
	} else {
		output = &colorized
	}
	t.FileContent = output

	// Attempt to page the content, regardless of colorization success
	if t.LessRSupported {
		err := t.Page()
		if err != nil {
			// If paging fails (less is not available), print the content directly
			fmt.Print(*output)
			log.Printf("error paging content: %v\n", err)
		}
	} else {
		// If less -R is not supported, print
		fmt.Print(*output)
		log.Println("system `less` does not support -R: not paging output")
	}
}
