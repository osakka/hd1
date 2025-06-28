# THD (The Holo-Deck) API Documentation

## Overview

The THD (The Holo-Deck) API enables creating, managing, and interacting with 3D virtual worlds through a RESTful interface. All endpoints are automatically generated from the OpenAPI 3.0.3 specification (`src/api.yaml`).

## Base URL
```
http://localhost:8080/api
```

## Authentication
*Currently, THD operates without authentication. All endpoints are publicly accessible.*

---

## Session Management

### Create Virtual World
Create a new isolated virtual world session with automatic world initialization.

```http
POST /sessions
```

**AUTOMATED BOOTSTRAPPING:** This endpoint now automatically initializes the 3D world coordinate system with a 25×25×25 grid and [-12, +12] bounds.

**Response:**
```json
{
    "success": true,
    "session_id": "session-abc123xyz",
    "created_at": "2025-06-28T15:30:00Z",
    "status": "active",
    "world": {
        "size": 25,
        "transparency": 0.1,
        "camera_x": 10,
        "camera_y": 10,
        "camera_z": 10
    },
    "bounds": {
        "min": -12,
        "max": 12
    },
    "coordinate_system": "fixed_grid",
    "message": "Session created with world ready - THD holo-deck activated"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/sessions
```

### List All Virtual Worlds
Retrieve all active virtual world sessions.

```http
GET /sessions
```

**Response:**
```json
{
    "sessions": [
        {
            "id": "session-abc123xyz",
            "created_at": "2025-06-28T15:30:00Z",
            "status": "active",
            "object_count": 3,
            "world_initialized": true,
            "world": {
                "size": 25,
                "transparency": 0.2,
                "camera_x": 15,
                "camera_y": 15,
                "camera_z": 15
            }
        }
    ],
    "total": 1,
    "timestamp": "2025-06-28T15:35:00Z"
}
```

### Get Virtual World Details
Retrieve comprehensive details about a specific virtual world.

```http
GET /sessions/{sessionId}
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Response:**
```json
{
    "id": "session-abc123xyz",
    "created_at": "2025-06-28T15:30:00Z",
    "status": "active",
    "object_count": 2,
    "world_initialized": true,
    "objects": [
        {
            "name": "cube1",
            "type": "cube",
            "x": 0,
            "y": 0,
            "z": 0,
            "scale": 1
        },
        {
            "name": "sphere1",
            "type": "sphere",
            "x": 5,
            "y": 3,
            "z": -2,
            "scale": 1.5,
            "color": "blue"
        }
    ],
    "world": {
        "size": 25,
        "transparency": 0.2,
        "camera_x": 15,
        "camera_y": 15,
        "camera_z": 15
    }
}
```

### Terminate Virtual World
Permanently delete a virtual world and all its data.

```http
DELETE /sessions/{sessionId}
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Response:**
```json
{
    "success": true,
    "message": "Session terminated and all data removed",
    "session_id": "session-abc123xyz"
}
```

---

## World Configuration

### Initialize Virtual World
Set up the 3D coordinate system, camera position, and world parameters.

```http
POST /sessions/{sessionId}/world
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Request Body:**
```json
{
    "size": 25,
    "transparency": 0.2,
    "camera_x": 15,
    "camera_y": 15,
    "camera_z": 15
}
```

**Request Fields:**
- `size` (integer, optional): Grid size (5-50, default: 25)
- `transparency` (float, optional): World transparency (0.0-1.0, default: 0.1)
- `camera_x` (float, optional): Camera X position (default: 10)
- `camera_y` (float, optional): Camera Y position (default: 10)
- `camera_z` (float, optional): Camera Z position (default: 10)

**Response:**
```json
{
    "success": true,
    "world": {
        "size": 25,
        "transparency": 0.2,
        "camera_x": 15,
        "camera_y": 15,
        "camera_z": 15
    },
    "bounds": {
        "min": -12,
        "max": 12
    },
    "coordinate_system": "fixed_grid",
    "message": "World initialized - ready for objects"
}
```

### Get World Specifications
Retrieve the current world configuration.

```http
GET /sessions/{sessionId}/world
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Response:**
```json
{
    "world": {
        "size": 25,
        "transparency": 0.2,
        "camera_x": 15,
        "camera_y": 15,
        "camera_z": 15
    },
    "bounds": {
        "min": -12,
        "max": 12
    },
    "coordinate_system": "fixed_grid"
}
```

---

## Object Management

### Create 3D Object
Add a new 3D object to the virtual world.

