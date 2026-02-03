# Başlarken (Getting Started)

Panel SDK ile dakikalar içinde modern bir admin paneli oluşturabilirsiniz.

## Kurulum

```bash
go get panel.go
```

## Hızlı Başlangıç

### 1. Veritabanı Modeli

GORM modelinizi tanımlayın:

```go
type User struct {
    ID        uint   `json:"id" gorm:"primaryKey"`
    FullName  string `json:"full_name"`
    Email     string `json:"email"`
}
```

### 2. Resource Tanımı

Model ile UI arasındaki köprüyü kurun:

```go
type UserResource struct{}

func (u *UserResource) Model() interface{} {
    return &User{}
}

func (u *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("İsim", "full_name"),
        fields.Email("E-Posta", "email"),
    }
}
```

### 3. Uygulamayı Başlatma

```go
func main() {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    db.AutoMigrate(&User{})

    app := panel.New(panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
    })
    
    app.Register("users", &UserResource{}) 
    
    app.Start()
}
```
