package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"
)

type tuiState int

const (
	tuiStateIdle tuiState = iota
	tuiStateStreaming
	tuiStateConfirming
)

type tuiModel struct {
	state       tuiState
	viewport    viewport.Model
	textInput   textinput.Model
	messages    []tuiMessage
	pendingActs []pendingAction
	convID      string
	memoryOff   bool
	width       int
	height      int
	streamCh    <-chan cadreen.ChatStreamEvent
	streamBuf   strings.Builder
	ready       bool
	err         error
}

type tuiMessage struct {
	role    string // "user", "cadreen", "system"
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
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 80

	vp := viewport.New(80, 20)

	return tuiModel{
		state:     tuiStateIdle,
		viewport:  vp,
		textInput: ti,
		memoryOff: memoryOff,
		messages: []tuiMessage{
			{role: "system", content: "Cadreen Chat — type 'exit' to quit, 'clear' to reset"},
		},
	}
}

func (m tuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		m.textInput.Width = msg.Width - 4
		m.ready = true
		m.renderViewport()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.state == tuiStateConfirming {
				return m.handleConfirmation()
			}
			return m.handleInput()
		}

	case streamChunkMsg:
		if msg.content != "" {
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
			m.messages = append(m.messages, tuiMessage{role: "system", content: formatConfirmationPrompt(msg.pendingActions)})
		} else {
			m.state = tuiStateIdle
		}
		m.renderViewport()
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *tuiModel) handleInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.textInput.Value())
	m.textInput.SetValue("")
	m.textInput.Reset()

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

func (m *tuiModel) handleConfirmation() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(strings.ToLower(m.textInput.Value()))
	m.textInput.SetValue("")
	m.textInput.Reset()

	if input == "yes" || input == "y" {
		m.messages = append(m.messages, tuiMessage{role: "user", content: "yes"})
		m.messages = append(m.messages, tuiMessage{role: "cadreen", content: ""})
		m.state = tuiStateStreaming
		m.pendingActs = nil
		m.streamBuf.Reset()
		m.renderViewport()
		return *m, m.startStream("yes")
	} else if input == "no" || input == "n" {
		m.messages = append(m.messages, tuiMessage{role: "user", content: "no"})
		m.messages = append(m.messages, tuiMessage{role: "cadreen", content: ""})
		m.state = tuiStateStreaming
		m.pendingActs = nil
		m.streamBuf.Reset()
		m.renderViewport()
		return *m, m.startStream("no")
	}

	m.messages = append(m.messages, tuiMessage{role: "system", content: "Please type 'yes' or 'no'."})
	m.renderViewport()
	return *m, nil
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

		if len(event.Chunk.Choices) > 0 {
			delta := event.Chunk.Choices[0].Delta
			if delta.Content != "" {
				tea.Printf("%s", delta.Content)()
			}
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
			b.WriteString(msg.content)
		case "system":
			b.WriteString(tuiSystemStyle.Render(msg.content))
		}
		b.WriteString("\n")
	}

	if m.state == tuiStateStreaming {
		b.WriteString("\n")
		b.WriteString(tuiStatusStyle.Render("streaming..."))
		b.WriteString("\n")
	}

	m.viewport.SetContent(b.String())
	m.viewport.GotoBottom()
}

func formatConfirmationPrompt(actions []pendingAction) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(tuiPromptStyle.Render("Before I can do this, I need your permission."))
	b.WriteString("\n\n")

	if len(actions) == 1 {
		b.WriteString(fmt.Sprintf("I want to: %s\n", humanToolName(actions[0].Tool)))
	} else {
		b.WriteString(fmt.Sprintf("I need your permission for %d things:\n\n", len(actions)))
		for i, a := range actions {
			b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, humanToolName(a.Tool)))
		}
	}

	b.WriteString("\n")
	b.WriteString(tuiPromptStyle.Render("Proceed? (yes/no)"))
	return b.String()
}

func (m tuiModel) View() string {
	if !ready(m) {
		return "Initializing..."
	}

	title := tuiTitleStyle.Render("Cadreen")
	status := m.statusBar()

	return fmt.Sprintf("%s\n%s\n%s\n%s",
		title,
		m.viewport.View(),
		status,
		m.textInput.View(),
	)
}

func ready(m tuiModel) bool {
	return m.ready
}

func (m tuiModel) statusBar() string {
	left := ""
	switch m.state {
	case tuiStateIdle:
		left = "ready"
	case tuiStateStreaming:
		left = "streaming..."
	case tuiStateConfirming:
		left = "waiting for confirmation"
	}

	right := ""
	if m.convID != "" {
		right = fmt.Sprintf("conversation: %s", m.convID[:8])
	}

	gap := m.width - len(left) - len(right) - 2
	if gap < 0 {
		gap = 0
	}

	return tuiStatusStyle.Render(left + strings.Repeat(" ", gap) + right)
}
