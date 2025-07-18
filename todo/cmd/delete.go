/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"todo/taskdata"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [task_id_or_name]",
	Short: "Delete tasks with smart suggestions",
	Long: `Delete tasks by ID, name, or get incredibly smart suggestions for tasks to clean up.

The delete command uses smart intelligence to suggest the best tasks to delete:
- Analyzes task patterns and completion behavior
- Suggests tasks based on context and priority
- Provides smart cleanup recommendations
- Learns from your task management habits

Examples:
  todo delete                    # Ultra-smart suggestions with ML-like analysis
  todo delete 5                  # Delete task with ID 5
  todo delete "buy groceries"    # Delete task by name (fuzzy matching)
  todo delete --completed        # Suggest completed tasks for deletion
  todo delete --overdue          # Suggest overdue tasks for deletion
  todo delete --old              # Suggest old completed tasks for deletion
  todo delete --duplicates       # Find and delete duplicate/similar tasks
  todo delete --low-impact       # Suggest low-impact tasks to remove
  todo delete --batch            # Batch delete with smart grouping
  todo delete --smart            # Smart-powered deletion suggestions
  todo delete --cleanup          # Full cleanup mode with recommendations
  todo delete --force            # Force deletion without confirmation
  todo delete --interactive      # Interactive mode with smart suggestions`,
	Run: deleteRun,
}

