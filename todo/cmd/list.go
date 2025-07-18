/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"todo/taskdata"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your tasks with smart filtering and insights",
	Long: `Display your todo tasks with various filtering options and smart insights.

By default, shows today's tasks. You can filter by:
- Time period: today (default), week, month, all
- Priority: low, normal, high (can use first letter: l, n, h)
- Completion status: pending, completed, all
- Smart filters: overdue, due-soon, no-date, productivity insights

Examples:
  todo list                      # Show today's tasks with insights
  todo list -w                   # Show this week's tasks  
  todo list -m                   # Show this month's tasks
  todo list -a                   # Show all tasks
  todo list -w -p h              # Show high priority tasks this week
  todo list -p high --all        # Show all high priority tasks
  todo list --completed          # Show completed tasks
  todo list --overdue            # Show only overdue tasks
  todo list --due-soon           # Show tasks due in next 3 days
  todo list --no-date            # Show tasks without due dates
  todo list --insights           # Show productivity insights
  todo list --smart              # Smart view with recommendations`,
	Run: listRun,
}

func listRun(cmd *cobra.Command, args []string) {
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

	// Get filter flags
	showWeek, _ := cmd.Flags().GetBool("week")
	showMonth, _ := cmd.Flags().GetBool("month")
	showAll, _ := cmd.Flags().GetBool("all")
	priority, _ := cmd.Flags().GetString("priority")
	showCompleted, _ := cmd.Flags().GetBool("completed")
	showPending, _ := cmd.Flags().GetBool("pending")
	showOverdue, _ := cmd.Flags().GetBool("overdue")
	showDueSoon, _ := cmd.Flags().GetBool("due-soon")
	showNoDate, _ := cmd.Flags().GetBool("no-date")
	showInsights, _ := cmd.Flags().GetBool("insights")
	showSmart, _ := cmd.Flags().GetBool("smart")
	showStats, _ := cmd.Flags().GetBool("stats")

	// Smart view mode
	if showSmart {
		displaySmartView(store)
		return
	}

	// Insights mode
	if showInsights {
		displayInsights(store)
		return
	}

	// Statistics mode
	if showStats {
		displayStatistics(store)
		return
	}

	// Filter tasks
	filteredTasks := filterTasks(store.Tasks, filterOptions{
		timeFilter:    getTimeFilter(showWeek, showMonth, showAll),
		priority:      priority,
		showCompleted: showCompleted,
		showPending:   showPending,
		showOverdue:   showOverdue,
		showDueSoon:   showDueSoon,
		showNoDate:    showNoDate,
	})

	if len(filteredTasks) == 0 {
		fmt.Println("No tasks match the specified filters.")
		return
	}

	// Display tasks
	displayTasks(filteredTasks, getTimeFilter(showWeek, showMonth, showAll))

	// Show quick insights if not in specific filter mode
	if !showOverdue && !showDueSoon && !showNoDate && !showCompleted {
		showQuickInsights(store, filteredTasks)
	}
}

type filterOptions struct {
	timeFilter    string
	priority      string
	showCompleted bool
	showPending   bool
	showOverdue   bool
	showDueSoon   bool
	showNoDate    bool
}

func getTimeFilter(week, month, all bool) string {
	if week {
		return "week"
	}
	if month {
		return "month"
	}
	if all {
		return "all"
	}
	return "today"
}

