# envdiff

> Diff `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-env> <compare-env> [additional-envs...]
```

### Example

```bash
envdiff .env.development .env.production
```

**Sample output:**

```
MISSING in .env.production:
  - DEBUG
  - LOG_LEVEL

MISMATCHED values:
  - DATABASE_URL  (.env.development) postgres://localhost/dev
                  (.env.production)  postgres://prod-host/app

OK: 12 keys match across both files.
```

### Flags

| Flag | Description |
|------|-------------|
| `--ignore-values` | Only check for missing keys, skip value comparison |
| `--strict` | Exit with non-zero code if any differences are found |
| `--format json` | Output results as JSON |

---

## Why envdiff?

Managing environment variables across staging, production, and local setups is error-prone. `envdiff` makes it easy to catch missing keys before they cause runtime failures.

---

## License

MIT © [yourusername](https://github.com/yourusername)