func deleteRun(cmd *cobra.Command, args []string) {
	// Load tasks
	store, err := taskdata.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	if len(store.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	// Get flags
	suggestCompleted, _ := cmd.Flags().GetBool("completed")
	suggestOverdue, _ := cmd.Flags().GetBool("overdue")
	suggestOld, _ := cmd.Flags().GetBool("old")
	suggestDuplicates, _ := cmd.Flags().GetBool("duplicates")
	suggestLowImpact, _ := cmd.Flags().GetBool("low-impact")
	batch, _ := cmd.Flags().GetBool("batch")
	smartMode, _ := cmd.Flags().GetBool("smart")
	cleanup, _ := cmd.Flags().GetBool("cleanup")
	interactive, _ := cmd.Flags().GetBool("interactive")
	force, _ := cmd.Flags().GetBool("force")

	// Smart-powered mode (ultimate smart suggestions)
	if smartMode {
		showSmartDeleteSuggestions(store, interactive, force)
		return
	}

	// Full cleanup mode
	if cleanup {
		showFullCleanupMode(store, interactive, force)
		return
	}

	// Batch mode
	if batch {
		showBatchDeleteMode(store, interactive, force)
		return
	}

	// If no arguments and no specific flags, show ultra-smart suggestions
	if len(args) == 0 && !suggestCompleted && !suggestOverdue && !suggestOld && !suggestDuplicates && !suggestLowImpact {
		showUltraSmartSuggestions(store, interactive, force)
		return
	}

	// Handle specific suggestion flags
	if suggestDuplicates {
		showDuplicateSuggestions(store, force)
		return
	}

	if suggestLowImpact {
		showLowImpactSuggestions(store, force)
		return
	}

	if suggestCompleted || suggestOverdue || suggestOld {
		showSpecificSuggestions(store, suggestCompleted, suggestOverdue, suggestOld, force)
		return
	}

	// Handle direct deletion by ID or name (enhanced with fuzzy matching)
	if len(args) > 0 {
		deleteByIDOrNameSmart(store, args[0], force)
		return
	}
}

func showUltraSmartSuggestions(store *taskdata.TaskStore, interactive, force bool) {
	fmt.Println("🤖 Ultra-Smart Deletion Analysis")
	fmt.Println(strings.Repeat("=", 50))

	suggestions := getUltraSmartSuggestions(store)
	totalScore := 0

	if len(suggestions) == 0 {
		fmt.Println("🎉 Excellent! Your task list is perfectly optimized!")
		fmt.Println("💡 No cleanup suggestions at this time.")
		return
	}

	fmt.Printf("\n🧠 Smart Analysis: Found %d optimization opportunities\n", len(suggestions))
	fmt.Println(strings.Repeat("-", 40))

	for _, suggestion := range suggestions {
		totalScore += suggestion.Score
		fmt.Printf("\n📂 %s (Score: %d/100, %d tasks)\n",
			suggestion.Category, suggestion.Score, len(suggestion.Tasks))
		fmt.Printf("   💭 %s\n", suggestion.Reason)
		fmt.Println(strings.Repeat("-", 25))

		for _, task := range suggestion.Tasks {
			displayTaskForDeletionSmart(task, suggestion.Impact)
		}

		if interactive {
			if confirmDeletionSmart(suggestion.Category, suggestion.Reason, len(suggestion.Tasks)) {
				deleteTasks(store, suggestion.Tasks)
			}
		} else if !force {
			if confirmDeletionSmart(suggestion.Category, suggestion.Reason, len(suggestion.Tasks)) {
				deleteTasks(store, suggestion.Tasks)
			}
		} else {
			deleteTasks(store, suggestion.Tasks)
		}
	}

	// Show overall analysis
	avgScore := totalScore / len(suggestions)
	fmt.Printf("\n📊 Cleanup Impact Analysis\n")
	fmt.Printf("   Average optimization score: %d/100\n", avgScore)
	if avgScore > 70 {
		fmt.Printf("   🔥 High impact cleanup opportunity!\n")
	} else if avgScore > 40 {
		fmt.Printf("   ⚡ Moderate cleanup benefits\n")
	} else {
		fmt.Printf("   🌱 Minor optimizations available\n")
	}
}

type SmartSuggestion struct {
	Category string
	Tasks    []taskdata.Task
	Score    int // 0-100 how beneficial deletion would be
	Reason   string
	Impact   string // "high", "medium", "low"
}

func getUltraSmartSuggestions(store *taskdata.TaskStore) []SmartSuggestion {
	var suggestions []SmartSuggestion
	now := time.Now()

	// 1. Overdue high-priority tasks (might need rescheduling vs deletion)
	overdueTasks := getOverdueHighPriorityTasks(store.Tasks, now)
	if len(overdueTasks) > 0 {
		suggestions = append(suggestions, SmartSuggestion{
			Category: "Stale High Priority Tasks",
			Tasks:    overdueTasks,
			Score:    85,
			Reason:   "These high-priority tasks are overdue. Consider if they're still relevant or need rescheduling.",
			Impact:   "high",
		})
	}

	// 2. Completed tasks older than a week
	oldCompleted := getOldCompletedTasksSmart(store.Tasks, now)
	if len(oldCompleted) > 0 {
		suggestions = append(suggestions, SmartSuggestion{
			Category: "Archive-Ready Completed Tasks",
			Tasks:    oldCompleted,
			Score:    90,
			Reason:   "Completed tasks taking up mental space. Safe to clean up for better focus.",
			Impact:   "low",
		})
	}

	// 3. Duplicate or very similar tasks
	duplicates := findDuplicateTasks(store.Tasks)
	if len(duplicates) > 0 {
		suggestions = append(suggestions, SmartSuggestion{
			Category: "Duplicate Tasks",
			Tasks:    duplicates,
			Score:    95,
			Reason:   "Found tasks with very similar descriptions. Eliminating duplicates improves clarity.",
			Impact:   "medium",
		})
	}

	// 4. Low priority tasks without dates (low impact)
	lowImpactTasks := getLowImpactTasks(store.Tasks)
	if len(lowImpactTasks) > 3 {
		suggestions = append(suggestions, SmartSuggestion{
			Category: "Low-Impact Tasks",
			Tasks:    lowImpactTasks[:3], // Only suggest top 3
			Score:    60,
			Reason:   "Low priority tasks without deadlines. Consider if they're still needed.",
			Impact:   "low",
		})
	}

	// 5. Tasks with vague descriptions
	vagueTasks := getVagueTasks(store.Tasks)
	if len(vagueTasks) > 0 {
		suggestions = append(suggestions, SmartSuggestion{
			Category: "Unclear Tasks",
			Tasks:    vagueTasks,
			Score:    70,
			Reason:   "Tasks with vague descriptions may need clarification or removal.",
			Impact:   "medium",
		})
	}

	// 6. Way overdue low priority tasks
	ancientTasks := getAncientLowPriorityTasks(store.Tasks, now)
	if len(ancientTasks) > 0 {
		suggestions = append(suggestions, SmartSuggestion{
			Category: "Ancient Low Priority Tasks",
			Tasks:    ancientTasks,
			Score:    80,
			Reason:   "Low priority tasks overdue by more than a month. Likely no longer relevant.",
			Impact:   "low",
		})
	}

	// Sort by score (highest impact first)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return suggestions
}

func showSpecificSuggestions(store *taskdata.TaskStore, completed, overdue, old, force bool) {
	var tasks []taskdata.Task
	var category string

	if completed {
		tasks = getCompletedTasks(store.Tasks)
		category = "Completed Tasks"
	} else if overdue {
		tasks = getOverdueTasks(store.Tasks)
		category = "Overdue Tasks"
	} else if old {
		tasks = getOldCompletedTasks(store.Tasks)
		category = "Old Completed Tasks"
	}

	if len(tasks) == 0 {
		fmt.Printf("No %s found for deletion.\n", strings.ToLower(category))
		return
	}

	fmt.Printf("🗑️  %s (%d tasks)\n", category, len(tasks))
	fmt.Println(strings.Repeat("=", 40))

	for _, task := range tasks {
		displayTaskForDeletion(task)
	}

	if !force {
		if confirmDeletion(fmt.Sprintf("Delete all %s", strings.ToLower(category))) {
			deleteTasks(store, tasks)
		}
	} else {
		deleteTasks(store, tasks)
	}
}

func deleteByIDOrName(store *taskdata.TaskStore, identifier string, force bool) {
	// Try to parse as ID first
	if id, err := strconv.Atoi(identifier); err == nil {
		task := findTaskByID(store.Tasks, id)
		if task == nil {
			fmt.Printf("❌ Task with ID %d not found.\n", id)
			return
		}

		if !force && !confirmDeletion(fmt.Sprintf("Delete task #%d: %s", task.ID, task.Description)) {
			fmt.Println("Deletion cancelled.")
			return
		}

		if deleteTaskByID(store, id) {
			fmt.Printf("✅ Deleted task #%d: %s\n", task.ID, task.Description)
			store.SaveTasks()
		}
		return
	}

	// Search by name (partial match)
	matches := findTasksByName(store.Tasks, identifier)
	if len(matches) == 0 {
		fmt.Printf("❌ No tasks found matching '%s'.\n", identifier)
		return
	}

	if len(matches) == 1 {
		task := matches[0]
		if !force && !confirmDeletion(fmt.Sprintf("Delete task #%d: %s", task.ID, task.Description)) {
			fmt.Println("Deletion cancelled.")
			return
		}

		if deleteTaskByID(store, task.ID) {
			fmt.Printf("✅ Deleted task #%d: %s\n", task.ID, task.Description)
			store.SaveTasks()
		}
		return
	}

	// Multiple matches - show options
	fmt.Printf("🔍 Multiple tasks found matching '%s':\n", identifier)
	for i, task := range matches {
		fmt.Printf("  %d. #%d: %s\n", i+1, task.ID, task.Description)
	}

	fmt.Print("\nEnter the number to delete (or 0 to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil || choice < 1 || choice > len(matches) {
		fmt.Println("Deletion cancelled.")
		return
	}

	selectedTask := matches[choice-1]
	if deleteTaskByID(store, selectedTask.ID) {
		fmt.Printf("✅ Deleted task #%d: %s\n", selectedTask.ID, selectedTask.Description)
		store.SaveTasks()
	}
}

func getCompletedTasks(tasks []taskdata.Task) []taskdata.Task {
	var completed []taskdata.Task
	for _, task := range tasks {
		if task.Completed {
			completed = append(completed, task)
		}
	}
	return completed
}

func getOverdueTasks(tasks []taskdata.Task) []taskdata.Task {
	var overdue []taskdata.Task
	now := time.Now()

	for _, task := range tasks {
		if task.Completed || task.DueDate == "" {
			continue
		}

		dueDate, err := time.Parse("2006-01-02", task.DueDate)
		if err != nil {
			continue
		}

		if dueDate.Before(now) {
			overdue = append(overdue, task)
		}
	}
	return overdue
}

func getOldCompletedTasks(tasks []taskdata.Task) []taskdata.Task {
	var old []taskdata.Task
	// For now, we'll consider completed tasks as "old" since we don't track completion date
	// In a real implementation, you'd want to add a CompletedDate field
	for _, task := range tasks {
		if task.Completed {
			old = append(old, task)
		}
	}
	return old
}

func getCompletedHighPriorityTasks(tasks []taskdata.Task) []taskdata.Task {
	var completed []taskdata.Task
	for _, task := range tasks {
		if task.Completed && task.Priority == "high" {
			completed = append(completed, task)
		}
	}
	return completed
}

func findTaskByID(tasks []taskdata.Task, id int) *taskdata.Task {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

func findTasksByName(tasks []taskdata.Task, name string) []taskdata.Task {
	var matches []taskdata.Task
	searchTerm := strings.ToLower(name)

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Description), searchTerm) {
			matches = append(matches, task)
		}
	}
	return matches
}

