/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"todo/taskdata"

	"github.com/spf13/cobra"
)

// markCmd represents the mark command
var markCmd = &cobra.Command{
	Use:   "mark [task_id_or_name] [flags]",
	Short: "Smart task management - mark, edit, and analyze tasks",
	Long: `Intelligently mark tasks as done/undone and edit task properties with smart suggestions.

Features:
- Mark tasks as completed/incomplete by ID or name
- Smart completion suggestions based on due dates and priority
- Edit task properties: description, due date, priority
- Auto-detect overdue tasks and suggest actions
- Batch operations with smart filtering
- Integration with delete command for cleanup suggestions

Examples:
  todo mark                      # Smart suggestions for completion
  todo mark 5                    # Mark task #5 as done
  todo mark 5 --undone           # Mark task #5 as not done
  todo mark "buy groceries"      # Mark task by name as done
  todo mark --overdue            # Show overdue tasks for action
  todo mark --smart              # Smart-powered task analysis
  todo mark 5 --edit             # Edit task #5 properties
  todo mark 5 --due "2025-07-20" # Change due date
  todo mark 5 --priority high    # Change priority
  todo mark 5 --desc "New desc"  # Change description
  todo mark --batch              # Batch mark multiple tasks
  todo mark --cleanup            # Mark and suggest cleanup`,
	Run: markRun,
}

func markRun(cmd *cobra.Command, args []string) {
	// Load tasks
	store, err := taskdata.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	if len(store.Tasks) == 0 {
		fmt.Println("No tasks found. Use 'todo add \"task description\"' to add a task.")
		return
	}

	// Get flags
	undone, _ := cmd.Flags().GetBool("undone")
	showOverdue, _ := cmd.Flags().GetBool("overdue")
	showSmart, _ := cmd.Flags().GetBool("smart")
	editMode, _ := cmd.Flags().GetBool("edit")
	batchMode, _ := cmd.Flags().GetBool("batch")
	cleanupMode, _ := cmd.Flags().GetBool("cleanup")
	newDue, _ := cmd.Flags().GetString("due")
	newPriority, _ := cmd.Flags().GetString("priority")
	newDesc, _ := cmd.Flags().GetString("desc")
	force, _ := cmd.Flags().GetBool("force")

	// Smart mode - smart-powered analysis
	if showSmart {
		showSmartTaskAnalysis(store)
		return
	}

	// Overdue mode - show overdue tasks for action
	if showOverdue {
		showOverdueTaskActions(store)
		return
	}

	// Batch mode - mark multiple tasks
	if batchMode {
		performBatchOperations(store, undone, force)
		return
	}

	// Cleanup mode - mark and suggest cleanup
	if cleanupMode {
		performCleanupOperations(store)
		return
	}

	// No arguments - show smart suggestions
	if len(args) == 0 {
		showMarkSuggestions(store)
		return
	}

	// Handle task marking/editing by ID or name
	identifier := args[0]

	// Check if we're editing properties
	if editMode || newDue != "" || newPriority != "" || newDesc != "" {
		editTaskProperties(store, identifier, newDue, newPriority, newDesc)
		return
	}

	// Mark task as done/undone
	markTask(store, identifier, undone, force)
}

func showSmartTaskAnalysis(store *taskdata.TaskStore) {
	fmt.Println("🧠 Smart Task Analysis")
	fmt.Println(strings.Repeat("=", 50))

	now := time.Now()

	// Analyze task patterns
	analyzeTaskPatterns(store, now)

	// Productivity insights
	showProductivityInsights(store, now)

	// Completion recommendations
	showCompletionRecommendations(store, now)

	// Integration with delete suggestions
	suggestCleanupTasks(store, now)
}

