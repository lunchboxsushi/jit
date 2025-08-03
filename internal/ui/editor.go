package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Editor provides functionality to open and edit files
type Editor struct {
	editor string
}

// NewEditor creates a new editor instance
func NewEditor() *Editor {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default fallback
	}

	return &Editor{
		editor: editor,
	}
}

// EditFile opens a file in the user's preferred editor
func (e *Editor) EditFile(filePath string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
	}

	// Open the editor
	cmd := exec.Command(e.editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// EditTemplate opens a template file for editing
func (e *Editor) EditTemplate(templatePath, outputPath string) error {
	// Read the template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %v", err)
	}

	// Ensure the output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write template content to output file
	if err := os.WriteFile(outputPath, templateContent, 0644); err != nil {
		return fmt.Errorf("failed to write template: %v", err)
	}

	// Open the file in editor
	return e.EditFile(outputPath)
}

// ReadFile reads the content of a file
func (e *Editor) ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(content), nil
}

// ParseMarkdownTicket parses markdown content into ticket fields
func (e *Editor) ParseMarkdownTicket(content string) (string, string, error) {
	lines := strings.Split(content, "\n")

	var title string
	var description strings.Builder

	inDescription := false
	skipNext := false

	for i, line := range lines {
		if skipNext {
			skipNext = false
			continue
		}

		line = strings.TrimSpace(line)

		// Skip empty lines at the beginning
		if title == "" && line == "" {
			continue
		}

		// Extract title (first non-empty line after #)
		if title == "" && strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			continue
		}

		// Start description after title
		if title != "" && !inDescription {
			inDescription = true
		}

		// Add to description
		if inDescription {
			// Skip section headers
			if strings.HasPrefix(line, "## ") {
				description.WriteString("\n\n")
				description.WriteString(line)
				description.WriteString("\n")
				continue
			}

			// Skip template placeholders
			if strings.Contains(line, "[Enter") || strings.Contains(line, "[Describe") ||
				strings.Contains(line, "[List") || strings.Contains(line, "[Additional") ||
				strings.Contains(line, "[Criterion") || strings.Contains(line, "[Technical") ||
				strings.Contains(line, "[Testing") || strings.Contains(line, "[Specific") {
				continue
			}

			// Add the line to description
			if line != "" || (i < len(lines)-1 && strings.TrimSpace(lines[i+1]) != "") {
				description.WriteString(line)
				description.WriteString("\n")
			}
		}
	}

	// Clean up description
	desc := strings.TrimSpace(description.String())

	// Validate
	if title == "" {
		return "", "", fmt.Errorf("title is required")
	}

	if desc == "" {
		return "", "", fmt.Errorf("description is required")
	}

	return title, desc, nil
}
