package tui

import (
	"go_jira_logger/pkg/api"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	page = int
	size = int
)

const (
	splashPage page = iota
	dailyStatusPage
)

type JiraTicket struct {
	api.JiraIssue
	Worklog textinput.Model
}

type state struct {
	splash      SplashState
	dailyStatus DailyStatusState
}

type model struct {
	state   state
	page    page
	tickets []JiraTicket
}

func (m model) Init() tea.Cmd {
	return m.SplashInit()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case []JiraTicket:
		m.tickets = msg
	}

	var cmd tea.Cmd
	switch m.page {
	case splashPage:
		m, cmd = m.SplashUpdate(msg)
	case dailyStatusPage:
		m, cmd = m.DailyStatusUpdate(msg)
	}

	cmds := []tea.Cmd{}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.page {
	case splashPage:
		return m.SplashView()
	case dailyStatusPage:
		return m.DailyStatusView()
	default:
		return m.getContent()
	}
}

func (m model) getContent() string {
	page := "unknown"
	switch m.page {
	case splashPage:
		page = m.SplashView()
	case dailyStatusPage:
		page = m.DailyStatusView()
	}
	return page
}

func (m model) SwitchPage(page page) model {
	m.page = page
	return m
}

func NewModel() (tea.Model, error) {
	result := model{
		page: splashPage,
		state: state{
			splash:      SplashState{},
			dailyStatus: DailyStatusState{},
		},
	}

	return result, nil
}
