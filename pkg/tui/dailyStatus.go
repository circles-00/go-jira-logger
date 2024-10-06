package tui

import (
	"fmt"
	"strings"

	"go_jira_logger/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DailyStatusState struct {
	cursor int
}

func (m model) DailyStatusUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.state.dailyStatus.cursor > 0 {
				m.tickets[m.state.dailyStatus.cursor].Worklog.Blur()
				m.state.dailyStatus.cursor--
				m.tickets[m.state.dailyStatus.cursor].Worklog.Focus()
			}
		case "down":
			if m.state.dailyStatus.cursor < len(m.tickets)-1 {
				m.tickets[m.state.dailyStatus.cursor].Worklog.Blur()
				m.state.dailyStatus.cursor++
				m.tickets[m.state.dailyStatus.cursor].Worklog.Focus()
			}
		case "backspace":
			m.tickets[m.state.dailyStatus.cursor].Worklog, cmd = m.tickets[m.state.dailyStatus.cursor].Worklog.Update(msg)
			return m, cmd
		default:
			m.tickets[m.state.dailyStatus.cursor].Worklog.Focus()
		}
	}

	worklogValue := m.tickets[m.state.dailyStatus.cursor].Worklog
	if len(worklogValue.Value()) < 5 {
		m.tickets[m.state.dailyStatus.cursor].Worklog, cmd = m.tickets[m.state.dailyStatus.cursor].Worklog.Update(msg)
	}
	return m, cmd
}

func (m model) DailyStatusView() string {
	// Header and cell styles
	ticketTypeStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Width(10)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Width(30)
	cellStyle := lipgloss.NewStyle().Width(30)
	selectedStyle := cellStyle.Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#c1c9d6")).Bold(true)

	// Separator style
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(strings.Repeat("-", 120))

	// Table header
	header := fmt.Sprintf("%s\t%s\t%s\t%s",
		ticketTypeStyle.Render("Type"),
		headerStyle.Render("Ticket URL"),
		headerStyle.Render("Status"),
		headerStyle.Render("Worklog"),
	)

	// Table rows with separators
	var rows string
	for i, ticket := range m.tickets {
		baseUrl := config.GetBoardUrl()
		url := fmt.Sprintf("%s/browse/%s", baseUrl, ticket.Key)

		urlCell := cellStyle.Render(fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, ticket.Key))
		statusCell := cellStyle.Render(ticket.Fields.Status.Name)
		worklogCell := ticket.Worklog.View()

		typeIcon := "✅"

		if ticket.Fields.Issuetype.Name == "Bug" {
			typeIcon = "🐜"
		}
		// Highlight current row if selected
		if i == m.state.dailyStatus.cursor {
			worklogCell = selectedStyle.Render(ticket.Worklog.View())
		}

		rows += fmt.Sprintf("%s\t%s\t%s\t%s\n%s\n", ticketTypeStyle.Render(typeIcon), urlCell, statusCell, worklogCell, separator)
	}

	// Instructions for the user
	instructions := "\nUse 'up'/'down' to select, 'enter' to submit worklog, 'q' to quit."

	return header + "\n" + separator + "\n" + rows + instructions
}