func analyzeTaskPatterns(store *taskdata.TaskStore, now time.Time) {
	fmt.Printf("\n📊 Task Pattern Analysis\n")
	fmt.Println(strings.Repeat("-", 30))

	totalTasks := len(store.Tasks)
	completedTasks := 0
	overdueTasks := 0
	todayTasks := 0
	highPriorityPending := 0

	for _, task := range store.Tasks {
		if task.Completed {
			completedTasks++
		} else {
			if task.Priority == "high" {
				highPriorityPending++
			}
			if isTaskOverdue(task, now) {
				overdueTasks++
			}
			if isTaskDueToday(task, now) {
				todayTasks++
			}
		}
	}

	// Calculate completion rate
	completionRate := float64(completedTasks) / float64(totalTasks) * 100

	fmt.Printf("📈 Completion Rate: %.1f%% (%d/%d)\n", completionRate, completedTasks, totalTasks)

	if overdueTasks > 0 {
		fmt.Printf("⚠️  Overdue Tasks: %d (needs immediate attention)\n", overdueTasks)
	}

	if todayTasks > 0 {
		fmt.Printf("🎯 Due Today: %d tasks\n", todayTasks)
	}

	if highPriorityPending > 0 {
		fmt.Printf("🔴 High Priority Pending: %d tasks\n", highPriorityPending)
	}

	// Task health score
	healthScore := calculateTaskHealthScore(totalTasks, completedTasks, overdueTasks, highPriorityPending)
	fmt.Printf("💚 Task Health Score: %d/100\n", healthScore)
}

func showProductivityInsights(store *taskdata.TaskStore, now time.Time) {
	fmt.Printf("\n💡 Productivity Insights\n")
	fmt.Println(strings.Repeat("-", 30))

	// Analyze task distribution
	priorityDist := analyzePriorityDistribution(store)
	fmt.Printf("Priority Distribution: High:%d, Normal:%d, Low:%d\n",
		priorityDist["high"], priorityDist["normal"], priorityDist["low"])

	// Time-based insights
	upcomingDeadlines := getUpcomingDeadlines(store, now)
	if len(upcomingDeadlines) > 0 {
		fmt.Printf("📅 Upcoming Deadlines (%d tasks in next 7 days)\n", len(upcomingDeadlines))
	}

	// Suggest optimal focus
	suggestOptimalFocus(store, now)
}

func showCompletionRecommendations(store *taskdata.TaskStore, now time.Time) {
	fmt.Printf("\n🎯 Completion Recommendations\n")
	fmt.Println(strings.Repeat("-", 30))

	// Quick wins (easy completions)
	quickWins := getQuickWinTasks(store)
	if len(quickWins) > 0 {
		fmt.Printf("⚡ Quick Wins (%d tasks):\n", len(quickWins))
		for _, task := range quickWins[:min(3, len(quickWins))] {
			fmt.Printf("  • #%d: %s\n", task.ID, task.Description)
		}
	}

	// High-impact tasks
	highImpact := getHighImpactTasks(store, now)
	if len(highImpact) > 0 {
		fmt.Printf("🎯 High Impact (%d tasks):\n", len(highImpact))
		for _, task := range highImpact[:min(3, len(highImpact))] {
			fmt.Printf("  • #%d: %s\n", task.ID, task.Description)
		}
	}

	// Overdue recovery
	overdue := getOverdueTasksForRecovery(store, now)
	if len(overdue) > 0 {
		fmt.Printf("🚨 Overdue Recovery (%d tasks):\n", len(overdue))
		for _, task := range overdue[:min(3, len(overdue))] {
			fmt.Printf("  • #%d: %s (due %s)\n", task.ID, task.Description, task.DueDate)
		}
	}
}

func suggestCleanupTasks(store *taskdata.TaskStore, now time.Time) {
	fmt.Printf("\n🧹 Cleanup Integration\n")
	fmt.Println(strings.Repeat("-", 30))

	// Old completed tasks
	oldCompleted := getOldCompletedTasksForCleanup(store)
	if len(oldCompleted) > 0 {
		fmt.Printf("🗑️  Consider deleting %d old completed tasks\n", len(oldCompleted))
		fmt.Printf("   Run: todo delete --completed\n")
	}

	// Stale tasks
	staleTasks := getStaleTasksForReview(store, now)
	if len(staleTasks) > 0 {
		fmt.Printf("📋 Review %d stale tasks (no due date, low priority)\n", len(staleTasks))
		fmt.Printf("   Run: todo delete --old\n")
	}
}

