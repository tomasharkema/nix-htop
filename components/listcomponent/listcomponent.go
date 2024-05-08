package listcomponent

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	items []string

	list list.Model
}

type Props struct {
	SetItems func([]string)
}

func New() *Component {
	return &Component{
		list: list.Model{},
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return nil
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	if msg.Type == tea.KeyEnter {
	// 		// Lifted state power! Woohooo
	// 		c.Props().SetText(c.textinput.Value())

	// 		reactea.SetRoute("/displayname")

	// 		return nil
	// 	}
	// }

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return cmd
}

// Here we are not using width and height, but you can!
// Using lipgloss styles for example
func (c *Component) Render(int, int) string {
	return c.list.View()
}
