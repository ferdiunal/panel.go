# Başlarken (Getting Started)

Panel SDK ile dakikalar içinde modern bir admin paneli oluşturabilirsiniz.

## Kurulum

```bash
go get github.com/ferdiunal/panel.go
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
import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

type UserResource struct{
    resource.Base
}

func GetUserResource() resource.Resource {
    return &UserResource{
        Base: resource.Base{
            DataModel: &User{},
            Label:     "Kullanıcılar",
            FieldsVal: []fields.Element{
                fields.ID(),
                fields.Text("İsim", "full_name"),
                fields.Email("E-Posta", "email"),
            },
        },
    }
}
```

### 3. Uygulamayı Başlatma

```go
package main

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/ferdiunal/panel.go/pkg/panel"
)

func main() {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    db.AutoMigrate(&User{})

    cfg := panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
        Server:   panel.ServerConfig{Port: "8080", Host: "localhost"},
        Environment: "production",
    }

    app := panel.New(cfg)
    
    app.RegisterResource(GetUserResource()) 
    
    app.Start()
}
```
