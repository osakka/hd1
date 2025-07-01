# Contributing to HD1

## Development Principles
- **API-first development** from OpenAPI specification
- **Quality solutions only** - no regressions ever
- **Single source of truth** - no parallel implementations
- **Specification-driven development** - all changes start with api.yaml

## Development Setup
```bash
cd src
make validate    # Validate OpenAPI specification
make generate    # Generate routing from spec
make build       # Build HD1 daemon
make test        # Run tests
```

## Making Changes
1. Update `src/api.yaml` for API changes
2. Run `make generate` to update code
3. Implement handlers in `src/api/` directories
4. Update documentation
5. Test thoroughly

## Code Standards
- Use absolute paths only
- Long flags only for clarity
- Structured logging with timestamps
- No emojis in system code
- Thread-safe implementations

## Commit Guidelines
- Clear, descriptive commit messages
- Reference relevant ADRs
- Include test results
- No breaking changes without ADR

## Architecture Decision Records
All significant architecture changes require an ADR in `docs/adr/` following the established format and numbering scheme.