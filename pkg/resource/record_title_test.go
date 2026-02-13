package resource

import (
	"testing"
)

// TestRecordTitle, RecordTitle metodunun temel işlevselliğini test eder
type TestUser struct {
	ID   int
	Name string
}

type TestPost struct {
	ID    int
	Title string
}

func TestOptimizedBaseRecordTitle(t *testing.T) {
	t.Run("varsayılan title key (id)", func(t *testing.T) {
		r := &OptimizedBase{}
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "1" {
			t.Errorf("Beklenen '1', alınan '%s'", title)
		}
	})

	t.Run("özel title key (name)", func(t *testing.T) {
		r := &OptimizedBase{}
		r.SetRecordTitleKey("name")
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "John Doe" {
			t.Errorf("Beklenen 'John Doe', alınan '%s'", title)
		}
	})

	t.Run("özel title key (title)", func(t *testing.T) {
		r := &OptimizedBase{}
		r.SetRecordTitleKey("title")
		post := &TestPost{ID: 1, Title: "Hello World"}

		title := r.RecordTitle(post)
		if title != "Hello World" {
			t.Errorf("Beklenen 'Hello World', alınan '%s'", title)
		}
	})

	t.Run("özel title fonksiyonu", func(t *testing.T) {
		r := &OptimizedBase{}
		r.SetRecordTitleFunc(func(record any) string {
			user := record.(*TestUser)
			return "User: " + user.Name
		})
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "User: John Doe" {
			t.Errorf("Beklenen 'User: John Doe', alınan '%s'", title)
		}
	})

	t.Run("nil kayıt", func(t *testing.T) {
		r := &OptimizedBase{}
		title := r.RecordTitle(nil)
		if title != "" {
			t.Errorf("Beklenen boş string, alınan '%s'", title)
		}
	})

	t.Run("case-insensitive field arama", func(t *testing.T) {
		r := &OptimizedBase{}
		r.SetRecordTitleKey("NAME") // Büyük harfle
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "John Doe" {
			t.Errorf("Beklenen 'John Doe', alınan '%s'", title)
		}
	})
}

func TestBaseRecordTitle(t *testing.T) {
	t.Run("varsayılan title key (id)", func(t *testing.T) {
		r := &Base{}
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "1" {
			t.Errorf("Beklenen '1', alınan '%s'", title)
		}
	})

	t.Run("özel title key (name)", func(t *testing.T) {
		r := &Base{}
		r.SetRecordTitleKey("name")
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "John Doe" {
			t.Errorf("Beklenen 'John Doe', alınan '%s'", title)
		}
	})

	t.Run("özel title fonksiyonu", func(t *testing.T) {
		r := &Base{}
		r.SetRecordTitleFunc(func(record any) string {
			user := record.(*TestUser)
			return "User: " + user.Name
		})
		user := &TestUser{ID: 1, Name: "John Doe"}

		title := r.RecordTitle(user)
		if title != "User: John Doe" {
			t.Errorf("Beklenen 'User: John Doe', alınan '%s'", title)
		}
	})
}

func TestGetRecordTitleKey(t *testing.T) {
	t.Run("varsayılan değer", func(t *testing.T) {
		r := &OptimizedBase{}
		key := r.GetRecordTitleKey()
		if key != "id" {
			t.Errorf("Beklenen 'id', alınan '%s'", key)
		}
	})

	t.Run("özel değer", func(t *testing.T) {
		r := &OptimizedBase{}
		r.SetRecordTitleKey("name")
		key := r.GetRecordTitleKey()
		if key != "name" {
			t.Errorf("Beklenen 'name', alınan '%s'", key)
		}
	})
}

func TestSetRecordTitleKey(t *testing.T) {
	t.Run("method chaining", func(t *testing.T) {
		r := &OptimizedBase{}
		result := r.SetRecordTitleKey("name")
		if result != r {
			t.Error("SetRecordTitleKey method chaining çalışmıyor")
		}
	})
}
