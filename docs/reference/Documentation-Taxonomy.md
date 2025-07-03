# HD1 Documentation Taxonomy

**Documentation Classification System for HD1 v5.0.1**

## 📂 Primary Categories

### 1. **Root Documentation** (`/opt/hd1/`)
- `README.md` - Project overview and quick start
- `CLAUDE.md` - Development context and configuration system
- `CHANGELOG.md` - Version history and changes

### 2. **User Documentation** (`docs/`)
- **Getting Started** (`getting-started/`) - Installation, quick start, tutorials
- **User Guides** (`user-guides/`) - End-user functionality and workflows
- **API Reference** (`reference/`) - Complete API documentation and schemas

### 3. **Developer Documentation** (`docs/`)
- **Developer Guides** (`developer-guide/`) - Development procedures and standards
- **Architecture** (`architecture/`) - System design and technical architecture
- **Decisions** (`decisions/`) - Architectural Decision Records (ADRs)

### 4. **Operations Documentation** (`docs/`)
- **Operations** (`operations/`) - Deployment, monitoring, administration
- **Troubleshooting** (`troubleshooting/`) - Common issues and solutions

### 5. **Source Documentation** (`src/`, `lib/`, `share/`)
- Technical implementation documentation co-located with code
- Auto-generated documentation from code comments
- Template and configuration documentation

## 🎯 Documentation Standards

### **File Naming Convention**
- Use kebab-case for directories: `getting-started/`, `user-guides/`
- Use sentence case for files: `README.md`, `Quick-Start.md`
- Use descriptive names: `API-Specification.md` not `api.md`

### **Content Classification**
1. **Factual Only** - 100% accurate to codebase
2. **Audience-Specific** - Clear user vs developer vs operations content
3. **Single Source of Truth** - No duplication across categories
4. **Version Consistent** - All references to v5.0.1

### **Navigation Structure**
- Each category has its own `README.md` as navigation hub
- Cross-references use relative paths
- Clear breadcrumb navigation in headers

## 📊 Current vs Target Structure

### **Current Issues**
- Mixed audience content in single files
- Inconsistent version references (v5.0.0 vs v5.0.1)
- API endpoint count inconsistencies (59 vs 82)
- Duplication between files
- No clear taxonomy or navigation

### **Target Structure**
```
docs/
├── README.md                    # Master navigation hub
├── getting-started/
│   ├── README.md               # Getting started guide
│   ├── Installation.md         # Installation procedures
│   └── Quick-Start.md          # Quick start tutorial
├── user-guides/
│   ├── README.md               # User guide navigation
│   ├── API-Usage.md            # API usage examples
│   └── WebSocket-Events.md     # Real-time events guide
├── developer-guide/
│   ├── README.md               # Developer navigation
│   ├── Contributing.md         # Contribution guidelines
│   ├── Code-Standards.md       # Coding standards
│   └── Build-System.md         # Build and development
├── architecture/
│   ├── README.md               # Architecture overview
│   ├── Overview.md             # System architecture
│   ├── Design-Principles.md    # Core principles
│   └── API-Design.md           # API design patterns
├── decisions/
│   ├── README.md               # ADR navigation
│   └── adr/                    # Individual ADRs
├── reference/
│   ├── README.md               # Reference navigation
│   ├── API-Specification.md    # Complete API reference
│   └── Configuration.md        # Configuration reference
├── operations/
│   ├── README.md               # Operations navigation
│   ├── Deployment.md           # Deployment procedures
│   └── Monitoring.md           # Monitoring and logging
└── troubleshooting/
    ├── README.md               # Troubleshooting navigation
    └── Common-Issues.md        # Common problems and solutions
```

## 🔄 Migration Strategy

### **Phase 1: Critical Fixes** ✅
- [x] Fix API endpoint count inconsistencies (82 total routes)
- [x] Standardize version references (v5.0.1)
- [x] Create taxonomy documentation

### **Phase 2: Restructure** 
- [ ] Create category directories with proper naming
- [ ] Move existing content to appropriate categories
- [ ] Create navigation README.md files
- [ ] Eliminate duplication

### **Phase 3: Content Audit**
- [ ] Verify factual accuracy against codebase
- [ ] Update CHANGELOG completeness
- [ ] Audit ADR consistency
- [ ] Create missing documentation

### **Phase 4: Single Source of Truth**
- [ ] Establish authoritative sources for each topic
- [ ] Remove duplicate content
- [ ] Create cross-references where needed
- [ ] Final validation

## 📋 Quality Criteria

### **Mandatory Requirements**
1. **100% Factual** - All claims verified against codebase
2. **No Exaggerations** - Conservative, accurate language only
3. **Clear and Crisp** - Concise, professional writing
4. **Consistent** - Uniform style and terminology
5. **Single Source of Truth** - No conflicting information
6. **Bar-Raising Solutions** - Industry-standard organization

### **Content Validation**
- API endpoint counts: 82 total routes (59 unique paths)
- Version references: v5.0.1 throughout
- Configuration system: Complete with environment variables
- Architecture: PlayCanvas-based API-first game engine
- Performance claims: Verified against actual implementation

---

**HD1 v5.0.1** - API-First Game Engine Platform  
**Documentation Taxonomy** - Industry-standard organization with single source of truth