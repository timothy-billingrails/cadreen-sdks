package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"
)

type tuiState int

const (
	tuiStateIdle tuiState = iota
	tuiStateStreaming
	tuiStateConfirming
)

type tuiKeyMap struct {
	Send  key.Binding
	Quit  key.Binding
	Help  key.Binding
	Clear key.Binding
}

func (k tuiKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.Quit}
}

func (k tuiKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Send, k.Clear},
		{k.Help, k.Quit},
	}
}

var tuiKeys = tuiKeyMap{
	Send: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Clear: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "clear"),
	),
}

type tuiModel struct {
	state         tuiState
	viewport      viewport.Model
	textarea      textarea.Model
	spinner       spinner.Model
	help          help.Model
	keys          tuiKeyMap
	glamourRender *glamour.TermRenderer
	messages      []tuiMessage
	pendingActs   []pendingAction
	confirmCursor int
	convID        string
	memoryOff     bool
	width         int
	height        int
	streamCh      <-chan cadreen.ChatStreamEvent
	streamBuf     strings.Builder
	ready         bool
}

type tuiMessage struct {
	role    string
	content string
}

var (
	tuiUserStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true)

	tuiCadreenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5"))

	tuiSystemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Italic(true)

	tuiPromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)

	tuiStatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(lipgloss.Color("0")).
			Padding(0, 1)

	tuiTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true)

	tuiHeaderSep = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	tuiFooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	tuiHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	tuiSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")).
				Bold(true)

	tuiUnselectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))
)

type streamDoneMsg struct {
	convID         string
	pendingActions []pendingAction
	err            error
}

type streamChunkMsg struct {
	content string
	rawJSON []byte
}

func initialTUIModel(memoryOff bool) tuiModel {
	ta := textarea.New()
	ta.Placeholder = "Ask Cadreen anything..."
	ta.MaxHeight = 10
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.Focus()

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	h := help.New()
	h.ShowAll = false

	vp := viewport.New(80, 20)

	return tuiModel{
		state:     tuiStateIdle,
		viewport:  vp,
		textarea:  ta,
		spinner:   sp,
		help:      h,
		keys:      tuiKeys,
		memoryOff: memoryOff,
		messages: []tuiMessage{
			{role: "system", content: "Cadreen Chat — type a message to start"},
		},
	}
}

func (m tuiModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		m.ready = true
		m.renderViewport()
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && !msg.Alt {
			if m.state == tuiStateConfirming {
				return m.handleConfirmKeys(msg)
			}
			return m.handleInput()
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.updateLayout()
			return m, nil
		case key.Matches(msg, m.keys.Clear):
			m.messages = []tuiMessage{
				{role: "system", content: "Conversation cleared."},
			}
			m.convID = ""
			m.renderViewport()
			return m, nil
		}

		if m.state == tuiStateConfirming {
			return m.handleConfirmKeys(msg)
		}

		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case streamChunkMsg:
		if msg.content != "" && len(m.messages) > 0 {
			m.streamBuf.WriteString(msg.content)
			m.messages[len(m.messages)-1].content = m.streamBuf.String()
			m.renderViewport()
		}
		if len(msg.rawJSON) > 0 {
			var extra streamExtra
			if json.Unmarshal(msg.rawJSON, &extra) == nil {
				if extra.ConversationID != "" {
					m.convID = extra.ConversationID
				}
			}
		}
		return m, readStreamChunk(m.streamCh)

	case streamDoneMsg:
		if msg.err != nil {
			m.messages = append(m.messages, tuiMessage{role: "system", content: fmt.Sprintf("Error: %s", msg.err)})
		}
		if msg.convID != "" {
			m.convID = msg.convID
		}
		if len(msg.pendingActions) > 0 {
			m.state = tuiStateConfirming
			m.pendingActs = msg.pendingActions
			m.confirmCursor = 0
		} else {
			m.state = tuiStateIdle
		}
		m.renderViewport()
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *tuiModel) handleInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.textarea.Value())
	m.textarea.SetValue("")

	if input == "" {
		return *m, nil
	}

	if input == "exit" || input == "quit" {
		return *m, tea.Quit
	}

	if input == "clear" {
		m.messages = []tuiMessage{
			{role: "system", content: "Conversation cleared."},
		}
		m.convID = ""
		m.renderViewport()
		return *m, nil
	}

	if m.state == tuiStateConfirming {
		m.messages = append(m.messages, tuiMessage{role: "system", content: "Please answer yes/no above first."})
		m.renderViewport()
		return *m, nil
	}

	m.messages = append(m.messages, tuiMessage{role: "user", content: input})
	m.messages = append(m.messages, tuiMessage{role: "cadreen", content: ""})
	m.state = tuiStateStreaming
	m.streamBuf.Reset()
	m.renderViewport()

	return *m, m.startStream(input)
}