func filterTasks(tasks []taskdata.Task, opts filterOptions) []taskdata.Task {
	var filtered []taskdata.Task
	now := time.Now()

	for _, task := range tasks {
		// Special filters first
		if opts.showOverdue && !isOverdue(task, now) {
			continue
		}
		if opts.showDueSoon && !isDueSoon(task, now) {
			continue
		}
		if opts.showNoDate && task.DueDate != "" {
			continue
		}

		// Time filter (only apply if no special filters)
		if !opts.showOverdue && !opts.showDueSoon && !opts.showNoDate {
			if !matchesTimeFilter(task, opts.timeFilter, now) {
				continue
			}
		}

		// Priority filter
		if opts.priority != "" && !matchesPriority(task, opts.priority) {
			continue
		}

		// Completion status filter
		if opts.showCompleted && !opts.showPending {
			if !task.Completed {
				continue
			}
		} else if opts.showPending && !opts.showCompleted {
			if task.Completed {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	// Sort by priority and due date
	sort.Slice(filtered, func(i, j int) bool {
		// First sort by completion status (pending first)
		if filtered[i].Completed != filtered[j].Completed {
			return !filtered[i].Completed
		}

		// Then by priority (high > normal > low)
		priorityOrder := map[string]int{"high": 3, "normal": 2, "low": 1}
		if priorityOrder[filtered[i].Priority] != priorityOrder[filtered[j].Priority] {
			return priorityOrder[filtered[i].Priority] > priorityOrder[filtered[j].Priority]
		}

		// Finally by due date (earlier dates first)
		if filtered[i].DueDate != "" && filtered[j].DueDate != "" {
			return filtered[i].DueDate < filtered[j].DueDate
		}
		if filtered[i].DueDate != "" {
			return true
		}
		return false
	})

	return filtered
}

func matchesTimeFilter(task taskdata.Task, timeFilter string, now time.Time) bool {
	if timeFilter == "all" {
		return true
	}

	if task.DueDate == "" {
		// Tasks without due dates (untracked tasks) are always shown
		return true
	}

	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return false
	}

	switch timeFilter {
	case "today":
		return isSameDay(dueDate, now)
	case "week":
		return isInWeekRange(dueDate, now)
	case "month":
		return isSameMonth(dueDate, now)
	default:
		return true
	}
}

func matchesPriority(task taskdata.Task, priority string) bool {
	priority = strings.ToLower(priority)

	// Handle single letter shortcuts
	switch priority {
	case "h":
		priority = "high"
	case "n":
		priority = "normal"
	case "l":
		priority = "low"
	}

	return task.Priority == priority
}

func isOverdue(task taskdata.Task, now time.Time) bool {
	if task.Completed || task.DueDate == "" {
		return false
	}

	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return false
	}

	// A task is overdue only if due date is before today (not including today)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	taskDate := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())

	return taskDate.Before(today)
}

func isDueSoon(task taskdata.Task, now time.Time) bool {
	if task.Completed || task.DueDate == "" {
		return false
	}

	dueDate, err := time.Parse("2006-01-02", task.DueDate)
	if err != nil {
		return false
	}

	// Due soon = within next 3 days
	threeDaysFromNow := now.AddDate(0, 0, 3)
	return !dueDate.Before(now) && !dueDate.After(threeDaysFromNow)
}

func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isInWeekRange(date, now time.Time) bool {
	// Week range is today + 6 days (7 days total)
	startOfRange := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfRange := startOfRange.AddDate(0, 0, 6)
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return !normalizedDate.Before(startOfRange) && !normalizedDate.After(endOfRange)
}

func isSameMonth(date1, date2 time.Time) bool {
	y1, m1, _ := date1.Date()
	y2, m2, _ := date2.Date()
	return y1 == y2 && m1 == m2
}

func displayTasks(tasks []taskdata.Task, timeFilter string) {
	// Display header
	switch timeFilter {
	case "today":
		fmt.Printf("📅 Today's Tasks (%s)\n", time.Now().Format("2006-01-02"))
	case "week":
		fmt.Printf("📅 This Week's Tasks\n")
	case "month":
		fmt.Printf("📅 This Month's Tasks\n")
	case "all":
		fmt.Printf("📅 All Tasks\n")
	}
	fmt.Println(strings.Repeat("=", 50))

	// Group by completion status
	var pendingTasks, completedTasks []taskdata.Task
	for _, task := range tasks {
		if task.Completed {
			completedTasks = append(completedTasks, task)
		} else {
			pendingTasks = append(pendingTasks, task)
		}
	}

	// Display pending tasks
	if len(pendingTasks) > 0 {
		fmt.Printf("\n🔲 Pending Tasks (%d)\n", len(pendingTasks))
		fmt.Println(strings.Repeat("-", 30))
		for _, task := range pendingTasks {
			displayTask(task)
		}
	}

	// Display completed tasks
	if len(completedTasks) > 0 {
		fmt.Printf("\n✅ Completed Tasks (%d)\n", len(completedTasks))
		fmt.Println(strings.Repeat("-", 30))
		for _, task := range completedTasks {
			displayTask(task)
		}
	}

	fmt.Printf("\nTotal: %d tasks\n", len(tasks))
}

func displayTask(task taskdata.Task) {
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

	// Format due date
	dueDateStr := ""
	if task.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", task.DueDate)
		if err == nil {
			now := time.Now()
			if isSameDay(dueDate, now) {
				dueDateStr = " 📅 Today"
			} else if dueDate.Before(now) && !task.Completed {
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

func displaySmartView(store *taskdata.TaskStore) {
	fmt.Println("🧠 Smart Task View")
	fmt.Println(strings.Repeat("=", 50))

	now := time.Now()

	// Critical tasks (overdue + high priority)
	criticalTasks := getCriticalTasks(store.Tasks, now)
	if len(criticalTasks) > 0 {
		fmt.Printf("\n🚨 Critical Tasks (%d)\n", len(criticalTasks))
		fmt.Println(strings.Repeat("-", 30))
		for _, task := range criticalTasks {
			displayTask(task)
		}
	}

	// Today's focus
	todayTasks := getTodayTasks(store.Tasks, now)
	if len(todayTasks) > 0 {
		fmt.Printf("\n🎯 Today's Focus (%d)\n", len(todayTasks))
		fmt.Println(strings.Repeat("-", 30))
		for _, task := range todayTasks {
			displayTask(task)
		}
	}

	// Due soon
	dueSoonTasks := getDueSoonTasks(store.Tasks, now)
	if len(dueSoonTasks) > 0 {
		fmt.Printf("\n⏰ Due Soon (Next 3 Days) (%d)\n", len(dueSoonTasks))
		fmt.Println(strings.Repeat("-", 30))
		for _, task := range dueSoonTasks {
			displayTask(task)
		}
	}

	// Quick wins (low priority, easy tasks)
	quickWins := getQuickWins(store.Tasks)
	if len(quickWins) > 0 && len(quickWins) <= 3 {
		fmt.Printf("\n⚡ Quick Wins (%d)\n", len(quickWins))
		fmt.Println(strings.Repeat("-", 30))
		for _, task := range quickWins {
			displayTask(task)
		}
	}

	// Show recommendations
	showSmartRecommendations(store, now)
}

func getCriticalTasks(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var critical []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && (isOverdue(task, now) || (task.Priority == "high" && isDueSoon(task, now))) {
			critical = append(critical, task)
		}
	}
	return critical
}

func getTodayTasks(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var today []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && task.DueDate != "" {
			dueDate, err := time.Parse("2006-01-02", task.DueDate)
			if err == nil && isSameDay(dueDate, now) {
				today = append(today, task)
			}
		}
	}
	return today
}

func getDueSoonTasks(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var dueSoon []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && isDueSoon(task, now) {
			dueSoon = append(dueSoon, task)
		}
	}
	return dueSoon
}

func getQuickWins(tasks []taskdata.Task) []taskdata.Task {
	var quickWins []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && task.Priority == "low" && task.DueDate == "" {
			quickWins = append(quickWins, task)
		}
	}
	return quickWins
}

