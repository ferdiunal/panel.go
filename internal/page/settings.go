package page

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/domain/setting"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/widget"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Settings struct {
	Base
	Elements         []fields.Element
	HideInNavigation bool
}

func (p *Settings) Slug() string {
	return "settings"
}

func (p *Settings) Title() string {
	return "Settings"
}

func (p *Settings) Group() string {
	return "System"
}

func (p *Settings) NavigationOrder() int {
	return 100 // System items usually at the bottom
}

func (p *Settings) Visible() bool {
	return !p.HideInNavigation
}

func (p *Settings) Cards() []widget.Card {
	return []widget.Card{}
}

func (p *Settings) Fields() []fields.Element {
	return p.Elements
}

func (p *Settings) Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error {
	for key, value := range data {
		val := map[string]interface{}{
			"value": value,
		}

		s := setting.Setting{
			Key:   key,
			Value: val,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
		}).Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}
