# HD1 Getting Started Guide

## Quick Start

### Prerequisites
- Go 1.19+ installed
- Linux/macOS environment
- Network access to `localhost:8080`

### Installation
```bash
# Clone repository
git clone https://git.uk.home.arpa/itdlabs/holo-deck.git
cd holo-deck

# Build HD1
cd src && make all

# Start daemon
make start
```

### First Steps
1. Navigate to http://localhost:8080
2. Create a session using the console
3. Load a scene or create objects
4. Use WASD to move, mouse to look around

### Basic API Usage
```bash
# Create session
SESSION_ID=$(./build/bin/hd1-client create-session | jq -r '.session_id')

# Create object
./build/bin/hd1-client create-object "$SESSION_ID" '{"name": "cube1", "type": "cube", "x": 0, "y": 1, "z": 0}'

# List objects
./build/bin/hd1-client list-objects "$SESSION_ID"
```

## Next Steps
- Read the [User Manual](user-manual.md)
- Explore [API Documentation](../api/README.md)
- Review [Architecture Overview](../architecture/system-architecture.md)