package nodelist

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
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
type builderItem struct {
	User     nixbuilders.ActiveUser
	Progress progress.Model
	Info     nixbuilders.ProcessInfo
}

func (m Model) getItems() []list.Item {
	items := []list.Item{}

	for _, item := range *m.status {

		item := builderItem{
			User:     item,
			Progress: progress.New(),
			Info:     nixbuilders.NewInfo(item),
		}

		items = append(items, item)
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

func (a builderItem) FilterValue() string {
	return a.User.User
}

func (a builderItem) Title() string {

	// cmd, _ := a.processes[0].Cmdline()

	sinceCreation := time.Since(a.Info.CreateTime).Round(time.Second)
	sinceCreationStr := sinceCreation.String()

	pkgs := lipgloss.NewStyle().Bold(true).Render(a.User.DirName())

	// dir := a.User.Dir.Name()

	return fmt.Sprintf("%s %s %s", pkgs, a.User.User, sinceCreationStr)
}

func (a builderItem) Description() string {

	if len(a.User.Processes) == 0 {
		return ""
	}

	// prog := a.Progress.ViewAs(a.Info.CpuPercent)
	// prog = ""
	return fmt.Sprintf("CPU %.0f%% | MEM %s %.0f%% | PROCS %d", a.Info.CpuPercent*100, a.Info.MemoryBytes.Floor().String(), a.Info.MemoryPercent*100, len(a.User.Processes)) //, prog)
}
