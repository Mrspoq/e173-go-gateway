# Contributing to E173 Go Gateway

This guide is designed for multiple AI agents and developers working on the E173 Go Gateway project.

## Multi-Agent Coordination

### Agent Specialization Areas

1. **Frontend Agent**: Templates, HTMX, Tailwind CSS
2. **Backend Agent**: API endpoints, business logic
3. **Database Agent**: Schema, migrations, repositories
4. **Integration Agent**: AMI, OpenSIPS, external services

### Communication Protocol

- Always update relevant documentation when making changes
- Leave clear comments in code for other agents
- Test your changes before committing
- Update CODEBASE_OVERVIEW.md for major architectural changes

## Development Workflow

### 1. Before Starting

```bash
# Always pull latest changes
git pull origin develop

# Check current status
make test
make run
```

### 2. Making Changes

#### Frontend Changes
- Templates are in `templates/`
- All templates must include `{{template "base" .}}`
- Use HTMX attributes for dynamic updates
- Follow existing Tailwind patterns

#### Backend Changes
- Add new handlers in `internal/handlers/`
- Implement repository pattern for data access
- Update routes in `cmd/server/main.go`
- Add tests for new functionality

#### Database Changes
- Create migration files in `migrations/`
- Never modify existing migrations
- Test migrations up and down
- Update models in `internal/database/models/`

### 3. Code Standards

#### Go Code
```go
// Good: Clear function names and error handling
func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
    // Implementation
}

// Bad: Unclear names and no error handling
func (r *CustomerRepository) Get(id int64) *models.Customer {
    // Implementation
}
```

#### Templates
```html
<!-- Good: Using HTMX for dynamic updates -->
<div hx-get="/api/stats" hx-trigger="every 5s" hx-swap="innerHTML">
    <!-- Content -->
</div>

<!-- Bad: Inline JavaScript -->
<div onclick="updateStats()">
    <!-- Content -->
</div>
```

## Testing

### Running Tests
```bash
make test              # Run all tests
go test ./pkg/...      # Test specific package
go test -v ./...       # Verbose output
```

### Writing Tests
- Use table-driven tests
- Mock external dependencies
- Test error cases
- Aim for >80% coverage

## Git Workflow

### Branch Naming
- `feature/description` - New features
- `fix/description` - Bug fixes
- `refactor/description` - Code improvements
- `docs/description` - Documentation

### Commit Messages
```bash
# Good
git commit -m "feat: add customer balance tracking"
git commit -m "fix: resolve template rendering issue"
git commit -m "docs: update API documentation"

# Bad
git commit -m "updated stuff"
git commit -m "fix"
```

### Pull Request Process
1. Create feature branch from `develop`
2. Make changes and test
3. Update documentation
4. Submit PR with clear description
5. Address review comments

## Common Tasks

### Adding a New UI Section

1. Create template directory:
```bash
mkdir -p templates/newsection
```

2. Create template files:
```go
// templates/newsection/list.tmpl
{{template "base" .}}
{{define "content"}}
    <!-- Your content -->
{{end}}
```

3. Add route handler:
```go
// internal/handlers/newsection.go
func NewSectionHandler(c *gin.Context) {
    c.HTML(200, "newsection/list.tmpl", gin.H{
        "Title": "New Section",
    })
}
```

4. Register route:
```go
// cmd/server/main.go
router.GET("/newsection", handlers.NewSectionHandler)
```

### Adding a New API Endpoint

1. Define repository method:
```go
// pkg/repository/interfaces.go
type NewRepository interface {
    GetAll(ctx context.Context) ([]*models.Item, error)
}
```

2. Implement repository:
```go
// internal/repository/new_repository.go
func (r *PostgresNewRepository) GetAll(ctx context.Context) ([]*models.Item, error) {
    // Implementation
}
```

3. Create handler:
```go
// internal/handlers/api_new.go
func GetItemsHandler(repo repository.NewRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        items, err := repo.GetAll(c.Request.Context())
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, items)
    }
}
```

4. Register API route:
```go
// cmd/server/main.go
api := router.Group("/api")
api.GET("/items", handlers.GetItemsHandler(newRepo))
```

## Troubleshooting

### Template Not Rendering
- Check `{{template "base" .}}` at top of file
- Verify template name matches file path
- Check for syntax errors in template

### Database Connection Issues
- Verify `.env` credentials
- Check PostgreSQL is running
- Ensure database exists
- Run migrations: `make migrate`

### Build Failures
- Run `go mod tidy`
- Check for syntax errors: `go vet ./...`
- Clear cache: `go clean -cache`

## Documentation

### Inline Documentation
```go
// CustomerRepository handles all database operations for customers.
// It implements the repository pattern for clean separation of concerns.
type CustomerRepository struct {
    db *sqlx.DB
}

// GetByID retrieves a customer by their unique identifier.
// Returns ErrNotFound if the customer doesn't exist.
func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
    // Implementation
}
```

### API Documentation
Document all API endpoints in code:
```go
// GetStatsHandler returns current system statistics
// @Summary Get system statistics
// @Description Returns counts of modems, SIMs, calls, and spam
// @Tags stats
// @Success 200 {object} StatsResponse
// @Router /api/stats [get]
func GetStatsHandler(repo repository.StatsRepository) gin.HandlerFunc {
    // Implementation
}
```

## Resources

- [Go Style Guide](https://github.com/golang/go/wiki/CodeReviewComments)
- [HTMX Documentation](https://htmx.org/docs/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Gin Framework](https://gin-gonic.com/docs/)

## Questions?

If you're an AI agent and need clarification:
1. Check CODEBASE_OVERVIEW.md first
2. Look for similar patterns in existing code
3. Test your assumptions with small examples
4. Document your decisions for other agents
