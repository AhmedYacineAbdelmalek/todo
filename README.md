# 🚀 Smart Todo CLI

A powerful, intelligent command-line todo application built with Go and Cobra. Features smart-powered task management, intelligent suggestions, and seamless productivity insights.

## ✨ Features

### 🎯 **Core Task Management**
- **Add tasks** with due dates, priorities, and smart validation
- **List tasks** with advanced filtering (today, week, month, all)
- **Mark tasks** as complete/incomplete with intelligent suggestions
- **Delete tasks** with smart cleanup recommendations
- **Edit tasks** - modify descriptions, due dates, and priorities

### 🧠 **Smart Intelligence**
- **Smart-powered analysis** - task health scores and productivity insights
- **Smart suggestions** - optimal focus recommendations based on priorities and deadlines
- **Overdue detection** - automatic identification and recovery suggestions
- **Quick wins identification** - find easy tasks to build momentum
- **Cleanup automation** - intelligent suggestions for task maintenance

### 📊 **Advanced Filtering & Views**
- **Time filters**: today, week, month, all tasks
- **Priority filters**: high, normal, low (with shortcuts: h, n, l)
- **Status filters**: pending, completed, overdue, due soon
- **Smart views**: insights, statistics, productivity analysis
- **Untracked tasks**: always visible regardless of time filters

### 🎨 **Beautiful Output**
- **Rich icons** and colored output for better visualization
- **Progress tracking** with completion rates and health scores
- **Intuitive displays** with clear task organization
- **Smart formatting** for dates, priorities, and status

## 📦 Installation

### Quick Install (Recommended)
```bash
# One-line installer (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/AhmedYacineAbdelmalek/todo/main/install.sh | bash

# Or download and run
curl -L https://raw.githubusercontent.com/AhmedYacineAbdelmalek/todo/main/install.sh -o install.sh
chmod +x install.sh
./install.sh
```

### Manual Download
```bash
# Linux (x64)
curl -L https://github.com/AhmedYacineAbdelmalek/todo/releases/latest/download/todo-linux-amd64 -o todo
chmod +x todo
sudo mv todo /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/AhmedYacineAbdelmalek/todo/releases/latest/download/todo-darwin-amd64 -o todo
chmod +x todo
sudo mv todo /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/AhmedYacineAbdelmalek/todo/releases/latest/download/todo-darwin-arm64 -o todo
chmod +x todo
sudo mv todo /usr/local/bin/

# Windows
# Download todo-windows-amd64.exe from releases and add to PATH
```

### Build from Source
```bash
git clone https://github.com/AhmedYacineAbdelmalek/todo.git
cd todo/todo
go build -o todo
```

## 🚀 Quick Start

```bash
# Add your first task
todo add "Buy groceries" --due "2025-07-20" --priority high

# View today's tasks
todo list

# View all tasks
todo list -a

# Mark task as complete
todo mark 1

# Get smart suggestions
todo mark --smart

# Clean up completed tasks
todo delete --smart
```

## 📖 Usage Guide

### Adding Tasks
```bash
# Basic task
todo add "Complete project"

# Task with due date and priority
todo add "Meeting with client" --due "2025-07-25" --priority high

# Multiple tasks at once
todo add "Task 1" "Task 2" "Task 3" --priority normal
```

### Listing Tasks
```bash
# Today's tasks (default)
todo list

# This week's tasks
todo list -w

# All tasks
todo list -a

# High priority tasks this week
todo list -w -p h

# Overdue tasks only
todo list --overdue

# Tasks without due dates
todo list --no-date

# Smart insights
todo list --insights
```

### Managing Tasks
```bash
# Mark task as complete
todo mark 1

# Mark task as incomplete
todo mark 1 --undone

# Edit task properties
todo mark 1 --due "2025-07-30" --priority normal

# Batch mark multiple tasks
todo mark --batch

# Smart analysis
todo mark --smart
```

### Deleting Tasks
```bash
# Smart cleanup suggestions
todo delete

# Delete specific task
todo delete 1

# Delete by name
todo delete "groceries"

# Delete completed tasks
todo delete --completed

# Interactive cleanup
todo delete --interactive
```

## 🎯 Command Reference

### `todo add [tasks...] [flags]`
Add new tasks with smart validation.

**Flags:**
- `-d, --due string`: Due date (YYYY-MM-DD)
- `-p, --priority string`: Priority (low, normal, high)

### `todo list [flags]`
List tasks with advanced filtering and insights.

**Flags:**
- `-w, --week`: Show this week's tasks
- `-m, --month`: Show this month's tasks  
- `-a, --all`: Show all tasks
- `-p, --priority string`: Filter by priority (l/low, n/normal, h/high)
- `--completed`: Show only completed tasks
- `--pending`: Show only pending tasks
- `--overdue`: Show only overdue tasks
- `--due-soon`: Show tasks due in next 3 days
- `--no-date`: Show tasks without due dates
- `-i, --insights`: Show productivity insights
- `-s, --smart`: Smart view with recommendations
- `--stats`: Show detailed statistics

### `todo mark [task_id_or_name] [flags]`
Mark tasks and edit properties with smart suggestions.

**Flags:**
- `-u, --undone`: Mark task as incomplete
- `-f, --force`: Skip confirmation prompts
- `--overdue`: Show overdue tasks for action
- `-s, --smart`: Smart task analysis
- `--batch`: Batch mark multiple tasks
- `--cleanup`: Mark and suggest cleanup
- `-e, --edit`: Edit task properties
- `--due string`: Change due date
- `-p, --priority string`: Change priority
- `-d, --desc string`: Change description

### `todo delete [task_id_or_name] [flags]`
Delete tasks with intelligent cleanup suggestions.

**Flags:**
- `--completed`: Suggest completed tasks for deletion
- `--overdue`: Suggest overdue tasks for deletion
- `--old`: Suggest old completed tasks for deletion
- `-i, --interactive`: Interactive mode
- `-f, --force`: Force deletion without confirmation
- `--smart`: Smart analysis
- `--pattern`: Pattern-based cleanup
- `--health`: Health-based suggestions

## 💡 Pro Tips

### Smart Workflows
```bash
# Morning routine
todo list                    # Check today's tasks
todo mark --smart           # Get smart recommendations

# End of day cleanup
todo mark --batch           # Complete finished tasks
todo delete --smart         # Clean up old tasks

# Weekly review
todo list --stats           # View productivity stats
todo list --insights        # Analyze patterns
```

### Power User Features
```bash
# Combine filters for precise results
todo list -w -p h --due-soon    # High priority, due soon, this week

# Use shortcuts for speed
todo list -w -p h               # Week view, high priority
todo mark 1 -f                  # Force mark without confirmation

# Smart editing
todo mark 1 --due "2025-08-01" --priority low --desc "Updated task"
```

## 🏗️ Architecture

```
todo/
├── cmd/                    # Cobra commands
│   ├── add.go             # Task creation
│   ├── list.go            # Task listing & filtering
│   ├── mark.go            # Task completion & editing
│   ├── delete.go          # Task deletion & cleanup
│   └── root.go            # Root command
├── taskdata/              # Data layer
│   └── task.go            # Task struct & storage
├── main.go                # Application entry point
└── go.mod                 # Go modules
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Inspired by modern productivity methodologies
- Designed for developers who love the command line

---

**Made by yasssin29 with ❤️ for productive developers**