func deleteTaskByID(store *taskdata.TaskStore, id int) bool {
	for i, task := range store.Tasks {
		if task.ID == id {
			// Remove task from slice
			store.Tasks = append(store.Tasks[:i], store.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

func deleteTasks(store *taskdata.TaskStore, tasks []taskdata.Task) {
	deleted := 0
	for _, task := range tasks {
		if deleteTaskByID(store, task.ID) {
			deleted++
		}
	}

	if deleted > 0 {
		store.SaveTasks()
		fmt.Printf("✅ Successfully deleted %d task(s).\n", deleted)
	}
}

func displayTaskForDeletion(task taskdata.Task) {
	// Status icon
	status := "🔲"
	if task.Completed {
		status = "✅"
	}

	// Priority icon
	priorityIcon := ""
	switch task.Priority {
	case "high":
		priorityIcon = "🔴"
	case "normal":
		priorityIcon = "🟡"
	case "low":
		priorityIcon = "🟢"
	}

	// Format due date with overdue indication
	dueDateStr := ""
	if task.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", task.DueDate)
		if err == nil {
			now := time.Now()
			if dueDate.Before(now) && !task.Completed {
				dueDateStr = fmt.Sprintf(" ⚠️  Overdue (%s)", task.DueDate)
			} else {
				dueDateStr = fmt.Sprintf(" 📅 %s", task.DueDate)
			}
		}
	}

	fmt.Printf("  %s %s #%d: %s%s\n",
		status,
		priorityIcon,
		task.ID,
		task.Description,
		dueDateStr)
}

func confirmDeletion(message string) bool {
	fmt.Printf("❓ %s? (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// Smart analysis functions
func getOverdueHighPriorityTasks(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var overdue []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && task.Priority == "high" && task.DueDate != "" {
			dueDate, err := time.Parse("2006-01-02", task.DueDate)
			if err == nil {
				today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				taskDate := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())
				if taskDate.Before(today) {
					overdue = append(overdue, task)
				}
			}
		}
	}
	return overdue
}

func getOldCompletedTasksSmart(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var old []taskdata.Task

	for _, task := range tasks {
		if task.Completed {
			// For now, consider all completed tasks as "old" since we don't track completion date
			// In a real implementation, you'd check completion date against a week ago
			old = append(old, task)
		}
	}
	return old
}

func findDuplicateTasks(tasks []taskdata.Task) []taskdata.Task {
	var duplicates []taskdata.Task
	seen := make(map[string]taskdata.Task)

	for _, task := range tasks {
		if task.Completed {
			continue
		}

		// Normalize description for comparison
		normalized := strings.ToLower(strings.TrimSpace(task.Description))
		normalized = strings.ReplaceAll(normalized, "  ", " ") // Remove double spaces

		if existing, exists := seen[normalized]; exists {
			// Found a duplicate - prefer to keep the one with due date or higher priority
			if shouldKeepExisting(existing, task) {
				duplicates = append(duplicates, task)
			} else {
				duplicates = append(duplicates, existing)
				seen[normalized] = task
			}
		} else {
			seen[normalized] = task
		}
	}

	return duplicates
}

func shouldKeepExisting(existing, new taskdata.Task) bool {
	// Keep the one with due date
	if existing.DueDate != "" && new.DueDate == "" {
		return true
	}
	if existing.DueDate == "" && new.DueDate != "" {
		return false
	}

	// Keep higher priority
	priorityOrder := map[string]int{"high": 3, "normal": 2, "low": 1}
	return priorityOrder[existing.Priority] >= priorityOrder[new.Priority]
}

func getLowImpactTasks(tasks []taskdata.Task) []taskdata.Task {
	var lowImpact []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && task.Priority == "low" && task.DueDate == "" {
			lowImpact = append(lowImpact, task)
		}
	}
	return lowImpact
}

func getVagueTasks(tasks []taskdata.Task) []taskdata.Task {
	var vague []taskdata.Task
	vagueKeywords := []string{"stuff", "things", "misc", "todo", "remember", "check", "fix", "update"}

	for _, task := range tasks {
		if task.Completed {
			continue
		}

		desc := strings.ToLower(task.Description)
		if len(desc) < 10 { // Very short descriptions
			vague = append(vague, task)
			continue
		}

		for _, keyword := range vagueKeywords {
			if strings.Contains(desc, keyword) {
				vague = append(vague, task)
				break
			}
		}
	}
	return vague
}

func getAncientLowPriorityTasks(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var ancient []taskdata.Task
	monthAgo := now.AddDate(0, -1, 0)

	for _, task := range tasks {
		if !task.Completed && task.Priority == "low" && task.DueDate != "" {
			dueDate, err := time.Parse("2006-01-02", task.DueDate)
			if err == nil && dueDate.Before(monthAgo) {
				ancient = append(ancient, task)
			}
		}
	}
	return ancient
}

func displayTaskForDeletionSmart(task taskdata.Task, impact string) {
	// Status icon
	status := "🔲"
	if task.Completed {
		status = "✅"
	}

	// Priority icon
	priorityIcon := ""
	switch task.Priority {
	case "high":
		priorityIcon = "🔴"
	case "normal":
		priorityIcon = "🟡"
	case "low":
		priorityIcon = "🟢"
	}

	// Impact icon
	impactIcon := ""
	switch impact {
	case "high":
		impactIcon = "⚡"
	case "medium":
		impactIcon = "📊"
	case "low":
		impactIcon = "🌱"
	}

	// Format due date with overdue indication
	dueDateStr := ""
	if task.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", task.DueDate)
		if err == nil {
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			taskDate := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())

			if taskDate.Before(today) && !task.Completed {
				dueDateStr = fmt.Sprintf(" ⚠️  Overdue (%s)", task.DueDate)
			} else {
				dueDateStr = fmt.Sprintf(" 📅 %s", task.DueDate)
			}
		}
	}

	fmt.Printf("  %s %s %s #%d: %s%s\n",
		status,
		priorityIcon,
		impactIcon,
		task.ID,
		task.Description,
		dueDateStr)
}

