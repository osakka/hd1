# ADR-012: Professional Status Bar System

## Status
ACCEPTED (2025-06-30)

## Context
The HD1 holodeck console required a professional status indication system to provide clear visual feedback about connection health, session state, and user interaction modes. The previous implementation mixed different status types in confusing ways, leading to user uncertainty about system state.

## Problem Statement
1. **Status Confusion**: Single status LED mixed connection health with pointer lock state
2. **Inconsistent UI Heights**: Title bar, scene bar, controls bar, and status bar had different heights
3. **Verbose Mouse Lock Text**: "Press ESC to release mouse" was unnecessarily long
4. **Session ID Truncation**: Removing "session-" prefix made debugging difficult
5. **Unprofessional Styling**: Inconsistent colors, alignment, and emoji usage in system output

## Decision
Implement a comprehensive professional status bar system with:

### Status Separation
- **Connection Health**: LED indicator with states (connecting, connected, receiving, error, disconnected)
- **Mouse Lock State**: Separate ESC indicator with orange theme (dim/active states)
- **Session Management**: Full session ID display with click-to-copy functionality

### Layout Standardization
- **Consistent Heights**: All bars use `min-height: 20px` and `padding: 6px 8px`
- **Professional Positioning**: Connection status left, mouse lock and session ID right
- **Flex Layout**: Proper alignment with `display: flex` and `align-items: center`

### Color Scheme
- **Connection LED**: Green (connected), orange (connecting), cyan (receiving), red (error)
- **Mouse Lock**: Orange theme - `rgba(255, 149, 0, 0.4)` dim, `rgba(255, 149, 0, 1.0)` active
- **Session ID**: Very faint `rgba(0, 255, 255, 0.3)` for minimal visual impact

### Professional Standards
- **No Emojis**: Removed from all auto-generated shell libraries
- **Concise Text**: Mouse lock shows only "ESC" when active
- **Full Session IDs**: No truncation for better debugging experience

## Implementation
```css
#debug-status-bar {
    padding: 6px 8px;
    min-height: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

#status-lock-indicator.locked {
    color: rgba(255, 149, 0, 1.0);
    text-shadow: 0 0 6px rgba(255, 149, 0, 0.8);
}
```

### JavaScript Logic
```javascript
function setLockStatus(status, message) {
    statusLockIndicator.className = status;
    statusLockIndicator.textContent = status === 'locked' ? 'ESC' : '';
}

function updateDebugSession(sessionId) {
    sessionIdTagStatus.textContent = sessionId; // Full ID, no truncation
}
```

## Consequences

### Positive
- **Clear Status Separation**: Users understand connection vs interaction state
- **Professional Appearance**: Consistent heights and professional color scheme
- **Better UX**: Concise ESC indicator, full session ID for debugging
- **Maintainable Code**: Standardized CSS patterns across all bars

### Negative
- **Increased Complexity**: More CSS rules and JavaScript state management
- **Breaking Change**: Full session ID display changes user expectations

## Compliance
- ✅ Professional engineering standards (no emojis, consistent styling)
- ✅ Single source of truth (centralized status management functions)
- ✅ Bar-raising performance (efficient CSS transitions, minimal DOM updates)
- ✅ Clean separation of concerns (connection health vs user interaction state)

## Related ADRs
- ADR-003: Professional UI Enhancement (builds upon this foundation)
- ADR-001: A-Frame WebXR Integration (status system supports VR/AR modes)

## Future Considerations
- Hover tooltips for status indicators could provide additional context
- Keyboard shortcuts could be indicated in the status bar
- VR mode status could be integrated into the same system