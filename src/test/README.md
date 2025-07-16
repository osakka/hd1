# HD1 (Holodeck One) Test Suite

This directory contains the comprehensive test suite for HD1 v7.0.0, the universal 3D interface platform.

## Overview

The test suite is organized into phases that correspond to the platform's development phases:

1. **Phase 1: Foundation** - Core authentication, sessions, and service registry
2. **Phase 2: Collaboration** - WebRTC, operational transforms, and asset management  
3. **Phase 3: AI Integration** - LLM avatars and content generation
4. **Phase 4: Universal Platform** - Cross-platform clients, plugins, and enterprise features
5. **Integration Tests** - End-to-end platform testing

## Test Files

### Core Test Files

- `phase1_foundation_test.go` - Foundation component tests
- `phase2_collaboration_test.go` - Collaboration feature tests  
- `phase3_ai_integration_test.go` - AI integration tests
- `phase4_universal_platform_test.go` - Universal platform and enterprise tests
- `integration_test.go` - Full platform integration tests

### Supporting Files

- `fixtures.go` - Test data and fixtures
- `run_tests.sh` - Test runner script
- `README.md` - This documentation file

## Quick Start

### Prerequisites

- Go 1.19 or later
- PostgreSQL 13 or later
- Make sure HD1 dependencies are installed

### Running Tests

```bash
# Run all tests
./run_tests.sh

# Run specific phase
./run_tests.sh phase1
./run_tests.sh phase2
./run_tests.sh phase3
./run_tests.sh phase4
./run_tests.sh integration

# Run performance tests
./run_tests.sh benchmarks
./run_tests.sh race
./run_tests.sh memory
```

## Test Structure

### Phase 1: Foundation Tests

Tests the core platform foundation including:

- **Authentication System**
  - User registration and login
  - JWT token validation
  - Refresh token handling
  - Password security

- **Session Management**
  - Session creation and lifecycle
  - User joining and leaving
  - Session permissions
  - Session discovery

- **Service Registry**
  - Service registration and discovery
  - Health monitoring
  - Capability management
  - Service lifecycle

- **Database Operations**
  - Connection management
  - Transaction support
  - Schema migrations
  - Data integrity

**Test Coverage**: Authentication flows, session operations, service management, database functionality

### Phase 2: Collaboration Tests

Tests real-time collaborative features including:

- **WebRTC Integration**
  - Session creation and management
  - Offer/answer signaling
  - ICE candidate exchange
  - Connection establishment

- **Operational Transforms**
  - Document creation and management
  - Concurrent operation handling
  - Conflict resolution
  - Version control

- **Asset Management**
  - File upload and storage
  - Asset processing and optimization
  - Version control
  - Usage tracking

- **WebSocket Synchronization**
  - Real-time message broadcasting
  - Client connection management
  - Message ordering
  - Connection recovery

**Test Coverage**: Real-time collaboration, asset handling, WebSocket communication, conflict resolution

### Phase 3: AI Integration Tests

Tests AI-powered features including:

- **LLM Avatars**
  - Avatar creation and management
  - Multi-provider support (OpenAI, Claude, Gemini)
  - Personality configuration
  - Avatar interaction

- **Content Generation**
  - Scene generation
  - Object creation
  - Animation sequences
  - Template-based generation

- **AI Interaction**
  - Natural language processing
  - Scene understanding
  - Context awareness
  - Multi-modal communication

- **Usage Tracking**
  - Token consumption monitoring
  - Cost calculation
  - Performance metrics
  - Rate limiting

**Test Coverage**: AI avatar functionality, content generation, natural language processing, usage monitoring

### Phase 4: Universal Platform Tests

Tests universal platform and enterprise features including:

- **Client Management**
  - Multi-platform client registration
  - Capability management
  - Heartbeat monitoring
  - Message broadcasting

- **Plugin System**
  - Plugin installation and management
  - Hook system
  - Configuration management
  - Plugin lifecycle

- **Enterprise Features**
  - Organization management
  - Role-based access control (RBAC)
  - Analytics and reporting
  - Security and compliance

- **Cross-Platform Integration**
  - Web client support
  - Mobile client support
  - Desktop client support
  - VR client support

**Test Coverage**: Client management, plugin system, enterprise features, cross-platform support

### Integration Tests

Tests end-to-end platform functionality including:

- **Complete Platform Workflow**
  - Full user journey from registration to collaboration
  - Multi-component integration
  - Real-time collaboration scenarios
  - AI-powered content creation

- **Multi-User Scenarios**
  - Concurrent user sessions
  - Collaborative editing
  - Asset sharing
  - Real-time synchronization

- **Enterprise Integration**
  - Role-based workflows
  - Analytics tracking
  - Security compliance
  - Multi-organization support

- **Performance Testing**
  - Load testing
  - Stress testing
  - Scalability testing
  - Resource usage monitoring

**Test Coverage**: End-to-end workflows, multi-user scenarios, enterprise integration, performance validation

## Test Data and Fixtures

The test suite uses comprehensive fixtures defined in `fixtures.go`:

### User Fixtures
- Admin, Designer, Developer, Manager, Viewer roles
- Complete user profiles with authentication data
- Linked to organizations and sessions

### Organization Fixtures
- Enterprise, Professional, Education tiers
- Complete organization settings
- Linked to users and roles

### Session Fixtures
- Various session types (collaboration, presentation, training)
- Different visibility levels
- Linked to users and organizations

### Service Fixtures
- Renderer, Physics, Audio, AI services
- Complete capability definitions
- Service health monitoring

