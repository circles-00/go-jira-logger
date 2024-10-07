package tui

import (
	"fmt"
	"strings"
	"time"

	"go_jira_logger/pkg/api"
	"go_jira_logger/pkg/config"
	"go_jira_logger/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DailyStatusState struct {
	cursor     int
	timeLogged bool
}

func (m model) SubmitWorklog(issueKey string, worklogPayload api.WorklogPayload) tea.Cmd {
	cmd := func() tea.Msg {
		fmt.Sprintln("FUNC RUNNING", issueKey, worklogPayload.Started, worklogPayload.TimeSpent)
		// TODO: Change hardcoded data
		// tags := []string{"Execution", "Planning_&_Analyzing"}
		// return api.SubmitWorklog(issueKey, tags, worklogPayload)

		return api.AddWorklogResponse{}
	}

	return cmd
}

func (m model) DailyStatusUpdate(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case api.AddWorklogResponse:
		m.state.dailyStatus.timeLogged = true
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
		case "enter":
			var cmds []tea.Cmd
			for _, ticket := range m.tickets {
				worklogPayload := api.WorklogPayload{
					TimeSpent: ticket.Worklog.Value(),
					Started:   time.Now().Format(time.RFC3339),
				}
				cmd := m.SubmitWorklog(ticket.Key, worklogPayload)

				cmds = append(cmds, cmd)
			}

			result := tea.Batch(cmds...)

			return m, result
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

		urlCell := cellStyle.Render(utils.ConstructOsc8Hyperlink(url, ticket.Key))
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
	loggedTimeText := ""

	if m.state.dailyStatus.timeLogged {
		loggedTimeText = "\nTime is all logged for today, good job!"
	}

	return header + "\n" + separator + "\n" + rows + instructions + loggedTimeText
}
