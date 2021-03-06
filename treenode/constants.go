package treenode

import (
	"strings"
)

const (
	NameAsc  = sortBy("nameAsc")
	NameDesc = sortBy("nameDesc")

	Any          = nodeType("any") //used for results filtering only
	Folder       = nodeType("folder")
	Document     = nodeType("document")
	ProjectSpace = nodeType("projectSpace")
)

type sortBy string
type nodeType string

func SortBy(sb string) sortBy {
	switch strings.ToLower(sb) {
	case "namedesc":
		return NameDesc
	default:
		return NameAsc
	}
}

func NodeType(nt string) nodeType {
	switch strings.ToLower(nt) {
	case "folder":
		return Folder
	case "document":
		return Document
	case "projectspace":
		return ProjectSpace
	default:
		return Any
	}
}