### Plugin Fixtures
- Physics, Animation, Lighting plugins
- Configuration and hook definitions
- Plugin lifecycle states

### Asset Fixtures
- Models, Textures, Audio, Environment assets
- Complete metadata and relationships
- Usage tracking data

### Other Fixtures
- Documents, Avatars, Clients, Roles
- Comprehensive test data coverage
- Realistic relationships and constraints

## Test Configuration

### Environment Variables

```bash
# Test mode configuration
HD1_TEST_MODE=true
HD1_TEST_DB_URL=postgres://test:test@localhost:5432/hd1_test?sslmode=disable
HD1_LOG_LEVEL=DEBUG
HD1_STATIC_DIR=/opt/hd1/share/htdocs/static

# Coverage configuration
COVERAGE_THRESHOLD=80
TIMEOUT=300s
PARALLEL_TESTS=4
```

### Database Setup

Tests use a dedicated PostgreSQL test database:

```bash
# Create test database
createdb hd1_test

# Run migrations
psql -d hd1_test -f database/migrations/schema.sql
```

## Reports and Coverage

### Test Reports

Test execution generates comprehensive reports:

- **JUnit XML Reports** - Machine-readable test results
- **HTML Coverage Reports** - Visual coverage analysis
- **Benchmark Reports** - Performance metrics
- **Race Detection Reports** - Concurrency safety
- **Memory Profile Reports** - Memory usage analysis

### Coverage Analysis

The test suite tracks code coverage across all phases:

- **Phase Coverage** - Individual phase coverage
- **Combined Coverage** - Overall platform coverage
- **Function Coverage** - Function-level analysis
- **Line Coverage** - Line-level analysis

### Report Locations

```
build/
├── coverage/
│   ├── phase1_coverage.out
│   ├── phase2_coverage.out
│   ├── phase3_coverage.out
│   ├── phase4_coverage.out
│   ├── integration_coverage.out
│   └── combined_coverage.out
└── reports/
    ├── test_report.html
    ├── phase1_coverage.html
    ├── phase2_coverage.html
    ├── phase3_coverage.html
    ├── phase4_coverage.html
    ├── integration_coverage.html
    ├── combined_coverage.html
    ├── benchmarks.txt
    ├── race_test.log
    ├── memory_test.log
    └── test_run.log
```

## Continuous Integration

### CI/CD Integration

The test suite integrates with CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
name: HD1 Test Suite
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.19
      - name: Run tests
        run: ./src/test/run_tests.sh
      - name: Upload coverage
        uses: codecov/codecov-action@v1
        with:
          file: ./build/coverage/combined_coverage.out
```

### Test Automation

Automated test execution includes:

- **Pre-commit hooks** - Run tests before commits
- **Pull request validation** - Validate changes
- **Nightly builds** - Full test suite execution
- **Performance monitoring** - Track performance trends

## Best Practices

### Writing Tests

1. **Test Organization**
   - Group related tests in sub-tests
   - Use descriptive test names
   - Follow table-driven test patterns

2. **Test Data**
   - Use fixtures for consistent test data
   - Clean up resources after tests
   - Avoid test interdependencies

3. **Assertions**
   - Use testify/assert for readable assertions
   - Check both success and error conditions
   - Validate all relevant response fields

4. **Performance**
   - Use parallel tests where appropriate
   - Set reasonable timeouts
   - Monitor resource usage

### Test Maintenance

1. **Regular Updates**
   - Keep tests synchronized with code changes
   - Update fixtures for new features
   - Maintain test documentation

2. **Coverage Monitoring**
   - Aim for 80%+ code coverage
   - Focus on critical path coverage
   - Monitor coverage trends

3. **Performance Monitoring**
   - Track test execution time
   - Monitor memory usage
   - Identify slow tests

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   ```bash
   # Check PostgreSQL status
   systemctl status postgresql
   
   # Verify test database exists
   psql -l | grep hd1_test
   
   # Reset test database
   dropdb hd1_test && createdb hd1_test
   ```

2. **Test Timeouts**
   ```bash
   # Increase timeout
   export TIMEOUT=600s
   
   # Run tests with verbose output
   ./run_tests.sh -v
   ```

3. **Coverage Issues**
   ```bash
   # Clean coverage data
   rm -rf build/coverage/*
   
   # Run with fresh coverage
   ./run_tests.sh
   ```

### Debug Mode

Enable detailed logging:

```bash
# Set debug environment
export HD1_LOG_LEVEL=DEBUG
export HD1_TRACE_MODULES=websocket,sync,threejs

# Run with debug output
./run_tests.sh
```

## Contributing

### Adding New Tests

1. **Create test file** in appropriate phase directory
2. **Use existing patterns** from other test files
3. **Add fixtures** for new test data
4. **Update documentation** with new test coverage

### Test Guidelines

1. **Test Naming** - Use descriptive names that explain what is being tested
2. **Test Structure** - Follow Arrange-Act-Assert pattern
3. **Error Handling** - Test both success and failure scenarios
4. **Resource Cleanup** - Always clean up resources after tests

### Review Process

1. **Code Review** - All test changes require code review
2. **Coverage Check** - Ensure new code maintains coverage levels
3. **Performance Impact** - Verify tests don't impact performance
4. **Documentation** - Update documentation for new tests

## Support

For test-related questions or issues:

1. **Documentation** - Check this README and code comments
2. **Issue Tracking** - Create GitHub issues for bugs
3. **Development Team** - Contact the HD1 development team
4. **Community** - Join the HD1 community discussions

---

*HD1 v7.0.0 Test Suite - Comprehensive testing for the universal 3D interface platform*