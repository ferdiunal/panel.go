// Package fields, test senaryoları için yardımcı fonksiyonlar ve yapılar sağlar.
//
// Bu paket, field sisteminin test edilmesi için gerekli olan veritabanı bağlantıları,
// test verileri ve model yapılarını içerir. SQLite in-memory veritabanı kullanarak
// hızlı ve izole test ortamları oluşturur.
//
// # Temel Özellikler
//
// - In-memory SQLite veritabanı desteği
// - Otomatik tablo oluşturma ve veri doldurma
// - İlişkisel veri modelleri (HasOne, HasMany, BelongsTo, BelongsToMany, MorphTo)
// - Test için hazır örnek veriler
//
// # Kullanım Senaryoları
//
// - Unit testlerde veritabanı bağlantısı gerektiğinde
// - İlişkisel sorguların test edilmesinde
// - Field sisteminin entegrasyon testlerinde
// - Performans testlerinde tutarlı veri seti oluşturmada
//
// # Örnek Kullanım
//
// ```go
// func TestMyFeature(t *testing.T) {
//     // Test veritabanı oluştur
//     testDB := NewTestDB(t)
//     defer testDB.Close()
//
//     // Veritabanını kullan
//     var users []User
//     testDB.DB.Query("SELECT * FROM users")
// }
// ```
package fields

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestDB, test senaryoları için veritabanı bağlantısı sağlayan ana yapıdır.
//
// Bu yapı, SQLite in-memory veritabanı kullanarak izole test ortamları oluşturur.
// Her test için temiz bir veritabanı örneği sağlar ve test sonunda otomatik temizlik yapılmasını kolaylaştırır.
//
// # Özellikler
//
// - In-memory SQLite veritabanı (disk I/O yok, hızlı testler)
// - Otomatik tablo oluşturma
// - Önceden tanımlı test verileri
// - Test context'i ile entegrasyon
//
// # Avantajlar
//
// - **Hız**: In-memory veritabanı sayesinde çok hızlı test çalıştırma
// - **İzolasyon**: Her test kendi veritabanı örneğine sahip
// - **Temizlik**: Test sonunda otomatik temizleme
// - **Tutarlılık**: Her test aynı başlangıç durumu ile çalışır
//
// # Dezavantajlar
//
// - SQLite sınırlamaları (bazı SQL özellikleri desteklenmez)
// - Gerçek veritabanı motorlarından farklı davranışlar gösterebilir
// - Büyük veri setleri için bellek kullanımı artabilir
//
// # Önemli Notlar
//
// - Her test için yeni bir TestDB örneği oluşturun
// - Test sonunda mutlaka Close() metodunu çağırın (defer kullanın)
// - In-memory veritabanı, bağlantı kapandığında tüm verileri kaybeder
//
// # Örnek Kullanım
//
// ```go
// func TestUserQueries(t *testing.T) {
//     testDB := NewTestDB(t)
//     defer testDB.Close()
//
//     // Test kodunuz
//     rows, err := testDB.DB.Query("SELECT * FROM users WHERE id = ?", 1)
//     if err != nil {
//         t.Fatal(err)
//     }
//     defer rows.Close()
// }
// ```
type TestDB struct {
	// DB, aktif SQLite veritabanı bağlantısını tutar.
	// Bu bağlantı üzerinden tüm SQL işlemleri gerçekleştirilir.
	DB *sql.DB

	// t, test context'ini tutar.
	// Hata durumlarında test başarısızlığını bildirmek için kullanılır.
	t  *testing.T
}