func showOverdueTaskActions(store *taskdata.TaskStore) {
	fmt.Println("⚠️  Overdue Task Actions")
	fmt.Println(strings.Repeat("=", 50))

	now := time.Now()
	overdueTasks := getOverdueTasksForRecovery(store, now)

	if len(overdueTasks) == 0 {
		fmt.Println("🎉 Great! No overdue tasks found.")
		return
	}

	fmt.Printf("Found %d overdue task(s):\n\n", len(overdueTasks))

	for i, task := range overdueTasks {
		fmt.Printf("%d. #%d: %s\n", i+1, task.ID, task.Description)
		fmt.Printf("   Due: %s (overdue by %s)\n", task.DueDate, getOverdueDuration(task, now))
		fmt.Printf("   Priority: %s\n", task.Priority)
		fmt.Println()
	}

	// Suggest actions
	fmt.Println("🎯 Suggested Actions:")
	fmt.Println("1. Complete overdue tasks immediately")
	fmt.Println("2. Reschedule to realistic dates")
	fmt.Println("3. Mark as done if already completed")
	fmt.Println("4. Delete if no longer relevant")

	if confirmAction("Would you like to take action on overdue tasks?") {
		handleOverdueTaskActions(store, overdueTasks)
	}
}

func performBatchOperations(store *taskdata.TaskStore, undone, force bool) {
	fmt.Println("📦 Batch Task Operations")
	fmt.Println(strings.Repeat("=", 50))

	pendingTasks := getPendingTasks(store)
	if len(pendingTasks) == 0 {
		fmt.Println("No pending tasks found.")
		return
	}

	fmt.Printf("Found %d pending task(s):\n\n", len(pendingTasks))

	// Show tasks for selection
	for i, task := range pendingTasks {
		status := "🔲"
		if task.Completed {
			status = "✅"
		}
		fmt.Printf("%d. %s #%d: %s", i+1, status, task.ID, task.Description)
		if task.DueDate != "" {
			fmt.Printf(" (due: %s)", task.DueDate)
		}
		fmt.Println()
	}

	if !force {
		fmt.Print("\nEnter task numbers to mark (comma-separated, or 'all'): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "all" {
			markAllTasks(store, pendingTasks, undone)
		} else {
			markSelectedTasks(store, pendingTasks, input, undone)
		}
	} else {
		markAllTasks(store, pendingTasks, undone)
	}
}

func performCleanupOperations(store *taskdata.TaskStore) {
	fmt.Println("🧹 Cleanup Operations")
	fmt.Println(strings.Repeat("=", 50))

	// Mark obvious completions
	autoMarkObviousCompletions(store)

	// Suggest cleanup
	fmt.Println("\n🗑️  Cleanup Suggestions:")
	fmt.Println("Run 'todo delete --smart' for intelligent cleanup options")
}

func showMarkSuggestions(store *taskdata.TaskStore) {
	fmt.Println("🎯 Smart Mark Suggestions")
	fmt.Println(strings.Repeat("=", 50))

	now := time.Now()

	// Today's tasks
	todayTasks := getTodayTasksForCompletion(store, now)
	if len(todayTasks) > 0 {
		fmt.Printf("\n📅 Due Today (%d tasks):\n", len(todayTasks))
		for _, task := range todayTasks {
			displayTaskForMarking(task)
		}
	}

	// High priority tasks
	highPriorityTasks := getHighPriorityPendingTasks(store)
	if len(highPriorityTasks) > 0 {
		fmt.Printf("\n🔴 High Priority (%d tasks):\n", len(highPriorityTasks))
		for _, task := range highPriorityTasks[:min(3, len(highPriorityTasks))] {
			displayTaskForMarking(task)
		}
	}

	// Quick wins
	quickWins := getQuickWinTasks(store)
	if len(quickWins) > 0 {
		fmt.Printf("\n⚡ Quick Wins (%d tasks):\n", len(quickWins))
		for _, task := range quickWins[:min(3, len(quickWins))] {
			displayTaskForMarking(task)
		}
	}

	// Overdue tasks
	overdue := getOverdueTasksForRecovery(store, now)
	if len(overdue) > 0 {
		fmt.Printf("\n⚠️  Overdue (%d tasks):\n", len(overdue))
		for _, task := range overdue[:min(3, len(overdue))] {
			displayTaskForMarking(task)
		}
	}

	fmt.Printf("\n💡 Use 'todo mark <id>' to mark tasks as complete\n")
	fmt.Printf("💡 Use 'todo mark --smart' for detailed analysis\n")
}

