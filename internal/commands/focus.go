package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/lunchboxsushi/jit/internal/storage"
	"github.com/lunchboxsushi/jit/internal/utils"
	"github.com/spf13/cobra"
)

var (
	typeFlag string
	listFlag bool
)

var focusCmd = &cobra.Command{
	Use:   "focus <query>",
	Short: "Focus on a ticket using fuzzy search",
	Long: `Switch your working context to a ticket using fuzzy search on keys and titles.
	
Examples:
  jit focus "5344"              # Focus on SRE-5344
  jit focus "bug"               # Search for tickets with "bug" in title
  jit focus "SRE" --type epic   # Search only epics
  jit focus "task" --list       # List matches without switching`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Configuration error: %v\n", err)
			fmt.Println("Run 'jit init' to create a configuration file")
			return
		}

		// Initialize storage
		storageInstance, err := storage.NewJSONStorage(cfg.App.DataDir)
		if err != nil {
			fmt.Printf("Storage error: %v\n", err)
			return
		}

		// Search for tickets
		results, err := searchTickets(query, storageInstance, typeFlag)
		if err != nil {
			fmt.Printf("Search error: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Printf("No tickets found matching '%s'\n", query)
			return
		}

		// If just listing, show results and exit
		if listFlag {
			displaySearchResults(results)
			return
		}

		// Select ticket to focus on
		selectedTicket, err := selectTicket(results)
		if err != nil {
			fmt.Printf("❌ Selection error: %v\n", err)
			return
		}

		// Update context
		contextManager := storage.NewContextManager(storageInstance)
		if err := contextManager.SetFocus(selectedTicket.Key, selectedTicket.Type); err != nil {
			fmt.Printf("❌ Failed to set focus: %v\n", err)
			return
		}

		// Add to recent tickets
		if err := contextManager.AddToRecent(selectedTicket.Key); err != nil {
			fmt.Printf("⚠️  Warning: Failed to add to recent tickets: %v\n", err)
		}

		fmt.Printf("🎯 Focused on %s (%s)\n", selectedTicket.Key, selectedTicket.Title)
	},
}

func init() {
	focusCmd.Flags().StringVar(&typeFlag, "type", "", "Filter by ticket type (epic|task|subtask)")
	focusCmd.Flags().BoolVar(&listFlag, "list", false, "List matches without switching focus")
}

// searchTickets searches for tickets matching the query
func searchTickets(query string, storageInstance storage.Storage, ticketType string) ([]utils.SearchResult, error) {
	// Get all ticket keys
	ticketKeys, err := storageInstance.ListTickets()
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %v", err)
	}

	// Load ticket information
	var tickets []utils.TicketInfo
	for _, key := range ticketKeys {
		ticket, err := storageInstance.LoadTicket(key)
		if err != nil {
			// Skip tickets that can't be loaded
			continue
		}

		tickets = append(tickets, utils.TicketInfo{
			Key:   ticket.Key,
			Title: ticket.Title,
			Type:  ticket.Type,
		})
	}

	// Perform fuzzy search
	results := utils.FuzzySearch(query, tickets)

	// Filter by type if specified
	if ticketType != "" {
		results = utils.FilterByType(results, ticketType)
	}

	return results, nil
}

// displaySearchResults displays search results
func displaySearchResults(results []utils.SearchResult) {
	fmt.Printf("🔍 Found %d matching tickets:\n\n", len(results))

	for i, result := range results {
		fmt.Printf("%d. %s (%s) - %s\n", i+1, result.Key, result.Type, result.Title)
		if result.Matched != "" {
			fmt.Printf("   Matched: %s\n", result.Matched)
		}
		fmt.Println()
	}
}

// selectTicket prompts user to select a ticket from results
func selectTicket(results []utils.SearchResult) (*utils.SearchResult, error) {
	if len(results) == 1 {
		// Auto-select if only one result
		return &results[0], nil
	}

	// Display results
	displaySearchResults(results)

	// Prompt for selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select ticket (1-", len(results), "): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %v", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("no selection made")
	}

	// Parse selection
	var selection int
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil {
		return nil, fmt.Errorf("invalid selection: %v", err)
	}

	if selection < 1 || selection > len(results) {
		return nil, fmt.Errorf("selection out of range")
	}

	return &results[selection-1], nil
}

// GetFocusCmd returns the focus command
func GetFocusCmd() *cobra.Command {
	return focusCmd
}
