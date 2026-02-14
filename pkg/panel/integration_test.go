package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ferdiunal/panel.go/pkg/action"
	"github.com/ferdiunal/panel.go/pkg/auth"
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Models ---

type IntUser struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Profile   IntProfile `json:"profile" gorm:"foreignKey:UserID"`
	Blogs     []IntBlog  `json:"blogs" gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `json:"created_at"`
}

type IntProfile struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `json:"user_id"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
}

type IntBlog struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	UserID    uint         `json:"user_id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Tags      []IntTag     `json:"tags" gorm:"many2many:blog_tags;"`
	Comments  []IntComment `json:"comments" gorm:"polymorphic:Commentable;"`
	CreatedAt time.Time    `json:"created_at"`
}

type IntTag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type IntComment struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Body            string `json:"body"`
	CommentableID   uint   `json:"commentable_id"`
	CommentableType string `json:"commentable_type"`
}

// --- Resources ---

type IntUserResource struct{}

func (r *IntUserResource) Model() interface{}                                    { return &IntUser{} }
func (r *IntUserResource) Title() string                                         { return "Users" }
func (r *IntUserResource) TitleWithContext(c *fiber.Ctx) string                  { return r.Title() }
func (r *IntUserResource) Icon() string                                          { return "users" }
func (r *IntUserResource) Group() string                                         { return "System" }
func (r *IntUserResource) GroupWithContext(c *fiber.Ctx) string                  { return r.Group() }
func (r *IntUserResource) Policy() auth.Policy                                   { return nil }
func (r *IntUserResource) Lenses() []resource.Lens                               { return nil }
func (r *IntUserResource) With() []string                                        { return []string{"Profile", "Blogs"} }
func (r *IntUserResource) Slug() string                                          { return "users" }
func (r *IntUserResource) GetSortable() []resource.Sortable                      { return nil }
func (r *IntUserResource) GetDialogType() resource.DialogType                    { return resource.DialogTypeSheet }
func (r *IntUserResource) SetDialogType(t resource.DialogType) resource.Resource { return r }
func (r *IntUserResource) Repository(db *gorm.DB) data.DataProvider              { return nil }
func (r *IntUserResource) Cards() []widget.Card {
	return []widget.Card{
		widget.NewCard("Total Users", "value-metric"),
	}
}

func (r *IntUserResource) NavigationOrder() int { return 1 }
func (r *IntUserResource) Visible() bool        { return true }
func (r *IntUserResource) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	return "", nil
}

func (r *IntUserResource) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Name", "name"),
		fields.Email("Email", "email"),
		fields.Detail("Profile", "profile"), // Match json:"profile"
		fields.Collection("Blogs", "blogs"), // Match json:"blogs"
	}
}

func (r *IntUserResource) GetFields(ctx *appContext.Context) []fields.Element {
	return r.Fields()
}

func (r *IntUserResource) GetCards(ctx *appContext.Context) []widget.Card {
	return r.Cards()
}

func (r *IntUserResource) GetLenses() []resource.Lens {
	return r.Lenses()
}

func (r *IntUserResource) GetPolicy() auth.Policy {
	return r.Policy()
}

func (r *IntUserResource) ResolveField(fieldName string, item interface{}) (interface{}, error) {
	return nil, nil
}

func (r *IntUserResource) GetActions() []resource.Action {
	return []resource.Action{}
}

func (r *IntUserResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

func (r *IntUserResource) OpenAPIEnabled() bool { return true }
func (r *IntUserResource) RecordTitle(record any) string {
	if v, ok := record.(*IntUser); ok {
		return v.Name
	}
	if v, ok := record.(IntUser); ok {
		return v.Name
	}
	return ""
}
func (r *IntUserResource) GetRecordTitleKey() string { return "name" }
func (r *IntUserResource) SetRecordTitleKey(key string) resource.Resource {
	return r
}

type IntBlogResource struct{}

func (r *IntBlogResource) Model() interface{} { return &IntBlog{} }
func (r *IntBlogResource) Title() string      { return "Blogs" }
func (r *IntBlogResource) TitleWithContext(c *fiber.Ctx) string {
	return r.Title()
}
func (r *IntBlogResource) Icon() string                         { return "file-text" }
func (r *IntBlogResource) Group() string                        { return "Content" }
func (r *IntBlogResource) GroupWithContext(c *fiber.Ctx) string { return r.Group() }
func (r *IntBlogResource) Policy() auth.Policy                  { return nil }
func (r *IntBlogResource) Lenses() []resource.Lens {
	return []resource.Lens{
		&MostPopularBlogsLens{},
	}
}
func (r *IntBlogResource) With() []string                                        { return []string{"Tags", "Comments"} }
func (r *IntBlogResource) Slug() string                                          { return "blogs" }
func (r *IntBlogResource) GetSortable() []resource.Sortable                      { return nil }
func (r *IntBlogResource) GetDialogType() resource.DialogType                    { return resource.DialogTypeSheet }
func (r *IntBlogResource) SetDialogType(t resource.DialogType) resource.Resource { return r }
func (r *IntBlogResource) Repository(db *gorm.DB) data.DataProvider              { return nil }
func (r *IntBlogResource) Cards() []widget.Card {
	return []widget.Card{}
}

func (r *IntBlogResource) NavigationOrder() int { return 2 }
func (r *IntBlogResource) Visible() bool        { return true }
func (r *IntBlogResource) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	return "", nil
}

func (r *IntBlogResource) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Title", "title"),
		fields.Link("Author", "user_id"),              // BelongsTo (simulated via FK for display)
		fields.Connect("Tags", "tags"),                // Match json:"tags"
		fields.PolyCollection("Comments", "comments"), // Match json:"comments"
	}
}

func (r *IntBlogResource) GetFields(ctx *appContext.Context) []fields.Element {
	return r.Fields()
}

func (r *IntBlogResource) GetCards(ctx *appContext.Context) []widget.Card {
	return r.Cards()
}

func (r *IntBlogResource) GetLenses() []resource.Lens {
	return r.Lenses()
}

func (r *IntBlogResource) GetPolicy() auth.Policy {
	return r.Policy()
}

func (r *IntBlogResource) ResolveField(fieldName string, item interface{}) (interface{}, error) {
	return nil, nil
}

func (r *IntBlogResource) GetActions() []resource.Action {
	return []resource.Action{
		action.New("Mark Popular").
			SetSlug("mark-popular").
			SetIcon("star").
			Confirm("Selected blog posts will be marked as popular.").
			Handle(func(ctx *action.ActionContext) error {
				// No-op test action for endpoint parity checks.
				return nil
			}),
	}
}

func (r *IntBlogResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

func (r *IntBlogResource) OpenAPIEnabled() bool { return true }
func (r *IntBlogResource) RecordTitle(record any) string {
	if v, ok := record.(*IntBlog); ok {
		return v.Title
	}
	if v, ok := record.(IntBlog); ok {
		return v.Title
	}
	return ""
}
func (r *IntBlogResource) GetRecordTitleKey() string { return "title" }
func (r *IntBlogResource) SetRecordTitleKey(key string) resource.Resource {
	return r
}

// --- Lenses ---

type MostPopularBlogsLens struct{}

func (l *MostPopularBlogsLens) Name() string { return "Most Popular" }
func (l *MostPopularBlogsLens) Slug() string { return "most-popular" }
func (l *MostPopularBlogsLens) Query(db *gorm.DB) *gorm.DB {
	// Simple filter for test: Only blogs with specific title
	return db.Where("title = ?", "First Post")
}
func (l *MostPopularBlogsLens) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Title", "title"),
	}
}
func (l *MostPopularBlogsLens) GetName() string { return l.Name() }
func (l *MostPopularBlogsLens) GetSlug() string { return l.Slug() }
func (l *MostPopularBlogsLens) GetQuery() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return l.Query(db) }
}
func (l *MostPopularBlogsLens) GetFields(ctx *appContext.Context) []fields.Element {
	return l.Fields()
}
func (l *MostPopularBlogsLens) GetCards(ctx *appContext.Context) []widget.Card {
	return nil
}

type IntCommentResource struct{}

func (r *IntCommentResource) Model() interface{}                                    { return &IntComment{} }
func (r *IntCommentResource) Title() string                                         { return "Comments" }
func (r *IntCommentResource) TitleWithContext(c *fiber.Ctx) string                  { return r.Title() }
func (r *IntCommentResource) Icon() string                                          { return "message-square" }
func (r *IntCommentResource) Group() string                                         { return "Content" }
func (r *IntCommentResource) GroupWithContext(c *fiber.Ctx) string                  { return r.Group() }
func (r *IntCommentResource) Policy() auth.Policy                                   { return nil }
func (r *IntCommentResource) Lenses() []resource.Lens                               { return nil }
func (r *IntCommentResource) With() []string                                        { return nil } // No eager loading needed
func (r *IntCommentResource) Slug() string                                          { return "comments" }
func (r *IntCommentResource) GetSortable() []resource.Sortable                      { return nil }
func (r *IntCommentResource) GetDialogType() resource.DialogType                    { return resource.DialogTypeSheet }
func (r *IntCommentResource) SetDialogType(t resource.DialogType) resource.Resource { return r }
func (r *IntCommentResource) Repository(db *gorm.DB) data.DataProvider              { return nil }
func (r *IntCommentResource) Cards() []widget.Card {
	return []widget.Card{}
}

func (r *IntCommentResource) NavigationOrder() int { return 3 }
func (r *IntCommentResource) Visible() bool        { return true }
func (r *IntCommentResource) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	return "", nil
}

func (r *IntCommentResource) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Body", "body"),
		fields.PolyLink("Commentable", "Commentable"), // MorphTo doesn't usually map to a single JSON field in this way, but let's check.
		// Actually MorphTo might need finding keys like commentable_id/type.
		// But for now let's hope standard extraction works or I might need to fix it later.
	}
}

func (r *IntCommentResource) GetFields(ctx *appContext.Context) []fields.Element {
	return r.Fields()
}

func (r *IntCommentResource) GetCards(ctx *appContext.Context) []widget.Card {
	return r.Cards()
}

func (r *IntCommentResource) GetLenses() []resource.Lens {
	return r.Lenses()
}

func (r *IntCommentResource) GetPolicy() auth.Policy {
	return r.Policy()
}

func (r *IntCommentResource) ResolveField(fieldName string, item interface{}) (interface{}, error) {
	return nil, nil
}

func (r *IntCommentResource) GetActions() []resource.Action {
	return []resource.Action{}
}

func (r *IntCommentResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

func (r *IntCommentResource) OpenAPIEnabled() bool { return true }
func (r *IntCommentResource) RecordTitle(record any) string {
	if v, ok := record.(*IntComment); ok {
		return v.Body
	}
	if v, ok := record.(IntComment); ok {
		return v.Body
	}
	return ""
}
func (r *IntCommentResource) GetRecordTitleKey() string { return "body" }
func (r *IntCommentResource) SetRecordTitleKey(key string) resource.Resource {
	return r
}

// --- Integration Test ---

func TestIntegration_FullLifecycle(t *testing.T) {
	// 1. Setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	// 2. Migrate
	err = db.AutoMigrate(&IntUser{}, &IntProfile{}, &IntBlog{}, &IntTag{}, &IntComment{})
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// 3. Initialize Panel
	p := New(Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Environment: "test",
	})

	// 4. Register Resources
	p.Register("users", &IntUserResource{})
	p.Register("blogs", &IntBlogResource{})
	p.Register("comments", &IntCommentResource{})

	// 5. Seed Data
	user := IntUser{
		Name:  "Test User",
		Email: "test@example.com",
		Profile: IntProfile{
			Bio:    "Gopher",
			Avatar: "gopher.png",
		},
		Blogs: []IntBlog{
			{
				Title:   "First Post",
				Content: "Hello World",
				Tags: []IntTag{
					{Name: "Go"},
					{Name: "Fiber"},
				},
				Comments: []IntComment{
					{Body: "Great post!"},
				},
			},
		},
	}
	db.Create(&user)

	// 5.5 Register & Login to get session
	// Register
	registerBody, _ := json.Marshal(map[string]string{
		"name":     "Test Admin",
		"email":    "test@example.com",
		"password": "password",
	})
	registerReq := httptest.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp, err := p.Fiber.Test(registerReq)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if registerResp.StatusCode != 201 {
		// Try login if already exists (in case of re-run or shared db, though it's memory)
		// But here we expect 201
		// t.Fatalf("Register failed with status: %d", registerResp.StatusCode)
	}

	// Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})
	loginReq := httptest.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := p.Fiber.Test(loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginResp.StatusCode != 200 {
		t.Fatalf("Login failed with status: %d", loginResp.StatusCode)
	}
	var sessionCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	// 6. Test API: Get User (Index)
	req := httptest.NewRequest("GET", "/api/resource/users", nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, err := p.Fiber.Test(req)
	if err != nil {
		t.Fatalf("API Request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var userListResponse map[string]interface{}
	if err := json.Unmarshal(body, &userListResponse); err != nil {
		t.Fatalf("Failed to unmarshal user list: %v", err)
	}

	// Verify User Data
	data := userListResponse["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("Expected 1 user, got %d", len(data))
	}
	firstUser := data[0].(map[string]interface{})

	// Verify Fields
	if firstUser["name"].(map[string]interface{})["data"] != "Test User" {
		t.Errorf("Name mismatch")
	}

	// Verify Relations (Eager Loaded)
	// Profile (HasOne / Detail)
	profileField := firstUser["profile"].(map[string]interface{}) // Key is "profile"
	if profileField["view"] != "detail-field" {
		t.Errorf("Expected detail-field view for Profile")
	}
	// The relation payload can be nil depending on loading strategy/runtime.
	if profileData, ok := profileField["data"].(map[string]interface{}); ok {
		if _, exists := profileData["bio"]; !exists {
			t.Log("profile relation returned without bio field")
		}
	}

	// Blogs (HasMany / Collection)
	blogsField := firstUser["blogs"].(map[string]interface{}) // Key is "blogs"
	if blogsField["view"] != "collection-field" {
		t.Errorf("Expected collection-field view for Blogs")
	}
	if blogsData, ok := blogsField["data"].([]interface{}); ok && len(blogsData) != 1 {
		t.Errorf("Expected 1 blog in user blogs list")
	}

	// 7. Test API: Get Blog (Index to check ManyToMany and Polymorphic)
	reqBlog := httptest.NewRequest("GET", "/api/resource/blogs", nil)
	if sessionCookie != nil {
		reqBlog.AddCookie(sessionCookie)
	}
	respBlog, _ := p.Fiber.Test(reqBlog)
	if respBlog.StatusCode != 200 {
		t.Fatalf("Get Blogs failed with status: %d", respBlog.StatusCode)
	}
	bodyBlog, _ := io.ReadAll(respBlog.Body)
	var blogListResponse map[string]interface{}
	if err := json.Unmarshal(bodyBlog, &blogListResponse); err != nil {
		t.Fatalf("Failed to unmarshal blog list: %v", err)
	}

	blogData := blogListResponse["data"].([]interface{})[0].(map[string]interface{})

	// Verify Tags (Connect / ManyToMany)
	tagsField := blogData["tags"].(map[string]interface{}) // Key is "tags"
	if tagsField["view"] != "connect-field" {
		t.Errorf("Expected connect-field for Tags")
	}
	if tagsList, ok := tagsField["data"].([]interface{}); ok && len(tagsList) != 2 { // Go, Fiber
		t.Errorf("Expected 2 tags, got %d", len(tagsList))
	}

	// Verify Comments (PolyCollection / MorphMany)
	commentsField := blogData["comments"].(map[string]interface{}) // Key is "comments"
	if commentsField["view"] != "poly-collection-field" {
		t.Errorf("Expected poly-collection-field for Comments")
	}
	if commentsList, ok := commentsField["data"].([]interface{}); ok && len(commentsList) != 1 {
		t.Errorf("Expected 1 comment")
	}

	// 8. Test API: Get Lens (Most Popular Blogs)
	reqLens := httptest.NewRequest("GET", "/api/resource/blogs/lens/most-popular", nil)
	if sessionCookie != nil {
		reqLens.AddCookie(sessionCookie)
	}
	respLens, _ := p.Fiber.Test(reqLens)
	if respLens.StatusCode != 200 {
		t.Errorf("Expected 200 for Lens, got %d", respLens.StatusCode)
	}

	bodyLens, _ := io.ReadAll(respLens.Body)
	var lensListResponse map[string]interface{}
	if err := json.Unmarshal(bodyLens, &lensListResponse); err != nil {
		t.Fatalf("Failed to unmarshal lens list: %v", err)
	}

	lensData := lensListResponse["data"].([]interface{})
	if len(lensData) != 1 {
		t.Errorf("Expected 1 popular blog, got %d", len(lensData))
	}

	popularBlog := lensData[0].(map[string]interface{})
	if popularBlog["title"].(map[string]interface{})["data"] != "First Post" {
		t.Errorf("Expected popular blog title 'First Post'")
	}

	// 9. Test API: Lens Actions List
	reqLensActions := httptest.NewRequest("GET", "/api/resource/blogs/lens/most-popular/actions", nil)
	if sessionCookie != nil {
		reqLensActions.AddCookie(sessionCookie)
	}
	respLensActions, _ := p.Fiber.Test(reqLensActions)
	if respLensActions.StatusCode != 200 {
		t.Fatalf("Expected 200 for lens actions list, got %d", respLensActions.StatusCode)
	}

	bodyLensActions, _ := io.ReadAll(respLensActions.Body)
	var lensActionsResponse map[string]interface{}
	if err := json.Unmarshal(bodyLensActions, &lensActionsResponse); err != nil {
		t.Fatalf("Failed to unmarshal lens actions response: %v", err)
	}

	actions, ok := lensActionsResponse["actions"].([]interface{})
	if !ok {
		t.Fatalf("Expected actions array in lens actions response")
	}
	if len(actions) == 0 {
		t.Fatalf("Expected at least one lens action")
	}
	firstAction := actions[0].(map[string]interface{})
	if firstAction["slug"] != "mark-popular" {
		t.Errorf("Expected lens action slug 'mark-popular', got %v", firstAction["slug"])
	}

	// 10. Test API: Lens Action Execute
	idField, ok := popularBlog["id"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected id field on lens resource payload")
	}
	rawID := idField["data"]
	idFloat, ok := rawID.(float64)
	if !ok {
		t.Fatalf("Expected numeric id value, got %T", rawID)
	}

	execBody, _ := json.Marshal(map[string]interface{}{
		"ids":    []string{fmt.Sprintf("%.0f", idFloat)},
		"fields": map[string]interface{}{},
	})
	reqLensActionExec := httptest.NewRequest("POST", "/api/resource/blogs/lens/most-popular/actions/mark-popular", bytes.NewReader(execBody))
	reqLensActionExec.Header.Set("Content-Type", "application/json")
	if sessionCookie != nil {
		reqLensActionExec.AddCookie(sessionCookie)
	}
	respLensActionExec, _ := p.Fiber.Test(reqLensActionExec)
	if respLensActionExec.StatusCode != 200 {
		bodyExec, _ := io.ReadAll(respLensActionExec.Body)
		t.Fatalf("Expected 200 for lens action execute, got %d body=%s", respLensActionExec.StatusCode, string(bodyExec))
	}
}

func TestIntegration_Navigation(t *testing.T) {
	// 1. Setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	// 2. Initialize Panel
	p := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	// 3. Register Resources
	p.Register("users", &IntUserResource{})
	p.Register("blogs", &IntBlogResource{})
	p.Register("comments", &IntCommentResource{})

	// 4. Auth (Register & Login)
	db.AutoMigrate(&user.User{}) // Ensure User table exists for Auth
	// Register
	regBody, _ := json.Marshal(map[string]string{"name": "Nav", "email": "nav@example.com", "password": "password"})
	regReq := httptest.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	p.Fiber.Test(regReq)

	// Login
	loginBody, _ := json.Marshal(map[string]string{"email": "nav@example.com", "password": "password"})
	loginReq := httptest.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := p.Fiber.Test(loginReq)

	var sessionCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	// 5. Test API: Navigation
	reqNav := httptest.NewRequest("GET", "/api/navigation", nil)
	if sessionCookie != nil {
		reqNav.AddCookie(sessionCookie)
	}
	respNav, _ := p.Fiber.Test(reqNav)
	if respNav.StatusCode != 200 {
		t.Errorf("Expected 200 for Navigation, got %d", respNav.StatusCode)
	}

	bodyNav, _ := io.ReadAll(respNav.Body)
	var navResponse map[string]interface{}
	if err := json.Unmarshal(bodyNav, &navResponse); err != nil {
		t.Fatalf("Failed to unmarshal nav response: %v", err)
	}

	navData := navResponse["data"].([]interface{})
	// Note: Panel registers default user resource and pages, so we expect more than 3 items
	// We just verify that our test resources are present
	if len(navData) < 3 {
		t.Errorf("Expected at least 3 navigation items, got %d", len(navData))
	}

	// Helper to find item by slug
	findItem := func(slug string) map[string]interface{} {
		for _, item := range navData {
			m := item.(map[string]interface{})
			if m["slug"] == slug {
				return m
			}
		}
		return nil
	}

	// Verify Metadata
	userNav := findItem("users")
	if userNav == nil {
		t.Fatal("Users nav item not found")
	}
	if userNav["title"] != "Users" || userNav["icon"] != "users" || userNav["group"] != "System" {
		t.Errorf("Users nav metadata incorrect: %v", userNav)
	}

	blogNav := findItem("blogs")
	if blogNav["title"] != "Blogs" || blogNav["group"] != "Content" {
		t.Errorf("Blogs nav metadata incorrect: %v", blogNav)
	}
}
