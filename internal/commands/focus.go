package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

		// Initialize command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "Failed to initialize")
			return
		}

		// Search for tickets
		results, err := searchTickets(query, ctx, typeFlag)
		if err != nil {
			HandleError(err, "Search failed")
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
			HandleError(err, "Selection failed")
			return
		}

		// Update context and recent tickets
		if err := ctx.UpdateContextAndRecent(selectedTicket.Key, selectedTicket.Type); err != nil {
			HandleError(err, "Failed to update context")
			return
		}

		fmt.Printf("Focused on %s (%s)\n", selectedTicket.Key, selectedTicket.Title)
	},
}

func init() {
	focusCmd.Flags().StringVar(&typeFlag, "type", "", "Filter by ticket type (epic|task|subtask)")
	focusCmd.Flags().BoolVar(&listFlag, "list", false, "List matches without switching focus")
}

// searchTickets searches for tickets matching the query
func searchTickets(query string, ctx *CommandContext, ticketType string) ([]utils.SearchResult, error) {
	// Get all ticket keys
	ticketKeys, err := ctx.Storage.ListTickets()
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %v", err)
	}

	// Load ticket information
	var tickets []utils.TicketInfo
	for _, key := range ticketKeys {
		ticket, err := ctx.Storage.LoadTicket(key)
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
	fmt.Printf("Found %d matching tickets:\n\n", len(results))

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
