package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomasharkema/nix-htop/nixbuilders"
	"github.com/tomasharkema/nix-htop/tui/nodelist"
)

type Model struct {
	statusLoaded statusLoaded

	lastUpdated time.Time

	selectedNodeID string
	isLoading      bool
	err            error
	nodelist       nodelist.Model
	// nodedetails    nodedetails.Model
	spinner spinner.Model
	w, h    int
}

type statusRequest struct {
	status nixbuilders.ActiveBuildersResponse
	err    error
}
type statusLoaded nixbuilders.ActiveBuildersResponse
type statusError error

func getBuilderStatusAsync(status chan statusRequest) {
	ctx := context.Background()
	users, err := nixbuilders.GetActiveBuilders(ctx)

	if err != nil {
		status <- statusRequest{status: nil, err: err}
	} else {
		status <- statusRequest{status: users, err: nil}
	}
}

func getBuilderStatus() tea.Cmd {
	return func() tea.Msg {
		c := make(chan statusRequest)

		go getBuilderStatusAsync(c)

		statusReq := <-c
		if statusReq.err != nil {
			return statusError(statusReq.err)
		}
		return statusLoaded(statusReq.status)
	}
}

type TickMsg time.Time

func (m Model) Tick() tea.Msg {
	return TickMsg(time.Now())
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		getBuilderStatus(),
		m.spinner.Tick,
		m.nodelist.Init(),
		m.Tick,
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case statusLoaded:
		// m.tsStatus = tailscale.Status(msg)
		m.statusLoaded = msg
		m.nodelist = nodelist.New(m.statusLoaded, m.w, m.h)
		m.isLoading = false
		m.lastUpdated = time.Now()

	case statusError:
		m.isLoading = false
		m.err = msg
		m.lastUpdated = time.Now()

		return m, tea.Quit

	// case nodedetails.BackMsg:
	// m.viewState = viewStateList

	case nodelist.RefreshMsg:
		cmd = getBuilderStatus()
		cmds = append(cmds, cmd)

	case TickMsg:
		cmd = getBuilderStatus()
		cmds = append(cmds, cmd)
		cmd = tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})

		cmds = append(cmds, cmd)
		// case nodelist.NodeSelectedMsg:
	// m.selectedNodeID = string(msg)
	// m.nodedetails = nodedetails.New(&m.tsStatus, m.selectedNodeID, m.w, m.h)
	// m.viewState = viewStateDetails

	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		if !m.isLoading {
			m.nodelist.SetSize(msg.Width, msg.Height)
		}
	}

	// switch m.viewState {
	// case viewStateDetails:
	// 	m.nodedetails, cmd = m.nodedetails.Update(msg)
	// case viewStateList:
	if m.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
	} else {
		m.nodelist, cmd = m.nodelist.Update(msg)
	}
	// }

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// switch m.viewState {
	// case viewStateDetails:
	// 	return m.nodedetails.View()
	// default:
	if m.isLoading {
		return fmt.Sprintf("\n\n   %s Loading ...\n\n", m.spinner.View())
	}
	return m.nodelist.View()
	// }
}

func New() Model {
	m := Model{
		isLoading: true,
		spinner:   spinner.New(),
	}
	m.spinner.Spinner = spinner.Dot
	// m.spinner.Style = constants.SpinnerStyle
	return m
}
