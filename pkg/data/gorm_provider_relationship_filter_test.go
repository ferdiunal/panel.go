package data

import (
	"fmt"
	"sort"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RelFilterUser struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Posts []RelFilterPost `gorm:"foreignKey:UserID"`
	Roles []RelFilterRole `gorm:"many2many:rel_filter_user_roles;"`
}

type RelFilterPost struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	Title  string
}

type RelFilterRole struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Users []RelFilterUser `gorm:"many2many:rel_filter_user_roles;"`
}

type RelPolyPost struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Comments []RelPolyComment `gorm:"polymorphic:Commentable;"`
}

type RelPolyComment struct {
	ID              uint `gorm:"primaryKey"`
	Body            string
	CommentableID   uint
	CommentableType string
}

type RelLegacyCategory struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type RelLegacyProduct struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	CategoryID uint
	Category   RelLegacyCategory `gorm:"foreignKey:CategoryID"`
}

type RelMorphTag struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type RelMorphProduct struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type RelMorphTaggable struct {
	TagID        uint   `gorm:"column:tag_id"`
	TaggableID   uint   `gorm:"column:taggable_id"`
	TaggableType string `gorm:"column:taggable_type"`
}

func (RelMorphTaggable) TableName() string {
	return "rel_morph_taggables"
}

func newRelationshipFilterTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect sqlite in-memory db: %v", err)
	}

	return db
}

func tableNameForModel(t *testing.T, db *gorm.DB, model interface{}) string {
	t.Helper()

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil || stmt.Schema == nil {
		t.Fatalf("failed to parse model schema: %v", err)
	}
	return stmt.Schema.Table
}

