# gx - Modern Go Script Runner

Run Go files like scripts with smart caching and auto-reload.

## Performance

Benchmarks show significant speedup over `go run`:

- **Local scripts:** 100ms vs go run's 275ms → **2.5x faster**
- **Scripts with dependencies:** 571ms vs go run's 1,759ms → **3x faster**

## Features

-  Smart caching - compile once, run instantly
-  Auto-reload - watch mode for development  
-  Cache management
-  Proper exit codes

## Commands
```bash
gx run <script> [args]    # Run a Go script (with caching)
gx watch <script> [args]  # Auto-reload on file changes
gx clean                  # Clear cache
gx version                # Show version
```

## Installation
```bash
go install github.com/slyt3/gx@latest
```

## Examples

### Run a script
```bash
gx run server.go --port 8080
```

### Development with auto-reload
```bash
gx watch server.go --port 8080
# Edit server.go and save → automatically restarts!
```

### Clear cache
```bash
gx clean
```

## How it works

1. **First run:** Compiles and caches the binary
2. **Subsequent runs:** Uses cached binary (validated by file size + modTime)
3. **Watch mode:** Monitors file changes and auto-restarts

## Why gx?

- Faster than `go run` for repeated executions
- No manual recompilation during development
- Works great for scripts, servers, and CLI tools

## License

MIT
