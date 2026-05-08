package sysfields

import "time"

// SystemFields lists the standard meta-columns every Gogal document table has.
var SystemFields = []string{
	"name", "owner", "creation", "modified", "modified_by", "docstatus", "idx",
}

// ProtectedFields are system fields that must never be overwritten by a user update.
var ProtectedFields = []string{
	"name", "creation", "owner", "docstatus",
}

// IsSystemField returns true when fieldname is one of the standard system columns.
func IsSystemField(fieldname string) bool {
	for _, f := range SystemFields {
		if f == fieldname {
			return true
		}
	}
	return false
}

// IsProtectedField returns true when fieldname must not be changed via a user update.
func IsProtectedField(fieldname string) bool {
	for _, f := range ProtectedFields {
		if f == fieldname {
			return true
		}
	}
	return false
}

// InjectCreate fills all system fields for a new document. Mutates doc in place.
// Call after the document name has been resolved, before the INSERT query.
func InjectCreate(doc map[string]any, name, user string) {
	now := time.Now().UTC()
	doc["name"] = name
	doc["owner"] = user
	doc["creation"] = now
	doc["modified"] = now
	doc["modified_by"] = user
	if _, ok := doc["docstatus"]; !ok {
		doc["docstatus"] = 0
	}
	if _, ok := doc["idx"]; !ok {
		doc["idx"] = 0
	}
}

// InjectUpdate sets modified/modified_by for an existing document update.
// Call just before executing the UPDATE query.
func InjectUpdate(doc map[string]any, user string) {
	doc["modified"] = time.Now().UTC()
	doc["modified_by"] = user
}
