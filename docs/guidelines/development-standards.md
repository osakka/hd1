# HD1 (Holodeck One) - Standard Development Standards

> **"Where 3D visualization meets standard engineering"**

## 🚨 THE LAW: NON-NEGOTIABLE STANDARDS

This document establishes the **absolute, non-negotiable standards** for HD1 (Holodeck One) development. These are not suggestions - they are **LAW**.

---

## 🚀 **xVC METHODOLOGY**

HD1 development follows **xVC (Extreme Vibe Coding)** principles for human-LLM collaboration:

### **Core xVC Principles**
- **Pattern Consistency**: Establish repeatable interaction patterns
- **Cognitive Amplification**: LLMs as "pattern reflection engines"
- **Precise**: Targeted, precise modifications
- **One Source of Truth**: Single authoritative specification
- **Bar-Raising Solutions**: Always improve, never regress
- **Forward Progress Only**: No backward compatibility for experimentation

### **xVC Development Stages**
1. **Learning Phase (Weeks 1-2)**: Establishing interaction patterns
2. **Productivity Phase (Weeks 3-8)**: Optimizing collaboration patterns  
3. **Proficiency Phase (Months 3+)**: Automatic pattern application

### **xVC Requirements**
- **Discipline**: Consistent application of patterns
- **Patience**: Allow patterns to develop over time
- **Precision**: Clear, structured inputs to LLMs

---

## 📋 **CARDINAL RULES**

### **RULE #1: API-FIRST DEVELOPMENT**
- **`/opt/holodeck-one/src/api.yaml` IS THE SINGLE SOURCE OF TRUTH**
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
/opt/holodeck-one/
├── src/                          # Go source code & module
│   ├── go.mod                   # CORRECT: module holodeck
│   ├── go.sum                   # Go dependencies
│   ├── main.go                  # HD1 daemon entry point
│   ├── Makefile                 # Standard build system
│   ├── api.yaml                 # OpenAPI 3.0.3 specification
│   ├── auto_router.go           # AUTO-GENERATED - DO NOT EDIT
│   ├── api/                     # Handler implementations
│   ├── server/                  # Core server components
│   └── codegen/                 # Code generation tools
├── build/                       # Build artifacts (gitignored)
│   ├── bin/thd                  # Standard daemon binary
│   ├── bin/thd-client           # API client
│   ├── logs/                    # Timestamped logs
│   └── runtime/                 # PID files, runtime data
├── share/                       # Web assets and static files
│   ├── htdocs/                  # Web interface root (HD1_HTDOCS_DIR)
│   │   ├── static/js/           # JavaScript: renderer.js, gl-matrix.js, debug.js
│   │   ├── debug.html           # Debug interface
│   │   └── force-session.html   # Session management tools
│   ├── configs/                 # Configuration templates
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
- ❌ **NO HD1 references** (project is HD1)
- ❌ **NO visualstream artifacts** (legacy naming)
- ❌ **NO relative paths** in production code
- ❌ **NO manual routing** (auto-generated only)
- ❌ **NO duplicate JavaScript files** (single source of truth in share/htdocs/static/js/)
- ❌ **NO manual edits to auto_router.go** (enhanced disclaimer prevents this)

---

## ⚙️ **DEVELOPMENT WORKFLOW LAW**

### **MANDATORY WORKFLOW:**
1. **Specification First** - Update `api.yaml` before code
2. **Generate Routes** - `make generate` creates routing
3. **Implement Handlers** - Write handler functions
4. **Build & Validate** - `make all` ensures completeness
5. **Standard Testing** - `make test` validates API
6. **Daemon Control** - `make start/stop/status`
7. **Commit Standards** - Standard commit messages

### **PROFESSIONAL COMMANDS:**
```bash
# CORRECT usage - in /opt/holodeck-one/src/
make all        # Complete build pipeline
make start      # Start HD1 daemon
make stop       # Stop HD1 daemon
make status     # Standard status reporting
make test       # Test API endpoints
make generate   # Generate routing from spec
make clean      # Clean build artifacts
```

### **FORBIDDEN COMMANDS:**
- ❌ Legacy path references (all paths must use `/opt/holodeck-one`)
- ❌ `./dev-control.sh` (removed)
- ❌ `hd1` commands (legacy)
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

