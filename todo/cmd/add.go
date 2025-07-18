/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"todo/taskdata"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new task to your todo list",
	Long: `The add command allows you to create a new task in your todo list.
You can specify the task description, and optionally set a due date or priority level. 
By default, the task will be added with no due date or priority.`,
	Run: addRun,
}

func addRun(cmd *cobra.Command, args []string) {
	// Check if no arguments provided
	if len(args) == 0 {
		fmt.Println("Please provide a task description.")
		return
	}

	// Get flags
	dueDate, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetString("priority")

	// Load existing tasks
	store, err := taskdata.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	// Handle multiple tasks
	if len(args) > 1 {
		fmt.Println("Multiple tasks detected. Adding each task separately:")
	}

	// Add each task
	successCount := 0
	for _, taskDesc := range args {
		if taskDesc == "" {
			fmt.Println("Skipping empty task description.")
			continue
		}

		// Add task to store with validation
		task, err := store.AddTask(taskDesc, dueDate, priority)
		if err != nil {
			fmt.Printf("Error adding task '%s': %v\n", taskDesc, err)
			continue
		}

		// Display success message
		fmt.Printf("✓ Added task #%d: %s\n", task.ID, task.Description)
		if task.DueDate != "" {
			fmt.Printf("  Due date: %s\n", task.DueDate)
		}
		fmt.Printf("  Priority: %s\n", task.Priority)
		fmt.Printf("  Status: %s\n", func() string {
			if task.Completed {
				return "Completed"
			}
			return "Pending"
		}())
		fmt.Println()

		successCount++
	}

	// Save tasks to file if any were added successfully
	if successCount > 0 {
		if err := store.SaveTasks(); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			return
		}
		fmt.Printf("Successfully added %d task(s) and saved to file.\n", successCount)
	}
} // Add the add command to the root command
// This allows the add command to be executed as a subcommand of the main todo command
// You can also define persistent flags that will work for this command and all subcommands
// (Persistent flag definition moved to init())

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.
	addCmd.Flags().StringP("due", "d", "", "Due date for the task (format: YYYY-MM-DD)")
	addCmd.Flags().StringP("priority", "p", "normal", "Priority level of the task (low, normal, high)")
}
