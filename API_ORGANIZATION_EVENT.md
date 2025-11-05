# Organization and Event Management API Documentation

## Overview

Ticketly Backend now supports organization and event management with role-based access control. Users can create organizations, invite members with different roles (admin/member), and create events.

## Key Features

- **Organization Management**: Create and manage organizations
- **Role-Based Access Control**: Admin and Member roles
- **Event Management**: Create and manage events within organizations
- **Public Events**: Browse and search public events without authentication

## Database Schema

### Organization
- Organization entity with owner
- Members with roles (admin/member)
- Active/inactive status

### Event
- Events belong to organizations
- Status: draft, published, ongoing, completed, cancelled
- Ticket management system
- Public/private visibility

## API Endpoints

### Organization Management

#### Create Organization
```http
POST /api/organizations
Authorization: Bearer {token}

Request Body:
{
  "name": "My Organization",
  "description": "Description of organization",
  "logo_url": "https://example.com/logo.png"
}

Response: 201 Created
{
  "message": "Organization created successfully",
  "organization": {
    "id": "uuid",
    "name": "My Organization",
    "description": "Description of organization",
    "logo_url": "https://example.com/logo.png",
    "owner_id": "uuid",
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### Get My Organizations
```http
GET /api/organizations/my
Authorization: Bearer {token}

Response: 200 OK
{
  "organizations": [
    {
      "id": "uuid",
      "name": "My Organization",
      "description": "Description",
      "logo_url": "https://example.com/logo.png",
      "owner_id": "uuid",
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z",
      "user_role": "admin"
    }
  ]
}
```

#### Get Organization by ID
```http
GET /api/organizations/:id
Authorization: Bearer {token}

Response: 200 OK
{
  "organization": {
    "id": "uuid",
    "name": "My Organization",
    ...
  }
}
```

#### Update Organization (Admin Only)
```http
PUT /api/organizations/:id
Authorization: Bearer {token}

Request Body:
{
  "name": "Updated Name",
  "description": "Updated description",
  "logo_url": "https://example.com/new-logo.png"
}

Response: 200 OK
{
  "message": "Organization updated successfully"
}
```

#### Delete Organization (Owner Only)
```http
DELETE /api/organizations/:id
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Organization deleted successfully"
}
```

### Member Management

#### Get Organization Members
```http
GET /api/organizations/:id/members
Authorization: Bearer {token}

