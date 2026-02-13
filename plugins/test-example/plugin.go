package test_example

import (
	_ "test-example/resources/address"
	_ "test-example/resources/billing_info"
	_ "test-example/resources/category"
	_ "test-example/resources/organization"
	_ "test-example/resources/product"
	_ "test-example/resources/shipment"
	_ "test-example/resources/shipment_row"

	"github.com/ferdiunal/panel.go/pkg/plugin"
)

// Plugin, test-example plugin'i.
type Plugin struct {
	plugin.BasePlugin
}

// init, plugin'i global registry'ye kaydeder.
func init() {
	plugin.Register(&Plugin{})
}

// Name, plugin adını döndürür.
func (p *Plugin) Name() string {
	return "test-example"
}

// Version, plugin versiyonunu döndürür.
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Register, plugin'i Panel'e kaydeder.
func (p *Plugin) Register(panel interface{}) error {
	// Tüm resource'lar init() fonksiyonları ile otomatik kayıt edilir
	return nil
}

// Boot, plugin'i boot eder.
func (p *Plugin) Boot(panel interface{}) error {
	// Plugin boot logic
	return nil
}
