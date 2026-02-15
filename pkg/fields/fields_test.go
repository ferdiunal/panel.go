package fields

import (
	"encoding/json"
	"testing"
)

func TestFieldSerialization(t *testing.T) {
	f := Text("Full Name", "full_name").
		Sortable().
		Required().
		Placeholder("Enter your name").
		HelpText("This is your full name").
		OnList()

	data := f.JsonSerialize()

	if data["name"] != "Full Name" {
		t.Errorf("Expected name to be 'Full Name', got %v", data["name"])
	}
	if data["key"] != "full_name" {
		t.Errorf("Expected key to be 'full_name', got %v", data["key"])
	}
	if data["view"] != "text-field-index" {
		t.Errorf("Expected view to be 'text-field-index', got %v", data["view"])
	}
	if data["sortable"] != true {
		t.Error("Expected sortable to be true")
	}
	if data["required"] != true {
		t.Error("Expected required to be true")
	}
	if data["placeholder"] != "Enter your name" {
		t.Errorf("Expected placeholder to be 'Enter your name', got %v", data["placeholder"])
	}
	if data["context"] != SHOW_ON_LIST {
		t.Errorf("Expected context to be %v, got %v", SHOW_ON_LIST, data["context"])
	}

	// Test JSON marshaling
	bytes, err := json.Marshal(f)
	if err != nil {
		t.Errorf("Failed to marshal field: %v", err)
	}
	if len(bytes) == 0 {
		t.Error("Expected non-empty JSON output")
	}
}

func TestIDInstantiation(t *testing.T) {
	id := ID()
	if id.Name != "ID" {
		t.Errorf("Expected default name 'ID', got %v", id.Name)
	}
	if id.Key != "id" {
		t.Errorf("Expected default key 'id', got %v", id.Key)
	}
	if id.View != "id-field" {
		t.Errorf("Expected view 'id-field', got %v", id.View)
	}
}

func TestPasswordType(t *testing.T) {
	pwd := Password("Secret")
	if pwd.Type != TYPE_PASSWORD {
		t.Errorf("Expected type password, got %v", pwd.Type)
	}
}

func TestNewFieldTypes(t *testing.T) {
	date := Date("Birth Date")
	if date.Type != TYPE_DATE || date.View != "date-field" {
		t.Errorf("Date field mismatch")
	}

	dateTime := DateTime("Created At")
	if dateTime.Type != TYPE_DATETIME || dateTime.View != "datetime-field" {
		t.Errorf("DateTime field mismatch")
	}

	file := File("Avatar")
	if file.Type != TYPE_FILE || file.View != "file-field" {
		t.Errorf("File field mismatch")
	}

	kv := KeyValue("Settings")
	if kv.Type != TYPE_KEY_VALUE || kv.View != "key-value-field" {
		t.Errorf("KeyValue field mismatch")
	}

	money := Money("Price")
	if money.Type != TYPE_MONEY || money.View != "money-field" {
		t.Errorf("Money field mismatch")
	}
	if money.Props["currency"] != string(CurrencyUSD) {
		t.Errorf("Money default currency mismatch")
	}
}

func TestRelationshipFields(t *testing.T) {
	bt := Link("User", "users")
	if bt.Type != TYPE_LINK || bt.Props["resource"] != "users" {
		t.Errorf("Link field mismatch")
	}

	hm := Collection("Posts", "posts")
	if hm.Type != TYPE_COLLECTION || hm.Props["resource"] != "posts" || hm.Context != HIDE_ON_LIST {
		t.Errorf("Collection field mismatch")
	}

	mt := PolyLink("Commentable")
	if mt.Type != TYPE_POLY_LINK {
		t.Errorf("PolyLink field mismatch")
	}
}

// TestSchemaImplementsElement verifies that Schema implements core.Element interface
func TestSchemaImplementsElement(t *testing.T) {
	// Compile-time check - if this compiles, Schema implements Element
	var _ Element = (*Schema)(nil)

	// Runtime check
	schema := &Schema{
		Key:  "test",
		Name: "Test Field",
	}

	if schema.GetKey() != "test" {
		t.Error("Schema does not properly implement GetKey()")
	}

	// Test fluent interface returns Element
	var element Element = schema.SetName("New Name")
	if element == nil {
		t.Error("Fluent methods should return Element interface")
	}
}