```http
POST /sessions/{sessionId}/objects
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Request Body:**
```json
{
    "name": "cube1",
    "type": "cube",
    "x": 0,
    "y": 0,
    "z": 0,
    "color": "red",
    "scale": 1.5
}
```

**Request Fields:**
- `name` (string, required): Unique object identifier within session
- `type` (string, required): Object geometry type (`cube`, `sphere`, etc.)
- `x` (float, required): X coordinate [-12, +12]
- `y` (float, required): Y coordinate [-12, +12]
- `z` (float, required): Z coordinate [-12, +12]
- `color` (string, optional): Color specification
- `scale` (float, optional): Size multiplier (default: 1.0)

**Response:**
```json
{
    "success": true,
    "object": {
        "name": "cube1",
        "type": "cube",
        "x": 0,
        "y": 0,
        "z": 0,
        "color": "red",
        "scale": 1.5
    },
    "message": "Object created successfully"
}
```

**Error Response (Invalid Coordinates):**
```json
{
    "error": "Coordinates must be within [-12, +12] bounds"
}
```

### List All Objects
Retrieve all 3D objects in the virtual world.

```http
GET /sessions/{sessionId}/objects
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Response:**
```json
{
    "objects": [
        {
            "name": "cube1",
            "type": "cube",
            "x": 0,
            "y": 0,
            "z": 0,
            "scale": 1
        },
        {
            "name": "sphere1",
            "type": "sphere",
            "x": 5,
            "y": 3,
            "z": -2,
            "scale": 1.5,
            "color": "blue"
        }
    ],
    "total": 2,
    "session_id": "session-abc123xyz"
}
```

### Get Object Details
Retrieve detailed information about a specific 3D object.

```http
GET /sessions/{sessionId}/objects/{objectName}
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier
- `objectName` (string, required): Unique object identifier

**Response:**
```json
{
    "object": {
        "name": "cube1",
        "type": "cube",
        "x": 0,
        "y": 0,
        "z": 0,
        "scale": 1,
        "color": "red"
    },
    "session_id": "session-abc123xyz"
}
```

### Update Object Properties
Modify properties of an existing 3D object.

```http
PUT /sessions/{sessionId}/objects/{objectName}
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier
- `objectName` (string, required): Unique object identifier

**Request Body:**
```json
{
    "x": 5,
    "y": 2,
    "color": "blue",
    "scale": 2.0
}
```

**Available Fields:**
- `x` (float): X coordinate [-12, +12]
- `y` (float): Y coordinate [-12, +12]
- `z` (float): Z coordinate [-12, +12]
- `color` (string): Color specification
- `scale` (float): Size multiplier
- `rotation` (string): Rotation specification

**Response:**
```json
{
    "success": true,
    "object": {
        "name": "cube1",
        "type": "cube",
        "x": 5,
        "y": 2,
        "z": 0,
        "scale": 2.0,
        "color": "blue"
    },
    "message": "Object updated successfully"
}
```

### Delete 3D Object
Remove an object from the virtual world.

```http
DELETE /sessions/{sessionId}/objects/{objectName}
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier
- `objectName` (string, required): Unique object identifier

**Response:**
```json
{
    "success": true,
    "message": "Object deleted successfully",
    "object_name": "cube1"
}
```

---

## Camera Controls

### Set Camera Position
Position the camera at specific 3D coordinates.

```http
PUT /sessions/{sessionId}/camera/position
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Request Body:**
```json
{
    "x": 20,
    "y": 15,
    "z": 10,
    "look_at_x": 0,
    "look_at_y": 0,
    "look_at_z": 0
}
```

**Request Fields:**
- `x` (float, required): Camera X position
- `y` (float, required): Camera Y position
- `z` (float, required): Camera Z position
- `look_at_x` (float, optional): Look-at target X (default: 0)
- `look_at_y` (float, optional): Look-at target Y (default: 0)
- `look_at_z` (float, optional): Look-at target Z (default: 0)

**Response:**
```json
{
    "success": true,
    "camera": {
        "x": 20,
        "y": 15,
        "z": 10,
        "look_at_x": 0,
        "look_at_y": 0,
        "look_at_z": 0
    },
    "message": "Camera position updated"
}
```

### Start Camera Orbit
Begin automatic orbital motion around a focal point.

```http
POST /sessions/{sessionId}/camera/orbit
```

**Parameters:**
- `sessionId` (string, required): Unique session identifier

**Request Body:**
```json
{
    "center_x": 0,
    "center_y": 0,
    "center_z": 0,
    "radius": 20,
    "speed": 1.0,
    "elevation": 15
}
```

**Request Fields:**
- `center_x` (float, optional): Orbit center X (default: 0)
- `center_y` (float, optional): Orbit center Y (default: 0)
- `center_z` (float, optional): Orbit center Z (default: 0)
- `radius` (float, optional): Orbit radius (default: 20)
- `speed` (float, optional): Rotation speed (default: 1.0)
- `elevation` (float, optional): Camera elevation (default: 15)