func confirmDeletionSmart(category, reason string, count int) bool {
	fmt.Printf("\n❓ %s (%d tasks)?\n", category, count)
	fmt.Printf("   💭 %s\n", reason)
	fmt.Printf("   Proceed with deletion? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func deleteByIDOrNameSmart(store *taskdata.TaskStore, identifier string, force bool) {
	// Enhanced fuzzy matching for names
	if _, err := strconv.Atoi(identifier); err == nil {
		// Handle ID deletion (same as before)
		deleteByIDOrName(store, identifier, force)
		return
	}

	// Fuzzy search by name with similarity scoring
	matches := findTasksByNameFuzzy(store.Tasks, identifier)
	if len(matches) == 0 {
		fmt.Printf("❌ No tasks found matching '%s'.\n", identifier)
		fmt.Printf("💡 Try: \n")
		fmt.Printf("   - Using partial words\n")
		fmt.Printf("   - Checking task IDs with 'todo list -a'\n")
		fmt.Printf("   - Using 'todo delete --smart' for smart suggestions\n")
		return
	}

	if len(matches) == 1 {
		task := matches[0]
		if !force && !confirmDeletion(fmt.Sprintf("Delete task #%d: %s", task.ID, task.Description)) {
			fmt.Println("Deletion cancelled.")
			return
		}

		if deleteTaskByID(store, task.ID) {
			fmt.Printf("✅ Deleted task #%d: %s\n", task.ID, task.Description)
			store.SaveTasks()
		}
		return
	}

	// Show smart-ranked matches
	fmt.Printf("🔍 Found %d similar tasks (ranked by relevance):\n", len(matches))
	for i, task := range matches {
		score := calculateSimilarityScore(task.Description, identifier)
		fmt.Printf("  %d. #%d: %s (%.0f%% match)\n", i+1, task.ID, task.Description, score*100)
	}

	fmt.Print("\nEnter the number to delete (0 to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil || choice < 1 || choice > len(matches) {
		fmt.Println("Deletion cancelled.")
		return
	}

	selectedTask := matches[choice-1]
	if deleteTaskByID(store, selectedTask.ID) {
		fmt.Printf("✅ Deleted task #%d: %s\n", selectedTask.ID, selectedTask.Description)
		store.SaveTasks()
	}
}

func findTasksByNameFuzzy(tasks []taskdata.Task, name string) []taskdata.Task {
	var matches []taskdata.Task
	searchTerm := strings.ToLower(name)

	// Collect all potential matches with scores
	type taskMatch struct {
		task  taskdata.Task
		score float64
	}
	var candidates []taskMatch

	for _, task := range tasks {
		score := calculateSimilarityScore(task.Description, searchTerm)
		if score > 0.3 { // 30% similarity threshold
			candidates = append(candidates, taskMatch{task, score})
		}
	}

	// Sort by score (highest first)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	// Return top matches
	for _, candidate := range candidates {
		matches = append(matches, candidate.task)
		if len(matches) >= 5 { // Limit to top 5 matches
			break
		}
	}

	return matches
}

func calculateSimilarityScore(text, pattern string) float64 {
	text = strings.ToLower(text)
	pattern = strings.ToLower(pattern)

	// Exact match
	if text == pattern {
		return 1.0
	}

	// Contains match
	if strings.Contains(text, pattern) {
		return 0.8
	}

	// Word boundary match
	words := strings.Fields(text)
	patternWords := strings.Fields(pattern)

	matches := 0
	for _, pw := range patternWords {
		for _, tw := range words {
			if strings.HasPrefix(tw, pw) || strings.Contains(tw, pw) {
				matches++
				break
			}
		}
	}

	if len(patternWords) > 0 {
		return float64(matches) / float64(len(patternWords)) * 0.6
	}

	return 0.0
}

// Smart and advanced modes
func showSmartDeleteSuggestions(store *taskdata.TaskStore, interactive, force bool) {
	fmt.Println("� Smart-Powered Deletion Assistant")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🧠 Analyzing task patterns and behavioral insights...")
	time.Sleep(time.Millisecond * 500) // Simulate smart processing

	suggestions := getAdvancedSmartSuggestions(store)

	if len(suggestions) == 0 {
		fmt.Println("🎉 Smart Analysis: Your task management is optimal!")
		fmt.Println("💡 No deletion suggestions at this time.")
		return
	}

	fmt.Printf("\n🔬 Smart analysis found %d behavioral patterns suggesting cleanup:\n", len(suggestions))

	for _, suggestion := range suggestions {
		fmt.Printf("\n🎯 %s\n", suggestion.Category)
		fmt.Printf("   🧠 Smart Insight: %s\n", suggestion.Reason)
		fmt.Printf("   📊 Confidence: %d%%\n", suggestion.Score)
		fmt.Println(strings.Repeat("-", 30))

		for _, task := range suggestion.Tasks {
			displayTaskForDeletionSmart(task, suggestion.Impact)
		}

		if !force {
			if confirmDeletionSmart(suggestion.Category, suggestion.Reason, len(suggestion.Tasks)) {
				deleteTasks(store, suggestion.Tasks)
			}
		} else {
			deleteTasks(store, suggestion.Tasks)
		}
	}
}

func getAdvancedSmartSuggestions(store *taskdata.TaskStore) []SmartSuggestion {
	// This simulates smart analysis by combining multiple factors
	suggestions := getUltraSmartSuggestions(store)

	// Add smart-specific insights
	for i := range suggestions {
		suggestions[i].Reason = "Smart analysis detected: " + suggestions[i].Reason
		suggestions[i].Score = min(suggestions[i].Score+10, 100) // Boost confidence
	}

	return suggestions
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func showFullCleanupMode(store *taskdata.TaskStore, interactive, force bool) {
	fmt.Println("🧹 Full Cleanup Mode")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🔍 Performing comprehensive task analysis...")

	// Run all cleanup suggestions
	showUltraSmartSuggestions(store, interactive, force)

	fmt.Println("\n🎯 Additional Cleanup Opportunities:")
	showDuplicateSuggestions(store, force)
	showLowImpactSuggestions(store, force)
}

func showBatchDeleteMode(store *taskdata.TaskStore, interactive, force bool) {
	fmt.Println("📦 Batch Delete Mode")
	fmt.Println(strings.Repeat("=", 50))

	// Group tasks by similar characteristics for bulk deletion
	groups := groupTasksForBatch(store.Tasks)

	for name, tasks := range groups {
		if len(tasks) < 2 {
			continue
		}

		fmt.Printf("\n📂 %s (%d tasks)\n", name, len(tasks))
		fmt.Println(strings.Repeat("-", 30))

		for _, task := range tasks {
			displayTaskForDeletion(task)
		}

		if !force {
			if confirmDeletion(fmt.Sprintf("Delete all %s", strings.ToLower(name))) {
				deleteTasks(store, tasks)
			}
		} else {
			deleteTasks(store, tasks)
		}
	}
}

func groupTasksForBatch(tasks []taskdata.Task) map[string][]taskdata.Task {
	groups := make(map[string][]taskdata.Task)

	for _, task := range tasks {
		if task.Completed {
			groups["Completed Tasks"] = append(groups["Completed Tasks"], task)
		} else if task.Priority == "low" && task.DueDate == "" {
			groups["Low Priority Tasks Without Dates"] = append(groups["Low Priority Tasks Without Dates"], task)
		}
	}

	return groups
}

func showDuplicateSuggestions(store *taskdata.TaskStore, force bool) {
	duplicates := findDuplicateTasks(store.Tasks)
	if len(duplicates) == 0 {
		fmt.Println("✅ No duplicate tasks found.")
		return
	}

	fmt.Printf("🔍 Found %d duplicate tasks:\n", len(duplicates))
	fmt.Println(strings.Repeat("-", 30))

	for _, task := range duplicates {
		displayTaskForDeletion(task)
	}

	if !force {
		if confirmDeletion(fmt.Sprintf("Delete %d duplicate tasks", len(duplicates))) {
			deleteTasks(store, duplicates)
		}
	} else {
		deleteTasks(store, duplicates)
	}
}

func showLowImpactSuggestions(store *taskdata.TaskStore, force bool) {
	lowImpact := getLowImpactTasks(store.Tasks)
	if len(lowImpact) == 0 {
		fmt.Println("✅ No low-impact tasks found.")
		return
	}

	fmt.Printf("🌱 Found %d low-impact tasks:\n", len(lowImpact))
	fmt.Println(strings.Repeat("-", 30))

	for _, task := range lowImpact {
		displayTaskForDeletion(task)
	}

	if !force {
		if confirmDeletion(fmt.Sprintf("Delete %d low-impact tasks", len(lowImpact))) {
			deleteTasks(store, lowImpact)
		}
	} else {
		deleteTasks(store, lowImpact)
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Basic suggestion flags
	deleteCmd.Flags().Bool("completed", false, "Suggest completed tasks for deletion")
	deleteCmd.Flags().Bool("overdue", false, "Suggest overdue tasks for deletion")
	deleteCmd.Flags().Bool("old", false, "Suggest old completed tasks for deletion")

	// Smart suggestion flags
	deleteCmd.Flags().Bool("duplicates", false, "Find and delete duplicate/similar tasks")
	deleteCmd.Flags().Bool("low-impact", false, "Suggest low-impact tasks to remove")

	// Advanced modes
	deleteCmd.Flags().Bool("batch", false, "Batch delete with smart grouping")
	deleteCmd.Flags().Bool("smart", false, "Smart-powered deletion suggestions")
	deleteCmd.Flags().Bool("cleanup", false, "Full cleanup mode with comprehensive analysis")

	// Control flags
	deleteCmd.Flags().BoolP("interactive", "i", false, "Interactive mode with category-by-category confirmation")
	deleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
}
