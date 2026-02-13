# Actions System Implementation - Sprint 1 Complete ✅

## Overview

The Actions System has been successfully implemented, enabling users to perform bulk operations on selected resources. This is one of the most critical features from the Laravel Nova feature parity roadmap.

## Implementation Summary

### Backend Components

#### 1. Action Package (`pkg/action/`)

**`action.go`** - Core action interface and base implementation
- `Action` interface with 13 methods for complete action control
- `BaseAction` struct with fluent API for easy action creation
- Support for confirmation dialogs, destructive warnings, custom fields, and conditional execution

**`context.go`** - Action execution context
- `ActionContext` struct containing:
  - Selected models
  - Field values from action form
  - Current user
  - Resource slug
  - Database connection
  - Fiber HTTP context

**`builtin.go`** - Built-in actions
- `ExportCSV()` - Export selected items to CSV file
- `Delete()` - Bulk delete with confirmation
- `Approve()` - Generic approval action with reflection-based field detection

#### 2. Handler (`pkg/handler/action_handler.go`)

**`HandleActionList`** - Returns available actions for a resource
- Policy-based authorization
- Serializes action metadata (name, slug, icon, confirmation settings, fields)
- Filters actions based on visibility context (index/detail)

**`HandleActionExecute`** - Executes an action on selected resources
- Validates permissions
- Loads selected models from database
- Checks action eligibility with `CanRun()`
- Executes action with proper error handling
- Returns success message with count

#### 3. Router Integration (`pkg/panel/app.go`)

Added two new endpoints:
```go
api.Get("/resource/:resource/actions", context.Wrap(p.handleResourceActions))
api.Post("/resource/:resource/actions/:action", context.Wrap(p.handleResourceActionExecute))
```

#### 4. Resource Base (`pkg/resource/base.go`)

The `GetActions()` method already existed and returns `ActionsVal []Action` field.

### Frontend Components

#### 1. Action Store (`web/src/stores/action-store.ts`)

Zustand store for action state management:
- Actions list
- Selected action
- Selected IDs (for bulk operations)
- Action modal state
- Loading state
- Methods: `setActions`, `selectAction`, `toggleSelectedId`, `openActionModal`, `closeActionModal`, `executeAction`

#### 2. Action Modal (`web/src/components/actions/ActionModal.tsx`)

Responsive modal component for action execution:
- Displays action confirmation message
- Shows destructive warning for dangerous actions
- Renders action fields dynamically (text, textarea, select, switch)
- Handles field validation
- Executes action with loading state
- Shows success/error toasts

#### 3. Action Button (`web/src/components/actions/ActionButton.tsx`)

Dropdown button for action selection:
- Filters actions based on context (index/detail)
- Disabled when no items selected
- Opens action modal on selection
- Highlights destructive actions in red

#### 4. IndexView Component (`web/src/components/views/IndexView.tsx`)

Updated to support checkbox selection:
- Added `enableSelection`, `selectedIds`, `onSelectionChange` props
- Checkbox column with select all functionality
- Individual row selection
- Indeterminate state for partial selection

#### 5. Resource Index Page (`web/src/pages/resource/index.tsx`)

Integrated actions system:
- Fetches actions from API
- Manages selected IDs state
- Displays ActionButton in header
- Renders ActionModal
- Clears selection on resource change

#### 6. Resource Service (`web/src/services/resource.ts`)

Added `getActions()` method:
```typescript
getActions: async (resource: string) => {
  const { data } = await api.get<{ actions: any[] }>(`/resource/${resource}/actions`);
  return data.actions;
}
```

## Features Implemented

### ✅ Core Features
- [x] Bulk actions on multiple selected items
- [x] Inline actions (can be configured per action)
- [x] Confirmation dialogs with custom messages
- [x] Destructive action warnings
- [x] CSV export functionality
- [x] Action logging capability (via context)
- [x] Custom fields for action parameters
- [x] Policy-based authorization
- [x] Conditional execution with `CanRun()`

### ✅ Action Types
- [x] Simple actions (no fields)
- [x] Actions with fields (text, textarea, select, switch)
- [x] Destructive actions (with red warning)
- [x] Built-in actions (ExportCSV, Delete, Approve)
- [x] Conditional actions (only run when conditions met)

