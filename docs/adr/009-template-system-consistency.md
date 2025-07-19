# ADR-009: Template System Consistency

**Status**: Accepted  
**Date**: 2025-07-19  
**Authors**: HD1 Development Team  

## Context

HD1's auto-generation system suffered from template regression where architectural decisions implemented in source code were being overwritten during the build process. The JavaScript API client template contained hardcoded `client_id` references that violated the unified HD1 ID system established in ADR-008.

### Problems Identified

1. **Template Regression**: Auto-generation reverting unified `hd1_id` back to `client_id`
2. **Architectural Violation**: Templates not reflecting architectural decisions
3. **Inconsistent Builds**: Manual fixes being overwritten by code generation
4. **Source Confusion**: Unclear what was authoritative - source or templates

### Root Cause Analysis

```javascript
// Template hardcoded old system
constructor(baseURL = '/api', clientId = null) {
    this.clientId = clientId;  // WRONG - violates ADR-008
}
```

This created a cycle where:
1. Developer implements unified `hd1_id` system
2. Build process regenerates code from templates  
3. Templates contain old `client_id` references
4. Generated code violates architectural decisions

## Decision

We implement **Template System Consistency** ensuring auto-generation maintains architectural purity.

### Core Principles

1. **Template Authority**: Templates must reflect current architectural decisions
2. **Regression Prevention**: Auto-generation cannot violate established ADRs
3. **Consistency Validation**: Build process validates architectural compliance
4. **Single Source Truth**: Templates serve architectural decisions, not legacy patterns

### Implementation Strategy

#### Template Updates
```javascript
// Fixed template reflects ADR-008
constructor(baseURL = '/api', hd1Id = null) {
    this.hd1Id = hd1Id;  // Unified HD1 ID system
}

setHd1Id(hd1Id) {
    this.hd1Id = hd1Id;  // Consistent naming
}

headers: {
    'Content-Type': 'application/json',
    'X-Client-ID': this.hd1Id  // Header name preserved for server compatibility
}
```

#### Template Validation Process
1. **Pre-Build Verification**: Templates checked against current ADRs
2. **Post-Build Validation**: Generated code validated for architectural compliance
3. **Regression Detection**: Automated detection of ADR violations
4. **Clean Build Guarantee**: Zero warnings policy maintained

### Breaking Changes

- **Template Source**: JavaScript client template updated to reflect unified ID system
- **Build Consistency**: Auto-generation now maintains architectural decisions
- **Validation Requirements**: Build process includes architectural compliance checks

## Implementation Timeline

- **2025-07-19 10:06**: Template system fixes preventing `client_id` regression
- **2025-07-19 10:42**: Surgical precision template updates with clean build verification
- **2025-07-19 11:44**: Complete template system consistency with validation

## Architectural Impact

This ADR closes the architectural loop by ensuring that:

1. **Decisions Persist**: ADR implementations survive build processes
2. **Templates Reflect Reality**: Auto-generation aligned with architectural decisions
3. **No Silent Regressions**: Build process cannot violate established patterns
4. **Maintainable Architecture**: Templates serve architectural vision consistently

## Consequences

### Positive

- **Architectural Integrity**: Templates maintain established architectural decisions
- **Build Reliability**: Auto-generation preserves manual architectural work
- **Developer Confidence**: Changes persist through build cycles
- **Documentation Accuracy**: Generated code matches architectural documentation

### Negative

- **Template Maintenance**: Templates require updates when architectural decisions change
- **Build Complexity**: Additional validation steps in build process

## Validation

Template system consistency verified through:
- **Clean Build**: Zero warnings with unified ID system
- **Protocol Testing**: Client-server unified `hd1_id` message exchange
- **Regression Prevention**: Automated detection of architectural violations

## Compliance

This ADR ensures HD1's auto-generation system serves architectural excellence rather than undermining it, establishing a foundation for maintainable template-driven development that respects architectural decisions.