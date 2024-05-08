package main

import (
	"context"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/samber/lo"
)

type Model struct {
	Users []string

	list     list.Model
	quitting bool
	choice   string
}

// func initialModel(ctx context.Context) Model {
// 	model := Model{}

// 	group, err := buildUser(ctx)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	model.Users = group.Users

// 	items := lo.Map(model.Users, func(user string, index int) list.Item {
// 		return item(user)
// 	})

// 	model.list = list.New(items, list.DefaultDelegate{}, 0, 0)

// 	return model
// }

// func (m Model) Init() tea.Cmd {
// 	// Just return `nil`, which means "no I/O right now, please."
// 	return nil
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmds []tea.Cmd

// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		h, v := appStyle.GetFrameSize()
// 		m.list.SetSize(msg.Width-h, msg.Height-v)

// 	case tea.KeyMsg:
// 		switch keypress := msg.String(); keypress {
// 		case "q", "ctrl+c":
// 			m.quitting = true
// 			return m, tea.Quit

// 		case "enter":
// 			i, ok := m.list.SelectedItem().(item)
// 			if ok {
// 				m.choice = string(i)
// 			}
// 			return m, tea.Quit
// 		}
// 	}

// 	// This will also call our delegate's update function.
// 	newListModel, cmd := m.list.Update(msg)
// 	m.list = newListModel
// 	cmds = append(cmds, cmd)

// 	return m, tea.Batch(cmds...)

// }

// func (m Model) View() string {
// 	// s := "Builders\n"

// 	// s += m.list.View()

// 	return m.list.View()
// }

// func main() {
// 	ctx := context.Background()

// 	// fmt.Println(group.Users)
// 	// fmt.Println(activeBuildUsers(ctx))

// 	p := tea.NewProgram(initialModel(ctx))
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }

type item string

func (i item) FilterValue() string { return string(i) }

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	users []string
	list  list.Model
}

func newApp() *Component {
	ctx := context.Background()
	group, err := buildUser(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	items := lo.Map(group.Users, func(user string, index int) list.Item {
		return item(user)
	})

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	return &Component{
		users: group.Users,
		list:  l,
	}
}

func (c *Component) Render(width, height int) string {
	return c.list.View()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return cmd
}

func main() {
	ctx := context.Background()
	group, err := buildUser(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(group)

	// reactea.NewProgram initializes program with
	// "translation layer", so Reactea components work
	// program := reactea.NewProgram(newApp())

	// if _, err := program.Run(); err != nil {
	// 	panic(err)
	// }
}