Response: 200 OK
{
  "members": [
    {
      "id": "uuid",
      "organization_id": "uuid",
      "user_id": "uuid",
      "role": "admin",
      "joined_at": "2025-01-01T00:00:00Z",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### Add Member (Admin Only)
```http
POST /api/organizations/:id/members
Authorization: Bearer {token}

Request Body:
{
  "user_id": "uuid",
  "role": "member"  // "admin" or "member"
}

Response: 201 Created
{
  "message": "Member added successfully"
}
```

#### Remove Member (Admin Only)
```http
DELETE /api/organizations/:id/members/:userId
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Member removed successfully"
}
```

#### Update Member Role (Admin Only)
```http
PUT /api/organizations/:id/members/:userId
Authorization: Bearer {token}

Request Body:
{
  "role": "admin"  // "admin" or "member"
}

Response: 200 OK
{
  "message": "Member role updated successfully"
}
```

### Event Management

#### Create Event (Admin Only)
```http
POST /api/organizations/:orgId/events
Authorization: Bearer {token}

Request Body:
{
  "title": "My Event",
  "description": "Event description",
  "location": "Seoul, Korea",
  "venue": "COEX Hall",
  "start_time": "2025-03-01T18:00:00Z",
  "end_time": "2025-03-01T22:00:00Z",
  "total_tickets": 1000,
  "ticket_price": 50000.0,
  "currency": "KRW",
  "thumbnail_url": "https://example.com/thumbnail.png",
  "is_public": true
}

Response: 201 Created
{
  "message": "Event created successfully",
  "event": {
    "id": "uuid",
    "organization_id": "uuid",
    "title": "My Event",
    "description": "Event description",
    "location": "Seoul, Korea",
    "venue": "COEX Hall",
    "start_time": "2025-03-01T18:00:00Z",
    "end_time": "2025-03-01T22:00:00Z",
    "total_tickets": 1000,
    "available_tickets": 1000,
    "ticket_price": 50000.0,
    "currency": "KRW",
    "thumbnail_url": "https://example.com/thumbnail.png",
    "status": "draft",
    "is_public": true,
    "created_by": "uuid",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### Get Event by ID
```http
GET /api/events/:id
Authorization: Bearer {token}

Response: 200 OK
{
  "event": {
    "id": "uuid",
    "organization_id": "uuid",
    "title": "My Event",
    ...
  }
}
```

#### Get Organization Events
```http
GET /api/organizations/:orgId/events
Authorization: Bearer {token}

Response: 200 OK
{
  "events": [
    {
      "id": "uuid",
      "organization_id": "uuid",
      "title": "My Event",
      ...
    }
  ]
}
```

#### Update Event (Admin Only)
```http
PUT /api/events/:id
Authorization: Bearer {token}

Request Body:
{
  "title": "Updated Event Title",
  "description": "Updated description",
  "location": "Updated location",
  "venue": "Updated venue",
  "start_time": "2025-03-01T18:00:00Z",
  "end_time": "2025-03-01T22:00:00Z",
  "total_tickets": 1500,
  "ticket_price": 60000.0,
  "currency": "KRW",
  "thumbnail_url": "https://example.com/new-thumbnail.png",
  "status": "published",
  "is_public": true
}

Response: 200 OK
{
  "message": "Event updated successfully"
}
```

#### Delete Event (Admin Only)
```http
DELETE /api/events/:id
Authorization: Bearer {token}

Response: 200 OK
{
  "message": "Event deleted successfully"
}
```

### Public Event Endpoints (No Authentication Required)

#### Get All Public Events
```http
GET /public/events

Response: 200 OK
{
  "events": [
    {
      "id": "uuid",
      "organization_id": "uuid",
      "title": "Public Event",
      "description": "Event description",
      "location": "Seoul, Korea",
      "venue": "COEX Hall",
      "start_time": "2025-03-01T18:00:00Z",
      "end_time": "2025-03-01T22:00:00Z",
      "total_tickets": 1000,
      "available_tickets": 850,
      "ticket_price": 50000.0,
      "currency": "KRW",
      "thumbnail_url": "https://example.com/thumbnail.png",
      "status": "published",
      "is_public": true,
      "created_by": "uuid",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z",
      "organization_name": "My Organization"
    }
  ]
}
```

#### Get Upcoming Events
```http
GET /public/events/upcoming

Response: 200 OK
{
  "events": [...]
}
```

#### Search Events
```http
GET /public/events/search?q=concert

Response: 200 OK
{
  "events": [...]
}
```

## Role-Based Access Control

### Roles

1. **Owner**: The user who created the organization
   - All admin permissions
   - Can delete the organization
   - Cannot be removed from the organization
   - Role cannot be changed

2. **Admin**: Organization administrators
   - Create, update, and delete events
   - Add and remove members
   - Update member roles
   - Update organization information

3. **Member**: Regular organization members
   - View organization information
   - View organization members
   - View organization events
   - Cannot create or modify events

### Permission Matrix

| Action | Owner | Admin | Member |
|--------|-------|-------|--------|
| View organization | ✓ | ✓ | ✓ |
| Update organization | ✓ | ✓ | ✗ |
| Delete organization | ✓ | ✗ | ✗ |
| View members | ✓ | ✓ | ✓ |
| Add members | ✓ | ✓ | ✗ |
| Remove members | ✓ | ✓ | ✗ |
| Update member roles | ✓ | ✓ | ✗ |
| Create events | ✓ | ✓ | ✗ |
| Update events | ✓ | ✓ | ✗ |
| Delete events | ✓ | ✓ | ✗ |
| View events | ✓ | ✓ | ✓ |

## Event Status

- **draft**: Event is being prepared, not visible to public
- **published**: Event is published and visible to public
- **ongoing**: Event is currently happening
- **completed**: Event has finished
- **cancelled**: Event has been cancelled

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 403 Forbidden
```json
{
  "error": "permission denied: admin role required"
}
```

### 404 Not Found
```json
{
  "error": "Organization not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error message"
}
```

## Examples

### Creating an Organization and Event Flow

1. **Create Organization**
```bash
curl -X POST http://localhost:3000/api/organizations \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tech Conference Organizers",
    "description": "We organize tech conferences"
  }'
```

2. **Add Admin Member**
```bash
curl -X POST http://localhost:3000/api/organizations/{org_id}/members \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "member-uuid",
    "role": "admin"
  }'
```

3. **Create Event**
```bash
curl -X POST http://localhost:3000/api/organizations/{org_id}/events \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Tech Summit 2025",
    "description": "Annual tech summit",
    "location": "Seoul",
    "venue": "COEX",
    "start_time": "2025-06-01T09:00:00Z",
    "end_time": "2025-06-01T18:00:00Z",
    "total_tickets": 500,
    "ticket_price": 100000.0,
    "currency": "KRW",
    "is_public": true
  }'
```

4. **Publish Event**
```bash
curl -X PUT http://localhost:3000/api/events/{event_id} \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "published"
  }'
```

## Notes

- All authenticated endpoints require a valid JWT token in the Authorization header
- User ID is automatically extracted from the JWT token (stored in `c.Locals("userID")`)
- When an organization is created, the creator is automatically added as an admin member
- Available tickets are automatically set to equal total tickets when an event is created
- When updating total tickets, available tickets are adjusted proportionally
