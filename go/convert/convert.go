package convert

import (
	"fmt"
	evt "github.com/saichler/l8events/go/types/l8events"
	"google.golang.org/protobuf/proto"
)

// Parser converts an EventRecord into a typed category-specific protobuf message.
type Parser interface {
	Parse(record *evt.EventRecord) (proto.Message, error)
}

// Converter dispatches EventRecord instances to category-specific parsers.
type Converter struct {
	parsers map[evt.EventCategory]Parser
}

// New creates a Converter pre-loaded with all 16 built-in category parsers.
func New() *Converter {
	c := &Converter{
		parsers: make(map[evt.EventCategory]Parser),
	}
	registerBuiltins(c)
	return c
}

// Register adds or replaces a parser for the given category.
func (c *Converter) Register(category evt.EventCategory, parser Parser) {
	c.parsers[category] = parser
}

// Convert dispatches the record to the appropriate category parser.
// Returns an error for nil records, UNSPECIFIED category, or unregistered categories.
// Returns (nil, nil) for CUSTOM category (no struct for custom events).
func (c *Converter) Convert(record *evt.EventRecord) (proto.Message, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}
	if record.Category == evt.EventCategory_EVENT_CATEGORY_UNSPECIFIED {
		return nil, fmt.Errorf("record has UNSPECIFIED category")
	}
	if record.Category == evt.EventCategory_EVENT_CATEGORY_CUSTOM {
		return nil, nil
	}
	parser, ok := c.parsers[record.Category]
	if !ok {
		return nil, fmt.Errorf("no parser registered for category %s", record.Category.String())
	}
	return parser.Parse(record)
}