func editTaskProperties(store *taskdata.TaskStore, identifier, newDue, newPriority, newDesc string) {
	task := findTaskByIDOrName(store, identifier)
	if task == nil {
		fmt.Printf("❌ Task not found: %s\n", identifier)
		return
	}

	fmt.Printf("📝 Editing Task #%d: %s\n", task.ID, task.Description)
	fmt.Println(strings.Repeat("=", 40))

	changes := make(map[string]string)
	updated := false

	// Update due date
	if newDue != "" {
		if err := taskdata.ValidateDate(newDue); err != nil {
			fmt.Printf("❌ Invalid due date: %v\n", err)
			return
		}
		changes["Due Date"] = fmt.Sprintf("%s → %s", task.DueDate, newDue)
		updateTaskDueDate(store, task.ID, newDue)
		updated = true
	}

	// Update priority
	if newPriority != "" {
		if err := taskdata.ValidatePriority(newPriority); err != nil {
			fmt.Printf("❌ Invalid priority: %v\n", err)
			return
		}
		changes["Priority"] = fmt.Sprintf("%s → %s", task.Priority, strings.ToLower(newPriority))
		updateTaskPriority(store, task.ID, strings.ToLower(newPriority))
		updated = true
	}

	// Update description
	if newDesc != "" {
		changes["Description"] = fmt.Sprintf("%s → %s", task.Description, newDesc)
		updateTaskDescription(store, task.ID, newDesc)
		updated = true
	}

	if updated {
		if err := store.SaveTasks(); err != nil {
			fmt.Printf("❌ Error saving changes: %v\n", err)
			return
		}

		fmt.Println("✅ Task updated successfully!")
		for field, change := range changes {
			fmt.Printf("  %s: %s\n", field, change)
		}
	} else {
		fmt.Println("ℹ️  No changes specified. Use --due, --priority, or --desc flags to edit.")
	}
}

func markTask(store *taskdata.TaskStore, identifier string, undone, force bool) {
	task := findTaskByIDOrName(store, identifier)
	if task == nil {
		fmt.Printf("❌ Task not found: %s\n", identifier)
		return
	}

	action := "complete"
	if undone {
		action = "mark as incomplete"
	}

	if !force {
		if !confirmAction(fmt.Sprintf("%s task #%d: %s", strings.Title(action), task.ID, task.Description)) {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	// Update task status
	if err := updateTaskCompletion(store, task.ID, !undone); err != nil {
		fmt.Printf("❌ Error updating task: %v\n", err)
		return
	}

	if err := store.SaveTasks(); err != nil {
		fmt.Printf("❌ Error saving changes: %v\n", err)
		return
	}

	status := "✅ completed"
	if undone {
		status = "🔲 marked as incomplete"
	}

	fmt.Printf("✅ Task #%d %s: %s\n", task.ID, status, task.Description)

	// Show smart suggestions after marking
	if !undone {
		showPostCompletionSuggestions(store, task)
	}
}

// Helper functions

func isTaskOverdue(task taskdata.Task, now time.Time) bool {
	if task.Completed || task.DueDate == "" {
		return false
	}
	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return false
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	taskDate := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())
	return taskDate.Before(today)
}

func isTaskDueToday(task taskdata.Task, now time.Time) bool {
	if task.DueDate == "" {
		return false
	}
	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return false
	}
	return dueDate.Year() == now.Year() && dueDate.YearDay() == now.YearDay()
}

