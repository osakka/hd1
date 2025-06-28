# THD (The Holo-Deck) - Professional Development Standards

> **"Where 3D visualization meets professional engineering"**

## 🚨 THE LAW: NON-NEGOTIABLE STANDARDS

This document establishes the **absolute, non-negotiable standards** for THD (The Holo-Deck) development. These are not suggestions - they are **LAW**.

---

## 📋 **CARDINAL RULES**

### **RULE #1: API-FIRST DEVELOPMENT**
- **`/opt/holo-deck/src/api.yaml` IS THE SINGLE SOURCE OF TRUTH**
- ALL routing is auto-generated from specification
- NO manual route configuration EVER
- Change the spec = change the API automatically
- **VIOLATION = IMMEDIATE ROLLBACK**

### **RULE #2: PROFESSIONAL STANDARDS ONLY**
- **NO EMOJIS** in system output, logs, or error messages
- **ABSOLUTE PATHS ONLY** - no relative paths in production code
- **LONG FLAGS ONLY** - no short flags to eliminate confusion
- **PROFESSIONAL ERROR MESSAGES** - clear, actionable, no decorative elements
- **VIOLATION = CODE REJECTION**

### **RULE #3: ZERO REGRESSIONS**
- ALL changes must maintain backward compatibility
- API endpoints cannot be removed without deprecation cycle
- Breaking changes require major version bump
- **VIOLATION = DEPLOYMENT BLOCKED**

### **RULE #4: BUILD SYSTEM VALIDATION**
- Code generation MUST pass before deployment
- Missing handlers = build failure
- API specification MUST validate
- **VIOLATION = RELEASE BLOCKED**

---

## 🏗️ **PROJECT STRUCTURE LAW**

### **CORRECT STRUCTURE:**
```
/opt/holo-deck/
├── src/                          # Go source code & module
│   ├── go.mod                   # CORRECT: module holodeck
│   ├── go.sum                   # Go dependencies
│   ├── main.go                  # THD daemon entry point
│   ├── Makefile                 # Professional build system
│   ├── api.yaml                 # OpenAPI 3.0.3 specification
│   ├── auto_router.go           # AUTO-GENERATED - DO NOT EDIT
│   ├── api/                     # Handler implementations
│   ├── server/                  # Core server components
│   └── codegen/                 # Code generation tools
├── build/                       # Build artifacts (gitignored)
│   ├── bin/thd                  # Professional daemon binary
│   ├── bin/thd-client           # API client
│   ├── logs/                    # Timestamped logs
│   └── runtime/                 # PID files, runtime data
├── share/                       # Static assets
│   ├── htdocs/                  # Web interface
│   ├── static/                  # WebGL renderer
│   └── templates/               # HTML templates
├── docs/                        # Documentation
│   ├── README.md                # THIS FILE - THE LAW
│   ├── adrs/                    # Architecture Decision Records
│   └── api/                     # API documentation
├── CLAUDE.md                    # Development context
├── CHANGELOG.md                 # Project history
└── README.md                    # Project overview
```

### **FORBIDDEN ITEMS:**
- ❌ **NO go.mod in project root** (belongs in src/)
- ❌ **NO legacy scripts** (dev-control.sh, etc.)
- ❌ **NO VWS references** (project is THD)
- ❌ **NO visualstream artifacts** (legacy naming)
- ❌ **NO relative paths** in production code
- ❌ **NO manual routing** (auto-generated only)

---

## ⚙️ **DEVELOPMENT WORKFLOW LAW**

### **MANDATORY WORKFLOW:**
1. **Specification First** - Update `api.yaml` before code
2. **Generate Routes** - `make generate` creates routing
3. **Implement Handlers** - Write handler functions
4. **Build & Validate** - `make all` ensures completeness
5. **Professional Testing** - `make test` validates API
6. **Daemon Control** - `make start/stop/status`
7. **Commit Standards** - Professional commit messages

### **PROFESSIONAL COMMANDS:**
```bash
# CORRECT usage - in /opt/holo-deck/src/
make all        # Complete build pipeline
make start      # Start THD daemon
make stop       # Stop THD daemon
make status     # Professional status reporting
make test       # Test API endpoints
make generate   # Generate routing from spec
make clean      # Clean build artifacts
```

### **FORBIDDEN COMMANDS:**
- ❌ `cd /home/claude-3/` (wrong path)
- ❌ `./dev-control.sh` (removed)
- ❌ `vws` commands (legacy)
- ❌ Short flags (`-d`, `-h`) (use `--daemon`, `--help`)

---

## 🎯 **TECHNICAL STANDARDS LAW**

### **COORDINATE SYSTEM:**
- **UNIVERSAL BOUNDARIES:** `[-12, +12]` on all axes
- **GRID SIZE:** 25×25×25 fixed grid
- **VALIDATION:** All coordinates MUST be within bounds
- **NO EXCEPTIONS**

### **SESSION MANAGEMENT:**
- **THREAD-SAFE:** All session operations use mutex
- **ISOLATION:** Each session has independent object store
- **PERSISTENCE:** Named objects with full lifecycle
- **REAL-TIME:** WebSocket broadcasts for live updates

### **API STANDARDS:**
- **OpenAPI 3.0.3** specification compliance
- **JSON responses** with consistent error formats
- **RESTful patterns** with proper HTTP status codes
- **Validation** at API boundary level

