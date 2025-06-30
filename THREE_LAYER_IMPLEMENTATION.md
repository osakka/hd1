# HD1 Three-Layer Architecture Implementation Plan

**Status**: 🚀 **ACTIVE IMPLEMENTATION**  
**Started**: 2025-06-30  
**Target**: Complete three-layer separation (Environment + Props + Scenes)  

## 🎯 Project Overview

Converting HD1 from monolithic scenes to industry-standard three-layer architecture:
- **Environment Layer**: Physical laws and world settings
- **Prop Layer**: Reusable object collections with cohesion physics  
- **Scene Layer**: Orchestration combining environment + prop placement

**Industry Validation**: Matches Unity (World Settings + Prefabs + Scenes) and Unreal (World Settings + Blueprints + Levels) patterns but with cleaner API separation.

---

## 📋 Implementation Phases

### 🔧 Phase 1: Environment System Foundation
**Goal**: Establish environment API and basic functionality  
**Status**: 📋 **PLANNED**

#### 1.1 OpenAPI Specification Extension
- [ ] Add `/environments` endpoints to `src/api.yaml`
- [ ] Add `/environments/{environmentId}` POST endpoint
- [ ] Define Environment schema in OpenAPI spec
- [ ] Regenerate API handlers with `make generate`

#### 1.2 Environment API Implementation  
- [ ] Create `src/api/environments/` directory
- [ ] Implement `list.go` (GET /environments)
- [ ] Implement `apply.go` (POST /environments/{id})  
- [ ] Create `types.go` (Environment data structures)
- [ ] Add environment validation logic

#### 1.3 Environment Scripts Creation
- [ ] Create `share/environments/` directory
- [ ] Convert `earth-surface.sh` from existing patterns
- [ ] Create `molecular-scale.sh` environment
- [ ] Create `space-vacuum.sh` environment  
- [ ] Create `underwater.sh` environment
- [ ] Add environment script discovery system

#### 1.4 UI Integration
- [ ] Add environment dropdown to `hd1-console.js`
- [ ] Create environment selection handler
- [ ] Add environment status display
- [ ] Update object counter to show environment
- [ ] Test environment switching in browser

**Phase 1 Success Criteria**: Environment dropdown works, applies physics/scale to session

---

### 📦 Phase 2: Prop System Implementation  
**Goal**: Convert objects to cohesive props with physics  
**Status**: 📋 **PLANNED**

#### 2.1 Prop API Development
- [ ] Add `/props` endpoints to `src/api.yaml` 
- [ ] Add `/props/{propId}` POST endpoint
- [ ] Define Prop schema in OpenAPI spec
- [ ] Implement `src/api/props/list.go`
- [ ] Implement `src/api/props/spawn.go`
- [ ] Create prop cohesion physics system

#### 2.2 Prop Conversion
- [ ] Create `share/props/` directory
- [ ] Convert `rustic-log-chair` from scene to prop
- [ ] Update chair script for prop behavior (spawn at origin)
- [ ] Create additional prop examples
- [ ] Add prop script discovery system

#### 2.3 Cohesion Physics A-Frame Component
- [ ] Create `prop-cohesion` A-Frame component
- [ ] Implement anchor point system
- [ ] Add cohesion force calculations
- [ ] Test prop movement as unified object
- [ ] Add prop rotation around anchor

#### 2.4 Prop UI Integration  
- [ ] Add prop dropdown to UI
- [ ] Create prop spawn handler
- [ ] Add prop positioning controls
- [ ] Update object counter for prop objects
- [ ] Test prop spawning and movement

**Phase 2 Success Criteria**: Rustic log chair spawns as cohesive prop, moves as unit

---

### 🎬 Phase 3: Scene Orchestration
**Goal**: Scenes become environment + prop orchestrators  
**Status**: 📋 **PLANNED**

#### 3.1 Scene Refactoring
- [ ] Update scene API to call environment + prop APIs
- [ ] Create scene definition schema (environment + prop placements)
- [ ] Refactor existing scenes to use new pattern
- [ ] Add scene validation (props fit in environment)
- [ ] Test scene loading with three-layer system

#### 3.2 Migration Strategy
- [ ] Implement backward compatibility for old scenes
- [ ] Create migration tools for existing scenes  
- [ ] Add version detection for scene formats
- [ ] Document migration path for users
- [ ] Test mixed old/new scene operation

#### 3.3 Scene Orchestration Logic
- [ ] Implement scene clear → environment → props workflow
- [ ] Add prop placement validation
- [ ] Create scene composition tools
- [ ] Add scene export/import functionality
- [ ] Test complex multi-prop scenes

**Phase 3 Success Criteria**: Complete scenes load with environment + multiple props

---

### 🎨 Phase 4: Polish & Optimization
**Goal**: Production-ready three-layer system  
**Status**: 📋 **PLANNED**  