// NewTestDB, yeni bir test veritabanı oluşturur ve başlatır.
//
// Bu fonksiyon, SQLite in-memory veritabanı kullanarak tamamen izole bir test ortamı sağlar.
// Veritabanı otomatik olarak tüm gerekli tabloları oluşturur ve test verileri ile doldurur.
//
// # İşlem Adımları
//
// 1. In-memory SQLite veritabanı bağlantısı açar
// 2. Tüm test tablolarını oluşturur (users, posts, comments, tags, profiles, vb.)
// 3. Örnek test verilerini ekler
// 4. Hazır TestDB örneğini döndürür
//
// # Parametreler
//
// - `t`: Test context'i. Hata durumlarında test başarısızlığını bildirmek için kullanılır.
//
// # Dönüş Değeri
//
// Kullanıma hazır TestDB örneği döndürür. Bu örnek üzerinden tüm veritabanı işlemleri yapılabilir.
//
// # Hata Yönetimi
//
// Herhangi bir hata durumunda (veritabanı açılamaz, tablolar oluşturulamaz, veri eklenemez)
// test otomatik olarak başarısız sayılır (t.Fatalf kullanılır).
//
// # Önemli Notlar
//
// - **Bellek Kullanımı**: In-memory veritabanı RAM'de çalışır, disk I/O yoktur
// - **İzolasyon**: Her çağrı tamamen bağımsız bir veritabanı oluşturur
// - **Temizlik**: Mutlaka defer ile Close() çağrılmalıdır
// - **Performans**: Disk tabanlı veritabanlarından çok daha hızlıdır
//
// # Kullanım Örnekleri
//
// ```go
// // Temel kullanım
// func TestBasicQuery(t *testing.T) {
//     testDB := NewTestDB(t)
//     defer testDB.Close()
//
//     var count int
//     err := testDB.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
//     if err != nil {
//         t.Fatal(err)
//     }
//     if count != 3 {
//         t.Errorf("Expected 3 users, got %d", count)
//     }
// }
//
// // İlişkisel sorgu testi
// func TestRelationships(t *testing.T) {
//     testDB := NewTestDB(t)
//     defer testDB.Close()
//
//     // User'ın post'larını getir
//     rows, err := testDB.DB.Query(`
//         SELECT p.title FROM posts p
//         JOIN users u ON p.user_id = u.id
//         WHERE u.email = ?
//     `, "john@example.com")
//     if err != nil {
//         t.Fatal(err)
//     }
//     defer rows.Close()
// }
// ```
//
// # Test Verileri
//
// Fonksiyon aşağıdaki örnek verileri otomatik olarak oluşturur:
// - 3 kullanıcı (John Doe, Jane Smith, Bob Johnson)
// - 3 gönderi (farklı kullanıcılara ait)
// - 3 yorum
// - 4 etiket (golang, database, testing, api)
// - 2 profil
// - Çeşitli ilişkisel veriler (post-tag, taggables)
func NewTestDB(t *testing.T) *TestDB {
	// Use in-memory SQLite database for tests
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	if err := createTestTables(db); err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Seed test data
	if err := seedTestData(db); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	return &TestDB{
		DB: db,
		t:  t,
	}
}

// Close, test veritabanı bağlantısını kapatır ve kaynakları serbest bırakır.
//
// Bu metod, test sonunda veritabanı bağlantısını düzgün bir şekilde kapatır.
// In-memory veritabanı kullanıldığı için, bağlantı kapandığında tüm veriler kaybolur.
//
// # Dönüş Değeri
//
// Bağlantı kapatma işlemi sırasında oluşan hata varsa döndürür, yoksa nil döner.
//
// # Önemli Notlar
//
// - **Mutlaka Çağrılmalı**: Her NewTestDB çağrısından sonra mutlaka Close çağrılmalıdır
// - **Defer Kullanımı**: En iyi pratik, defer ile kullanmaktır
// - **Veri Kaybı**: In-memory veritabanı olduğu için, Close sonrası tüm veriler kaybolur
// - **Kaynak Sızıntısı**: Close çağrılmazsa veritabanı bağlantısı açık kalır
//
// # Kullanım Örnekleri
//
// ```go
// // Önerilen kullanım (defer ile)
// func TestExample(t *testing.T) {
//     testDB := NewTestDB(t)
//     defer testDB.Close() // Test bitince otomatik kapanır
//
//     // Test kodunuz...
// }
//
// // Manuel kullanım (önerilmez)
// func TestManual(t *testing.T) {
//     testDB := NewTestDB(t)
//
//     // Test kodunuz...
//
//     if err := testDB.Close(); err != nil {
//         t.Errorf("Failed to close database: %v", err)
//     }
// }
// ```
func (tdb *TestDB) Close() error {
	return tdb.DB.Close()
}

