# HD1 Architecture Overview

HD1 (Holodeck One) is an API-first, specification-driven Three.js game engine platform designed for real-time 3D applications with ultra-minimal overhead.

## Core Principles

### 1. API-First Design
- **Single Source of Truth**: Complete system defined in `src/api.yaml` OpenAPI specification
- **15 REST Endpoints**: Three.js platform with avatar lifecycle management
- **Auto-Generation**: All routing, client libraries, and documentation generated from spec
- **Zero Manual Configuration**: Routes, handlers, and validation automatically maintained

### 2. Real-Time Synchronization
- **WebSocket Hub**: TCP-simple sequence-based synchronization protocol
- **Client State Management**: Intelligent reconnection with exponential backoff
- **Session Isolation**: Multi-tenant support with session-based entity management
- **High-Frequency Updates**: Optimized for real-time avatar position updates

### 3. Three.js Integration
- **Direct WebGL**: Zero abstraction layers for maximum performance
- **Ultra-Minimal Console**: Essential debugging and monitoring only
- **Progressive Enhancement**: Graceful degradation for different client capabilities
- **Asset Pipeline**: Optimized loading and caching strategies

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │   HTTP Client   │    │   Mobile App    │
│   (Three.js)    │    │   (REST API)    │    │   (WebSocket)   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ WebSocket            │ HTTP/REST            │ WebSocket
          │                      │                      │
    ┌─────▼──────────────────────▼──────────────────────▼─────┐
    │                HD1 Server Core                          │
    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
    │  │ WebSocket   │  │ HTTP Router │  │ Static File │     │
    │  │ Hub         │  │ (Generated) │  │ Server      │     │
    │  └─────────────┘  └─────────────┘  └─────────────┘     │
    │                                                         │
    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
    │  │ API         │  │ Config      │  │ Logging     │     │
    │  │ Handlers    │  │ Management  │  │ System      │     │
    │  └─────────────┘  └─────────────┘  └─────────────┘     │
    └─────────────────────────────────────────────────────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Code Generation     │
                    │   (from api.yaml)     │
                    └───────────────────────┘
```

## Data Flow

### HTTP Request Flow
1. **Client Request** → HTTP Router (auto-generated)
2. **Router** → API Handler (with validation)
3. **Handler** → Business Logic + WebSocket Broadcast
4. **Response** → Client (JSON + CORS headers)

### WebSocket Flow
1. **Client Connect** → WebSocket Upgrader
2. **Registration** → Hub with Session Association
3. **Message Flow** → Bidirectional real-time communication
4. **Broadcast** → All clients in session (filtered)

### Code Generation Flow
1. **api.yaml Changes** → Trigger regeneration
2. **Generator** → Parse specification + validate
3. **Templates** → Generate routing, handlers, clients
4. **Build** → Compile and deploy updated system

## Key Components

### Server Core (`main.go`)
- **Configuration Loading**: Environment, flags, .env files
- **HTTP Server Setup**: Static files, API routing, WebSocket endpoints
- **Daemon Management**: Process forking, PID files, signal handling
- **Graceful Shutdown**: Clean resource cleanup and connection termination

### WebSocket Hub (`server/hub.go`)
- **Client Registry**: Active connection management
- **Session Isolation**: Multi-tenant message routing
- **Broadcast Engine**: Efficient message distribution
- **Connection Lifecycle**: Registration, heartbeat, cleanup

### API Router (`auto_router.go`)
- **Auto-Generated**: Never manually edited
- **Specification-Driven**: Complete route mapping from api.yaml
- **CORS Handling**: Cross-origin request support
- **Context Injection**: Hub and configuration availability

### Configuration System (`config/`)
- **Priority Order**: Flags > Environment > .env > Defaults
- **Type Safety**: Structured configuration with validation
- **Hot Reload**: Runtime configuration updates
- **Environment Isolation**: Development, staging, production profiles

### Logging System (`logging/`)
- **Structured JSON**: Machine-readable log format
- **Module Tracing**: Fine-grained debug control
- **Performance**: Zero-overhead disabled log levels
- **Context Enrichment**: Request IDs, session data, timing

## Performance Characteristics

### Scalability
- **Concurrent Connections**: Thousands of WebSocket clients
- **Memory Efficiency**: Object pooling for high-frequency operations
- **CPU Optimization**: Zero-allocation hot paths
- **Network**: Optimized JSON serialization and WebSocket frames

### Reliability
- **Error Handling**: Comprehensive error recovery and logging
- **Health Monitoring**: Built-in health checks and metrics
- **Graceful Degradation**: Fallback mechanisms for component failures
- **Data Consistency**: Transaction-like operations for critical updates

### Security
- **Input Validation**: Automatic validation from OpenAPI schemas
- **CORS Protection**: Configurable cross-origin policies
- **Session Management**: Secure session isolation and cleanup
- **Rate Limiting**: Configurable request throttling

## Technology Stack

### Backend
- **Go 1.21+**: High-performance concurrent server
- **Gorilla WebSocket**: Production-grade WebSocket implementation
- **Gorilla Mux**: HTTP routing with parameter extraction
- **YAML v3**: Configuration and specification parsing

### Frontend
- **Three.js r150+**: WebGL 3D rendering engine
- **Native WebSocket**: Direct browser WebSocket API
- **Vanilla JavaScript**: Zero framework dependencies
- **CSS3**: Modern styling with animations

### Development
- **OpenAPI 3.0**: API specification and documentation
- **Make**: Build automation and development workflow
- **Git**: Version control with comprehensive history
- **JSON**: Structured logging and API communication

---

*Architecture documentation for HD1 v0.7.0*