package project

import (
	"database/sql"
	"github.com/modelhub/core/util"
	"github.com/modelhub/vada"
	"github.com/robsix/golog"
	"strings"
)

func NewSqlProjectStore(db *sql.DB, vada vada.VadaClient, ossBucketPrefix string, ossBucketPolicy vada.BucketPolicy, log golog.Log) ProjectStore {

	getter := func(query string, colLen int, args ...interface{}) ([]*Project, error) {
		ps := make([]*Project, 0, colLen)
		rowsScan := func(rows *sql.Rows) error {
			p := Project{}
			if err := rows.Scan(&p.Id, &p.Name, &p.Created, &p.ImageFileExtension); err != nil {
				return err
			}
			ps = append(ps, &p)
			return nil
		}
		return ps, util.SqlQuery(db, rowsScan, query, args...)
	}

	getterString := func(query string, args ...interface{}) (string, error) {
		str := ""
		rowsScan := func(rows *sql.Rows) error {
			if err := rows.Scan(&str); err != nil {
				return err
			}
			return nil
		}
		return str, util.SqlQuery(db, rowsScan, query, args...)
	}

	offsetGetter := func(query string, args ...interface{}) ([]*Project, int, error) {
		ps := make([]*Project, 0, util.DefaultSqlOffsetQueryLimit)
		totalResults := 0
		rowsScan := func(rows *sql.Rows) error {
			if util.RowsContainsOnlyTotalResults(&totalResults, rows) {
				return nil
			}
			p := Project{}
			if err := rows.Scan(&totalResults, &p.Id, &p.Name, &p.Created, &p.ImageFileExtension); err != nil {
				return err
			}
			ps = append(ps, &p)
			return nil
		}
		return ps, totalResults, util.SqlQuery(db, rowsScan, query, args...)
	}

	offsetGetterMembership := func(query string, args ...interface{}) ([]*Membership, int, error) {
		ms := make([]*Membership, 0, util.DefaultSqlOffsetQueryLimit)
		totalResults := 0
		rowsScan := func(rows *sql.Rows) error {
			if util.RowsContainsOnlyTotalResults(&totalResults, rows) {
				return nil
			}
			m := Membership{}
			if err := rows.Scan(&totalResults, &m.User, &m.Role); err != nil {
				return err
			}
			ms = append(ms, &m)
			return nil
		}
		return ms, totalResults, util.SqlQuery(db, rowsScan, query, args...)
	}

	offsetGetterInUserContext := func(query string, args ...interface{}) ([]*ProjectInUserContext, int, error) {
		ps := make([]*ProjectInUserContext, 0, util.DefaultSqlOffsetQueryLimit)
		totalResults := 0
		rowsScan := func(rows *sql.Rows) error {
			if util.RowsContainsOnlyTotalResults(&totalResults, rows) {
				return nil
			}
			p := ProjectInUserContext{}
			if err := rows.Scan(&totalResults, &p.Id, &p.Name, &p.Created, &p.ImageFileExtension, &p.Role); err != nil {
				return err
			}
			ps = append(ps, &p)
			return nil
		}
		return ps, totalResults, util.SqlQuery(db, rowsScan, query, args...)
	}

	create := func(forUser string, id string, name string, description string, imageFileExtension string) (*Project, error) {
		if ps, err := getter("CALL projectCreate(?, ?, ?, ?, ?)", 1, forUser, id, name, description, imageFileExtension); len(ps) == 1 {
			return ps[0], err
		} else {
			return nil, err
		}
	}

	delete := func(forUser string, id string) error {
		return util.SqlExec(db, "CALL projectDelete(?, ?)", forUser, id)
	}

	setName := func(forUser string, id string, newName string) error {
		return util.SqlExec(db, "CALL projectSetName(?, ?, ?)", forUser, id, newName)
	}

	setDescription := func(forUser string, id string, newDescription string) error {
		return util.SqlExec(db, "CALL projectSetDescription(?, ?, ?)", forUser, id, newDescription)
	}

	setImageFileExtension := func(forUser string, id string, newImageFileExtension string) error {
		return util.SqlExec(db, "CALL projectSetImageFileExtension(?, ?, ?)", forUser, id, newImageFileExtension)
	}

	addUsers := func(forUser string, id string, role role, users []string) error {
		return util.SqlExec(db, "CALL projectAddUsers(?, ?, ?, ?)", forUser, id, string(role), strings.Join(users, ","))
	}

	removeUsers := func(forUser string, id string, users []string) error {
		return util.SqlExec(db, "CALL projectRemoveUsers(?, ?, ?)", forUser, id, strings.Join(users, ","))
	}

	acceptInvite := func(forUser string, id string) error {
		return util.SqlExec(db, "CALL projectAcceptInvite(?, ?)", forUser, id)
	}

	declineInvite := func(forUser string, id string) error {
		return util.SqlExec(db, "CALL projectDeclineInvite(?, ?)", forUser, id)
	}

	getMemberships := func(forUser string, id string, role role, offset int, limit int, sortBy sortBy) ([]*Membership, int, error) {
		return offsetGetterMembership("CALL projectGetMemberships(?, ?, ?, ?, ?, ?)", forUser, id, string(role), offset, limit, string(sortBy))
	}

	getMembershipInvites := func(forUser string, id string, role role, offset int, limit int, sortBy sortBy) ([]*Membership, int, error) {
		return offsetGetterMembership("CALL projectGetMembershipInvites(?, ?, ?, ?, ?, ?)", forUser, id, string(role), offset, limit, string(sortBy))
	}

	getDescription := func(forUser string, id string) (string, error) {
		return getterString("CALL projectGetDescription(?, ?)", forUser, id)
	}

	get := func(forUser string, ids []string) ([]*Project, error) {
		return getter("CALL projectGet(?, ?)", len(ids), forUser, strings.Join(ids, ","))
	}

	getInUserContext := func(forUser string, user string, role role, offset int, limit int, sortBy sortBy) ([]*ProjectInUserContext, int, error) {
		return offsetGetterInUserContext("CALL projectGetInUserContext(?, ?, ?, ?, ?, ?)", forUser, user, string(role), offset, limit, string(sortBy))
	}

	getInUserInviteContext := func(forUser string, user string, role role, offset int, limit int, sortBy sortBy) ([]*ProjectInUserContext, int, error) {
		return offsetGetterInUserContext("CALL projectGetInUserInviteContext(?, ?, ?, ?, ?, ?)", forUser, user, string(role), offset, limit, string(sortBy))
	}

	search := func(forUser string, search string, offset int, limit int, sortBy sortBy) ([]*Project, int, error) {
		return offsetGetter("CALL projectSearch(?, ?, ?, ?, ?)", forUser, search, offset, limit, string(sortBy))
	}

	return newProjectStore(create, delete, setName, setDescription, setImageFileExtension, addUsers, removeUsers, acceptInvite, declineInvite, util.GetRoleFunc(db), getMemberships, getMembershipInvites, getDescription, get, getInUserContext, getInUserInviteContext, search, vada, ossBucketPrefix, ossBucketPolicy, log)
}
