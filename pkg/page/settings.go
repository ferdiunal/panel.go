package page

import (
	"encoding/json"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/domain/setting"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
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

func (p *Settings) Description() string {
	return "Sistem ayarlarını yönetin"
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
		// Convert value to string
		var strValue string
		if v, ok := value.(string); ok {
			strValue = v
		} else {
			// For non-string values, convert to JSON string
			b, _ := json.Marshal(value)
			strValue = string(b)
		}

		s := setting.Setting{
			Key:   key,
			Value: strValue,
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
