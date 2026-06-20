package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"
)

type wizardStep int

const (
	wizardStepAuth wizardStep = iota
	wizardStepHealth
	wizardStepDoctor
)

var wizardStepNames = []string{"Authenticate", "Check Health", "Run Doctor"}

type wizardKeyMap struct {
	Next  key.Binding
	Prev  key.Binding
	Quit  key.Binding
	Help  key.Binding
	Run   key.Binding
}

func (k wizardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Run, k.Quit}
}

func (k wizardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Run},
		{k.Help, k.Quit},
	}
}

var wizardKeys = wizardKeyMap{
	Next: key.NewBinding(
		key.WithKeys("tab", "right"),
		key.WithHelp("tab/→", "next step"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab", "left"),
		key.WithHelp("shift+tab/←", "prev step"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "run step"),
	),
}

type wizardModel struct {
	steps    []wizardStep
	active   wizardStep
	results  []wizardResult
	progress progress.Model
	spinner  spinner.Model
	help     help.Model
	keys     wizardKeyMap
	width    int
	height   int
	running  bool
	done     bool
}

type wizardResult struct {
	step    wizardStep
	status  string // "ok", "fail", "skip", "pending"
	message string
}

type wizardStepDoneMsg struct {
	step    wizardStep
	status  string
	message string
}

func newWizardModel() wizardModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	h := help.New()

	return wizardModel{
		steps:    []wizardStep{wizardStepAuth, wizardStepHealth, wizardStepDoctor},
		active:   wizardStepAuth,
		progress: p,
		spinner:  sp,
		help:     h,
		keys:     wizardKeys,
		results: []wizardResult{
			{step: wizardStepAuth, status: "pending", message: ""},
			{step: wizardStepHealth, status: "pending", message: ""},
			{step: wizardStepDoctor, status: "pending", message: ""},
		},
	}
}