**Response:**
```json
{
    "success": true,
    "orbit": {
        "center_x": 0,
        "center_y": 0,
        "center_z": 0,
        "radius": 20,
        "speed": 1.0,
        "elevation": 15,
        "active": true
    },
    "message": "Camera orbit started"
}
```

---

## WebSocket Real-time Events

THD provides real-time updates via WebSocket connection at:
```
ws://localhost:8080/ws
```

### Event Types

#### Session Events
```json
{
    "type": "session_created",
    "data": {
        "session_id": "session-abc123xyz",
        "created_at": "2025-06-28T15:30:00Z"
    },
    "timestamp": 1672531200
}
```

```json
{
    "type": "session_deleted",
    "data": {
        "session_id": "session-abc123xyz"
    },
    "timestamp": 1672531200
}
```

#### World Events
```json
{
    "type": "world_initialized",
    "data": {
        "session_id": "session-abc123xyz",
        "world": {
            "size": 25,
            "transparency": 0.2,
            "camera_x": 15,
            "camera_y": 15,
            "camera_z": 15
        }
    },
    "timestamp": 1672531200
}
```

#### Object Events
```json
{
    "type": "object_created",
    "data": {
        "session_id": "session-abc123xyz",
        "object": {
            "name": "cube1",
            "type": "cube",
            "x": 0,
            "y": 0,
            "z": 0,
            "scale": 1
        }
    },
    "timestamp": 1672531200
}
```

```json
{
    "type": "object_updated",
    "data": {
        "session_id": "session-abc123xyz",
        "object": {
            "name": "cube1",
            "type": "cube",
            "x": 5,
            "y": 2,
            "z": 0,
            "scale": 2.0,
            "color": "blue"
        }
    },
    "timestamp": 1672531200
}
```

---

## Error Handling

### HTTP Status Codes
- `200 OK` - Successful GET request
- `201 Created` - Successful POST request (resource created)
- `400 Bad Request` - Invalid request format or parameters
- `404 Not Found` - Session or object not found
- `500 Internal Server Error` - Server error

### Error Response Format
```json
{
    "error": "Descriptive error message",
    "code": "ERROR_CODE",
    "details": {
        "field": "Additional error context"
    }
}
```

### Common Errors

#### Invalid Coordinates
```json
{
    "error": "Coordinates must be within [-12, +12] bounds"
}
```

#### Session Not Found
```json
{
    "error": "Session not found"
}
```

#### Object Already Exists
```json
{
    "error": "Object with name 'cube1' already exists in session"
}
```

#### Missing Required Fields
```json
{
    "error": "Missing required fields: name, type"
}
```

---

## Complete Workflow Example

### 1. Create Virtual World
```bash
SESSION_ID=$(curl -s -X POST http://localhost:8080/api/sessions | jq -r '.session_id')
echo "Created session: $SESSION_ID"
```

### 2. Initialize World
```bash
curl -X POST http://localhost:8080/api/sessions/$SESSION_ID/world \
  -H "Content-Type: application/json" \
  -d '{
    "size": 25,
    "transparency": 0.2,
    "camera_x": 15,
    "camera_y": 15,
    "camera_z": 15
  }'
```

### 3. Create Objects
```bash
# Create a cube at origin
curl -X POST http://localhost:8080/api/sessions/$SESSION_ID/objects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "cube1",
    "type": "cube",
    "x": 0,
    "y": 0,
    "z": 0,
    "color": "red"
  }'

# Create a sphere
curl -X POST http://localhost:8080/api/sessions/$SESSION_ID/objects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sphere1",
    "type": "sphere",
    "x": 5,
    "y": 3,
    "z": -2,
    "color": "blue",
    "scale": 1.5
  }'
```

### 4. Position Camera
```bash
curl -X PUT http://localhost:8080/api/sessions/$SESSION_ID/camera/position \
  -H "Content-Type: application/json" \
  -d '{
    "x": 20,
    "y": 15,
    "z": 10,
    "look_at_x": 0,
    "look_at_y": 0,
    "look_at_z": 0
  }'
```

### 5. Get Complete World State
```bash
curl http://localhost:8080/api/sessions/$SESSION_ID | jq .
```

### 6. Start Camera Orbit
```bash
curl -X POST http://localhost:8080/api/sessions/$SESSION_ID/camera/orbit \
  -H "Content-Type: application/json" \
  -d '{
    "radius": 25,
    "speed": 0.5,
    "elevation": 20
  }'
```

---

## Related Resources

- **[OpenAPI Specification](../../src/api.yaml)** - Complete API definition
- **[Architecture Guide](../architecture/)** - System design details
- **[Development Guide](../development/)** - Implementation details
- **[WebSocket Documentation](../guides/websocket-integration.md)** - Real-time integration

---

*This documentation is automatically maintained in sync with the OpenAPI specification through THD's specification-driven development approach.*