func (m *tuiModel) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.String() == "y" || msg.String() == "Y":
		return m.sendConfirm("yes")
	case msg.String() == "n" || msg.String() == "N":
		return m.sendConfirm("no")
	case msg.String() == "up" || msg.String() == "k":
		if m.confirmCursor > 0 {
			m.confirmCursor--
		}
		m.renderViewport()
		return *m, nil
	case msg.String() == "down" || msg.String() == "j":
		if m.confirmCursor < 1 {
			m.confirmCursor++
		}
		m.renderViewport()
		return *m, nil
	case msg.Type == tea.KeyEnter:
		if m.confirmCursor == 0 {
			return m.sendConfirm("yes")
		}
		return m.sendConfirm("no")
	}
	return *m, nil
}

func (m *tuiModel) sendConfirm(answer string) (tea.Model, tea.Cmd) {
	m.messages = append(m.messages, tuiMessage{role: "user", content: answer})
	m.messages = append(m.messages, tuiMessage{role: "cadreen", content: ""})
	m.state = tuiStateStreaming
	m.pendingActs = nil
	m.streamBuf.Reset()
	return *m, m.startStream(answer)
}

func (m tuiModel) startStream(input string) tea.Cmd {
	return func() tea.Msg {
		req := cadreen.ChatCompletionRequest{
			Messages: []cadreen.ChatMessage{
				{Role: "user", Content: input},
			},
		}
		if m.memoryOff {
			req.Context = map[string]any{"memory": false}
		}
		if m.convID != "" {
			req.ConversationID = m.convID
		}

		client := newClient()
		ch, err := client.ChatCompletionsStream(context.Background(), req)
		if err != nil {
			return streamDoneMsg{err: err}
		}

		return streamDoneMsg(m.collectStream(ch))
	}
}

func (m tuiModel) collectStream(ch <-chan cadreen.ChatStreamEvent) streamDoneMsg {
	var convID string
	var pendingActions []pendingAction
	var lastErr error

	for event := range ch {
		if event.Error != nil {
			lastErr = event.Error
			break
		}
		if event.Chunk == nil {
			continue
		}

		if len(event.RawJSON) > 0 {
			var extra streamExtra
			if json.Unmarshal(event.RawJSON, &extra) == nil {
				if extra.ConversationID != "" {
					convID = extra.ConversationID
				}
				if len(extra.PendingActions) > 0 {
					pendingActions = extra.PendingActions
				}
			}
		}
	}

	return streamDoneMsg{convID: convID, pendingActions: pendingActions, err: lastErr}
}

func readStreamChunk(ch <-chan cadreen.ChatStreamEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-ch
		if !ok {
			return streamDoneMsg{}
		}
		if event.Error != nil {
			return streamDoneMsg{err: event.Error}
		}
		content := ""
		if event.Chunk != nil && len(event.Chunk.Choices) > 0 {
			content = event.Chunk.Choices[0].Delta.Content
		}
		return streamChunkMsg{content: content, rawJSON: event.RawJSON}
	}
}