// createTestTables, test veritabanında gerekli tüm tabloları oluşturur.
//
// Bu fonksiyon, field sisteminin tüm ilişki türlerini test edebilmek için
// kapsamlı bir veritabanı şeması oluşturur. Tablolar, gerçek dünya senaryolarını
// simüle edecek şekilde tasarlanmıştır.
//
// # Oluşturulan Tablolar
//
// 1. **users**: Kullanıcı bilgilerini tutar (id, name, email, created_at)
// 2. **posts**: Blog gönderilerini tutar (id, user_id, title, content, created_at)
// 3. **comments**: Gönderi yorumlarını tutar (id, post_id, user_id, content, created_at)
// 4. **tags**: Etiketleri tutar (id, name, created_at)
// 5. **post_tag**: Post-Tag çoka-çok ilişkisi için pivot tablo (post_id, tag_id)
// 6. **taggables**: Polimorfik ilişkiler için tablo (id, taggable_id, taggable_type, tag_id)
// 7. **profiles**: Kullanıcı profilleri için tablo (id, user_id, bio, avatar_url, created_at)
//
// # İlişki Türleri
//
// - **HasOne**: User -> Profile (bir kullanıcının bir profili var)
// - **HasMany**: User -> Posts (bir kullanıcının birden fazla gönderisi var)
// - **BelongsTo**: Post -> User (bir gönderi bir kullanıcıya ait)
// - **BelongsToMany**: Post <-> Tag (çoka-çok ilişki, pivot tablo ile)
// - **MorphTo**: Taggables (polimorfik ilişki, farklı modellere etiket eklenebilir)
//
// # Parametreler
//
// - `db`: Tabloların oluşturulacağı veritabanı bağlantısı
//
// # Dönüş Değeri
//
// Tablo oluşturma işlemi başarılı ise nil, hata varsa error döner.
//
// # Önemli Notlar
//
// - **Foreign Key Constraints**: Tüm ilişkiler foreign key ile tanımlanmıştır
// - **Unique Constraints**: Email ve user_id gibi alanlar unique olarak işaretlenmiştir
// - **Auto Increment**: Tüm primary key'ler otomatik artan değerlerdir
// - **Timestamps**: created_at alanları otomatik olarak CURRENT_TIMESTAMP ile doldurulur
// - **IF NOT EXISTS**: Tablolar zaten varsa hata vermez
//
// # Hata Yönetimi
//
// Herhangi bir tablo oluşturulamazsa, işlem durur ve hata döndürülür.
// Hata mesajı hangi tablonun oluşturulamadığını belirtir.
//
// # Kullanım Senaryoları
//
// Bu fonksiyon genellikle doğrudan çağrılmaz, NewTestDB tarafından otomatik olarak çalıştırılır.
// Ancak özel test senaryoları için manuel olarak da kullanılabilir:
//
// ```go
// db, _ := sql.Open("sqlite3", ":memory:")
// if err := createTestTables(db); err != nil {
//     log.Fatal(err)
// }
// ```
func createTestTables(db *sql.DB) error {
	tables := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Comments table
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Tags table
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Post-Tag pivot table (BelongsToMany)
		`CREATE TABLE IF NOT EXISTS post_tag (
			post_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)`,

		// Taggable polymorphic table
		`CREATE TABLE IF NOT EXISTS taggables (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			taggable_id INTEGER NOT NULL,
			taggable_type TEXT NOT NULL,
			tag_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)`,

		// Profiles table (HasOne)
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL UNIQUE,
			bio TEXT,
			avatar_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// seedTestData, test veritabanına örnek verileri ekler.
//
// Bu fonksiyon, field sisteminin tüm özelliklerini test edebilmek için
// gerçekçi ve ilişkisel test verileri oluşturur. Veriler, gerçek dünya
// senaryolarını simüle edecek şekilde tasarlanmıştır.
//
// # Eklenen Veriler
//
// **Kullanıcılar (3 adet)**:
// - John Doe (john@example.com)
// - Jane Smith (jane@example.com)
// - Bob Johnson (bob@example.com)
//
// **Gönderiler (3 adet)**:
// - "First Post" (John'a ait)
// - "Second Post" (John'a ait)
// - "Jane's Post" (Jane'e ait)
//
// **Yorumlar (3 adet)**:
// - Post 1'e Jane'den yorum
// - Post 1'e Bob'dan yorum
// - Post 2'ye Jane'den yorum
//
// **Etiketler (4 adet)**:
// - golang
// - database
// - testing
// - api
//
// **Profiller (2 adet)**:
// - John'un profili (Software developer)
// - Jane'in profili (Designer)
//
// **İlişkiler**:
// - Post-Tag ilişkileri (BelongsToMany test için)
// - Taggable ilişkileri (Polimorfik ilişki test için)
//
// # Parametreler
//
// - `db`: Verilerin ekleneceği veritabanı bağlantısı
//
// # Dönüş Değeri
//
// Veri ekleme işlemi başarılı ise nil, hata varsa error döner.
//
// # Önemli Notlar
//
// - **Sıralı Ekleme**: Veriler foreign key constraint'leri nedeniyle belirli sırada eklenir
// - **ID'ler**: Auto-increment kullanıldığı için ID'ler 1'den başlar
// - **İlişkiler**: Tüm foreign key ilişkileri doğru şekilde kurulur
// - **Tutarlılık**: Her test aynı veri seti ile başlar
//
// # Hata Yönetimi
//
// Herhangi bir veri eklenemezse, işlem durur ve hangi verinin eklenemediğini
// belirten bir hata döndürülür.
//
// # Test Senaryoları
//
// Bu veriler aşağıdaki test senaryolarını destekler:
// - HasOne: User -> Profile
// - HasMany: User -> Posts, Post -> Comments
// - BelongsTo: Post -> User, Comment -> User/Post
// - BelongsToMany: Post <-> Tag
// - MorphTo: Taggables (polimorfik ilişkiler)
//
// # Kullanım
//
// Bu fonksiyon genellikle doğrudan çağrılmaz, NewTestDB tarafından otomatik olarak çalıştırılır.
// Ancak özel test senaryoları için manuel olarak da kullanılabilir:
//
// ```go
// db, _ := sql.Open("sqlite3", ":memory:")
// createTestTables(db)
// if err := seedTestData(db); err != nil {
//     log.Fatal(err)
// }
// ```
func seedTestData(db *sql.DB) error {
	// Insert users
	users := []struct {
		name  string
		email string
	}{
		{"John Doe", "john@example.com"},
		{"Jane Smith", "jane@example.com"},
		{"Bob Johnson", "bob@example.com"},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.name, user.email)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	// Insert posts
	posts := []struct {
		userID  int
		title   string
		content string
	}{
		{1, "First Post", "This is the first post"},
		{1, "Second Post", "This is the second post"},
		{2, "Jane's Post", "Jane's first post"},
	}

	for _, post := range posts {
		_, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", post.userID, post.title, post.content)
		if err != nil {
			return fmt.Errorf("failed to insert post: %w", err)
		}
	}

	// Insert comments
	comments := []struct {
		postID  int
		userID  int
		content string
	}{
		{1, 2, "Great post!"},
		{1, 3, "Thanks for sharing"},
		{2, 2, "Interesting"},
	}

	for _, comment := range comments {
		_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", comment.postID, comment.userID, comment.content)
		if err != nil {
			return fmt.Errorf("failed to insert comment: %w", err)
		}
	}

	// Insert tags
	tags := []string{"golang", "database", "testing", "api"}
	for _, tag := range tags {
		_, err := db.Exec("INSERT INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
	}

	// Insert post-tag relationships
	postTags := []struct {
		postID int
		tagID  int
	}{
		{1, 1}, // First post - golang
		{1, 2}, // First post - database
		{2, 1}, // Second post - golang
		{3, 3}, // Jane's post - testing
	}

	for _, pt := range postTags {
		_, err := db.Exec("INSERT INTO post_tag (post_id, tag_id) VALUES (?, ?)", pt.postID, pt.tagID)
		if err != nil {
			return fmt.Errorf("failed to insert post-tag: %w", err)
		}
	}

	// Insert profiles
	profiles := []struct {
		userID    int
		bio       string
		avatarURL string
	}{
		{1, "Software developer", "https://example.com/avatar1.jpg"},
		{2, "Designer", "https://example.com/avatar2.jpg"},
	}

	for _, profile := range profiles {
		_, err := db.Exec("INSERT INTO profiles (user_id, bio, avatar_url) VALUES (?, ?, ?)", profile.userID, profile.bio, profile.avatarURL)
		if err != nil {
			return fmt.Errorf("failed to insert profile: %w", err)
		}
	}

	// Insert taggable relationships (polymorphic)
	taggables := []struct {
		taggableID   int
		taggableType string
		tagID        int
	}{
		{1, "post", 1},    // Post 1 tagged with golang
		{1, "post", 2},    // Post 1 tagged with database
		{1, "comment", 3}, // Comment 1 tagged with testing
	}

	for _, taggable := range taggables {
		_, err := db.Exec("INSERT INTO taggables (taggable_id, taggable_type, tag_id) VALUES (?, ?, ?)", taggable.taggableID, taggable.taggableType, taggable.tagID)
		if err != nil {
			return fmt.Errorf("failed to insert taggable: %w", err)
		}
	}

	return nil
}

// ============================================================================
// Test Model Yapıları
// ============================================================================
//
// Aşağıdaki struct'lar, field sisteminin test edilmesi için kullanılan
// model yapılarıdır. Bu modeller, gerçek dünya uygulamalarındaki veri
// modellerini simüle eder ve tüm ilişki türlerini test etmeyi sağlar.
//
// # Model İlişkileri
//
// - User -> Profile (HasOne)
// - User -> Posts (HasMany)
// - Post -> User (BelongsTo)
// - Post -> Tags (BelongsToMany)
// - Comment -> Post (BelongsTo)
// - Comment -> User (BelongsTo)
// - Taggables (MorphTo - Polimorfik ilişki)

// User, bir kullanıcıyı temsil eden model yapısıdır.
//
// Bu model, sistemdeki kullanıcı bilgilerini tutar ve diğer modeller ile
// ilişkilendirilir. HasOne ve HasMany ilişkilerinin test edilmesinde kullanılır.
//
// # İlişkiler
//
// - **HasOne**: Profile (bir kullanıcının bir profili vardır)
// - **HasMany**: Posts (bir kullanıcının birden fazla gönderisi olabilir)
// - **HasMany**: Comments (bir kullanıcının birden fazla yorumu olabilir)
//
// # JSON Serileştirme
//
// Tüm alanlar JSON olarak serileştirilebilir. API yanıtlarında kullanılabilir.
//
// # Örnek Kullanım
//
// ```go
// var user User
// err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", 1).
//     Scan(&user.ID, &user.Name, &user.Email)
// ```
type User struct {
	// ID, kullanıcının benzersiz kimliğidir (primary key).
	ID    int    `json:"id"`

	// Name, kullanıcının tam adıdır.
	Name  string `json:"name"`

	// Email, kullanıcının benzersiz e-posta adresidir.
	// Veritabanında UNIQUE constraint ile korunur.
	Email string `json:"email"`
}

// Post, bir blog gönderisini temsil eden model yapısıdır.
//
// Bu model, kullanıcılar tarafından oluşturulan içerikleri tutar ve
// BelongsTo, HasMany ve BelongsToMany ilişkilerinin test edilmesinde kullanılır.
//
// # İlişkiler
//
// - **BelongsTo**: User (her gönderi bir kullanıcıya aittir)
// - **HasMany**: Comments (bir gönderinin birden fazla yorumu olabilir)
// - **BelongsToMany**: Tags (bir gönderi birden fazla etikete sahip olabilir)
//
// # JSON Serileştirme
//
// Author alanı opsiyoneldir ve `omitempty` ile işaretlenmiştir.
// İlişki yüklenmemişse JSON'da görünmez.
//
// # Örnek Kullanım
//
// ```go
// var post Post
// err := db.QueryRow(`
//     SELECT p.id, p.user_id, p.title, p.content,
//            u.id, u.name, u.email
//     FROM posts p
//     LEFT JOIN users u ON p.user_id = u.id
//     WHERE p.id = ?
// `, 1).Scan(&post.ID, &post.UserID, &post.Title, &post.Content,
//           &post.Author.ID, &post.Author.Name, &post.Author.Email)
// ```
type Post struct {
	// ID, gönderinin benzersiz kimliğidir (primary key).
	ID      int    `json:"id"`

	// UserID, gönderinin sahibi olan kullanıcının ID'sidir (foreign key).
	UserID  int    `json:"user_id"`

	// Title, gönderinin başlığıdır.
	Title   string `json:"title"`

	// Content, gönderinin içeriğidir.
	Content string `json:"content"`

	// Author, gönderinin sahibi olan kullanıcıdır (BelongsTo ilişkisi).
	// Bu alan opsiyoneldir ve eager loading ile doldurulur.
	Author  *User  `json:"author,omitempty"`
}

// Comment, bir gönderi yorumunu temsil eden model yapısıdır.
//
// Bu model, kullanıcıların gönderilere yaptığı yorumları tutar ve
// birden fazla BelongsTo ilişkisinin test edilmesinde kullanılır.
//
// # İlişkiler
//
// - **BelongsTo**: Post (her yorum bir gönderiye aittir)
// - **BelongsTo**: User (her yorum bir kullanıcıya aittir)
//
// # JSON Serileştirme
//
// Post ve Author alanları opsiyoneldir ve `omitempty` ile işaretlenmiştir.
// İlişkiler yüklenmemişse JSON'da görünmez.
//
// # Örnek Kullanım
//
// ```go
// var comment Comment
// err := db.QueryRow(`
//     SELECT c.id, c.post_id, c.user_id, c.content,
//            p.id, p.title,
//            u.id, u.name
//     FROM comments c
//     LEFT JOIN posts p ON c.post_id = p.id
//     LEFT JOIN users u ON c.user_id = u.id
//     WHERE c.id = ?
// `, 1).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,
//           &comment.Post.ID, &comment.Post.Title,
//           &comment.Author.ID, &comment.Author.Name)
// ```
type Comment struct {
	// ID, yorumun benzersiz kimliğidir (primary key).
	ID      int    `json:"id"`

	// PostID, yorumun yapıldığı gönderinin ID'sidir (foreign key).
	PostID  int    `json:"post_id"`

	// UserID, yorumu yapan kullanıcının ID'sidir (foreign key).
	UserID  int    `json:"user_id"`

	// Content, yorumun içeriğidir.
	Content string `json:"content"`

	// Post, yorumun yapıldığı gönderidir (BelongsTo ilişkisi).
	// Bu alan opsiyoneldir ve eager loading ile doldurulur.
	Post    *Post  `json:"post,omitempty"`

	// Author, yorumu yapan kullanıcıdır (BelongsTo ilişkisi).
	// Bu alan opsiyoneldir ve eager loading ile doldurulur.
	Author  *User  `json:"author,omitempty"`
}

// Tag, bir etiketi temsil eden model yapısıdır.
//
// Bu model, içerikleri kategorize etmek için kullanılan etiketleri tutar ve
// BelongsToMany ilişkisinin test edilmesinde kullanılır.
//
// # İlişkiler
//
// - **BelongsToMany**: Posts (bir etiket birden fazla gönderiye ait olabilir)
// - **MorphTo**: Taggables (polimorfik ilişki ile farklı modellere etiket eklenebilir)
//
// # JSON Serileştirme
//
// Tüm alanlar JSON olarak serileştirilebilir.
//
// # Örnek Kullanım
//
// ```go
// var tag Tag
// err := db.QueryRow("SELECT id, name FROM tags WHERE id = ?", 1).
//     Scan(&tag.ID, &tag.Name)
// ```
type Tag struct {
	// ID, etiketin benzersiz kimliğidir (primary key).
	ID   int    `json:"id"`

	// Name, etiketin adıdır.
	// Veritabanında UNIQUE constraint ile korunur.
	Name string `json:"name"`
}

// Profile, bir kullanıcı profilini temsil eden model yapısıdır.
//
// Bu model, kullanıcıların ek bilgilerini tutar ve HasOne ilişkisinin
// test edilmesinde kullanılır. Her kullanıcının yalnızca bir profili olabilir.
//
// # İlişkiler
//
// - **BelongsTo**: User (her profil bir kullanıcıya aittir)
// - **HasOne** (ters yönde): User -> Profile
//
// # JSON Serileştirme
//
// User alanı opsiyoneldir ve `omitempty` ile işaretlenmiştir.
// İlişki yüklenmemişse JSON'da görünmez.
//
// # Önemli Notlar
//
// - UserID alanı UNIQUE constraint'e sahiptir
// - Bir kullanıcının yalnızca bir profili olabilir
// - Profil olmadan kullanıcı var olabilir (opsiyonel ilişki)
//
// # Örnek Kullanım
//
// ```go
// var profile Profile
// err := db.QueryRow(`
//     SELECT p.id, p.user_id, p.bio, p.avatar_url,
//            u.id, u.name, u.email
//     FROM profiles p
//     LEFT JOIN users u ON p.user_id = u.id
//     WHERE p.user_id = ?
// `, 1).Scan(&profile.ID, &profile.UserID, &profile.Bio, &profile.AvatarURL,
//           &profile.User.ID, &profile.User.Name, &profile.User.Email)
// ```
type Profile struct {
	// ID, profilin benzersiz kimliğidir (primary key).
	ID        int    `json:"id"`

	// UserID, profilin ait olduğu kullanıcının ID'sidir (foreign key, unique).
	UserID    int    `json:"user_id"`

	// Bio, kullanıcının biyografisidir.
	Bio       string `json:"bio"`

	// AvatarURL, kullanıcının profil resminin URL'sidir.
	AvatarURL string `json:"avatar_url"`

	// User, profilin ait olduğu kullanıcıdır (BelongsTo ilişkisi).
	// Bu alan opsiyoneldir ve eager loading ile doldurulur.
	User      *User  `json:"user,omitempty"`
}

// PostTag, gönderi-etiket ilişkisini temsil eden pivot model yapısıdır.
//
// Bu model, Post ve Tag arasındaki çoka-çok (many-to-many) ilişkiyi tutar.
// BelongsToMany ilişkisinin test edilmesinde kullanılır.
//
// # İlişki Türü
//
// Bu bir **pivot tablo** modelidir ve iki model arasındaki çoka-çok ilişkiyi sağlar:
// - Bir gönderi birden fazla etikete sahip olabilir
// - Bir etiket birden fazla gönderide kullanılabilir
//
// # Veritabanı Yapısı
//
// - Composite primary key (post_id, tag_id)
// - Foreign key: post_id -> posts.id
// - Foreign key: tag_id -> tags.id
//
// # Önemli Notlar
//
// - Bu tablo ek veri tutmaz, sadece ilişkiyi tanımlar
// - Composite primary key sayesinde aynı ilişki iki kez eklenemez
// - Pivot tablolar genellikle ek metadata tutabilir (created_at, vb.)
//
// # Örnek Kullanım
//
// ```go
// // Bir gönderiye etiket ekle
// _, err := db.Exec("INSERT INTO post_tag (post_id, tag_id) VALUES (?, ?)", 1, 2)
//
// // Bir gönderinin tüm etiketlerini getir
// rows, err := db.Query(`
//     SELECT t.id, t.name
//     FROM tags t
//     JOIN post_tag pt ON t.id = pt.tag_id
//     WHERE pt.post_id = ?
// `, 1)
// ```
type PostTag struct {
	// PostID, gönderinin ID'sidir (composite primary key'in bir parçası).
	PostID int `json:"post_id"`

	// TagID, etiketin ID'sidir (composite primary key'in bir parçası).
	TagID  int `json:"tag_id"`
}

// Taggable, polimorfik etiketleme ilişkisini temsil eden model yapısıdır.
//
// Bu model, farklı türdeki modellere (Post, Comment, vb.) etiket eklenmesini
// sağlar. MorphTo (polimorfik) ilişkisinin test edilmesinde kullanılır.
//
// # Polimorfik İlişki
//
// Polimorfik ilişki, bir modelin birden fazla farklı model türü ile ilişkilendirilmesini sağlar:
// - Bir etiket bir Post'a eklenebilir (taggable_type = "post")
// - Bir etiket bir Comment'e eklenebilir (taggable_type = "comment")
// - Gelecekte başka model türleri de eklenebilir
//
// # Veritabanı Yapısı
//
// - taggable_id: İlişkilendirilen kaydın ID'si
// - taggable_type: İlişkilendirilen modelin türü (string)
// - tag_id: Etiketin ID'si
//
// # Avantajlar
//
// - **Esneklik**: Yeni model türleri kolayca eklenebilir
// - **Tekrar Kullanılabilirlik**: Aynı etiket sistemi farklı modellerde kullanılabilir
// - **Genişletilebilirlik**: Yeni özellikler kolayca eklenebilir
//
// # Dezavantajlar
//
// - **Foreign Key Constraint Yok**: taggable_id için doğrudan foreign key tanımlanamaz
// - **Tip Güvenliği**: taggable_type string olduğu için yazım hataları olabilir
// - **Performans**: Sorgular daha karmaşık olabilir
//
// # Önemli Notlar
//
// - taggable_type değeri tutarlı olmalıdır (örn: "post", "comment")
// - Silme işlemlerinde cascade dikkatli yapılmalıdır
// - Sorgularda hem taggable_id hem de taggable_type kullanılmalıdır
//
// # Örnek Kullanım
//
// ```go
// // Bir gönderiye etiket ekle
// _, err := db.Exec(`
//     INSERT INTO taggables (taggable_id, taggable_type, tag_id)
//     VALUES (?, ?, ?)
// `, 1, "post", 2)
//
// // Bir gönderinin tüm etiketlerini getir (polimorfik sorgu)
// rows, err := db.Query(`
//     SELECT t.id, t.name
//     FROM tags t
//     JOIN taggables tg ON t.id = tg.tag_id
//     WHERE tg.taggable_id = ? AND tg.taggable_type = ?
// `, 1, "post")
//
// // Bir etiketin tüm gönderilerini getir
// rows, err := db.Query(`
//     SELECT p.id, p.title
//     FROM posts p
//     JOIN taggables tg ON p.id = tg.taggable_id
//     WHERE tg.tag_id = ? AND tg.taggable_type = ?
// `, 1, "post")
// ```
type Taggable struct {
	// ID, taggable kaydının benzersiz kimliğidir (primary key).
	ID           int    `json:"id"`

	// TaggableID, etiketlenen kaydın ID'sidir.
	// Bu, Post, Comment veya başka bir modelin ID'si olabilir.
	TaggableID   int    `json:"taggable_id"`

	// TaggableType, etiketlenen kaydın model türünü belirtir.
	// Örnek değerler: "post", "comment", "user", vb.
	TaggableType string `json:"taggable_type"`

	// TagID, etiketin ID'sidir (foreign key).
	TagID        int    `json:"tag_id"`
}

// GetTestDataDir, test verileri için geçici dizin yolunu döndürür.
//
// Bu fonksiyon, test sırasında geçici dosyaların saklanması gereken
// dizin yolunu sağlar. İşletim sisteminin standart geçici dizinini kullanır.
//
// # Dönüş Değeri
//
// İşletim sisteminin geçici dizin yolunu string olarak döndürür.
// - Unix/Linux/macOS: Genellikle `/tmp`
// - Windows: Genellikle `C:\Users\<username>\AppData\Local\Temp`
//
// # Kullanım Senaryoları
//
// - Test sırasında geçici dosya oluşturma
// - Test verilerini diske yazma
// - Dosya yükleme testleri
// - Önbellek testleri
//
// # Önemli Notlar
//
// - **Geçici Dizin**: Döndürülen dizin işletim sistemi tarafından yönetilir
// - **Temizlik**: Test sonunda oluşturulan dosyaları manuel olarak temizlemek gerekebilir
// - **İzinler**: Geçici dizine yazma izni genellikle mevcuttur
// - **Çakışma**: Benzersiz dosya adları kullanarak çakışmaları önleyin
//
// # Örnek Kullanım
//
// ```go
// func TestFileUpload(t *testing.T) {
//     testDir := GetTestDataDir()
//     testFile := filepath.Join(testDir, "test_upload.txt")
//
//     // Test dosyası oluştur
//     err := os.WriteFile(testFile, []byte("test data"), 0644)
//     if err != nil {
//         t.Fatal(err)
//     }
//     defer os.Remove(testFile) // Temizlik
//
//     // Test kodunuz...
// }
//
// // Benzersiz dosya adı ile kullanım
// func TestWithUniqueFile(t *testing.T) {
//     testDir := GetTestDataDir()
//     testFile := filepath.Join(testDir, fmt.Sprintf("test_%d.txt", time.Now().UnixNano()))
//
//     // Test kodunuz...
// }
// ```
func GetTestDataDir() string {
	return os.TempDir()
}