func showSmartRecommendations(store *taskdata.TaskStore, now time.Time) {
	fmt.Printf("\n💡 Smart Recommendations\n")
	fmt.Println(strings.Repeat("-", 30))

	overdueTasks := getOverdueTasksWithTime(store.Tasks, now)
	if len(overdueTasks) > 0 {
		fmt.Printf("• You have %d overdue task(s). Consider rescheduling or completing them.\n", len(overdueTasks))
	}

	noDateTasks := getNoDateTasks(store.Tasks)
	if len(noDateTasks) > 5 {
		fmt.Printf("• You have %d tasks without due dates. Consider adding dates for better planning.\n", len(noDateTasks))
	}

	highPriorityCount := getHighPriorityPendingCount(store.Tasks)
	if highPriorityCount > 3 {
		fmt.Printf("• You have %d high-priority tasks. Consider focusing on top 3 first.\n", highPriorityCount)
	}

	completedToday := getCompletedTodayCount(store.Tasks, now)
	if completedToday > 0 {
		fmt.Printf("• Great job! You've completed %d task(s) today! 🎉\n", completedToday)
	}
}

func getOverdueTasksWithTime(tasks []taskdata.Task, now time.Time) []taskdata.Task {
	var overdue []taskdata.Task
	for _, task := range tasks {
		if isOverdue(task, now) {
			overdue = append(overdue, task)
		}
	}
	return overdue
}

func getNoDateTasks(tasks []taskdata.Task) []taskdata.Task {
	var noDate []taskdata.Task
	for _, task := range tasks {
		if !task.Completed && task.DueDate == "" {
			noDate = append(noDate, task)
		}
	}
	return noDate
}

func getHighPriorityPendingCount(tasks []taskdata.Task) int {
	count := 0
	for _, task := range tasks {
		if !task.Completed && task.Priority == "high" {
			count++
		}
	}
	return count
}

func getCompletedTodayCount(tasks []taskdata.Task, now time.Time) int {
	// Since we don't track completion date, we'll return 0 for now
	// In a real implementation, you'd track when tasks were completed
	return 0
}