func calculateTaskHealthScore(total, completed, overdue, highPriorityPending int) int {
	if total == 0 {
		return 100
	}

	completionBonus := (completed * 100) / total
	overduePenalty := (overdue * 20)
	highPriorityPenalty := (highPriorityPending * 10)

	score := completionBonus - overduePenalty - highPriorityPenalty
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

func analyzePriorityDistribution(store *taskdata.TaskStore) map[string]int {
	dist := map[string]int{"high": 0, "normal": 0, "low": 0}
	for _, task := range store.Tasks {
		if !task.Completed {
			dist[task.Priority]++
		}
	}
	return dist
}

func getUpcomingDeadlines(store *taskdata.TaskStore, now time.Time) []taskdata.Task {
	var upcoming []taskdata.Task
	nextWeek := now.AddDate(0, 0, 7)

	for _, task := range store.Tasks {
		if !task.Completed && task.DueDate != "" {
			dueDate, err := time.Parse("2006-01-02", task.DueDate)
			if err == nil && dueDate.After(now) && dueDate.Before(nextWeek) {
				upcoming = append(upcoming, task)
			}
		}
	}
	return upcoming
}

func suggestOptimalFocus(store *taskdata.TaskStore, now time.Time) {
	overdue := getOverdueTasksForRecovery(store, now)
	today := getTodayTasksForCompletion(store, now)
	highPriority := getHighPriorityPendingTasks(store)

	if len(overdue) > 0 {
		fmt.Printf("🎯 Focus: Handle %d overdue task(s) first\n", len(overdue))
	} else if len(today) > 0 {
		fmt.Printf("🎯 Focus: Complete %d task(s) due today\n", len(today))
	} else if len(highPriority) > 0 {
		fmt.Printf("🎯 Focus: Work on %d high-priority task(s)\n", len(highPriority))
	} else {
		fmt.Println("🎯 Focus: Great job! Consider picking up some quick wins")
	}
}

func getQuickWinTasks(store *taskdata.TaskStore) []taskdata.Task {
	var quickWins []taskdata.Task
	for _, task := range store.Tasks {
		if !task.Completed && task.Priority == "low" && task.DueDate == "" {
			quickWins = append(quickWins, task)
		}
	}
	return quickWins
}

func getHighImpactTasks(store *taskdata.TaskStore, now time.Time) []taskdata.Task {
	var highImpact []taskdata.Task
	for _, task := range store.Tasks {
		if !task.Completed && task.Priority == "high" {
			highImpact = append(highImpact, task)
		}
	}
	return highImpact
}

func getOverdueTasksForRecovery(store *taskdata.TaskStore, now time.Time) []taskdata.Task {
	var overdue []taskdata.Task
	for _, task := range store.Tasks {
		if isTaskOverdue(task, now) {
			overdue = append(overdue, task)
		}
	}
	return overdue
}

func getOldCompletedTasksForCleanup(store *taskdata.TaskStore) []taskdata.Task {
	var oldCompleted []taskdata.Task
	for _, task := range store.Tasks {
		if task.Completed {
			oldCompleted = append(oldCompleted, task)
		}
	}
	return oldCompleted
}

func getStaleTasksForReview(store *taskdata.TaskStore, now time.Time) []taskdata.Task {
	var stale []taskdata.Task
	for _, task := range store.Tasks {
		if !task.Completed && task.DueDate == "" && task.Priority == "low" {
			stale = append(stale, task)
		}
	}
	return stale
}

func getOverdueDuration(task taskdata.Task, now time.Time) string {
	if task.DueDate == "" {
		return "unknown"
	}
	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return "unknown"
	}
	duration := now.Sub(dueDate)
	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

func confirmAction(message string) bool {
	fmt.Printf("❓ %s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func handleOverdueTaskActions(store *taskdata.TaskStore, overdueTasks []taskdata.Task) {
	fmt.Println("\n🎯 Taking action on overdue tasks...")

	for _, task := range overdueTasks {
		fmt.Printf("\nTask #%d: %s (due %s)\n", task.ID, task.Description, task.DueDate)
		fmt.Println("Actions: (c)omplete, (r)eschedule, (d)elete, (s)kip")
		fmt.Print("Choose action: ")

		reader := bufio.NewReader(os.Stdin)
		action, _ := reader.ReadString('\n')
		action = strings.TrimSpace(strings.ToLower(action))

		switch action {
		case "c", "complete":
			updateTaskCompletion(store, task.ID, true)
			fmt.Printf("✅ Marked task #%d as completed\n", task.ID)
		case "r", "reschedule":
			fmt.Print("New due date (YYYY-MM-DD): ")
			newDate, _ := reader.ReadString('\n')
			newDate = strings.TrimSpace(newDate)
			if err := taskdata.ValidateDate(newDate); err == nil {
				updateTaskDueDate(store, task.ID, newDate)
				fmt.Printf("📅 Rescheduled task #%d to %s\n", task.ID, newDate)
			} else {
				fmt.Printf("❌ Invalid date format\n")
			}
		case "d", "delete":
			if confirmAction(fmt.Sprintf("Delete task #%d", task.ID)) {
				// Remove task from slice
				for i, t := range store.Tasks {
					if t.ID == task.ID {
						store.Tasks = append(store.Tasks[:i], store.Tasks[i+1:]...)
						break
					}
				}
				fmt.Printf("🗑️  Deleted task #%d\n", task.ID)
			}
		default:
			fmt.Printf("⏭️  Skipped task #%d\n", task.ID)
		}
	}

	store.SaveTasks()
}

func getPendingTasks(store *taskdata.TaskStore) []taskdata.Task {
	var pending []taskdata.Task
	for _, task := range store.Tasks {
		if !task.Completed {
			pending = append(pending, task)
		}
	}
	return pending
}

func markAllTasks(store *taskdata.TaskStore, tasks []taskdata.Task, undone bool) {
	count := 0
	for _, task := range tasks {
		if err := updateTaskCompletion(store, task.ID, !undone); err == nil {
			count++
		}
	}

	store.SaveTasks()
	action := "completed"
	if undone {
		action = "marked as incomplete"
	}
	fmt.Printf("✅ %s %d task(s)\n", strings.Title(action), count)
}

func markSelectedTasks(store *taskdata.TaskStore, tasks []taskdata.Task, input string, undone bool) {
	selections := strings.Split(input, ",")
	count := 0

	for _, sel := range selections {
		sel = strings.TrimSpace(sel)
		if idx, err := strconv.Atoi(sel); err == nil && idx > 0 && idx <= len(tasks) {
			task := tasks[idx-1]
			if err := updateTaskCompletion(store, task.ID, !undone); err == nil {
				count++
			}
		}
	}

	store.SaveTasks()
	action := "completed"
	if undone {
		action = "marked as incomplete"
	}
	fmt.Printf("✅ %s %d task(s)\n", strings.Title(action), count)
}

func autoMarkObviousCompletions(store *taskdata.TaskStore) {
	// This could implement ML-based suggestions in the future
	fmt.Println("🤖 Scanning for obvious completions...")
	fmt.Println("   (No obvious completions detected)")
}

func getTodayTasksForCompletion(store *taskdata.TaskStore, now time.Time) []taskdata.Task {
	var today []taskdata.Task
	for _, task := range store.Tasks {
		if !task.Completed && isTaskDueToday(task, now) {
			today = append(today, task)
		}
	}
	return today
}

func getHighPriorityPendingTasks(store *taskdata.TaskStore) []taskdata.Task {
	var highPriority []taskdata.Task
	for _, task := range store.Tasks {
		if !task.Completed && task.Priority == "high" {
			highPriority = append(highPriority, task)
		}
	}
	return highPriority
}

func displayTaskForMarking(task taskdata.Task) {
	priorityIcon := ""
	switch task.Priority {
	case "high":
		priorityIcon = "🔴"
	case "normal":
		priorityIcon = "🟡"
	case "low":
		priorityIcon = "🟢"
	}

	dueDateStr := ""
	if task.DueDate != "" {
		dueDateStr = fmt.Sprintf(" (due: %s)", task.DueDate)
	}

	fmt.Printf("  %s #%d: %s%s\n", priorityIcon, task.ID, task.Description, dueDateStr)
}

func findTaskByIDOrName(store *taskdata.TaskStore, identifier string) *taskdata.Task {
	// Try to parse as ID first
	if id, err := strconv.Atoi(identifier); err == nil {
		for _, task := range store.Tasks {
			if task.ID == id {
				return &task
			}
		}
	}

	// Search by name (partial match)
	searchTerm := strings.ToLower(identifier)
	for _, task := range store.Tasks {
		if strings.Contains(strings.ToLower(task.Description), searchTerm) {
			return &task
		}
	}

	return nil
}

func updateTaskDueDate(store *taskdata.TaskStore, id int, newDue string) error {
	for i, task := range store.Tasks {
		if task.ID == id {
			store.Tasks[i].DueDate = newDue
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func updateTaskPriority(store *taskdata.TaskStore, id int, newPriority string) error {
	for i, task := range store.Tasks {
		if task.ID == id {
			store.Tasks[i].Priority = newPriority
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func updateTaskDescription(store *taskdata.TaskStore, id int, newDesc string) error {
	for i, task := range store.Tasks {
		if task.ID == id {
			store.Tasks[i].Description = newDesc
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func updateTaskCompletion(store *taskdata.TaskStore, id int, completed bool) error {
	for i, task := range store.Tasks {
		if task.ID == id {
			store.Tasks[i].Completed = completed
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func showPostCompletionSuggestions(store *taskdata.TaskStore, completedTask *taskdata.Task) {
	fmt.Printf("\n🎉 Great job completing: %s\n", completedTask.Description)

	// Check for related tasks or next actions
	now := time.Now()
	pending := getPendingTasks(store)

	if len(pending) > 0 {
		fmt.Printf("💡 Next suggestions:\n")

		// Show high priority tasks
		highPriority := getHighPriorityPendingTasks(store)
		if len(highPriority) > 0 {
			fmt.Printf("   🔴 High priority: %s\n", highPriority[0].Description)
		}

		// Show due today
		today := getTodayTasksForCompletion(store, now)
		if len(today) > 0 {
			fmt.Printf("   📅 Due today: %s\n", today[0].Description)
		}

		// Suggest cleanup if many completed
		completed := 0
		for _, task := range store.Tasks {
			if task.Completed {
				completed++
			}
		}

		if completed > 5 {
			fmt.Printf("   🧹 Consider running 'todo delete --completed' to clean up\n")
		}
	}
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Core operation flags
	markCmd.Flags().BoolP("undone", "u", false, "Mark task as incomplete/undone")
	markCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompts")

	// Analysis and view flags
	markCmd.Flags().Bool("overdue", false, "Show overdue tasks for action")
	markCmd.Flags().BoolP("smart", "s", false, "Smart-powered task analysis and suggestions")
	markCmd.Flags().Bool("batch", false, "Batch mark multiple tasks")
	markCmd.Flags().Bool("cleanup", false, "Mark and suggest cleanup operations")

	// Edit flags
	markCmd.Flags().BoolP("edit", "e", false, "Edit task properties")
	markCmd.Flags().String("due", "", "Change due date (YYYY-MM-DD)")
	markCmd.Flags().StringP("priority", "p", "", "Change priority (low, normal, high)")
	markCmd.Flags().StringP("desc", "d", "", "Change task description")
}
