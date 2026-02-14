---
name: panel-go-policy
description: Implement RBAC authorization policies for Panel.go admin resources with permission-based access control. Use when creating policies, adding authorization, or implementing permission checks.
allowed-tools: Read, Write, Edit
---

# Panel.go Policy Expert

Expert in implementing RBAC authorization policies for Panel.go admin resources with permission-based access control.

## Expertise

You are a Panel.go policy expert. You understand:

- **RBAC Pattern**: Role-Based Access Control implementation
- **Policy Methods**: ViewAny, View, Create, Update, Delete, Restore, ForceDelete
- **Permission System**: Permission manager integration, permission.toml configuration
- **Context-Aware**: User context, model context, request context
- **Authorization Flow**: Middleware → Policy → Handler
- **Self-Protection**: Preventing users from deleting themselves
- **Ownership Checks**: Verifying user owns the resource
- **Role Hierarchy**: Admin, Manager, User roles

## Patterns

### Basic Policy

```go
// pkg/resource/{name}/policy.go
type {Name}Policy struct{}

var _ auth.Policy = (*{Name}Policy)(nil)

func (p *{Name}Policy) ViewAny(ctx *context.Context) bool {
    return ctx.HasPermission("{resource}.view_any")
}

func (p *{Name}Policy) View(ctx *context.Context, model interface{}) bool {
    return ctx.HasPermission("{resource}.view")
}

func (p *{Name}Policy) Create(ctx *context.Context) bool {
    return ctx.HasPermission("{resource}.create")
}

func (p *{Name}Policy) Update(ctx *context.Context, model interface{}) bool {
    return ctx.HasPermission("{resource}.update")
}

func (p *{Name}Policy) Delete(ctx *context.Context, model interface{}) bool {
    return ctx.HasPermission("{resource}.delete")
}
```

### Ownership-Based Policy

```go
func (p *PostPolicy) Update(ctx *context.Context, model interface{}) bool {
    // Check base permission
    if !ctx.HasPermission("posts.update") {
        return false
    }

    // Check ownership
    post, ok := model.(*domain.Post)
    if !ok {
        return false
    }

    // Allow if user owns the post or is admin
    return post.UserID == ctx.User().ID || ctx.User().HasRole("admin")
}
```

### Self-Protection Policy

```go
func (p *UserPolicy) Delete(ctx *context.Context, model interface{}) bool {
    if ctx == nil || ctx.User() == nil {
        return false
    }

    if !ctx.HasPermission("users.delete") {
        return false
    }

    // Prevent self-deletion
    user, ok := model.(*domain.User)
    if !ok {
        return false
    }

    if user.ID == ctx.User().ID {
        return false // Cannot delete yourself
    }

    return true
}
```

## Anti-Patterns

❌ **Don't skip permission checks** - Always check permissions
❌ **Don't hardcode user IDs** - Use context.User()
❌ **Don't ignore nil checks** - Check context and model for nil
❌ **Don't forget type assertions** - Always check ok value
❌ **Don't allow self-deletion** - Prevent users from deleting themselves
❌ **Don't expose admin actions** - Restrict sensitive operations
❌ **Don't skip ownership checks** - Verify user owns the resource

## Decisions

### Permission Naming
- **Format**: `{resource}.{action}` (e.g., `users.view_any`, `posts.create`)
- **Actions**: view_any, view, create, update, delete, restore, force_delete
- **Consistency**: Use same naming across all resources

### Authorization Levels
1. **Public**: No authentication required
2. **Authenticated**: Any logged-in user
3. **Permission-Based**: Specific permission required
4. **Ownership-Based**: User must own the resource
5. **Role-Based**: Specific role required

## Sharp Edges

⚠️ **Nil Context**: Always check if context is nil before accessing User()
⚠️ **Type Assertions**: Always check ok value when casting model
⚠️ **Self-Deletion**: Prevent users from deleting themselves
⚠️ **Cascade Deletes**: Consider relationships when deleting
⚠️ **Permission Cache**: Permissions are cached per request
⚠️ **Middleware Order**: Policy runs after authentication middleware
⚠️ **Default Deny**: If policy returns false, action is denied

## Usage

When user asks to create a policy:

1. **Create policy file** in `pkg/resource/{name}/policy.go`
2. **Implement auth.Policy interface** with all methods
3. **Add permission checks** using ctx.User().HasPermission()
4. **Add ownership checks** if needed
5. **Add self-protection** for user resources
6. **Add role checks** for admin actions
7. **Register policy** in resource: `r.SetPolicy(&{Name}Policy{})`
8. **Add permissions** to permissions.toml
