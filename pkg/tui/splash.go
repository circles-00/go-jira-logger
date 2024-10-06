package tui

import (
	"fmt"
	"time"

	"go_jira_logger/pkg/api"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SplashState struct {
	dailyStatus bool
	delay       bool
}

type DelayCompleteMsg struct{}

func (m model) LoadCmds() []tea.Cmd {
	cmds := []tea.Cmd{}

	// Make sure the loading state shows for at least a couple seconds
	cmds = append(cmds, tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return DelayCompleteMsg{}
	}))

	cmds = append(cmds, func() tea.Msg {
		issues := api.FetchIssues()
		tickets := make([]JiraTicket, 0)

		for _, issue := range issues.Issues {
			tickets = append(tickets, JiraTicket{
				JiraIssue: issue,
				Worklog:   textinput.New(),
			})
		}

		return tickets
	})

	return cmds
}

func (m model) IsLoadingComplete() bool {
	return m.state.splash.dailyStatus && m.state.splash.delay
}

func (m model) SplashInit() tea.Cmd {
	return tea.Batch(m.LoadCmds()...)
}

func (m model) SplashUpdate(msg tea.Msg) (model, tea.Cmd) {
	switch msg.(type) {
	case DelayCompleteMsg:
		m.state.splash.delay = true
	case []JiraTicket:
		m.state.splash.dailyStatus = true
	}

	if m.IsLoadingComplete() {
		return m.SwitchPage(dailyStatusPage), nil
	}
	return m, nil
}

func (m model) SplashView() string {
	return fmt.Sprintln("Jira Logger")
}
