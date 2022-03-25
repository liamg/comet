package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultWidth = 20
const listHeight = 14

var (
	titleStyle           = lipgloss.NewStyle().MarginLeft(2)
	itemStyle            = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle    = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ff55ff"))
	itemDescriptionStyle = lipgloss.NewStyle().PaddingLeft(2).Faint(true)
	paginationStyle      = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle            = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle        = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(prefix)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Title())

	var output string
	if index == m.Index() {
		output = selectedItemStyle.Render("> " + str)
	} else {
		output = itemStyle.Render(str)
	}
	output += itemDescriptionStyle.PaddingLeft(12 - len(str)).Render(i.Description())

	_, _ = fmt.Fprint(w, output)
}

type model struct {
	chosenPrefix bool
	chosenScope  bool
	chosenMsg    bool
	chosenBody   bool
	specifyBody  bool
	prefix       string
	scope        string
	msg          string
	prefixList   list.Model
	msgInput     textinput.Model
	scopeInput   textinput.Model
	ynInput      textinput.Model
	items        []prefix
	quitting     bool
	err          error
}

func newModel(prefixes []list.Item) *model {

	// set up list
	prefixList := list.New(prefixes, itemDelegate{}, defaultWidth, listHeight)
	prefixList.Title = "What are you committing?"
	prefixList.SetShowStatusBar(false)
	prefixList.SetFilteringEnabled(true)
	prefixList.Styles.Title = titleStyle
	prefixList.Styles.PaginationStyle = paginationStyle
	prefixList.Styles.HelpStyle = helpStyle

	// set up scope prompt
	scopeInput := textinput.New()
	scopeInput.Placeholder = "Scope"
	scopeInput.CharLimit = 16
	scopeInput.Width = 20

	// set up commit message prompt
	commitInput := textinput.New()
	commitInput.Placeholder = "Commit message"
	commitInput.CharLimit = 100
	commitInput.Width = 50

	// set up add body confirmation
	bodyConfirmation := textinput.New()
	bodyConfirmation.Placeholder = "y/N"
	bodyConfirmation.CharLimit = 1
	bodyConfirmation.Width = 20

	return &model{
		prefixList: prefixList,
		scopeInput: scopeInput,
		msgInput:   commitInput,
		ynInput:    bodyConfirmation,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch {
	case !m.chosenPrefix:
		return m.updatePrefixList(msg)
	case !m.chosenScope:
		return m.updateScopeInput(msg)
	case !m.chosenMsg:
		return m.updateMsgInput(msg)
	case !m.chosenBody:
		return m.updateYNInput(msg)
	default:
		return m, tea.Quit
	}
}

func (m *model) Finished() bool {
	return m.chosenBody
}

func (m *model) CommitMessage() (string, bool) {
	prefix := m.prefix
	if m.scope != "" {
		prefix = fmt.Sprintf("%s(%s)", prefix, m.scope)
	}
	return fmt.Sprintf("%s: %s", prefix, m.msg), m.specifyBody
}

func (m *model) updatePrefixList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.prefixList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.prefixList.SelectedItem().(prefix)
			if ok {
				m.prefix = i.Title()
				m.chosenPrefix = true
				m.scopeInput.Focus()
			}
		}
	}

	var cmd tea.Cmd
	m.prefixList, cmd = m.prefixList.Update(msg)
	return m, cmd
}

func (m *model) updateScopeInput(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.chosenScope = true
			m.scope = m.scopeInput.Value()
			m.msgInput.Focus()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.scopeInput, cmd = m.scopeInput.Update(msg)
	return m, cmd
}

func (m *model) updateMsgInput(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.chosenMsg = true
			m.msg = m.msgInput.Value()
			m.ynInput.Focus()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.msgInput, cmd = m.msgInput.Update(msg)
	return m, cmd
}

func (m *model) updateYNInput(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.chosenMsg = true
			switch strings.ToLower(m.ynInput.Value()) {
			case "y":
				m.specifyBody = true
			}
			m.chosenBody = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.ynInput, cmd = m.ynInput.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	switch {
	case !m.chosenPrefix:
		return "\n" + m.prefixList.View()
	case !m.chosenScope:
		return titleStyle.Render(fmt.Sprintf(
			"\nEnter a scope (enter to skip):\n\n%s\n\n%s",
			m.scopeInput.View(),
			"(esc to cancel)",
		) + "\n")
	case !m.chosenMsg:
		return titleStyle.Render(fmt.Sprintf(
			"\nEnter a commit message:\n\n%s\n\n%s",
			m.msgInput.View(),
			"(esc to cancel)",
		) + "\n")
	case !m.chosenBody:
		return fmt.Sprintf("\nDo you need to specify a body/footer?\n\n%s\n", m.ynInput.View())
	case m.quitting:
		return quitTextStyle.Render("Aborted.\n")
	default:
		return "\nCreating commit...\n"
	}
}
