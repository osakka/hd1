# THD Library Architecture

**Professional Upstream API Libraries - Auto-Generated from Specification**

This directory contains the **upstream core libraries** for THD (The Holo-Deck), providing professional API wrapper functionality auto-generated directly from the OpenAPI specification.

## 📋 Upstream Libraries

### [thdlib.sh](thdlib.sh)
**Core Shell API Library** - Professional shell wrapper for all THD API endpoints

- **Auto-generated** from `api.yaml` specification
- **Single source of truth** - zero manual synchronization 
- **Professional HTTP client** with unified error handling
- **Complete API coverage** - every endpoint wrapped with validation

**Usage:**
```bash
source /opt/holo-deck/lib/thdlib.sh

# Core functions auto-generated from API spec
thd::create_object "cube1" "box" 0 1 0
thd::camera 5 5 5
thd::canvas_control "clear"
```

### [thdlib.js](../share/htdocs/static/js/thdlib.js)
**Core JavaScript API Client** - Professional web client for all THD API endpoints

- **Auto-generated** from `api.yaml` specification  
- **Identical API coverage** to shell library
- **Professional async/await** HTTP client
- **Type-safe parameter** handling

**Usage:**
```javascript
// Global instance automatically available
const result = await thdAPI.createObject('session-id', {
    name: 'cube1',
    type: 'box', 
    x: 0, y: 1, z: 0
});
```

## 🏗️ Architecture Principles

### Single Source of Truth
- **API Specification Drives Everything** - `api.yaml` is the authoritative source
- **Zero Manual Synchronization** - change spec, regenerate libraries automatically
- **Perfect Consistency** - shell and JavaScript libraries have identical coverage

### Professional Standards
- **Comprehensive Error Handling** - actionable error messages with context
- **Parameter Validation** - all inputs validated before API calls
- **Professional Logging** - structured, timestamped output
- **Enterprise Quality** - production-ready reliability

### Revolutionary Generation
Generated via THD's revolutionary code generator:
```bash
cd /opt/holo-deck/src && make generate
```

This command auto-generates:
- `lib/thdlib.sh` - Shell API wrapper
- `share/htdocs/static/js/thdlib.js` - JavaScript API client
- `src/auto_router.go` - Go routing layer
- Complete UI components and forms

## 🔗 Integration with Downstream

The upstream libraries serve as the foundation for **downstream integrations**:

### Downstream A-Frame Integration
```
lib/thdlib.sh (upstream core)
    ↓ imported by
lib/downstream/aframelib.sh (A-Frame integration)
```

See [downstream/README.md](downstream/README.md) for A-Frame specific capabilities.

## 📊 API Coverage

Current upstream library coverage:

| Category | Endpoints | Shell Functions | JS Methods |
|----------|-----------|----------------|------------|
| Sessions | 4 | ✅ Complete | ✅ Complete |
| Objects | 5 | ✅ Complete | ✅ Complete |
| Camera | 2 | ✅ Complete | ✅ Complete |
| Browser | 2 | ✅ Complete | ✅ Complete |
| World | 2 | ✅ Complete | ✅ Complete |
| **Total** | **28** | **✅ 100%** | **✅ 100%** |

## 🚀 Development Workflow

### For API Changes
1. **Update specification**: Edit `src/api.yaml`
2. **Regenerate libraries**: `make generate`
3. **Test integration**: Verify downstream compatibility
4. **Deploy**: Libraries automatically updated

### For New Endpoints
1. **Add to specification**: Define in `src/api.yaml`
2. **Implement handler**: Create Go handler file
3. **Regenerate**: `make generate` creates shell + JS functions
4. **Use immediately**: Functions available in both environments

## 🎯 Quality Assurance

### Automated Testing
- **Specification validation** - OpenAPI schema validation
- **Generation verification** - ensure all endpoints covered
- **Integration testing** - downstream compatibility maintained

### Professional Standards
- **No manual edits** - libraries are 100% generated
- **Consistent patterns** - identical function signatures across languages
- **Error handling** - professional error reporting throughout

---

**The upstream libraries represent the foundation of THD's revolutionary specification-driven architecture - where changing the API specification automatically updates all client libraries across all environments.**