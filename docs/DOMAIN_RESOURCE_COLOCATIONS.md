# Domain-Resource Co-location Guide

## Overview

Domain-Resource Co-location is an architectural pattern that organizes domain models and their corresponding resource definitions in the same directory. This approach improves code organization, maintainability, and makes it easier to understand the relationship between domain models and their API representations.

## Motivation

In traditional architectures, domain models and their API resources are often scattered across different directories:

```
pkg/
├── domain/
│   ├── user/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
├── resource/
│   ├── user_resource.go
│   ├── post_resource.go
│   └── comment_resource.go
```

This separation makes it difficult to:
- Understand the complete picture of a domain entity
- Maintain consistency between the model and its resource
- Navigate between related files
- Refactor domain entities

## Solution: Co-location

The co-location pattern places domain models and their resources in the same directory:

```
pkg/domain/
├── user/
│   ├── model.go           (Domain model)
│   ├── repository.go      (Data access)
│   ├── service.go         (Business logic)
│   └── resource.go        (API resource)
├── post/
│   ├── model.go
│   ├── repository.go
│   ├── service.go
│   └── resource.go
├── comment/
│   ├── model.go
│   ├── repository.go
│   ├── service.go
│   └── resource.go
└── tag/
    ├── model.go
    ├── repository.go
    ├── service.go
    └── resource.go
```

## Benefits

### 1. **Improved Navigation**
All related code for a domain entity is in one place, making it easier to find and understand.

### 2. **Better Maintainability**
Changes to a domain model can be immediately reflected in its resource definition without searching across multiple directories.

### 3. **Clear Dependencies**
The relationship between models and resources is explicit and easy to understand.

### 4. **Easier Refactoring**
Renaming or restructuring a domain entity is simpler when all related code is co-located.

### 5. **Scalability**
As the codebase grows, the co-location pattern scales better than scattered resources.

## Directory Structure

Each domain directory should contain:

```
pkg/domain/{entity}/
├── model.go              # Domain model definition
├── repository.go         # Data access layer (optional)
├── service.go            # Business logic (optional)
├── resource.go           # API resource definition
├── {entity}_test.go      # Tests for the domain
└── README.md             # Documentation (optional)
```

### File Descriptions

#### model.go
Contains the domain model struct and any model-specific methods:

```go
package user

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) IsAdmin() bool {
    // Business logic
    return false
}
```

#### repository.go
Contains data access methods:

```go
package user

import "gorm.io/gorm"

type Repository interface {
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}

type GormRepository struct {
    db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
    return &GormRepository{db: db}
}

// Implement Repository interface...
```

#### service.go
Contains business logic:

```go
package user

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) RegisterUser(name, email, password string) (*User, error) {
    // Business logic for user registration
    user := &User{
        Name:  name,
        Email: email,
    }
    return user, s.repo.Create(user)
}
```

#### resource.go
Contains the API resource definition:

```go
package user

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

type UserResource struct {
    *resource.Base
}

func (r *UserResource) Title() string {
    return "Users"
}

func (r *UserResource) Slug() string {
    return "users"
}

func (r *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Name", "name"),
        fields.Email("Email", "email"),
    }
}

func (r *UserResource) Model() interface{} {
    return &User{}
}
```

## Migration Guide

### Step 1: Create Domain Directory Structure

```bash
mkdir -p pkg/domain/user
mkdir -p pkg/domain/post
mkdir -p pkg/domain/comment
```

### Step 2: Move Model Files

Move existing model definitions to `model.go`:

```bash
# Before
pkg/domain/user.go

# After
pkg/domain/user/model.go
```

### Step 3: Move Resource Files

Move resource definitions to `resource.go`:

```bash
# Before
pkg/resource/user_resource.go

# After
pkg/domain/user/resource.go
```

### Step 4: Update Imports

Update all imports to reflect the new structure:

```go
// Before
import "github.com/ferdiunal/panel.go/pkg/domain"
import "github.com/ferdiunal/panel.go/pkg/resource"

// After
import "github.com/ferdiunal/panel.go/pkg/domain/user"
```

### Step 5: Update Package Initialization

Update the main application to register resources from the new locations:

```go
// Before
p.RegisterResource(resource.GetUserResource())

// After
p.RegisterResource(&user.UserResource{})
```

## Best Practices

### 1. **Keep Concerns Separated**
Even though files are co-located, maintain clear separation of concerns:
- Model: Data structure and model-specific methods
- Repository: Data access
- Service: Business logic
- Resource: API representation

### 2. **Use Interfaces**
Define interfaces for repositories and services to enable testing and flexibility:

```go
type Repository interface {
    FindByID(id uint) (*User, error)
    Create(user *User) error
}
```

### 3. **Document Relationships**
Add comments explaining the relationship between model and resource:

```go
// UserResource represents the API resource for User domain entity.
// It defines fields, policies, and actions available for users in the API.
type UserResource struct {
    *resource.Base
}
```

### 4. **Test Co-location**
Keep tests in the same directory as the code they test:

```
pkg/domain/user/
├── model.go
├── model_test.go
├── repository.go
├── repository_test.go
├── service.go
├── service_test.go
├── resource.go
└── resource_test.go
```

### 5. **Use Package-level Documentation**
Add a `README.md` or package comment explaining the domain:

```go
// Package user provides domain models, repositories, services, and API resources
// for user management. It implements the complete user lifecycle including
// registration, authentication, and profile management.
package user
```

## Example: Complete User Domain

```
pkg/domain/user/
├── model.go
│   └── User struct with methods
├── repository.go
│   └── UserRepository interface and implementation
├── service.go
│   └── UserService with business logic
├── resource.go
│   └── UserResource for API
├── user_test.go
│   └── Tests for the domain
└── README.md
    └── Documentation
```

## Refactoring Existing Code

### Phase 1: Create New Structure
1. Create new domain directories
2. Copy model and resource files
3. Update imports in new files
4. Keep old files temporarily

### Phase 2: Update References
1. Update all imports to use new locations
2. Update resource registration in main app
3. Run tests to ensure everything works

### Phase 3: Cleanup
1. Remove old files
2. Remove old directories
3. Verify all tests pass

## Advantages Over Alternatives

### vs. Monolithic Domain Directory
```
# Monolithic (harder to navigate)
pkg/domain/
├── user.go
├── post.go
├── comment.go
├── user_resource.go
├── post_resource.go
└── comment_resource.go

# Co-located (easier to navigate)
pkg/domain/
├── user/
│   ├── model.go
│   └── resource.go
├── post/
│   ├── model.go
│   └── resource.go
└── comment/
    ├── model.go
    └── resource.go
```

### vs. Scattered Resources
```
# Scattered (hard to maintain)
pkg/domain/user.go
pkg/resource/user_resource.go
pkg/service/user_service.go
pkg/repository/user_repository.go

# Co-located (easy to maintain)
pkg/domain/user/
├── model.go
├── resource.go
├── service.go
└── repository.go
```

## Conclusion

Domain-Resource Co-location is a powerful architectural pattern that improves code organization and maintainability. By placing related code in the same directory, developers can more easily understand, maintain, and refactor domain entities.

This pattern is particularly effective in Go applications where package organization is a key part of the architecture. It encourages developers to think about domains as cohesive units rather than scattered pieces across the codebase.