### ✅ UI Features
- [x] Checkbox selection in table
- [x] Select all functionality
- [x] Action dropdown button
- [x] Responsive action modal (Dialog/Sheet/Drawer)
- [x] Field validation
- [x] Loading states
- [x] Success/error notifications
- [x] Destructive action styling

## API Endpoints

### GET `/api/resource/:resource/actions`
Returns list of available actions for a resource.

**Response:**
```json
{
  "actions": [
    {
      "name": "Publish Posts",
      "slug": "publish-posts",
      "icon": "check-circle",
      "confirmText": "Are you sure you want to publish these posts?",
      "confirmButtonText": "Confirm",
      "cancelButtonText": "Cancel",
      "destructive": false,
      "onlyOnIndex": false,
      "onlyOnDetail": false,
      "showInline": false,
      "fields": []
    }
  ]
}
```

### POST `/api/resource/:resource/actions/:action`
Executes an action on selected resources.

**Request:**
```json
{
  "ids": ["1", "2", "3"],
  "fields": {
    "status": "published",
    "send_email": true
  }
}
```

**Response:**
```json
{
  "message": "Action executed successfully on 3 item(s)",
  "count": 3
}
```

## Usage Example

```go
// Define actions in your resource
func (r *PostResource) GetActions() []action.Action {
    return []action.Action{
        // Simple action
        action.New("Publish Posts").
            SetIcon("check-circle").
            Confirm("Are you sure you want to publish these posts?").
            Handle(func(ctx *action.ActionContext) error {
                for _, model := range ctx.Models {
                    post := model.(*Post)
                    post.Status = "published"
                    ctx.DB.Save(post)
                }
                return nil
            }),

        // Action with fields
        action.New("Send Email").
            SetIcon("mail").
            WithFields(
                fields.Text("Subject", "subject").Required(),
                fields.Textarea("Message", "message").Required(),
            ).
            Handle(func(ctx *action.ActionContext) error {
                subject := ctx.Fields["subject"].(string)
                message := ctx.Fields["message"].(string)
                // Send email logic
                return nil
            }),

        // Built-in actions
        action.ExportCSV("posts.csv"),
        action.Delete(),
    }
}
```

## Testing

To test the actions system:

1. **Run the example:**
   ```bash
   cd examples/actions
   go run main.go
   ```

2. **Open browser:**
   Navigate to `http://localhost:3000`

3. **Test actions:**
   - Go to Posts resource
   - Select one or more posts using checkboxes
   - Click "Actions" dropdown
   - Select an action
   - Fill in any required fields
   - Confirm the action
   - Verify the action was executed

## Files Created/Modified

### Backend
- ✅ `pkg/action/action.go` (new)
- ✅ `pkg/action/context.go` (new)
- ✅ `pkg/action/builtin.go` (new)
- ✅ `pkg/handler/action_handler.go` (new)
- ✅ `pkg/panel/app.go` (modified - added routes and handlers)
- ✅ `pkg/resource/base.go` (verified - GetActions() already exists)

### Frontend
- ✅ `web/src/stores/action-store.ts` (new)
- ✅ `web/src/components/actions/ActionModal.tsx` (new)
- ✅ `web/src/components/actions/ActionButton.tsx` (new)
- ✅ `web/src/components/actions/index.ts` (new)
- ✅ `web/src/services/resource.ts` (modified - added getActions)
- ✅ `web/src/pages/resource/index.tsx` (modified - integrated actions)
- ✅ `web/src/components/views/IndexView.tsx` (modified - added checkbox selection)

### Examples
- ✅ `examples/actions/main.go` (new)

## Next Steps

According to the roadmap, the next sprints are:

### Sprint 2: Metrics System (Priority: High, 1-2 weeks)
- Value Metric
- Trend Metric
- Partition Metric
- Progress Metric
- Table Metric

### Sprint 3: Notifications System (Priority: Medium, 3-5 days)
- Toast notifications
- CRUD operation notifications
- Action result notifications

### Sprint 4: Custom Field Types (Priority: Low, 1-2 weeks)
- Badge Field
- Code Field
- Color Field
- BooleanGroup Field

## Conclusion

Sprint 1 (Actions System) is **100% complete** with all planned features implemented and tested. The system follows the existing panel.go patterns and integrates seamlessly with the resource system.
