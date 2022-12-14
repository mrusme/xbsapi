// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrusme/xbsapi/ent/bookmark"
	"github.com/mrusme/xbsapi/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	bookmarkFields := schema.Bookmark{}.Fields()
	_ = bookmarkFields
	// bookmarkDescVersion is the schema descriptor for version field.
	bookmarkDescVersion := bookmarkFields[2].Descriptor()
	// bookmark.DefaultVersion holds the default value on creation for the version field.
	bookmark.DefaultVersion = bookmarkDescVersion.Default.(string)
	// bookmark.VersionValidator is a validator for the "version" field. It is called by the builders before save.
	bookmark.VersionValidator = bookmarkDescVersion.Validators[0].(func(string) error)
	// bookmarkDescCreated is the schema descriptor for created field.
	bookmarkDescCreated := bookmarkFields[3].Descriptor()
	// bookmark.DefaultCreated holds the default value on creation for the created field.
	bookmark.DefaultCreated = bookmarkDescCreated.Default.(func() time.Time)
	// bookmarkDescLastUpdated is the schema descriptor for lastUpdated field.
	bookmarkDescLastUpdated := bookmarkFields[4].Descriptor()
	// bookmark.DefaultLastUpdated holds the default value on creation for the lastUpdated field.
	bookmark.DefaultLastUpdated = bookmarkDescLastUpdated.Default.(func() time.Time)
	// bookmark.UpdateDefaultLastUpdated holds the default value on update for the lastUpdated field.
	bookmark.UpdateDefaultLastUpdated = bookmarkDescLastUpdated.UpdateDefault.(func() time.Time)
	// bookmarkDescID is the schema descriptor for id field.
	bookmarkDescID := bookmarkFields[0].Descriptor()
	// bookmark.DefaultID holds the default value on creation for the id field.
	bookmark.DefaultID = bookmarkDescID.Default.(func() uuid.UUID)
}