func displayInsights(store *taskdata.TaskStore) {
	fmt.Println("📊 Task Insights")
	fmt.Println(strings.Repeat("=", 50))

	now := time.Now()

	// Task breakdown
	total := len(store.Tasks)
	completed := 0
	pending := 0
	overdue := 0
	dueSoon := 0
	noDate := 0

	for _, task := range store.Tasks {
		if task.Completed {
			completed++
		} else {
			pending++
			if isOverdue(task, now) {
				overdue++
			} else if isDueSoon(task, now) {
				dueSoon++
			}
			if task.DueDate == "" {
				noDate++
			}
		}
	}

	// Display stats
	fmt.Printf("\n📈 Task Overview\n")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("Total Tasks: %d\n", total)
	fmt.Printf("Completed: %d (%.1f%%)\n", completed, float64(completed)/float64(total)*100)
	fmt.Printf("Pending: %d (%.1f%%)\n", pending, float64(pending)/float64(total)*100)

	if overdue > 0 {
		fmt.Printf("⚠️  Overdue: %d\n", overdue)
	}
	if dueSoon > 0 {
		fmt.Printf("⏰ Due Soon: %d\n", dueSoon)
	}
	if noDate > 0 {
		fmt.Printf("📝 No Due Date: %d\n", noDate)
	}

	// Priority breakdown
	fmt.Printf("\n🎯 Priority Breakdown\n")
	fmt.Println(strings.Repeat("-", 30))
	highCount, normalCount, lowCount := getPriorityBreakdown(store.Tasks)
	fmt.Printf("🔴 High: %d\n", highCount)
	fmt.Printf("🟡 Normal: %d\n", normalCount)
	fmt.Printf("🟢 Low: %d\n", lowCount)
}

func displayStatistics(store *taskdata.TaskStore) {
	fmt.Println("📊 Detailed Statistics")
	fmt.Println(strings.Repeat("=", 50))

	displayInsights(store)

	// Additional detailed stats
	fmt.Printf("\n📅 Time-based Analysis\n")
	fmt.Println(strings.Repeat("-", 30))

	now := time.Now()
	todayCount := 0
	thisWeekCount := 0
	thisMonthCount := 0

	for _, task := range store.Tasks {
		if !task.Completed && task.DueDate != "" {
			dueDate, err := time.Parse("2006-01-02", task.DueDate)
			if err == nil {
				if isSameDay(dueDate, now) {
					todayCount++
				}
				if isInWeekRange(dueDate, now) {
					thisWeekCount++
				}
				if isSameMonth(dueDate, now) {
					thisMonthCount++
				}
			}
		}
	}

	fmt.Printf("Due Today: %d\n", todayCount)
	fmt.Printf("Due This Week: %d\n", thisWeekCount)
	fmt.Printf("Due This Month: %d\n", thisMonthCount)
}

func getPriorityBreakdown(tasks []taskdata.Task) (int, int, int) {
	high, normal, low := 0, 0, 0
	for _, task := range tasks {
		if !task.Completed {
			switch task.Priority {
			case "high":
				high++
			case "normal":
				normal++
			case "low":
				low++
			}
		}
	}
	return high, normal, low
}

func showQuickInsights(store *taskdata.TaskStore, filteredTasks []taskdata.Task) {
	now := time.Now()

	// Show quick stats
	overdue := 0
	dueSoon := 0

	for _, task := range filteredTasks {
		if !task.Completed {
			if isOverdue(task, now) {
				overdue++
			} else if isDueSoon(task, now) {
				dueSoon++
			}
		}
	}

	if overdue > 0 || dueSoon > 0 {
		fmt.Printf("\n💡 Quick Insights: ")
		if overdue > 0 {
			fmt.Printf("%d overdue", overdue)
		}
		if overdue > 0 && dueSoon > 0 {
			fmt.Printf(", ")
		}
		if dueSoon > 0 {
			fmt.Printf("%d due soon", dueSoon)
		}
		fmt.Println()
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Time filter flags
	listCmd.Flags().BoolP("week", "w", false, "Show this week's tasks")
	listCmd.Flags().BoolP("month", "m", false, "Show this month's tasks")
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks")

	// Priority filter flag
	listCmd.Flags().StringP("priority", "p", "", "Filter by priority (low/l, normal/n, high/h)")

	// Completion status flags
	listCmd.Flags().Bool("completed", false, "Show only completed tasks")
	listCmd.Flags().Bool("pending", false, "Show only pending tasks")

	// Smart filter flags
	listCmd.Flags().Bool("overdue", false, "Show only overdue tasks")
	listCmd.Flags().Bool("due-soon", false, "Show tasks due in next 3 days")
	listCmd.Flags().Bool("no-date", false, "Show tasks without due dates")

	// Analysis flags
	listCmd.Flags().BoolP("insights", "i", false, "Show productivity insights")
	listCmd.Flags().BoolP("smart", "s", false, "Smart view with recommendations")
	listCmd.Flags().Bool("stats", false, "Show detailed statistics")
}
