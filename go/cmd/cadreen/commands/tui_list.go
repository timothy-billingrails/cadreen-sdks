package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type tuiListItem struct {
	title       string
	description string
}

func (i tuiListItem) Title() string       { return i.title }
func (i tuiListItem) Description() string { return i.description }
func (i tuiListItem) FilterValue() string { return i.title }

type tuiListKeyMap struct {
	Quit key.Binding
	Help key.Binding
}

func (k tuiListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k tuiListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

type tuiListModel struct {
	list   list.Model
	help   help.Model
	keys   tuiListKeyMap
	width  int
	height int
	title  string
}

func newTUIListModel(title string, items []tuiListItem) tuiListModel {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	h := help.New()

	return tuiListModel{
		list:  l,
		help:  h,
		title: title,
		keys: tuiListKeyMap{
			Quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
	}
}

func (m tuiListModel) Init() tea.Cmd {
	return nil
}

func (m tuiListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m tuiListModel) View() string {
	var b strings.Builder
	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.help.View(m.keys))
	return b.String()
}

func tuiListFromMemory(items []map[string]interface{}) []tuiListItem {
	var result []tuiListItem
	for _, item := range items {
		title := ""
		desc := ""
		if t, ok := item["title"].(string); ok {
			title = t
		}
		if d, ok := item["description"].(string); ok {
			desc = d
		}
		if title == "" {
			title = "Untitled"
		}
		if desc == "" {
			desc = "No description"
		}
		result = append(result, tuiListItem{title: title, description: desc})
	}
	return result
}

func tuiListFromPolicies(items []map[string]interface{}) []tuiListItem {
	var result []tuiListItem
	for _, item := range items {
		name := ""
		desc := ""
		if n, ok := item["name"].(string); ok {
			name = n
		}
		if d, ok := item["description"].(string); ok {
			desc = d
		}
		if name == "" {
			name = "Unnamed policy"
		}
		if desc == "" {
			desc = "No description"
		}
		result = append(result, tuiListItem{title: name, description: desc})
	}
	return result
}

func tuiListEmptyMessage(entity string) string {
	return fmt.Sprintf("No %s found. Try running 'cadreen %s add' first.", entity, entity)
}