func (m wizardModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 20
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Next):
			if !m.running && int(m.active) < len(m.steps)-1 {
				m.active++
			}
			return m, nil
		case key.Matches(msg, m.keys.Prev):
			if !m.running && m.active > 0 {
				m.active--
			}
			return m, nil
		case key.Matches(msg, m.keys.Run):
			if !m.running && m.results[m.active].status == "pending" {
				m.running = true
				return m, tea.Batch(m.spinner.Tick, m.runStep(m.active))
			}
		}

	case spinner.TickMsg:
		if m.running {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case wizardStepDoneMsg:
		m.results[msg.step] = wizardResult{
			step:    msg.step,
			status:  msg.status,
			message: msg.message,
		}
		m.running = false
		if msg.step == wizardStepDoctor {
			m.done = true
		} else if msg.status == "ok" {
			m.active = msg.step + 1
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m wizardModel) runStep(step wizardStep) tea.Cmd {
	return func() tea.Msg {
		switch step {
		case wizardStepAuth:
			return m.runAuthStep()
		case wizardStepHealth:
			return m.runHealthStep()
		case wizardStepDoctor:
			return m.runDoctorStep()
		}
		return wizardStepDoneMsg{step: step, status: "skip", message: "Unknown step"}
	}
}

func (m wizardModel) runAuthStep() tea.Msg {
	if cfg.IsAuthenticated() {
		return wizardStepDoneMsg{
			step:    wizardStepAuth,
			status:  "ok",
			message: fmt.Sprintf("Authenticated as %s", output.MaskKey(cfg.APIKeyResolved())),
		}
	}
	return wizardStepDoneMsg{
		step:    wizardStepAuth,
		status:  "fail",
		message: "Not authenticated. Run 'cadreen login' first.",
	}
}

func (m wizardModel) runHealthStep() tea.Msg {
	resp, err := rawDo("GET", "/api/v1/cadreen/health", nil)
	if err != nil {
		return wizardStepDoneMsg{step: wizardStepHealth, status: "fail", message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return wizardStepDoneMsg{step: wizardStepHealth, status: "fail", message: fmt.Sprintf("Health check returned %d", resp.StatusCode)}
	}
	return wizardStepDoneMsg{step: wizardStepHealth, status: "ok", message: "API is healthy"}
}

func (m wizardModel) runDoctorStep() tea.Msg {
	resp, err := rawDo("GET", "/api/v1/cadreen/health", nil)
	if err != nil {
		return wizardStepDoneMsg{step: wizardStepDoctor, status: "fail", message: err.Error()}
	}
	defer resp.Body.Close()

	var health struct {
		Status      string `json:"status"`
		Connections struct {
			Total   int `json:"total"`
			Healthy int `json:"healthy"`
		} `json:"connections"`
		Governance struct {
			Active int `json:"active"`
		} `json:"governance"`
		Memory struct {
			Status string `json:"status"`
		} `json:"memory"`
	}
	if err := readJSON(resp.Body, &health); err != nil {
		return wizardStepDoneMsg{step: wizardStepDoctor, status: "fail", message: "Could not parse health response"}
	}

	checks := 0
	passed := 0
	if health.Status == "ok" {
		passed++
	}
	checks++
	if health.Connections.Total > 0 {
		passed++
	}
	checks++
	if health.Governance.Active > 0 {
		passed++
	}
	checks++
	if health.Memory.Status != "" {
		passed++
	}
	checks++

	status := "ok"
	if passed < checks {
		status = "fail"
	}
	return wizardStepDoneMsg{
		step:    wizardStepDoctor,
		status:  status,
		message: fmt.Sprintf("%d/%d checks passed", passed, checks),
	}
}

func (m wizardModel) View() string {
	var b strings.Builder

	b.WriteString(wizardTitleStyle.Render("Cadreen Setup"))
	b.WriteString("\n\n")

	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	b.WriteString(m.renderStepContent())
	b.WriteString("\n\n")

	b.WriteString(m.renderProgressBar())
	b.WriteString("\n\n")

	b.WriteString(m.help.View(m.keys))

	return b.String()
}

func (m wizardModel) renderTabs() string {
	var tabs []string
	for i, name := range wizardStepNames {
		style := wizardTabInactive
		if wizardStep(i) == m.active {
			style = wizardTabActive
		}
		status := ""
		switch m.results[i].status {
		case "ok":
			status = " ✓"
		case "fail":
			status = " ✗"
		case "skip":
			status = " ○"
		}
		tabs = append(tabs, style.Render(name+status))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m wizardModel) renderStepContent() string {
	step := m.active
	result := m.results[step]

	var b strings.Builder
	b.WriteString(wizardStepTitle.Render(wizardStepNames[step]))
	b.WriteString("\n\n")

	switch step {
	case wizardStepAuth:
		b.WriteString("Cadreen needs an API key to connect.\n")
		b.WriteString("You can paste one, or get one from the dashboard.\n\n")
	case wizardStepHealth:
		b.WriteString("Checking if the API is reachable.\n\n")
	case wizardStepDoctor:
		b.WriteString("Running readiness checks.\n\n")
	}

	if m.running {
		b.WriteString(m.spinner.View() + " Running...")
	} else {
		switch result.status {
		case "ok":
			b.WriteString(wizardStatusOK.Render("✓ " + result.message))
		case "fail":
			b.WriteString(wizardStatusFail.Render("✗ " + result.message))
		case "pending":
			b.WriteString(wizardStatusPending.Render("○ Press Enter to run"))
		}
	}

	return b.String()
}

func (m wizardModel) renderProgressBar() string {
	completed := 0
	for _, r := range m.results {
		if r.status == "ok" || r.status == "fail" || r.status == "skip" {
			completed++
		}
	}
	pct := float64(completed) / float64(len(m.results))
	return m.progress.ViewAs(pct)
}

var (
	wizardTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("5")).
				Bold(true)

	wizardTabActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("6")).
			Padding(0, 2)

	wizardTabInactive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("8")).
				Padding(0, 2)

	wizardStepTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true)

	wizardStatusOK = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2"))

	wizardStatusFail = lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Bold(true)

	wizardStatusPending = lipgloss.NewStyle().
				Foreground(lipgloss.Color("3"))
)