#### 4.1 Performance Optimization
- [ ] Optimize prop cohesion physics
- [ ] Add prop LOD (Level of Detail) system
- [ ] Implement efficient environment switching
- [ ] Add prop instance pooling
- [ ] Performance test with complex scenes

#### 4.2 Advanced Features
- [ ] Add prop animation support
- [ ] Create prop composition tools
- [ ] Add environment transition effects
- [ ] Implement prop interaction system
- [ ] Add prop physics simulation

#### 4.3 Documentation & Testing
- [ ] Complete API documentation
- [ ] Add comprehensive test suite
- [ ] Create user guides for three-layer system
- [ ] Add developer examples
- [ ] Validate industry compliance

**Phase 4 Success Criteria**: Production-ready system with full documentation

---

## 🎯 Current Sprint

### 🔥 **Phase 1.1: OpenAPI Specification Extension** 
**Active Task**: Adding environment endpoints to API specification

#### Today's Objectives
- [ ] Analyze current `src/api.yaml` structure
- [ ] Add `/environments` and `/environments/{environmentId}` endpoints
- [ ] Define Environment data schema
- [ ] Test `make generate` with new endpoints
- [ ] Verify auto-generated handlers created

#### Working Session Notes
*Session started: 2025-06-30*

**Research Completed**:
- ✅ Analyzed game engine architecture patterns (Unity, Unreal, Three.js)
- ✅ Validated three-layer approach against industry standards  
- ✅ Confirmed architecture is bar-raising and simpler than alternatives
- ✅ Documented complete technical specification

**Ready to Begin**: OpenAPI specification extension for environment system

---

## 📊 Progress Tracking

### Overall Progress
```
Phase 1: Environment System     [⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜] 0%
Phase 2: Prop System           [⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜] 0%  
Phase 3: Scene Orchestration   [⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜] 0%
Phase 4: Polish & Optimization [⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜] 0%

Total Implementation: [⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜] 0%
```

### Milestones
- [ ] **M1**: Environment dropdown functional
- [ ] **M2**: First prop (rustic chair) working with cohesion  
- [ ] **M3**: First orchestrated scene loads environment + props
- [ ] **M4**: Complete system ready for production

### Key Metrics
- **API Endpoints**: 0/6 implemented
- **Environment Scripts**: 0/5 created
- **Props Converted**: 0/1 (rustic-log-chair pending)
- **Scenes Refactored**: 0/12 planned
- **Tests Passing**: TBD
- **Documentation**: Architecture complete ✅

---

## 🔧 Technical Decisions

### Architecture Choices
- ✅ **Three-layer separation**: Environment + Props + Scenes
- ✅ **OpenAPI-driven development**: Single source of truth
- ✅ **A-Frame integration**: Prop cohesion via components
- ✅ **Backward compatibility**: Gradual migration strategy

### Implementation Decisions  
- ✅ **Environment at session level**: One environment per session
- ✅ **Props spawn at origin**: Consistent placement behavior
- ✅ **Scene orchestration**: Clear → Environment → Props workflow
- ✅ **Progressive enhancement**: Old scenes work during transition

---

## 🚀 Next Actions

### Immediate (Today)
1. **Start Phase 1.1**: Add environment endpoints to `src/api.yaml`
2. **Test API generation**: Verify `make generate` creates handlers
3. **Create environment directory**: Set up `share/environments/`

### This Week
1. Complete Phase 1.1-1.2: Environment API foundation
2. Create first environment script (`earth-surface.sh`)
3. Add environment dropdown to UI
4. Test environment switching

### This Sprint (Next 2 weeks)
1. Complete Phase 1: Full environment system
2. Begin Phase 2: Start prop API development
3. Convert rustic-log-chair to prop
4. Validate architecture with working examples

---

## 📝 Session Log

### 2025-06-30 Session
**Objective**: Architecture planning and Phase 1 initiation

**Completed**:
- ✅ Researched game engine architecture patterns
- ✅ Validated three-layer approach against Unity/Unreal standards
- ✅ Created comprehensive technical documentation
- ✅ Planned complete implementation strategy
- ✅ Set up progress tracking system

**Insights**:
- Three-layer architecture is industry-validated and bar-raising
- HD1's approach is simpler than Unity/Unreal due to clean API separation
- OpenAPI-driven development maintains single source of truth
- Backward compatibility ensures zero regressions

**Next Session**: Begin OpenAPI specification extension for environment endpoints

---

## 🔗 Related Documents

- [Architecture Decision Records](docs/adr/README.md)
- [HD1 Development Context](CLAUDE.md)
- [API Specification](src/api.yaml)
- [Build System](Makefile)

---

**Last Updated**: 2025-06-30  
**Next Update**: After Phase 1.1 completion  
**Maintained By**: Claude + Development Team