func (m *tuiModel) renderViewport() {
	var b strings.Builder
	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			b.WriteString(tuiUserStyle.Render("You: "))
			b.WriteString(msg.content)
		case "cadreen":
			b.WriteString(tuiCadreenStyle.Render("Cadreen: "))
			if m.glamourRender != nil {
				rendered, err := m.glamourRender.Render(msg.content)
				if err == nil {
					b.WriteString(strings.TrimRight(rendered, "\n"))
				} else {
					b.WriteString(msg.content)
				}
			} else {
				b.WriteString(msg.content)
			}
		case "system":
			b.WriteString(tuiSystemStyle.Render(msg.content))
		}
		b.WriteString("\n")
	}

	if m.state == tuiStateConfirming {
		b.WriteString(m.renderConfirmation())
	}

	m.viewport.SetContent(b.String())
	m.viewport.GotoBottom()
}

func (m tuiModel) renderConfirmation() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(tuiPromptStyle.Render("Before I can do this, I need your permission."))
	b.WriteString("\n\n")

	if len(m.pendingActs) == 1 {
		b.WriteString(fmt.Sprintf("I want to: %s\n", humanToolName(m.pendingActs[0].Tool)))
	} else {
		b.WriteString(fmt.Sprintf("I need your permission for %d things:\n\n", len(m.pendingActs)))
		for i, a := range m.pendingActs {
			b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, humanToolName(a.Tool)))
		}
	}

	b.WriteString("\n")
	choices := []string{"Yes, proceed", "No, skip this"}
	for i, c := range choices {
		if m.confirmCursor == i {
			b.WriteString(tuiSelectedStyle.Render("(•) "))
		} else {
			b.WriteString(tuiUnselectedStyle.Render("( ) "))
		}
		b.WriteString(c + "\n")
	}
	b.WriteString("\n")
	b.WriteString(tuiHelpStyle.Render("↑/↓ to select • enter to confirm • y/n shortcut"))
	return b.String()
}

func (m *tuiModel) updateGlamour() {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.viewport.Width-4),
	)
	if err == nil {
		m.glamourRender = renderer
	}
}

func (m *tuiModel) updateLayout() {
	headerHeight := 1
	footerHeight := 1
	statusBarHeight := 1
	helpHeight := lipgloss.Height(m.help.View(m.keys))
	textareaHeight := m.textarea.Height()
	if textareaHeight < 1 {
		textareaHeight = 1
	}

	m.viewport.Width = m.width
	m.viewport.Height = m.height - headerHeight - footerHeight - statusBarHeight - helpHeight - textareaHeight - 2
	if m.viewport.Height < 1 {
		m.viewport.Height = 1
	}

	m.textarea.SetWidth(m.width - 4)

	m.updateGlamour()
}

func (m tuiModel) headerView() string {
	title := tuiTitleStyle.Render("Cadreen")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, tuiHeaderSep.Render(line))
}

func (m tuiModel) footerView() string {
	info := ""
	if len(m.convID) >= 8 {
		info = fmt.Sprintf("conv: %s", m.convID[:8])
	} else if m.convID != "" {
		info = fmt.Sprintf("conv: %s", m.convID)
	}
	if m.viewport.ScrollPercent() < 1.0 {
		pct := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
		if info != "" {
			info = pct + "  " + info
		} else {
			info = pct
		}
	}
	if info == "" {
		info = " "
	}
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, tuiHeaderSep.Render(line), tuiFooterStyle.Render(info))
}

func (m tuiModel) statusBar() string {
	left := ""
	switch m.state {
	case tuiStateIdle:
		left = "ready"
	case tuiStateStreaming:
		left = m.spinner.View() + " thinking..."
	case tuiStateConfirming:
		left = "waiting for confirmation"
	}

	right := ""
	if m.memoryOff {
		right = "memory: off"
	} else {
		right = "memory: on"
	}

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 0 {
		gap = 0
	}

	return tuiStatusStyle.Render(left + strings.Repeat(" ", gap) + right)
}

func (m tuiModel) View() string {
	if !m.ready {
		return "Initializing..."
	}

	helpView := m.help.View(m.keys)

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
		m.textarea.View(),
		helpView,
		m.statusBar(),
	)
}
