# Manuel Plugin Oluşturma Örneği

## 1. plugin.yaml

```yaml
name: course-management
version: 1.0.0
author: Your Name
description: Course management plugin for LMS
```

## 2. plugin.go

```go
package course_management

import (
    "github.com/ferdiunal/panel.go/pkg/plugin"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/gofiber/fiber/v2"
)

// CourseManagementPlugin - Kurs yönetimi plugin'i
type CourseManagementPlugin struct {
    plugin.BasePlugin
}

// Metadata
func (p *CourseManagementPlugin) Name() string {
    return "course-management"
}

func (p *CourseManagementPlugin) Version() string {
    return "1.0.0"
}

func (p *CourseManagementPlugin) Author() string {
    return "Your Name"
}

func (p *CourseManagementPlugin) Description() string {
    return "Course management plugin for LMS"
}

// Lifecycle
func (p *CourseManagementPlugin) Register(panel interface{}) error {
    // Plugin kaydı sırasında çağrılır
    return nil
}

func (p *CourseManagementPlugin) Boot(panel interface{}) error {
    // Panel başlatıldığında çağrılır
    return nil
}

// Resources - Custom resource'ları döndür
func (p *CourseManagementPlugin) Resources() []resource.Resource {
    return []resource.Resource{
        // Örnek: Course resource
        resource.New(
            &Course{},
            "course",
            "Courses",
        ).Icon("book-open").
            Group("Education").
            Fields(func(r *resource.Resource) []fields.Field {
                return []fields.Field{
                    fields.ID(),
                    fields.Text("Title", "title").Required(),
                    fields.Textarea("Description", "description"),
                    fields.Number("Duration", "duration").
                        Placeholder("Duration in hours"),
                    fields.Select("Level", "level").
                        Options(map[string]string{
                            "beginner":     "Beginner",
                            "intermediate": "Intermediate",
                            "advanced":     "Advanced",
                        }),
                    fields.DateTime("Start Date", "start_date"),
                    fields.Boolean("Published", "published"),
                }
            }),
    }
}

// Middleware - Custom middleware'leri döndür
func (p *CourseManagementPlugin) Middleware() []fiber.Handler {
    return []fiber.Handler{
        // Custom middleware ekleyebilirsin
    }
}

// Routes - Custom route'ları kaydet
func (p *CourseManagementPlugin) Routes(router fiber.Router) {
    // Custom API endpoint'leri ekleyebilirsin
    router.Get("/api/courses/stats", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "total_courses": 100,
            "active_courses": 75,
        })
    })
}

// Migrations - Database migration'ları döndür
func (p *CourseManagementPlugin) Migrations() []plugin.Migration {
    return []plugin.Migration{
        &CreateCoursesTable{},
    }
}

// Course model
type Course struct {
    ID          uint   `gorm:"primaryKey"`
    Title       string `gorm:"size:255;not null"`
    Description string `gorm:"type:text"`
    Duration    int    `gorm:"default:0"`
    Level       string `gorm:"size:50"`
    StartDate   string `gorm:"type:datetime"`
    Published   bool   `gorm:"default:false"`
}

// Migration
type CreateCoursesTable struct{}

func (m *CreateCoursesTable) Name() string {
    return "create_courses_table"
}

func (m *CreateCoursesTable) Up(db interface{}) error {
    gormDB := db.(*gorm.DB)
    return gormDB.AutoMigrate(&Course{})
}

func (m *CreateCoursesTable) Down(db interface{}) error {
    gormDB := db.(*gorm.DB)
    return gormDB.Migrator().DropTable(&Course{})
}

// Plugin'i global registry'ye kaydet
func init() {
    plugin.Register(&CourseManagementPlugin{})
}
```

## 3. main.go'da Import Et

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    _ "learning_management/plugins/course-management" // Plugin'i import et
)

func main() {
    config := panel.Config{
        Database: panel.DatabaseConfig{
            Driver: "sqlite",
            DSN:    "lms.db",
        },
        Server: panel.ServerConfig{
            Port: 8787,
        },
    }

    p := panel.New(config)
    p.Start()
}
```

## Frontend Plugin (Opsiyonel)

Eğer custom field component'i gerekiyorsa:

### 1. Frontend Dizini Oluştur

```bash
mkdir -p plugins/course-management/frontend/fields
```

### 2. Plugin Export

```typescript
// plugins/course-management/frontend/index.ts
import type { Plugin } from '@/plugins/types';
import { CourseField } from './fields/CourseField';

export const CourseManagementPlugin: Plugin = {
  name: 'course-management',
  version: '1.0.0',
  description: 'Course management plugin',
  author: 'Your Name',

  fields: [
    {
      type: 'course-field',
      component: CourseField,
    },
  ],

  init: async () => {
    console.log('CourseManagementPlugin initialized');
  },
};
```

### 3. Custom Field Component

```typescript
// plugins/course-management/frontend/fields/CourseField.tsx
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface CourseFieldProps {
  field: {
    key: string;
    label?: string;
    placeholder?: string;
  };
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function CourseField({
  field,
  value,
  onChange,
  error,
}: CourseFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={field.key}>{field.label || field.key}</Label>
      <Input
        id={field.key}
        value={value || ''}
        onChange={(e) => onChange(e.target.value)}
        placeholder={field.placeholder}
        className={error ? 'border-destructive' : ''}
      />
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  );
}
```

### 4. Frontend Plugin'i Kaydet

Panel.go projesinde frontend plugin'leri `web/src/plugins/` altında kayıtlı. Kendi projenizde frontend plugin kullanmak için:

1. Panel.go'nun web-ui'sini clone edin
2. Plugin frontend dosyalarını `web/src/plugins/course-management/frontend/` altına kopyalayın
3. `web/src/plugins/index.ts` dosyasına plugin'i ekleyin
4. Build alın: `panel plugin build`

## Özet

**Backend-only plugin için:**
1. `plugins/course-management/` dizini oluştur
2. `plugin.yaml` ve `plugin.go` dosyalarını oluştur
3. `main.go`'da import et
4. Çalıştır: `go run main.go`

**Frontend plugin ile:**
1. Backend adımlarını tamamla
2. Frontend dosyalarını oluştur
3. Panel CLI ile build al: `panel plugin build`
4. Çalıştır: `go run main.go`
