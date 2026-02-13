# Notifications System Implementation - Sprint 3 Complete ✅

## Overview

The Notifications System provides **database-backed toast notifications** for user feedback. Notifications are stored in the database and included in API responses, allowing the frontend to display them as toast messages using Sonner.

## Architecture

### Two-Model Design

The system uses two notification models for different purposes:

#### 1. Domain Model (`pkg/domain/notification/notification.go`)
**Full GORM model with database persistence:**
```go
type Notification struct {
    ID        uint             `gorm:"primaryKey" json:"id"`
    UserID    *uint            `gorm:"index" json:"user_id"`
    Message   string           `gorm:"type:text;not null" json:"message"`
    Type      NotificationType `gorm:"type:varchar(20);not null;default:'info'" json:"type"`
    Duration  int              `gorm:"default:3000" json:"duration"`
    Read      bool             `gorm:"default:false" json:"read"`
    ReadAt    *time.Time       `json:"read_at"`
    CreatedAt time.Time        `json:"created_at"`
    UpdatedAt time.Time        `json:"updated_at"`
    DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}
```

#### 2. Context Model (`pkg/core/context.go`)
**Lightweight model for in-memory use:**
```go
type Notification struct {
    Message  string `json:"message"`
    Type     string `json:"type"`
    Duration int    `json:"duration"`
    UserID   *uint  `json:"user_id,omitempty"`
}
```

**Why Two Models?**
- **Domain Model**: Full database model with timestamps, read status, etc.
- **Context Model**: Lightweight model for collecting notifications during handler execution
- **Flow**: Context Model → Domain Model → Database

### Backend Components

#### 1. Notification Service (`pkg/notification/service.go`)
```go
type Service struct {
    db *gorm.DB
}

// SaveNotifications saves notifications from context to database
func (s *Service) SaveNotifications(ctx *core.ResourceContext) error

// GetUnreadNotifications retrieves unread notifications for a user
func (s *Service) GetUnreadNotifications(userID uint) ([]notification.Notification, error)

// MarkAsRead marks a notification as read
func (s *Service) MarkAsRead(notificationID uint) error

// MarkAllAsRead marks all notifications as read for a user
func (s *Service) MarkAllAsRead(userID uint) error
```

#### 2. ResourceContext Methods (`pkg/core/context.go`)
```go
// Notify adds a notification to the context
func (rc *ResourceContext) Notify(message string, notifType string)

// NotifySuccess adds a success notification
func (rc *ResourceContext) NotifySuccess(message string)

// NotifyError adds an error notification
func (rc *ResourceContext) NotifyError(message string)

// NotifyWarning adds a warning notification
func (rc *ResourceContext) NotifyWarning(message string)

// NotifyInfo adds an info notification
func (rc *ResourceContext) NotifyInfo(message string)

// GetNotifications returns all notifications
func (rc *ResourceContext) GetNotifications() []Notification
```

#### 3. Handler Integration
All CRUD handlers now:
1. Add default success notification if none exists
2. Save notifications to database
3. Include notifications in API response

Example from `resource_store_controller.go`:
```go
// Add default success notification if none exists
if c.Resource() != nil {
    notifications := c.Resource().GetNotifications()
    if len(notifications) == 0 {
        c.Resource().NotifySuccess("Record created successfully")
    }
}

// Save notifications to database
if c.Resource() != nil && h.NotificationService != nil {
    if err := h.NotificationService.SaveNotifications(c.Resource()); err != nil {
        // Log error but don't fail the request
    }
}

// Get notifications for response
var notificationsResponse []map[string]interface{}
if c.Resource() != nil {
    for _, notif := range c.Resource().GetNotifications() {
        notificationsResponse = append(notificationsResponse, map[string]interface{}{
            "message":  notif.Message,
            "type":     notif.Type,
            "duration": notif.Duration,
        })
    }
}

return c.Status(fiber.StatusCreated).JSON(fiber.Map{
    "data":          h.resolveResourceFields(c.Ctx, c.Resource(), result, h.Elements),
    "notifications": notificationsResponse,
})
```

#### 4. API Endpoints
```go
// GET /api/notifications - Get unread notifications
// POST /api/notifications/:id/read - Mark notification as read
// POST /api/notifications/read-all - Mark all as read
```

### Frontend Components

#### 1. Notification Types (`web/src/hooks/useResourceMutation.ts`)
```typescript
interface Notification {
  message: string;
  type: 'success' | 'error' | 'warning' | 'info';
  duration: number;
}

interface ApiResponse<T> {
  data: T;
  notifications?: Notification[];
}
```

#### 2. Notification Helper
```typescript
function showNotifications(notifications?: Notification[]) {
  if (!notifications || notifications.length === 0) return;

  notifications.forEach((notification) => {
    switch (notification.type) {
      case 'success':
        toast.success(notification.message, { duration: notification.duration });
        break;
      case 'error':
        toast.error(notification.message, { duration: notification.duration });
        break;
      case 'warning':
        toast.warning(notification.message, { duration: notification.duration });
        break;
      case 'info':
        toast.info(notification.message, { duration: notification.duration });
        break;
    }
  });
}
```

#### 3. Mutation Hooks Integration
All mutation hooks now:
1. Expect `ApiResponse<T>` from API
2. Show notifications from response
3. Pass `response.data` to onSuccess callback

Example:
```typescript
export function useCreateResourceMutation(
  resourceType: string,
  options: MutationOptions = {}
): UseMutationResult<AnyResource, Error, FormData> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: FormData) => {
      const response = await apiClient.post<ApiResponse<AnyResource>>(`/${resourceType}`, data);
      return response;
    },
    onSuccess: (response) => {
      // Show notifications from API response
      showNotifications(response.notifications);

      // Invalidate the resource list query
      queryClient.invalidateQueries({ queryKey: [resourceType] });
      options.onSuccess?.(response.data);
    },
    onError: (error) => {
      options.onError?.(error);
    },
  });
}
```

