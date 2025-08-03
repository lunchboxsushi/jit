package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateManager handles AI prompt templates
type TemplateManager struct {
	templateDir string
	cache       map[string]*template.Template
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(templateDir string) *TemplateManager {
	return &TemplateManager{
		templateDir: templateDir,
		cache:       make(map[string]*template.Template),
	}
}

// LoadTemplate loads a template by name
func (tm *TemplateManager) LoadTemplate(name string) (*template.Template, error) {
	// Check cache first
	if cached, exists := tm.cache[name]; exists {
		return cached, nil
	}

	// Load from file
	templatePath := filepath.Join(tm.templateDir, name+".txt")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %v", name, err)
	}

	// Parse template
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %v", name, err)
	}

	// Cache the template
	tm.cache[name] = tmpl
	return tmpl, nil
}

// ProcessTemplate processes a template with given data
func (tm *TemplateManager) ProcessTemplate(name string, data interface{}) (string, error) {
	tmpl, err := tm.LoadTemplate(name)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %v", name, err)
	}

	return result.String(), nil
}

// GetEnrichmentPrompt generates an enrichment prompt for the given context
func (tm *TemplateManager) GetEnrichmentPrompt(content string, context *EnrichmentContext) (string, error) {
	// Create template data
	data := map[string]interface{}{
		"Content":      content,
		"TicketType":   context.TicketType,
		"Project":      context.Project,
		"CurrentEpic":  context.CurrentEpic,
		"CurrentTask":  context.CurrentTask,
		"UserEmail":    context.UserEmail,
		"CustomFields": context.CustomFields,
	}

	// Try to load specific template for ticket type
	templateName := fmt.Sprintf("enrich_%s", strings.ToLower(context.TicketType))

	// Fall back to generic template if specific one doesn't exist
	if _, err := tm.LoadTemplate(templateName); err != nil {
		templateName = "enrich_generic"
	}

	return tm.ProcessTemplate(templateName, data)
}

// GetCommentPrompt generates a comment enrichment prompt
func (tm *TemplateManager) GetCommentPrompt(comment string, context *EnrichmentContext) (string, error) {
	data := map[string]interface{}{
		"Comment":      comment,
		"TicketType":   context.TicketType,
		"Project":      context.Project,
		"CurrentEpic":  context.CurrentEpic,
		"CurrentTask":  context.CurrentTask,
		"UserEmail":    context.UserEmail,
		"CustomFields": context.CustomFields,
	}

	return tm.ProcessTemplate("enrich_comment", data)
}

// CreateDefaultTemplates creates default template files if they don't exist
func (tm *TemplateManager) CreateDefaultTemplates() error {
	templates := map[string]string{
		"enrich_generic.txt": `You are an expert software development assistant. Your task is to enhance the following {{.TicketType}} description for a {{.Project}} project.

Original content:
{{.Content}}

Please enhance this description by:
1. Adding more technical details and context
2. Improving clarity and structure
3. Adding acceptance criteria if missing
4. Including relevant technical considerations
5. Making it more professional and comprehensive

Enhanced description:`,

		"enrich_epic.txt": `You are an expert software development assistant. Your task is to enhance the following EPIC description for a {{.Project}} project.

Original content:
{{.Content}}

Please enhance this EPIC description by:
1. Adding strategic context and business value
2. Defining clear objectives and success metrics
3. Outlining major milestones and deliverables
4. Identifying key stakeholders and dependencies
5. Adding risk considerations and mitigation strategies
6. Making it comprehensive for project planning

Enhanced EPIC description:`,

		"enrich_task.txt": `You are an expert software development assistant. Your task is to enhance the following TASK description for a {{.Project}} project.

Original content:
{{.Content}}

{{if .CurrentEpic}}This task belongs to epic: {{.CurrentEpic}}{{end}}

Please enhance this TASK description by:
1. Adding detailed technical requirements
2. Defining clear acceptance criteria
3. Including implementation considerations
4. Adding testing requirements
5. Specifying dependencies and prerequisites
6. Making it actionable for developers

Enhanced TASK description:`,

		"enrich_subtask.txt": `You are an expert software development assistant. Your task is to enhance the following SUBTASK description for a {{.Project}} project.

Original content:
{{.Content}}

{{if .CurrentTask}}This subtask belongs to task: {{.CurrentTask}}{{end}}
{{if .CurrentEpic}}This is part of epic: {{.CurrentEpic}}{{end}}

Please enhance this SUBTASK description by:
1. Adding specific implementation details
2. Defining precise acceptance criteria
3. Including code-level considerations
4. Adding unit testing requirements
5. Specifying any configuration needed
6. Making it ready for immediate development

Enhanced SUBTASK description:`,

		"enrich_comment.txt": `You are an expert software development assistant. Your task is to enhance the following comment for a {{.TicketType}} in the {{.Project}} project.

Original comment:
{{.Comment}}

{{if .CurrentEpic}}This is for a ticket in epic: {{.CurrentEpic}}{{end}}
{{if .CurrentTask}}This is for a subtask in task: {{.CurrentTask}}{{end}}

Please enhance this comment by:
1. Making it more professional and clear
2. Adding relevant technical context
3. Including actionable insights
4. Maintaining the original intent
5. Making it helpful for other team members

Enhanced comment:`,
	}

	// Create template directory if it doesn't exist
	if err := os.MkdirAll(tm.templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %v", err)
	}

	// Create each template file
	for name, content := range templates {
		path := filepath.Join(tm.templateDir, name)

		// Skip if file already exists
		if _, err := os.Stat(path); err == nil {
			continue
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create template %s: %v", name, err)
		}
	}

	return nil
}
