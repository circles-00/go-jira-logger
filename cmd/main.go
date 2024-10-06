package main

import (
	"fmt"

	"go_jira_logger/pkg/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model, err := tui.NewModel()
	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
