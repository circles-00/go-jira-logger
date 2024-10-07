package tui

import (
	"math"

	"go_jira_logger/pkg/api"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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

const (
	undersized size = iota
	small
	medium
	large
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
	state           state
	page            page
	tickets         []JiraTicket
	viewportWidth   int
	viewportHeight  int
	widthContainer  int
	heightContainer int
	widthContent    int
	heightContent   int
	size            size
	accessToken     string
	viewport        viewport.Model
	hasScroll       bool
	ready           bool
}

func (m model) Init() tea.Cmd {
	return m.SplashInit()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.viewportHeight = msg.Height

		switch {
		case m.viewportWidth < 20 || m.viewportHeight < 10:
			m.size = undersized
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 40:
			m.size = small
			m.widthContainer = m.viewportWidth
			m.heightContainer = m.viewportHeight
		case m.viewportWidth < 60:
			m.size = medium
			m.widthContainer = 40
			m.heightContainer = int(math.Min(float64(msg.Height), 30))
		default:
			m.size = large
			m.widthContainer = 60
			m.heightContainer = int(math.Min(float64(msg.Height), 30))
		}

		m.widthContent = m.widthContainer - 4
		m = m.updateViewport()
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

func (m model) updateViewport() model {
	width := m.widthContainer - 4
	m.heightContent = m.heightContainer

	if !m.ready {
		m.viewport = viewport.New(width, m.heightContent)
		m.viewport.HighPerformanceRendering = false
		m.ready = true
	} else {
		m.viewport.Width = width
		m.viewport.Height = m.heightContent
		m.viewport.GotoTop()
	}

	m.viewport.KeyMap = viewport.DefaultKeyMap()

	m.hasScroll = m.viewport.VisibleLineCount() < m.viewport.TotalLineCount()

	if m.hasScroll {
		m.widthContent = m.widthContainer - 4
	} else {
		m.widthContent = m.widthContainer - 2
	}

	return m
}
