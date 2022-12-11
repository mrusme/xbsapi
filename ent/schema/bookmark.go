package schema

import (
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Bookmark holds the schema definition for the Bookmark entity.
type Bookmark struct {
	ent.Schema
}

// Fields of the Bookmark.
func (Bookmark) Fields() []ent.Field {
	// validate := validator.New()

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),

		field.String("bookmarks"),
		field.String("version").
			Default("2.0.0").
			Match(regexp.MustCompile("^[0-9]{1,2}\\.[0-9]{1,2}\\.[0-9]{1,2}$")),

		field.Time("created").
			Default(time.Now),

		field.Time("lastUpdated").
			Default(time.Now).
			UpdateDefault(time.Now),

		field.Time("deleted").
			Default(nil).
			Optional().
			Nillable(),
	}
}

// Edges of the Bookmark.
func (Bookmark) Edges() []ent.Edge {
	return nil
}