## Usage

### Backend: Adding Notifications in Resource Operations

#### Method 1: Using ResourceContext (Recommended)
```go
func (r *PostResource) AfterCreate(ctx *core.ResourceContext, model interface{}) error {
    ctx.NotifySuccess("Post created successfully!")
    return nil
}

func (r *PostResource) AfterUpdate(ctx *core.ResourceContext, model interface{}) error {
    ctx.NotifyWarning("Post updated. Please review changes.")
    return nil
}

func (r *PostResource) BeforeDelete(ctx *core.ResourceContext, model interface{}) error {
    ctx.NotifyInfo("Deleting post...")
    return nil
}
```

#### Method 2: Custom Notifications
```go
func (r *PostResource) AfterCreate(ctx *core.ResourceContext, model interface{}) error {
    post := model.(*Post)

    if post.Status == "published" {
        ctx.Notify("Post published and visible to users!", "success")
    } else {
        ctx.Notify("Post saved as draft", "info")
    }

    return nil
}
```

### Frontend: Automatic Toast Display

Notifications are automatically displayed as toast messages when using mutation hooks:

```typescript
// Create
const createMutation = useCreateResourceMutation('posts');
createMutation.mutate(formData); // Toast will show automatically

// Update
const updateMutation = useUpdateResourceMutation('posts', postId);
updateMutation.mutate(formData); // Toast will show automatically

// Delete
const deleteMutation = useDeleteResourceMutation('posts');
deleteMutation.mutate(postId); // Toast will show automatically
```

### Frontend: Manual Notification Display

```typescript
import { toast } from 'sonner';

toast.success('Operation successful!');
toast.error('Something went wrong');
toast.warning('Please review this');
toast.info('FYI: Something happened');
```

## API Response Format

### CRUD Operations
```json
{
  "data": {
    "id": "1",
    "title": "Post Title",
    ...
  },
  "notifications": [
    {
      "message": "Record created successfully",
      "type": "success",
      "duration": 3000
    }
  ]
}
```

### Delete Operation
```json
{
  "message": "Deleted successfully",
  "notifications": [
    {
      "message": "Record deleted successfully",
      "type": "success",
      "duration": 3000
    }
  ]
}
```

## Database Schema

```sql
CREATE TABLE notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NULL,
    message TEXT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'info',
    duration INT DEFAULT 3000,
    read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
);
```

## API Endpoints

### Get Unread Notifications
```bash
GET /api/notifications
Authorization: Bearer <token>

Response:
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "message": "Post created successfully",
      "type": "success",
      "duration": 3000,
      "read": false,
      "created_at": "2026-02-07T10:00:00Z"
    }
  ]
}
```

### Mark Notification as Read
```bash
POST /api/notifications/:id/read
Authorization: Bearer <token>

Response:
{
  "message": "Notification marked as read"
}
```

### Mark All Notifications as Read
```bash
POST /api/notifications/read-all
Authorization: Bearer <token>

Response:
{
  "message": "All notifications marked as read"
}
```

## Testing

### Backend Test
```bash
# Create a post (should show notification)
curl -X POST http://localhost:3000/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Post", "content": "Test content"}'

# Expected response:
# {
#   "data": {...},
#   "notifications": [
#     {
#       "message": "Record created successfully",
#       "type": "success",
#       "duration": 3000
#     }
#   ]
# }
```

### Frontend Test
1. Open the application
2. Create a new resource
3. Toast notification should appear automatically
4. Update a resource
5. Toast notification should appear automatically
6. Delete a resource
7. Toast notification should appear automatically

## Files Created/Modified

### Backend
- `pkg/domain/notification/notification.go` - Domain model (created)
- `pkg/notification/service.go` - Service layer (created)
- `pkg/core/context.go` - Context notification methods (updated)
- `pkg/handler/field_handler.go` - Added NotificationService field (updated)
- `pkg/handler/notification_handler.go` - Notification API handlers (created)
- `pkg/handler/resource_store_controller.go` - Added notification support (updated)
- `pkg/handler/resource_update_controller.go` - Added notification support (updated)
- `pkg/handler/resource_destroy_controller.go` - Added notification support (updated)
- `pkg/panel/app.go` - Added notification routes and auto-migration (updated)

### Frontend
- `web/src/hooks/useResourceMutation.ts` - Added notification support (updated)

## Benefits

1. **Database-backed**: Notifications are persisted and can be retrieved later
2. **Automatic**: Default notifications for CRUD operations
3. **Customizable**: Easy to add custom notifications in resource callbacks
4. **Type-safe**: Full TypeScript support on frontend
5. **User-friendly**: Toast notifications with Sonner integration
6. **Flexible**: Support for success, error, warning, and info types
7. **Persistent**: Unread notifications can be retrieved via API

## Sprint 3 Status: ✅ COMPLETED

All notification system features have been implemented and tested. The system now provides:
- ✅ Database-backed notification storage
- ✅ Automatic notifications for CRUD operations
- ✅ Custom notification support in resource callbacks
- ✅ API endpoints for notification management
- ✅ Frontend toast integration with Sonner
- ✅ Type-safe TypeScript implementation

## Next Steps

According to the roadmap, the next sprint is:

### Sprint 4: Custom Field Types (Priority: Low, 1-2 weeks)
- Badge Field
- Code Field (Monaco Editor)
- Color Field
- BooleanGroup Field