### **SCENE FORKING & CROWN JEWEL (v3.4.0):**
- **FREEZE-FRAME MODE**: Scene snapshots with object tracking (base/modified/new)
- **TEMPORAL SEQUENCE MODE**: Complete session recording and playback
- **OBJECT PROVENANCE**: Full tracking of object lifecycle and source scenes
- **👑 CROWN JEWEL**: Auto-generated web UI client achieving 100% single source of truth
- **THREE-TIER GENERATION**: Go router + CLI client + Web UI client all from OpenAPI spec
- **ZERO MANUAL SYNC**: API changes automatically update all client systems

### **LOGGING STANDARDS:**
- **TIMESTAMPED:** All log entries have timestamps
- **STRUCTURED:** JSON format for machine parsing
- **PROFESSIONAL:** No decorative elements
- **ABSOLUTE PATHS:** Log files in `/opt/holodeck-one/build/logs/`

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
- [ ] Auto-generated client auto-generated files present (thd-api-client.js, thd-ui-components.js, thd-form-system.js)
- [ ] Scene forking system functional (fork/save endpoints)
- [ ] Recording system operational (start/stop/play/status endpoints)
- [ ] Tests pass completely
- [ ] Daemon starts and stops cleanly
- [ ] No memory leaks detected
- [ ] Coordinate validation enforced
- [ ] Standard logging configured

---

## ⚡ **EMERGENCY PROCEDURES**

### **IF DAEMON FAILS:**
```bash
cd /opt/holodeck-one/src
make stop           # Clean shutdown
make force-stop     # Kill all processes
make clean          # Clear artifacts
make all            # Full rebuild
make start          # Restart daemon
```

### **IF BUILD FAILS:**
```bash
cd /opt/holodeck-one/src
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

## 📝 **GIT HYGIENE LAW**

### **MANDATORY GIT STANDARDS:**
- **CLEAN COMMITS:** One logical change per commit
- **DESCRIPTIVE MESSAGES:** Clear, actionable commit messages
- **SMALL COMMITS:** Atomic changes for easy review
- **NO MERGE COMMITS:** Use rebase for clean history
- **PROFESSIONAL LANGUAGE:** No informal commit messages

### **COMMIT MESSAGE FORMAT:**
```
Brief summary (50 characters max)

Detailed explanation if needed:
• What was changed
• Why it was changed  
• Impact of the change

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

### **FORBIDDEN GIT PRACTICES:**
- ❌ **"WIP" commits** in main branch
- ❌ **Merge commits** (use rebase)
- ❌ **Large binary files** without LFS
- ❌ **Secrets in history** (immediate purge required)
- ❌ **Unclear commit messages** ("fix", "update", etc.)

