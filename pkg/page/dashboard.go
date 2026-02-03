package page

import (
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

type Dashboard struct {
	Base
}

func (d *Dashboard) Slug() string {
	return "dashboard"
}

func (d *Dashboard) Title() string {
	return "Dashboard"
}

func (d *Dashboard) Icon() string {
	return "layout-dashboard"
}

func (d *Dashboard) Cards() []widget.Card {
	return []widget.Card{
		widget.NewCountWidget("Total Users", &user.User{}),
	}
}
