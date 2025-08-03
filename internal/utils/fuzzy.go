package utils

import (
	"strings"
	"unicode"
)

// SearchResult represents a search result with score
type SearchResult struct {
	Key     string
	Title   string
	Type    string
	Score   int
	Matched string // What part matched
}

// TicketInfo represents basic ticket information for search
type TicketInfo struct {
	Key   string
	Title string
	Type  string
}

// FuzzySearch performs fuzzy search on tickets
func FuzzySearch(query string, tickets []TicketInfo) []SearchResult {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return []SearchResult{}
	}

	var results []SearchResult

	for _, ticket := range tickets {
		score, matched := calculateScore(query, ticket.Key, ticket.Title)
		if score > 0 {
			results = append(results, SearchResult{
				Key:     ticket.Key,
				Title:   ticket.Title,
				Type:    ticket.Type,
				Score:   score,
				Matched: matched,
			})
		}
	}

	// Sort by score (higher is better)
	sortSearchResults(results)

	return results
}

// calculateScore calculates a fuzzy match score
func calculateScore(query, key, title string) (int, string) {
	keyLower := strings.ToLower(key)
	titleLower := strings.ToLower(title)

	// Exact key match gets highest score
	if keyLower == query {
		return 1000, key
	}

	// Key starts with query
	if strings.HasPrefix(keyLower, query) {
		return 800, key
	}

	// Key contains query
	if strings.Contains(keyLower, query) {
		return 600, key
	}

	// Title contains query
	if strings.Contains(titleLower, query) {
		return 400, title
	}

	// Fuzzy match on key
	if fuzzyMatch(query, keyLower) {
		return 300, key
	}

	// Fuzzy match on title
	if fuzzyMatch(query, titleLower) {
		return 200, title
	}

	return 0, ""
}

// fuzzyMatch performs a simple fuzzy match
func fuzzyMatch(query, text string) bool {
	if len(query) == 0 {
		return true
	}

	queryIdx := 0
	for _, char := range text {
		if queryIdx < len(query) && unicode.ToLower(char) == unicode.ToLower(rune(query[queryIdx])) {
			queryIdx++
		}
	}

	return queryIdx == len(query)
}

// sortSearchResults sorts results by score (descending)
func sortSearchResults(results []SearchResult) {
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

// FilterByType filters search results by ticket type
func FilterByType(results []SearchResult, ticketType string) []SearchResult {
	if ticketType == "" {
		return results
	}

	var filtered []SearchResult
	for _, result := range results {
		if result.Type == ticketType {
			filtered = append(filtered, result)
		}
	}

	return filtered
}
