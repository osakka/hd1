# ADR-003: Three.js Minimal Console

**Status**: Accepted  
**Date**: 2025-07-14  
**Authors**: HD1 Development Team  
**Related**: ADR-005 (Ultra-Minimal Build Strategy)

## Context

HD1 requires a debugging and monitoring interface for development and operations. Traditional approaches include complex dashboards, heavy framework dependencies, or feature-rich consoles that add significant overhead. The system needs real-time visibility into WebSocket connections, API performance, and system health while maintaining minimal resource footprint.

## Decision

We implement an **Ultra-Minimal Three.js Console** that provides essential debugging capabilities with zero framework dependencies and minimal resource overhead.

### Core Principles

1. **Minimal Surface Area**: Only essential features included
2. **Zero Dependencies**: No frameworks, libraries, or external dependencies
3. **Real-Time Monitoring**: Live WebSocket and API status updates
4. **Progressive Enhancement**: Graceful degradation for different capabilities
5. **Performance First**: Sub-millisecond UI response times

### Feature Set

```javascript
// Ultra-minimal feature set
const consoleFeatures = {
    websocket: {
        connection: 'live status',
        reconnection: 'automatic with exponential backoff',
        statistics: 'real-time metrics'
    },
    debugging: {
        panel: 'collapsible debug information',
        logs: 'client-side event logging',
        rebootstrap: 'storage clearing + page reload'
    },
    monitoring: {
        version: 'client/server version checking',
        performance: 'connection latency tracking',
        health: 'basic system health indicators'
    }
};
```

## Implementation

### Minimal HTML Structure
```html
<!-- Ultra-minimal DOM structure -->
<div id="holodeck-container">
    <canvas id="holodeck-canvas"></canvas>
    <div id="debug-panel">
        <div id="debug-header">HD1 Console</div>
        <div id="debug-content">
            <div id="websocket-status"></div>
            <div id="system-info"></div>
        </div>
    </div>
</div>
```

### Zero-Framework JavaScript
```javascript
// No frameworks - pure JavaScript implementation
class HD1Console {
    constructor() {
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 99;
        this.init();
    }
    
    // Minimal implementation focused on essential features only
    init() {
        this.setupWebSocket();
        this.setupDebugPanel();
        this.setupRebootstrap();
    }
}
```

### Minimal CSS
```css
/* Ultra-minimal styling - essential visual elements only */
#debug-panel {
    background: rgba(0, 0, 0, 0.7);
    border: 1px solid rgba(0, 255, 255, 0.3);
    border-radius: 4px;
    font-family: 'Courier New', monospace;
    font-size: 12px;
    color: #00ffff;
}
```

## Key Features

### Intelligent WebSocket Reconnection
```javascript
// Exponential backoff with automatic recovery
connectWebSocket() {
    const backoffDelay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        this.triggerRebootstrap();
        return;
    }
    
    setTimeout(() => this.attemptConnection(), backoffDelay);
}
```

### Rebootstrap Recovery
```javascript
// Nuclear option - clear all storage and reload
triggerRebootstrap() {
    this.addDebug('REBOOTSTRAP', 'Clearing storage and reloading page...');
    localStorage.clear();
    sessionStorage.clear();
    
    // Clear cookies
    document.cookie.split(";").forEach(cookie => {
        const eqPos = cookie.indexOf("=");
        const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/";
    });
    
    setTimeout(() => window.location.reload(true), 1000);
}
```

### Version Synchronization
```javascript
// Automatic client/server version checking
checkVersion() {
    this.ws.send(JSON.stringify({
        type: 'version_check',
        js_version: this.getJSVersion()
    }));
}

// Automatic refresh on version mismatch
handleVersionMismatch(data) {
    this.addDebug('VERSION', `Mismatch detected: ${data.client_version} != ${data.server_version}`);
    this.addDebug('VERSION', 'Refreshing page in 3 seconds...');
    setTimeout(() => window.location.reload(true), 3000);
}
```

## Consequences

### Positive
- **Ultra-Fast Load Times**: Minimal HTML/CSS/JS bundle size
- **Zero Dependencies**: No external libraries or frameworks to manage
- **Reliable Recovery**: Rebootstrap provides guaranteed recovery mechanism
- **Real-Time Monitoring**: Live system status without overhead
- **Easy Maintenance**: Simple codebase with minimal attack surface

### Negative
- **Limited Features**: Only essential debugging capabilities
- **Manual Implementation**: No framework conveniences or abstractions
- **Visual Simplicity**: Basic styling without rich UI components
- **Browser Compatibility**: Modern browser features required

### Neutral
- **Custom Implementation**: All features implemented specifically for HD1
- **Performance Focus**: Optimized for speed over feature richness

## Technical Specifications

### Resource Requirements
- **HTML**: ~500 bytes (minimal DOM structure)
- **CSS**: ~2KB (essential styling only) 
- **JavaScript**: ~8KB (core functionality)
- **Memory**: <1MB runtime footprint
- **CPU**: <1% utilization during normal operation

### Browser Compatibility
- **WebSocket**: Modern browsers with native WebSocket support
- **ES6**: Arrow functions, template literals, classes
- **CSS3**: Flexbox, animations, rgba colors
- **DOM**: Modern DOM APIs and event handling

### Performance Characteristics
- **Startup Time**: <100ms initialization
- **WebSocket Latency**: <10ms reconnection detection
- **UI Response**: <1ms for debug panel operations
- **Memory Efficiency**: Minimal object allocation

## Integration Points

### Server Communication
```javascript
// Simple message protocol with server
const messageTypes = {
    version_check: 'Client/server version synchronization',
    ping: 'Latency measurement and keepalive', 
    debug_info: 'System status and health metrics',
    rebootstrap: 'Recovery mechanism activation'
};
```

### Three.js Integration
```javascript
// Minimal Three.js setup for holodeck rendering
function initThreeJS() {
    // Ultra-minimal Three.js initialization
    // Only essential components for debugging
}
```

## Alternative Approaches Considered

### Rich Dashboard Framework
**Rejected**: Heavy dependency overhead, complex build process, slower load times, unnecessary features for HD1's use case.

### Server-Side Rendering
**Rejected**: Adds server complexity, slower initial load, unnecessary for debugging interface.

### WebGL-Based Console
**Rejected**: Overkill for text-based debugging, compatibility issues, higher resource usage.

### Terminal-Based Interface
**Rejected**: Poor user experience, limited visual debugging capabilities, deployment complexity.

## Success Metrics

- **Load Time**: Console ready in <500ms
- **Reconnection**: WebSocket recovery in <30s for 99% of failures
- **Resource Usage**: <1MB memory, <1% CPU during operation
- **Reliability**: Rebootstrap recovery success rate >99%

## Future Enhancements

### Performance Monitoring
- **Real-time Metrics**: CPU, memory, network usage
- **Historical Data**: Performance trends and patterns
- **Alert System**: Threshold-based notifications

### Advanced Debugging
- **Message Inspection**: WebSocket message logging and filtering
- **State Visualization**: Real-time system state display
- **Error Tracking**: Client-side error capture and reporting

---

*ADR-003 defines HD1's ultra-minimal debugging and monitoring console*