### **LOGGING STANDARDS:**
- **TIMESTAMPED:** All log entries have timestamps
- **STRUCTURED:** JSON format for machine parsing
- **PROFESSIONAL:** No decorative elements
- **ABSOLUTE PATHS:** Log files in `/opt/holo-deck/build/logs/`

---

## 🔒 **SECURITY LAW**

### **MANDATORY SECURITY:**
- **NO SECRETS** in code or logs
- **NO KEYS** committed to repository
- **VALIDATED INPUT** at all API endpoints
- **COORDINATE BOUNDS** enforced absolutely
- **CLEAN SHUTDOWN** procedures for all processes

### **FORBIDDEN PRACTICES:**
- ❌ Hardcoded credentials
- ❌ Unvalidated user input
- ❌ Unbounded coordinates
- ❌ Memory leaks in long-running processes
- ❌ Insecure WebSocket origins (in production)

---

## 🚀 **DEPLOYMENT LAW**

### **PRODUCTION REQUIREMENTS:**
- **DAEMON MODE:** `thd --daemon` with proper PID management
- **ABSOLUTE PATHS:** All paths must be absolute
- **PROFESSIONAL LOGGING:** Timestamped, structured output
- **CLEAN SHUTDOWN:** Proper resource cleanup
- **VALIDATION:** Build system prevents incomplete deployments

### **DEPLOYMENT CHECKLIST:**
- [ ] API specification validates
- [ ] All handlers implemented
- [ ] Auto-router generates successfully
- [ ] Tests pass completely
- [ ] Daemon starts and stops cleanly
- [ ] No memory leaks detected
- [ ] Coordinate validation enforced
- [ ] Professional logging configured

---

## ⚡ **EMERGENCY PROCEDURES**

### **IF DAEMON FAILS:**
```bash
cd /opt/holo-deck/src
make stop           # Clean shutdown
make force-stop     # Kill all processes
make clean          # Clear artifacts
make all            # Full rebuild
make start          # Restart daemon
```

### **IF BUILD FAILS:**
```bash
cd /opt/holo-deck/src
make validate       # Check API spec
make generate       # Regenerate router
make build          # Build with validation
```

### **IF API BROKEN:**
```bash
# Check specification
cat api.yaml | grep -i error

# Validate handlers exist
find api/ -name "*.go" | wc -l

# Test core endpoints
curl -s http://localhost:8080/api/sessions
```

---

## 📊 **MONITORING LAW**

### **MANDATORY MONITORING:**
- **Daemon Status:** `make status` shows complete health
- **API Health:** Core endpoints must respond
- **Resource Usage:** Memory and CPU within limits
- **Log Rotation:** Prevent disk space issues

### **ALERT CONDITIONS:**
- 🚨 Daemon not responding
- 🚨 API endpoints returning errors
- 🚨 Memory usage exceeding limits
- 🚨 Coordinate validation failures
- 🚨 WebSocket connection failures

---

## 📚 **DOCUMENTATION LAW**

### **MANDATORY DOCUMENTATION:**
- **API Changes:** Update OpenAPI specification
- **Architecture Decisions:** ADR documents required
- **Breaking Changes:** CHANGELOG.md updates
- **Development Context:** CLAUDE.md maintenance

### **DOCUMENTATION STANDARDS:**
- **PROFESSIONAL LANGUAGE:** No informal language
- **ACCURATE PATHS:** All paths must be current
- **COMPLETE EXAMPLES:** Working code samples
- **VERSION TRACKING:** Proper semantic versioning

---

## ⚖️ **VIOLATION CONSEQUENCES**

### **IMMEDIATE ACTIONS:**
- **Code Rejection** - Non-compliant code rejected
- **Build Failure** - Specification violations block builds
- **Deployment Block** - Professional standards not met
- **Rollback Required** - Breaking changes reverted

### **ZERO TOLERANCE:**
- Hardcoded paths
- Emoji in system output
- Manual routing
- Unvalidated coordinates
- Insecure practices

---

## 🎖️ **EXCELLENCE STANDARDS**

THD represents the evolution from innovative concept to **professional-grade engineering solution**. We maintain revolutionary capabilities while implementing proper software engineering practices.

### **CORE PHILOSOPHY:**
> **"Professional engineering excellence with zero compromise on innovation"**

### **SUCCESS METRICS:**
- ✅ 100% API-driven development
- ✅ Zero regressions ever
- ✅ Professional daemon control
- ✅ Specification-driven architecture
- ✅ Thread-safe session management
- ✅ Real-time 3D visualization

---

## 📞 **SUPPORT & ESCALATION**

### **DEVELOPMENT ISSUES:**
1. Check `make status` for daemon health
2. Review `/opt/holo-deck/build/logs/` for errors
3. Validate API specification syntax
4. Verify all handlers implemented
5. Rebuild with `make all`

### **EMERGENCY CONTACTS:**
- **Build Issues:** Check ADR documents
- **API Problems:** Validate against specification
- **Performance Issues:** Monitor resource usage
- **Security Concerns:** Review security checklist

---

**REMEMBER: These standards are non-negotiable. They ensure THD remains a professional, reliable, and innovative 3D visualization platform.**

---

*Last Updated: 2025-06-28*  
*THD Version: 2.0.0*  
*Authority: Professional Development Standards*