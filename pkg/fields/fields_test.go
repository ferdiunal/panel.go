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

func TestFieldSerialization_IncludesDependencies(t *testing.T) {
	f := Select("City", "city_id").
		DependsOn("country_id")

	data := f.JsonSerialize()

	rawDependsOn, exists := data["depends_on"]
	if !exists {
		t.Fatalf("expected depends_on to be present in serialized field")
	}

	dependsOn, ok := rawDependsOn.([]string)
	if !ok {
		t.Fatalf("expected depends_on to be []string, got %T", rawDependsOn)
	}
	if len(dependsOn) != 1 || dependsOn[0] != "country_id" {
		t.Fatalf("unexpected depends_on value: %v", dependsOn)
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

	matrix := Matrix("Variant Matrix")
	if matrix.Type != TYPE_KEY_VALUE || matrix.View != "matrix-field" {
		t.Errorf("Matrix field mismatch")
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

func TestHideOnApiSetsContext(t *testing.T) {
	f := Text("Secret", "secret").HideOnApi()
	if f.GetContext() != HIDE_ON_API {
		t.Fatalf("expected context %q, got %q", HIDE_ON_API, f.GetContext())
	}
}

func TestFieldSpan(t *testing.T) {
	field, ok := Text("First Name", "first_name").Span(6).(*Schema)
	if !ok {
		t.Fatal("expected span result to be *Schema")
	}
	if got := field.Props["span"]; got != 6 {
		t.Fatalf("expected span to be 6, got %v", got)
	}

	minClamped, ok := Text("Min", "min").Span(0).(*Schema)
	if !ok {
		t.Fatal("expected min span result to be *Schema")
	}
	if got := minClamped.Props["span"]; got != 1 {
		t.Fatalf("expected min-clamped span to be 1, got %v", got)
	}

	maxClamped, ok := Text("Max", "max").Span(99).(*Schema)
	if !ok {
		t.Fatal("expected max span result to be *Schema")
	}
	if got := maxClamped.Props["span"]; got != 12 {
		t.Fatalf("expected max-clamped span to be 12, got %v", got)
	}
}

func TestNumberControlVisibility(t *testing.T) {
	hidden, ok := Number("Price", "price").HideNumberControls().(*Schema)
	if !ok {
		t.Fatal("expected hidden controls result to be *Schema")
	}
	if got := hidden.Props["showControls"]; got != false {
		t.Fatalf("expected showControls to be false, got %v", got)
	}

	shown, ok := Number("Quantity", "quantity").ShowNumberControls(true).(*Schema)
	if !ok {
		t.Fatal("expected shown controls result to be *Schema")
	}
	if got := shown.Props["showControls"]; got != true {
		t.Fatalf("expected showControls to be true, got %v", got)
	}
}

func TestGridVisibilityRules(t *testing.T) {
	hideOnGrid, ok := Text("Image", "image").HideOnGrid().(*Schema)
	if !ok {
		t.Fatal("expected HideOnGrid chain to return *Schema")
	}
	if !hideOnGrid.IsVisibleInContext(ContextIndex) {
		t.Fatal("HideOnGrid field should remain visible in index/table")
	}
	if hideOnGrid.IsVisibleInContext(ContextGrid) {
		t.Fatal("HideOnGrid field should be hidden in grid")
	}

	showOnGrid, ok := Text("Secret", "secret").HideOnList().ShowOnGrid().(*Schema)
	if !ok {
		t.Fatal("expected ShowOnGrid chain to return *Schema")
	}
	if showOnGrid.IsVisibleInContext(ContextIndex) {
		t.Fatal("HideOnList+ShowOnGrid field should remain hidden in index/table")
	}
	if !showOnGrid.IsVisibleInContext(ContextGrid) {
		t.Fatal("HideOnList+ShowOnGrid field should be visible in grid")
	}

	showOnlyGrid, ok := Text("Grid", "grid_only").ShowOnlyGrid().(*Schema)
	if !ok {
		t.Fatal("expected ShowOnlyGrid chain to return *Schema")
	}
	if !showOnlyGrid.IsVisibleInContext(ContextIndex) {
		t.Fatal("ShowOnlyGrid field should be visible in index/table")
	}
	if !showOnlyGrid.IsVisibleInContext(ContextGrid) {
		t.Fatal("ShowOnlyGrid field should be visible in grid")
	}
	if showOnlyGrid.IsVisibleInContext(ContextDetail) {
		t.Fatal("ShowOnlyGrid field should be hidden in detail")
	}
	if showOnlyGrid.IsVisibleInContext(ContextCreate) {
		t.Fatal("ShowOnlyGrid field should be hidden in create")
	}
	if showOnlyGrid.IsVisibleInContext(ContextUpdate) {
		t.Fatal("ShowOnlyGrid field should be hidden in update")
	}
}