func TestGormDataProvider_IndexViaRelationship_HasMany(t *testing.T) {
	db := newRelationshipFilterTestDB(t)
	if err := db.AutoMigrate(&RelFilterUser{}, &RelFilterPost{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	user1 := RelFilterUser{Name: "u1"}
	user2 := RelFilterUser{Name: "u2"}
	if err := db.Create(&user1).Error; err != nil {
		t.Fatalf("failed to create user1: %v", err)
	}
	if err := db.Create(&user2).Error; err != nil {
		t.Fatalf("failed to create user2: %v", err)
	}

	posts := []RelFilterPost{
		{UserID: user1.ID, Title: "p1"},
		{UserID: user1.ID, Title: "p2"},
		{UserID: user2.ID, Title: "p3"},
	}
	if err := db.Create(&posts).Error; err != nil {
		t.Fatalf("failed to create posts: %v", err)
	}

	provider := NewGormDataProvider(db, &RelFilterPost{})
	resp, err := provider.Index(nil, QueryRequest{
		Page:            1,
		PerPage:         50,
		Sorts:           []Sort{{Column: "id", Direction: "asc"}},
		ViaResource:     tableNameForModel(t, db, &RelFilterUser{}),
		ViaResourceId:   fmt.Sprint(user1.ID),
		ViaRelationship: "posts",
		ViaParentModel:  &RelFilterUser{},
	})
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}

	if resp.Total != 2 {
		t.Fatalf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	for _, item := range resp.Items {
		post, ok := item.(*RelFilterPost)
		if !ok {
			t.Fatalf("expected *RelFilterPost item, got %T", item)
		}
		if post.UserID != user1.ID {
			t.Fatalf("expected post user_id=%d, got %d", user1.ID, post.UserID)
		}
	}
}

func TestGormDataProvider_IndexViaRelationship_ManyToMany(t *testing.T) {
	db := newRelationshipFilterTestDB(t)
	if err := db.AutoMigrate(&RelFilterUser{}, &RelFilterRole{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	roles := []RelFilterRole{
		{Name: "r1"},
		{Name: "r2"},
		{Name: "r3"},
	}
	if err := db.Create(&roles).Error; err != nil {
		t.Fatalf("failed to create roles: %v", err)
	}

	user1 := RelFilterUser{Name: "u1"}
	user2 := RelFilterUser{Name: "u2"}
	if err := db.Create(&user1).Error; err != nil {
		t.Fatalf("failed to create user1: %v", err)
	}
	if err := db.Create(&user2).Error; err != nil {
		t.Fatalf("failed to create user2: %v", err)
	}

	if err := db.Model(&user1).Association("Roles").Append(&roles[0], &roles[1]); err != nil {
		t.Fatalf("failed to append user1 roles: %v", err)
	}
	if err := db.Model(&user2).Association("Roles").Append(&roles[1], &roles[2]); err != nil {
		t.Fatalf("failed to append user2 roles: %v", err)
	}

	provider := NewGormDataProvider(db, &RelFilterRole{})
	resp, err := provider.Index(nil, QueryRequest{
		Page:            1,
		PerPage:         50,
		Sorts:           []Sort{{Column: "id", Direction: "asc"}},
		ViaResource:     tableNameForModel(t, db, &RelFilterUser{}),
		ViaResourceId:   fmt.Sprint(user1.ID),
		ViaRelationship: "roles",
		ViaParentModel:  &RelFilterUser{},
	})
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}

	if resp.Total != 2 {
		t.Fatalf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	gotIDs := make([]uint, 0, len(resp.Items))
	for _, item := range resp.Items {
		role, ok := item.(*RelFilterRole)
		if !ok {
			t.Fatalf("expected *RelFilterRole item, got %T", item)
		}
		gotIDs = append(gotIDs, role.ID)
	}
	sort.Slice(gotIDs, func(i, j int) bool { return gotIDs[i] < gotIDs[j] })
	if len(gotIDs) != 2 || gotIDs[0] != roles[0].ID || gotIDs[1] != roles[1].ID {
		t.Fatalf("expected role ids [%d %d], got %v", roles[0].ID, roles[1].ID, gotIDs)
	}
}

func TestGormDataProvider_IndexViaRelationship_PolymorphicHasMany(t *testing.T) {
	db := newRelationshipFilterTestDB(t)
	if err := db.AutoMigrate(&RelPolyPost{}, &RelPolyComment{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	post1 := RelPolyPost{Title: "p1"}
	post2 := RelPolyPost{Title: "p2"}
	if err := db.Create(&post1).Error; err != nil {
		t.Fatalf("failed to create post1: %v", err)
	}
	if err := db.Create(&post2).Error; err != nil {
		t.Fatalf("failed to create post2: %v", err)
	}

	polyType := tableNameForModel(t, db, &RelPolyPost{})
	comments := []RelPolyComment{
		{Body: "c1", CommentableID: post1.ID, CommentableType: polyType},
		{Body: "c2", CommentableID: post1.ID, CommentableType: polyType},
		{Body: "c3", CommentableID: post2.ID, CommentableType: polyType},
		{Body: "c4", CommentableID: post1.ID, CommentableType: "other_type"},
	}
	if err := db.Create(&comments).Error; err != nil {
		t.Fatalf("failed to create comments: %v", err)
	}

	provider := NewGormDataProvider(db, &RelPolyComment{})
	resp, err := provider.Index(nil, QueryRequest{
		Page:            1,
		PerPage:         50,
		Sorts:           []Sort{{Column: "id", Direction: "asc"}},
		ViaResource:     polyType,
		ViaResourceId:   fmt.Sprint(post1.ID),
		ViaRelationship: "comments",
		ViaParentModel:  &RelPolyPost{},
	})
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}

	if resp.Total != 2 {
		t.Fatalf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	for _, item := range resp.Items {
		comment, ok := item.(*RelPolyComment)
		if !ok {
			t.Fatalf("expected *RelPolyComment item, got %T", item)
		}
		if comment.CommentableID != post1.ID {
			t.Fatalf("expected commentable_id=%d, got %d", post1.ID, comment.CommentableID)
		}
		if comment.CommentableType != polyType {
			t.Fatalf("expected commentable_type=%s, got %s", polyType, comment.CommentableType)
		}
	}
}

func TestGormDataProvider_IndexViaRelationship_MorphToManyConfig(t *testing.T) {
	db := newRelationshipFilterTestDB(t)
	if err := db.AutoMigrate(&RelMorphTag{}, &RelMorphProduct{}, &RelMorphTaggable{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	tag1 := RelMorphTag{Name: "t1"}
	tag2 := RelMorphTag{Name: "t2"}
	if err := db.Create(&tag1).Error; err != nil {
		t.Fatalf("failed to create tag1: %v", err)
	}
	if err := db.Create(&tag2).Error; err != nil {
		t.Fatalf("failed to create tag2: %v", err)
	}

	product1 := RelMorphProduct{Name: "p1"}
	product2 := RelMorphProduct{Name: "p2"}
	product3 := RelMorphProduct{Name: "p3"}
	if err := db.Create(&product1).Error; err != nil {
		t.Fatalf("failed to create product1: %v", err)
	}
	if err := db.Create(&product2).Error; err != nil {
		t.Fatalf("failed to create product2: %v", err)
	}
	if err := db.Create(&product3).Error; err != nil {
		t.Fatalf("failed to create product3: %v", err)
	}

	entries := []RelMorphTaggable{
		{TagID: tag1.ID, TaggableID: product1.ID, TaggableType: "products"},
		{TagID: tag1.ID, TaggableID: product2.ID, TaggableType: "products"},
		{TagID: tag2.ID, TaggableID: product3.ID, TaggableType: "products"},
		{TagID: tag1.ID, TaggableID: product3.ID, TaggableType: "shipments"},
	}
	if err := db.Create(&entries).Error; err != nil {
		t.Fatalf("failed to create taggables: %v", err)
	}

	provider := NewGormDataProvider(db, &RelMorphProduct{})
	resp, err := provider.Index(nil, QueryRequest{
		Page:          1,
		PerPage:       50,
		Sorts:         []Sort{{Column: "id", Direction: "asc"}},
		ViaResource:   tableNameForModel(t, db, &RelMorphTag{}),
		ViaResourceId: fmt.Sprint(tag1.ID),
		ViaRelationshipConfig: &ViaRelationshipConfig{
			PivotTable:        "rel_morph_taggables",
			ParentPivotColumn: "tag_id",
			ChildPivotColumn:  "taggable_id",
			MorphTypeColumn:   "taggable_type",
			MorphTypeValue:    "products",
		},
	})
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}

	if resp.Total != 2 {
		t.Fatalf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	gotIDs := make([]uint, 0, len(resp.Items))
	for _, item := range resp.Items {
		product, ok := item.(*RelMorphProduct)
		if !ok {
			t.Fatalf("expected *RelMorphProduct item, got %T", item)
		}
		gotIDs = append(gotIDs, product.ID)
	}
	sort.Slice(gotIDs, func(i, j int) bool { return gotIDs[i] < gotIDs[j] })
	if len(gotIDs) != 2 || gotIDs[0] != product1.ID || gotIDs[1] != product2.ID {
		t.Fatalf("expected product ids [%d %d], got %v", product1.ID, product2.ID, gotIDs)
	}
}

func TestGormDataProvider_IndexViaRelationship_FallbackBelongsTo(t *testing.T) {
	db := newRelationshipFilterTestDB(t)
	if err := db.AutoMigrate(&RelLegacyCategory{}, &RelLegacyProduct{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	category1 := RelLegacyCategory{Name: "c1"}
	category2 := RelLegacyCategory{Name: "c2"}
	if err := db.Create(&category1).Error; err != nil {
		t.Fatalf("failed to create category1: %v", err)
	}
	if err := db.Create(&category2).Error; err != nil {
		t.Fatalf("failed to create category2: %v", err)
	}

	products := []RelLegacyProduct{
		{Name: "p1", CategoryID: category1.ID},
		{Name: "p2", CategoryID: category1.ID},
		{Name: "p3", CategoryID: category2.ID},
	}
	if err := db.Create(&products).Error; err != nil {
		t.Fatalf("failed to create products: %v", err)
	}

	provider := NewGormDataProvider(db, &RelLegacyProduct{})
	resp, err := provider.Index(nil, QueryRequest{
		Page:          1,
		PerPage:       50,
		Sorts:         []Sort{{Column: "id", Direction: "asc"}},
		ViaResource:   tableNameForModel(t, db, &RelLegacyCategory{}),
		ViaResourceId: fmt.Sprint(category1.ID),
	})
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}

	if resp.Total != 2 {
		t.Fatalf("expected total 2, got %d", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	for _, item := range resp.Items {
		product, ok := item.(*RelLegacyProduct)
		if !ok {
			t.Fatalf("expected *RelLegacyProduct item, got %T", item)
		}
		if product.CategoryID != category1.ID {
			t.Fatalf("expected category_id=%d, got %d", category1.ID, product.CategoryID)
		}
	}
}
