package treenode

import (
	"database/sql"
	"github.com/modelhub/core/util"
	"github.com/modelhub/core/vada"
	"github.com/robsix/golog"
	"github.com/robsix/json"
	"strings"
)

func NewSqlTreeNodeStore(db *sql.DB, vada vada.VadaClient, ossBucketPrefix string, log golog.Log) TreeNodeStore {

	getter := func(query string, colLen int, args ...interface{}) ([]*TreeNode, error) {
		tns := make([]*TreeNode, 0, colLen)
		rowsScan := func(rows *sql.Rows) error {
			tn := TreeNode{}
			scanNodeType := ""
			if err := rows.Scan(&tn.Id, &tn.Parent, &tn.Project, &tn.Name, &scanNodeType); err != nil {
				return err
			}
			tn.NodeType = nodeType(scanNodeType)
			tns = append(tns, &tn)
			return nil
		}
		return tns, util.SqlQuery(db, rowsScan, query, args...)
	}

	offsetGetter := func(query string, args ...interface{}) ([]*TreeNode, int, error) {
		tns := make([]*TreeNode, 0, util.DefaultSqlOffsetQueryLimit)
		totalResults := 0
		rowsScan := func(rows *sql.Rows) error {
			tn := TreeNode{}
			scanNodeType := ""
			if err := rows.Scan(&totalResults, &tn.Id, &tn.Parent, &tn.Project, &tn.Name, &scanNodeType); err != nil {
				return err
			}
			tn.NodeType = nodeType(scanNodeType)
			tns = append(tns, &tn)
			return nil
		}
		return tns, totalResults, util.SqlQuery(db, rowsScan, query, args...)
	}

	createFolder := func(forUser string, parent string, name string) (*TreeNode, error) {
		if tns, err := getter("CALL treeNodeCreateFolder(?, ?, ?)", 1, forUser, parent, name); len(tns) == 1 {
			return tns[0], err
		} else {
			return nil, err
		}
	}

	createDocument := func(forUser string, parent string, name string, documentVersion string, uploadComment string, fileExtension string, urn string, status string) (*TreeNode, error) {
		if tns, err := getter("CALL treeNodeCreateDocument(?, ?, ?, ?, ?, ?, ?, ?)", 1, forUser, parent, name, documentVersion, uploadComment, fileExtension, urn, status); len(tns) == 1 {
			return tns[0], err
		} else {
			return nil, err
		}
	}

	createViewerState := func(forUser string, parent string, name string, createComment string, definition *json.Json) (*TreeNode, error) {
		// TODO
		return nil, nil
	}

	setName := func(forUser string, id string, newName string) error {
		return util.SqlExec(db, "CALL treeNodeSetName(?, ?, ?)", forUser, id, newName)
	}

	move := func(forUser string, newParent string, ids []string) error {
		return util.SqlExec(db, "CALL treeNodeMove(?, ?, ?)", forUser, newParent, strings.Join(ids, ","))
	}

	get := func(forUser string, ids []string) ([]*TreeNode, error) {
		return getter("CALL treeNodeGet(?, ?)", len(ids), forUser, strings.Join(ids, ","))
	}

	getChildren := func(forUser string, id string, nt nodeType, offset int, limit int, sortBy sortBy) ([]*TreeNode, int, error) {
		return offsetGetter("CALL treeNodeGetChildren(?, ?, ?, ?, ?, ?)", forUser, id, string(nt), offset, limit, string(sortBy))
	}

	getParents := func(forUser string, id string) ([]*TreeNode, error) {
		return getter("CALL treeNodeGetParents(?, ?)", util.DefaultSqlOffsetQueryLimit, forUser, id)
	}

	globalSearch := func(forUser string, search string, nt nodeType, offset int, limit int, sortBy sortBy) ([]*TreeNode, int, error) {
		return offsetGetter("CALL treeNodeGlobalSearch(?, ?, ?, ?, ?, ?)", forUser, search, string(nt), offset, limit, string(sortBy))
	}

	projectSearch := func(forUser string, project string, search string, nt nodeType, offset int, limit int, sortBy sortBy) ([]*TreeNode, int, error) {
		return offsetGetter("CALL treeNodeProjectSearch(?, ?, ?, ?, ?, ?, ?)", forUser, project, search, string(nt), offset, limit, string(sortBy))
	}

	return newTreeNodeStore(createFolder, createDocument, createViewerState, setName, move, get, getChildren, getParents, globalSearch, projectSearch, util.GetRoleFunc(db), vada, ossBucketPrefix, log)
}