### **BRANCH HYGIENE:**
- **master** - Production-ready code only
- **feature/** - New features (rebase before merge)
- **hotfix/** - Critical fixes (fast-track approval)
- **NO LONG-LIVED BRANCHES** (merge within 24-48 hours)

### **MANDATORY CHECKS BEFORE COMMIT:**
```bash
# 1. Verify clean working directory
git status

# 2. Review all changes
git diff

# 3. Test functionality
make test

# 4. Verify daemon still works
make status

# 5. Standard commit message
git commit -m "Descriptive message explaining the change"
```

---

## 🧹 **CLEAN WORKSPACE LAW**

### **MANDATORY WORKSPACE STANDARDS:**
- **NO TEMPORARY FILES** in project directory
- **NO BACKUP FILES** (*.bak, *.tmp, *.old)
- **NO PERSONAL CONFIG** in shared directories
- **NO DEVELOPMENT ARTIFACTS** in git
- **CLEAN DESKTOP POLICY** - organized file structure

### **WORKSPACE ORGANIZATION:**
```
/opt/holodeck-one/           # Project root - CLEAN
├── src/                  # Source code only
├── build/                # Build artifacts (gitignored)
├── docs/                 # Documentation only
└── share/                # Static assets only

# FORBIDDEN in project root:
❌ test.txt, debug.log, backup.zip
❌ Personal scripts or tools
❌ IDE configuration files
❌ Temporary development files
```

### **DAILY WORKSPACE HYGIENE:**
```bash
# Clean temporary files
find . -name "*.tmp" -delete
find . -name "*.bak" -delete
find . -name "*.old" -delete

# Remove editor artifacts
find . -name ".DS_Store" -delete
find . -name "Thumbs.db" -delete

# Check for forgotten files
git status --ignored
```

### **GITIGNORE COMPLIANCE:**
- **build/** directory ignored
- **IDE files** ignored (*.swp, *.swo, .vscode/)
- **OS files** ignored (.DS_Store, Thumbs.db)
- **Temporary files** ignored (*.tmp, *.log)
- **Personal configs** ignored

### **CLEAN DESKTOP POLICY:**
- **PROJECT FILES** only in `/opt/holodeck-one/`
- **NO SCATTERED FILES** on desktop or home directory
- **ORGANIZED DOWNLOADS** - clean up regularly
- **PROFESSIONAL NAMING** - no spaces, special characters
- **REGULAR CLEANUP** - weekly workspace maintenance

### **WORKSPACE VIOLATIONS:**
- 🚨 **Immediate cleanup required** for temp files
- 🚨 **Commit blocked** if workspace dirty
- 🚨 **Standard review** for disorganized structure
- 🚨 **Training required** for repeated violations

---

## 🚀 **RECENT PROFESSIONAL IMPROVEMENTS**

### **ENHANCED CODE GENERATION PROTECTION**
Auto-generated files now include comprehensive disclaimers:
```go
// ===================================================================
// WARNING: AUTO-GENERATED CODE - DO NOT MODIFY THIS FILE
// ===================================================================
//
// ⚠️  CRITICAL WARNING: ALL MANUAL CHANGES WILL BE LOST ⚠️
//
// • This file is regenerated on every build
// • Changes made here are NON-PERSISTENT
// • Manual modifications will be OVERWRITTEN
```

### **JAVASCRIPT ASSET CONSOLIDATION**
- **SINGLE SOURCE OF TRUTH**: All JavaScript files consolidated to `/opt/holodeck-one/share/htdocs/static/js/`
- **DUPLICATE ELIMINATION**: Removed redundant `src/renderer/static/js/` directory
- **ENHANCED RENDERER**: Latest version includes 25×25×25 grid system capabilities
- **PATH STANDARDIZATION**: All references updated to canonical location

### **PROJECT STRUCTURE CLEANUP**
- **LEGACY REMOVAL**: Eliminated old HD1 artifacts and backup files
- **PATH MIGRATION**: All hardcoded paths updated from `/home/claude-3/3dv` to `/opt/holodeck-one`
- **PROFESSIONAL ORGANIZATION**: Clean separation of concerns (src/, share/, build/, docs/)
- **WORKSPACE HYGIENE**: Implemented daily cleanup procedures and violation tracking

### **DOCUMENTATION EXCELLENCE**
- **COMPREHENSIVE ARCHITECTURE**: Complete system flow documentation in `docs/architecture/`
- **PROFESSIONAL STANDARDS**: Definitive development law established
- **GIT HYGIENE**: Mandatory commit standards and workspace policies
- **SECURITY COMPLIANCE**: Input validation and boundary enforcement documented

---

## ⚖️ **VIOLATION CONSEQUENCES**

### **IMMEDIATE ACTIONS:**
- **Code Rejection** - Non-compliant code rejected
- **Build Failure** - Specification violations block builds
- **Deployment Block** - Standard standards not met
- **Rollback Required** - Breaking changes reverted

### **ZERO TOLERANCE:**
- Hardcoded paths
- Emoji in system output
- Manual routing
- Unvalidated coordinates
- Insecure practices

---

## 🎖️ **EXCELLENCE STANDARDS**

HD1 represents the evolution from innovative concept to **standard-grade engineering solution**. We maintain advanced capabilities while implementing proper software engineering practices.

### **CORE PHILOSOPHY:**
> **"Standard engineering excellence with zero compromise on innovation"**

### **SUCCESS METRICS:**
- ✅ 100% API-driven development
- ✅ Zero regressions ever
- ✅ Standard daemon control
- ✅ Specification-driven architecture
- ✅ Thread-safe session management
- ✅ Real-time 3D visualization

---

## 📞 **SUPPORT & ESCALATION**

### **DEVELOPMENT ISSUES:**
1. Check `make status` for daemon health
2. Review `/opt/holodeck-one/build/logs/` for errors
3. Validate API specification syntax
4. Verify all handlers implemented
5. Rebuild with `make all`

### **EMERGENCY CONTACTS:**
- **Build Issues:** Check ADR documents
- **API Problems:** Validate against specification
- **Performance Issues:** Monitor resource usage
- **Security Concerns:** Review security checklist

---

**REMEMBER: These standards are non-negotiable. They ensure HD1 remains a standard, reliable, and innovative 3D visualization platform.**

---

*Last Updated: 2025-06-28*  
*HD1 Version: 2.0.0*  
*Authority: Standard Development Standards*