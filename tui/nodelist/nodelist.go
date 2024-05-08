package nodelist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomasharkema/nix-htop/nixbuilders"
	"github.com/tomasharkema/nix-htop/tui/keymap"
)

type RefreshMsg bool

type Model struct {
	status   nixbuilders.ActiveBuildersResponse
	exitNode string
	list     list.Model
	keyMap   keymap.KeyMap
	msg      string
	w        int
	h        int
}

func (m Model) getItems() []list.Item {
	items := []list.Item{}

	for _, item := range *m.status {
		items = append(items, builderItem{item})
	}

	return items
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	var kcmds []tea.Cmd
	// 	m, kcmds = m.keyBindingsHandler(msg)
	// 	cmds = append(cmds, kcmds...)
	// default:
	// }

	// cmd = func() tea.Msg {
	// 	<-time.After(time.Second * 2)
	// 	return RefreshMsg(true)
	// }

	cmds = append(cmds, cmd)

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	// m.updateKeybindings()
	return m, tea.Batch(cmds...)
}

func New(status nixbuilders.ActiveBuildersResponse, w, h int) Model {

	d := list.NewDefaultDelegate()

	m := Model{
		status: status,
		list:   list.New([]list.Item{}, d, w, h),
		// keyMap:     keymap.NewKeyMap(),
		// tailStatus: status,
	}

	m.list.SetItems(m.getItems())

	headerHeight := lipgloss.Height(m.headerView())
	m.list.SetHeight(h - headerHeight)

	return m
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s", m.headerView(), m.list.View())
}

func (m Model) headerView() string {
	s := lipgloss.NewStyle().
		Bold(true).Render("Nix Builders")

	s += fmt.Sprintf("\n%d builders active...", len(*m.status))

	return lipgloss.NewStyle().Margin(1, 2).Render(s)
}

func (m *Model) SetSize(w int, h int) {
	m.w = w
	m.h = h
	headerHeight := lipgloss.Height(m.headerView())
	m.list.SetSize(w, h-headerHeight)
}

type builderItem struct {
	nixbuilders.ActiveUser
}

func (a builderItem) FilterValue() string {
	return a.User
}

func (a builderItem) Title() string {

	// cmd, _ := a.processes[0].Cmdline()

	return fmt.Sprintf("%s %s", lipgloss.NewStyle().Bold(true).Render(a.DirName()), a.User)
}

func (a builderItem) Description() string {

	percent, _ := a.Processes[0].CPUPercent()
	mem, _ := a.Processes[0].MemoryPercent()
	return fmt.Sprintf("CPU %.0f%% | MEM %.0f%% | PROCS %d", percent*100, mem*100, len(a.Processes))
}
