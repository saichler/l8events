# Plan: Event Parsing & Conversion Engine

## Context

l8events has 16 category-specific protobuf event structs (AuditEvent, SystemEvent, ... AutomationEvent) and a generic `EventRecord` that carries raw data in its `Attributes` map. Consumer projects receive raw events as `EventRecord` instances and need to convert them into the appropriate typed category struct. This plan adds a `convert` package that dispatches on `EventRecord.Category` and populates the typed struct from the record's fields and Attributes map.

## Package: `go/convert/`

**Architecture**: Constructor injection (matching `archive.go` pattern). `Converter` struct holds a parser registry. `New()` pre-loads all 16 built-in parsers. Consumers can register custom parsers via `Register()`.

**Return type**: `proto.Message` — all category structs are generated protobuf types.

**Error strategy**: Lenient — missing attributes yield zero values (no error). Malformed numeric/bool strings return an error.

## API

```go
type Parser interface {
    Parse(record *evt.EventRecord) (proto.Message, error)
}

type Converter struct { parsers map[evt.EventCategory]Parser }

func New() *Converter                                                    // Pre-loaded with 16 built-in parsers
func (c *Converter) Register(category evt.EventCategory, parser Parser)  // Add/replace parser
func (c *Converter) Convert(record *evt.EventRecord) (proto.Message, error)
```

**Convert() behavior**:
- `nil` record → error
- UNSPECIFIED category → error
- CUSTOM category → `(nil, nil)` (no struct for custom)
- Unregistered category → error
- Otherwise → delegates to the registered parser

## Field Mapping

**Common fields (all 16 parsers)**:
| Target Field | Source |
|---|---|
| `EventId` | `record.EventId` |
| `PropertyId` | `record.Attributes["propertyId"]` |
| `SourceId` | `record.SourceId` |
| `SourceType` | `record.SourceType` |
| `Message` | `record.Message` |

**SubCategory** (15 of 16 — NOT SyslogEvent): `int32(record.Attributes["subCategory"])` cast to the category's enum type.

**Domain fields**: All from `record.Attributes[camelCaseFieldName]` with type conversion helpers.

**Special cases**:
- `TrapEvent.Varbinds`: `extractMap(attrs, "varbinds.")` — keys with prefix `varbinds.` become entries in the map
- `SyslogEvent`: No SubCategory field; uses `Facility` and `SyslogSeverity` (int32) instead
- `PerformanceEvent.ThresholdType`: Second enum field besides SubCategory
- `AutomationEvent.Success`: bool field

## Files

| File | ~Lines | Content |
|------|--------|---------|
| `go/convert/convert.go` | 80 | Converter struct, Parser interface, New(), Convert(), Register() |
| `go/convert/helpers.go` | 70 | str(), i32(), i64(), f64(), boolean(), setCommon(), subCategory(), extractMap() |
| `go/convert/parsers_infra.go` | 180 | 7 parsers: Network, Kubernetes, Compute, Storage, Power, GPU, Topology |
| `go/convert/parsers_ops.go` | 180 | 9 parsers: Audit, System, Monitoring, Security, Integration, Performance, Syslog, Trap, Automation |
| `go/convert/builtins.go` | 30 | registerBuiltins() wiring |
| `go/convert/convert_test.go` | 450 | Full test coverage (split if >450) |

All source files well under 500 lines.

## Implementation Phases

### Phase 1: Helpers (`helpers.go`)
Type conversion utilities: `str`, `i32`, `i64`, `f64`, `boolean`, `setCommon`, `subCategory`, `extractMap`.

### Phase 2: Engine core (`convert.go`)
Converter struct, Parser interface, New(), Convert(), Register().

### Phase 3: Parsers (`parsers_ops.go` + `parsers_infra.go`)
16 parser structs implementing the Parser interface. Can be written in parallel.

### Phase 4: Wiring (`builtins.go`)
`registerBuiltins()` maps each EventCategory to its parser instance.

### Phase 5: Tests (`convert_test.go`)
- Nil/UNSPECIFIED/CUSTOM/unregistered error cases
- One test per category (16 tests): common fields, domain fields, type conversions, enum SubCategory
- Edge cases: bad numeric strings, missing attributes, TrapEvent varbinds, SyslogEvent no-SubCategory
- Custom parser registration

### Phase 6: Update `test.sh`
Add `./convert/...` to the `-coverpkg` flag.

## Verification

1. `go build ./convert/...` — compiles
2. `go vet ./convert/...` — clean
3. `go test -v ./convert/...` — all pass
4. `wc -l go/convert/*.go` — all under 500 lines
5. Each of 16 parsers exercised by at least one test
