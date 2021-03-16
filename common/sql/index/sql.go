/*
 * Copyright (c) 2018. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */
package index

import (
	"bytes"
	"context"
	"crypto/md5"
	databasesql "database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	json "github.com/pydio/cells/x/jsonx"

	"github.com/pydio/cells/common/utils/mtree"

	"github.com/pborman/uuid"
	"github.com/pydio/packr"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"

	"github.com/pydio/cells/common/log"
	"github.com/pydio/cells/common/proto/tree"
	"github.com/pydio/cells/common/sql"
	"github.com/pydio/cells/x/configx"
)

var (
	queries   = map[string]interface{}{}
	mu        atomic.Value
	inserting atomic.Value
	cond      *sync.Cond
)

const (
	batchLen = 20
	indexLen = 767
)

var (
	//	queries = map[string]string{}
	batch = "?" + strings.Repeat(", ?", batchLen-1)
)

// BatchSend sql structure
type BatchSend struct {
	in  chan *mtree.TreeNode
	out chan error
}

func init() {

	inserting.Store(make(map[string]bool))
	cond = sync.NewCond(&sync.Mutex{})

	queries["insertTree"] = func(dao sql.DAO, args ...interface{}) string {
		var num int
		if len(args) == 1 {
			num = args[0].(int)
		}

		columns := []string{"uuid", "level", "name", "leaf", "mtime", "etag", "size", "mode", "mpath1", "mpath2", "mpath3", "mpath4"}
		values := []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"}

		hash := dao.Hash("mpath1", "mpath2", "mpath3", "mpath4")
		if hash != "" {
			columns = append(columns, "hash")
			values = append(values, hash)
		}

		str := `
			insert into %%PREFIX%%_idx_tree (` +
			strings.Join(columns, ",") + `
			) values (` +
			strings.Join(values, ",") + `
			)`

		for i := 1; i < num; i++ {
			str = str + `, (` + strings.Join(values, ",") + `)`
		}

		return str
	}
	queries["updateTree"] = func(dao sql.DAO, mpathes ...string) string {

		columns := []string{"level", "name", "leaf", "mtime", "etag", "size", "mode", "mpath1", "mpath2", "mpath3", "mpath4", "uuid"}
		values := []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"}

		hash := dao.Hash("mpath1", "mpath2", "mpath3", "mpath4")
		if hash != "" {
			columns = append(columns, "hash")
			values = append(values, hash)
		}

		return `
			replace into %%PREFIX%%_idx_tree (` +
			strings.Join(columns, ",") + `
			) values (` +
			strings.Join(values, ",") + `
			)
		`
	}
	queries["updateReplace"] = func(dao sql.DAO, args ...string) (string, []interface{}) {
		whereSub, whereArgs := getMPathLike([]byte(args[0]))

		mapping := args[1:]
		mpath1 := mapping[0:2]
		mpath2 := mapping[2:4]
		mpath3 := mapping[4:6]
		mpath4 := mapping[6:8]

		var mpathSub []string

		if mpath1[0] != mpath1[1] {
			mpathSub = append(mpathSub, `mpath1 = `+dao.Concat(`"`+mpath1[1]+`."`, `SUBSTR(mpath1, `+fmt.Sprintf("%d", len(mpath1[0])+2)+`)`))
		}
		if mpath2[0] != mpath2[1] {
			mpathSub = append(mpathSub, `mpath2 = `+dao.Concat(`"`+mpath2[1]+`."`, `SUBSTR(mpath1, `+fmt.Sprintf("%d", len(mpath2[0])+2)+`)`))
		}
		if mpath3[0] != mpath3[1] {
			mpathSub = append(mpathSub, `mpath3 = `+dao.Concat(`"`+mpath3[1]+`."`, `SUBSTR(mpath1, `+fmt.Sprintf("%d", len(mpath3[0])+2)+`)`))
		}
		if mpath4[0] != mpath4[1] {
			mpathSub = append(mpathSub, `mpath4 = `+dao.Concat(`"`+mpath4[1]+`."`, `SUBSTR(mpath1, `+fmt.Sprintf("%d", len(mpath4[0])+2)+`)`))
		}

		hash := dao.Hash("mpath1", "mpath2", "mpath3", "mpath4")
		if hash != "" {
			mpathSub = append(mpathSub, `hash = `+hash)
		}

		return fmt.Sprintf(`
		update %%PREFIX%%_idx_tree set level = level + ?, %s
		where %s`, strings.Join(mpathSub, ", "), whereSub), whereArgs
	}
	queries["updateMeta"] = func(dao sql.DAO, mpathes ...string) string {
		return `UPDATE %%PREFIX%%_idx_tree set name = ?, leaf = ?, mtime = ?, etag=?, size=?, mode=? WHERE uuid = ?`
	}
	queries["updateEtag"] = func(dao sql.DAO, mpathes ...string) string {
		return `UPDATE %%PREFIX%%_idx_tree set etag = ? WHERE uuid = ?`
	}

	queries["selectNodeUuid"] = func(dao sql.DAO, mpathes ...string) string {
		return `
		select uuid, level, mpath1, mpath2, mpath3, mpath4, name, leaf, mtime, etag, size, mode
        from %%PREFIX%%_idx_tree where uuid = ?`
	}

	queries["updateNodes"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathesIn(mpathes...)

		return fmt.Sprintf(`
			update %%PREFIX%%_idx_tree set mtime = ?, etag = ?, size = size + ?
			where (%s)`, sub), args
	}

	queries["deleteTree"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathEqualsOrLike([]byte(mpathes[0]))

		return fmt.Sprintf(`
			delete from %%PREFIX%%_idx_tree
			where (%s)`, sub), args
	}

	queries["selectNode"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathEquals([]byte(mpathes[0]))

		return fmt.Sprintf(`
		SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
		FROM %%PREFIX%%_idx_tree
		WHERE %s`, sub), args
	}

	queries["selectNodes"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathesIn(mpathes...)

		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE (%s)	
			ORDER BY mpath1, mpath2, mpath3, mpath4`, sub), args
	}

	queries["tree"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))

		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE %s and level >= ?
			ORDER BY mpath1, mpath2, mpath3, mpath4`, sub), args
	}

	queries["children"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE %s AND level = ?
			ORDER BY name`, sub), args
	}

	queries["child"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE %s AND level = ? AND name like ?`, sub), args
	}

	queries["child_sqlite3"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE %s AND level = ? AND name like ? ESCAPE '\'`, sub), args
	}

	queries["lastChild"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE %s AND level = ?
			ORDER BY mpath4, mpath3, mpath2, mpath1 DESC LIMIT 1`, sub), args
	}

	queries["childrenEtags"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			SELECT etag
			FROM %%PREFIX%%_idx_tree
			WHERE %s AND level = ?
			ORDER BY name`, sub), args
	}

	queries["dirtyEtags"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathEqualsOrLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			SELECT uuid, level, mpath1, mpath2, mpath3, mpath4,  name, leaf, mtime, etag, size, mode
			FROM %%PREFIX%%_idx_tree
			WHERE etag = '-1' AND (%s) AND level >= ?
			ORDER BY level DESC`, sub), args
	}

	queries["childrenCount"] = func(dao sql.DAO, mpathes ...string) (string, []interface{}) {
		sub, args := getMPathLike([]byte(mpathes[0]))
		return fmt.Sprintf(`
			select leaf, count(leaf)
			FROM %%PREFIX%%_idx_tree
			WHERE %s AND level = ? AND name != '.pydio'
			GROUP BY leaf`, sub), args
	}
}

// IndexSQL implementation
type IndexSQL struct {
	*sql.Handler

	rootNodeId string
}

// Init handles the db version migration and prepare the statements
func (dao *IndexSQL) Init(options configx.Values) error {

	migrations := &sql.PackrMigrationSource{
		Box:         packr.NewBox("../../../common/sql/index/migrations"),
		Dir:         "./" + dao.Driver(),
		TablePrefix: dao.Prefix() + "_idx",
	}

	_, err := sql.ExecMigration(dao.DB(), dao.Driver(), migrations, migrate.Up, dao.Prefix()+"_idx_")
	if err != nil {
		return err
	}

	if options.Val("prepare").Default(true).Bool() {
		for key, query := range queries {
			if err := dao.Prepare(key, query); err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanResourcesOnDeletion revert the creation of the table for a datasource
func (dao *IndexSQL) CleanResourcesOnDeletion() (error, string) {

	migrations := &sql.PackrMigrationSource{
		Box:         packr.NewBox("../../../common/sql/index/migrations"),
		Dir:         "./" + dao.Driver(),
		TablePrefix: dao.Prefix() + "_idx",
	}

	_, err := sql.ExecMigration(dao.DB(), dao.Driver(), migrations, migrate.Down, dao.Prefix()+"_idx_")
	if err != nil {
		return err, ""
	}

	return nil, "Removed tables for index"
}

// AddNode to the underlying SQL DB.
func (dao *IndexSQL) AddNode(node *mtree.TreeNode) error {

	dao.Lock()
	defer dao.Unlock()

	mTime := node.GetMTime()
	if mTime == 0 {
		mTime = time.Now().Unix()
	}

	mpath1, mpath2, mpath3, mpath4 := prepareMPathParts(node)

	stmt, er := dao.GetStmt("insertTree")
	if er != nil {
		return er
	}
	if _, err := stmt.Exec(
		node.Uuid,
		node.Level,
		node.Name(),
		node.IsLeafInt(),
		mTime,
		node.GetEtag(),
		node.GetSize(),
		node.GetMode(),
		mpath1,
		mpath2,
		mpath3,
		mpath4,
	); err != nil {
		return err
	}

	return nil
}

// AddNodeStream creates a channel to write to the SQL database
func (dao *IndexSQL) AddNodeStream(max int) (chan *mtree.TreeNode, chan error) {

	c := make(chan *mtree.TreeNode)
	e := make(chan error)

	go func() {

		defer close(e)

		insert := func(num int, valsInsertTree []interface{}) error {
			dao.Lock()
			defer dao.Unlock()

			insertTree, er := dao.GetStmt("insertTree", num)
			if er != nil {
				return er
			}
			if _, err := insertTree.Exec(valsInsertTree...); err != nil {
				return err
			}

			return nil
		}

		valsInsertTree := []interface{}{}

		var count int
		for node := range c {

			//log.Logger(context.Background()).Info("SQL:AddNodeStream", node.ZapUuid(), node.ZapPath(), zap.String("MPath", node.MPath.String()))

			mTime := node.GetMTime()
			if mTime == 0 {
				mTime = time.Now().Unix()
			}

			mpath1, mpath2, mpath3, mpath4 := prepareMPathParts(node)

			valsInsertTree = append(valsInsertTree, node.Uuid, node.Level, node.Name(), node.IsLeafInt(), mTime, node.GetEtag(), node.GetSize(), node.GetMode(), mpath1, mpath2, mpath3, mpath4)

			count = count + 1

			if count >= max {

				if err := insert(max, valsInsertTree); err != nil {
					e <- err
				}

				count = 0
				valsInsertTree = []interface{}{}

			}
		}

		if count > 0 {
			if err := insert(count, valsInsertTree); err != nil {
				e <- err
			}
		}
	}()

	return c, e
}

// Flush the database in case of cached inserts
func (dao *IndexSQL) Flush(final bool) error {
	return nil
}

// SetNode in replacement of previous node
func (dao *IndexSQL) SetNode(node *mtree.TreeNode) error {

	dao.Lock()
	defer dao.Unlock()

	mpath1, mpath2, mpath3, mpath4 := prepareMPathParts(node)

	updateTree, er := dao.GetStmt("updateTree")
	if er != nil {
		return er
	}

	_, err := updateTree.Exec(
		node.Level,
		node.Name(),
		node.IsLeafInt(),
		node.MTime,
		node.Etag,
		node.Size,
		node.Mode,
		mpath1,
		mpath2,
		mpath3,
		mpath4,
		node.Uuid,
	)

	return err
}

// SetNode in replacement of previous node
func (dao *IndexSQL) SetNodeMeta(node *mtree.TreeNode) error {

	dao.Lock()
	defer dao.Unlock()

	updateMeta, er := dao.GetStmt("updateMeta")
	if er != nil {
		return er
	}

	_, err := updateMeta.Exec(
		node.Name(),
		node.IsLeafInt(),
		node.MTime,
		node.Etag,
		node.Size,
		node.Mode,
		node.Uuid,
	)

	return err
}

// etagFromChildren recompute ETag from children ETags
func (dao *IndexSQL) etagFromChildren(node *mtree.TreeNode) (string, error) {

	SEPARATOR := "."
	hasher := md5.New()
	dao.Lock()

	var rows *databasesql.Rows
	var err error
	defer func() {
		if rows != nil {
			rows.Close()
		}
		dao.Unlock()
	}()

	mpath := node.MPath

	// First we check if we already have an object with the same key
	if stmt, args, e := dao.GetStmtWithArgs("childrenEtags", mpath.String()); e == nil {
		rows, err = stmt.Query(append(args, len(mpath)+1)...)
		if err != nil {
			return "", err
		}
	} else {
		return "", e
	}

	first := true
	for rows.Next() {
		var etag string
		rows.Scan(&etag)
		if !first {
			hasher.Write([]byte(SEPARATOR))
		}
		hasher.Write([]byte(etag))
		first = false
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// ResyncDirtyEtags ensures that etags are rightly calculated
func (dao *IndexSQL) ResyncDirtyEtags(rootNode *mtree.TreeNode) error {

	dao.Lock()

	var rows *databasesql.Rows
	var err error

	mpath := rootNode.MPath
	if stmt, args, e := dao.GetStmtWithArgs("dirtyEtags", mpath.String()); e == nil {
		rows, err = stmt.Query(append(args, len(mpath))...) // Start at root level
		if err != nil {
			dao.Unlock()
			return err
		}
	} else {
		return e
	}
	var nodesToUpdate []*mtree.TreeNode
	for rows.Next() {
		node, e := dao.scanDbRowToTreeNode(rows)
		if e != nil {
			rows.Close()
			dao.Unlock()
			return e
		}
		nodesToUpdate = append(nodesToUpdate, node)
	}
	log.Logger(context.Background()).Info("Total Nodes Resynced", zap.Any("t", len(nodesToUpdate)))
	rows.Close()
	dao.Unlock()

	for _, node := range nodesToUpdate {
		log.Logger(context.Background()).Info("Resyncing Etag For Node", zap.Any("n", node))
		newEtag, eE := dao.etagFromChildren(node)
		if eE != nil {
			return eE
		}
		log.Logger(context.Background()).Info("Computed Etag For Node", zap.Any("etag", newEtag))
		stmt, er := dao.GetStmt("updateEtag")
		if er != nil {
			return er
		}
		if _, err = stmt.Exec(
			newEtag,
			node.Uuid,
		); err != nil {
			return err
		}
	}
	return nil

}

// SetNodes returns a channel and waits for arriving nodes before updating them in batch.
func (dao *IndexSQL) SetNodes(etag string, deltaSize int64) sql.BatchSender {

	b := NewBatchSend()

	go func() {
		dao.Lock()
		defer dao.Unlock()

		defer func() {
			close(b.out)
		}()

		insert := func(mpathes ...interface{}) {

			updateNodes, args, e := dao.GetStmtWithArgs("updateNodes", mpathes...)
			if e != nil {
				b.out <- e
			} else {
				if _, err := updateNodes.Exec(append([]interface{}{time.Now().Unix(), etag, deltaSize}, args...)...); err != nil {
					b.out <- err
				}
			}

		}

		all := make([]interface{}, 0, batchLen)

		for node := range b.in {
			all = append(all, node.MPath.String())
			if len(all) == cap(all) {
				insert(all...)
				all = all[:0]
			}
		}

		if len(all) > 0 {
			insert(all...)
		}

	}()

	return b
}

// DelNode from database
func (dao *IndexSQL) DelNode(node *mtree.TreeNode) error {

	dao.Lock()
	defer dao.Unlock()

	stmt, args, e := dao.GetStmtWithArgs("deleteTree", node.MPath.String())
	if e != nil {
		return e
	}
	if _, err := stmt.Exec(args...); err != nil {
		return err
	}

	/*
		if len(dao.commitsTableName) > 0 {
			if _, err = dao.GetStmt("deleteCommits").Exec(
				mpath, mpathLike,
			); err != nil {
				return err
			}
		}
	*/

	return nil
}

// GetNode from path
func (dao *IndexSQL) GetNode(path mtree.MPath) (*mtree.TreeNode, error) {

	dao.Lock()
	defer dao.Unlock()

	if len(path) == 0 {
		return nil, fmt.Errorf("Empty path")
	}

	node := mtree.NewTreeNode()
	node.SetMPath(path...)

	mpath := node.MPath.String()

	if stmt, args, e := dao.GetStmtWithArgs("selectNode", mpath); e == nil {
		row := stmt.QueryRow(args...)
		treeNode, err := dao.scanDbRowToTreeNode(row)
		if err != nil {
			return nil, err
		}
		return treeNode, nil
	} else {
		return nil, e
	}
}

// GetNodeByUUID returns the node stored with the unique uuid
func (dao *IndexSQL) GetNodeByUUID(uuid string) (*mtree.TreeNode, error) {

	dao.Lock()
	defer dao.Unlock()

	stmt, er := dao.GetStmt("selectNodeUuid")
	if er != nil {
		return nil, er
	}
	row := stmt.QueryRow(uuid)
	treeNode, err := dao.scanDbRowToTreeNode(row)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return treeNode, nil
}

// GetNodes List
func (dao *IndexSQL) GetNodes(mpathes ...mtree.MPath) chan *mtree.TreeNode {

	dao.Lock()

	c := make(chan *mtree.TreeNode)

	go func() {

		defer func() {
			close(c)
			dao.Unlock()
		}()

		get := func(mpathes ...interface{}) {
			if stmt, args, e := dao.GetStmtWithArgs("selectNodes", mpathes...); e == nil {
				rows, err := stmt.Query(args...)
				if err != nil {
					return
				}
				defer rows.Close()

				for rows.Next() {
					node, err := dao.scanDbRowToTreeNode(rows)
					if err != nil {
						break
					}

					c <- node
				}
			} else {
				log.Logger(context.Background()).Error("Error while getting statement in GetNodes", zap.Error(e))
				return
			}
		}

		all := make([]interface{}, 0, batchLen)

		for _, mpath := range mpathes {
			all = append(all, mpath.String())
			if len(all) == cap(all) {
				get(all...)
				all = all[:0]
			}
		}

		if len(all) > 0 {
			get(all...)
		}
	}()

	return c
}

// GetNodeChild from node path whose name matches
func (dao *IndexSQL) GetNodeChild(reqPath mtree.MPath, reqName string) (*mtree.TreeNode, error) {

	dao.Lock()
	defer dao.Unlock()

	node := mtree.NewTreeNode()
	node.SetMPath(reqPath...)

	mpath := node.MPath

	stmtName := "child"
	if dao.Driver() == "sqlite3" {
		stmtName = "child_sqlite3"
	}

	if stmt, args, e := dao.GetStmtWithArgs(stmtName, mpath.String()); e == nil {
		// Escape for LIKE query
		reqName = strings.ReplaceAll(reqName, "_", "\\_")
		reqName = strings.ReplaceAll(reqName, "%", "\\%")
		row := stmt.QueryRow(append(args, len(reqPath)+1, reqName)...)
		treeNode, err := dao.scanDbRowToTreeNode(row)
		if err != nil {
			return nil, err
		}
		return treeNode, nil
	} else {
		return nil, e
	}

}

// GetNodeLastChild from path
func (dao *IndexSQL) GetNodeLastChild(reqPath mtree.MPath) (*mtree.TreeNode, error) {

	dao.Lock()
	defer dao.Unlock()

	node := mtree.NewTreeNode()
	node.SetMPath(reqPath...)

	mpath := node.MPath

	if stmt, args, e := dao.GetStmtWithArgs("lastChild", mpath.String()); e == nil {
		row := stmt.QueryRow(append(args, len(reqPath)+1)...)
		treeNode, err := dao.scanDbRowToTreeNode(row)
		if err != nil {
			return nil, err
		}
		return treeNode, nil
	} else {
		return nil, e
	}

}

// GetNodeFirstAvailableChildIndex from path
func (dao *IndexSQL) GetNodeFirstAvailableChildIndex(reqPath mtree.MPath) (uint64, error) {

	all := []int{}

	for node := range dao.GetNodeChildren(reqPath) {
		all = append(all, int(node.MPath.Index()))
	}

	if len(all) == 0 {
		return 1, nil
	}

	sort.Ints(all)
	max := all[len(all)-1]

	for i := 1; i <= max; i++ {
		found := false
		for _, v := range all {
			if i == v {
				// We found the entry, so next one
				found = true
				break
			}
		}

		if !found {
			// This number is not present, returning it
			return uint64(i), nil
		}
	}

	return uint64(max + 1), nil
}

// GetNodeChildrenCounts List
func (dao *IndexSQL) GetNodeChildrenCounts(path mtree.MPath) (int, int) {

	dao.Lock()
	defer dao.Unlock()

	node := mtree.NewTreeNode()
	node.SetMPath(path...)

	mpath := node.MPath

	var folderCount, fileCount int

	// First we check if we already have an object with the same key
	if stmt, args, e := dao.GetStmtWithArgs("childrenCount", mpath.String()); e == nil {
		if rows, e := stmt.Query(append(args, len(path)+1)...); e == nil {
			defer rows.Close()
			for rows.Next() {
				var leaf bool
				var count int
				if sE := rows.Scan(&leaf, &count); sE == nil {
					if leaf {
						fileCount = count
					} else {
						folderCount = count
					}
				}
			}
		}
	}

	return folderCount, fileCount
}

// GetNodeChildren List
func (dao *IndexSQL) GetNodeChildren(path mtree.MPath) chan *mtree.TreeNode {

	dao.Lock()

	c := make(chan *mtree.TreeNode)

	go func() {
		var rows *databasesql.Rows
		var err error

		defer func() {
			if rows != nil {
				rows.Close()
			}
			close(c)
			dao.Unlock()
		}()

		node := mtree.NewTreeNode()
		node.SetMPath(path...)

		mpath := node.MPath

		// First we check if we already have an object with the same key
		if stmt, args, e := dao.GetStmtWithArgs("children", mpath.String()); e == nil {
			rows, err = stmt.Query(append(args, len(path)+1)...)
			if err != nil {
				return
			}

			for rows.Next() {
				treeNode, err := dao.scanDbRowToTreeNode(rows)
				if err != nil {
					break
				}
				c <- treeNode
			}
		}
	}()

	return c
}

// GetNodeTree List from the path
func (dao *IndexSQL) GetNodeTree(path mtree.MPath) chan *mtree.TreeNode {

	dao.Lock()

	c := make(chan *mtree.TreeNode)

	go func() {
		var rows *databasesql.Rows
		var err error

		defer func() {
			if rows != nil {
				rows.Close()
			}

			close(c)
			dao.Unlock()
		}()

		node := mtree.NewTreeNode()
		node.SetMPath(path...)

		mpath := node.MPath

		// First we check if we already have an object with the same key
		if stmt, args, e := dao.GetStmtWithArgs("tree", mpath.String()); e == nil {
			rows, err = stmt.Query(append(args, len(mpath)+1)...)
			if err != nil {
				return
			}

			for rows.Next() {
				treeNode, err := dao.scanDbRowToTreeNode(rows)
				if err != nil {
					break
				}
				c <- treeNode
			}
		}
	}()

	return c
}

// MoveNodeTree move all the nodes belonging to a tree by calculating the new mpathes
func (dao *IndexSQL) MoveNodeTree(nodeFrom *mtree.TreeNode, nodeTo *mtree.TreeNode) error {
	if nodeFrom == nil {
		return errors.New("Source node cannot be empty")
	}

	if nodeTo == nil {
		return errors.New("Target node cannot be empty")
	}

	pathFrom := nodeFrom.MPath
	pathTo := nodeTo.MPath

	mpath1From, mpath2From, mpath3From, mpath4From := prepareMPathParts(nodeFrom)
	mpath1To, mpath2To, mpath3To, mpath4To := prepareMPathParts(nodeTo)

	nodeFrom.SetName(nodeTo.Name())
	nodeFrom.SetMPath(pathTo...)

	// Start by updating the original node
	if err := dao.SetNode(nodeFrom); err != nil {
		return err
	}

	// Then replace the children mpaths
	updateChildren, updateChildrenArgs, err := dao.GetStmtWithArgs("updateReplace",
		pathFrom.String(),
		mpath1From, mpath1To,
		mpath2From, mpath2To,
		mpath3From, mpath3To,
		mpath4From, mpath4To,
	)
	if err != nil {
		return err
	}

	updateChildrenArgs = append([]interface{}{
		len(pathTo) - len(pathFrom),
	}, updateChildrenArgs...)

	res, err := updateChildren.Exec(
		updateChildrenArgs...,
	)

	if err != nil {
		return err
	}

	nrows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	t1 := time.Now()
	ctx := context.Background()

	log.Logger(ctx).Info("[MoveNodeTree] Finished moving", zap.Int64("nodes", nrows), zap.Duration("duration", time.Now().Sub(t1)))
	return nil
}

func (dao *IndexSQL) scanDbRowToTreeNode(row sql.Scanner) (*mtree.TreeNode, error) {
	var (
		uuid   string
		mpath1 string
		mpath2 string
		mpath3 string
		mpath4 string
		level  uint32
		name   string
		leaf   int32
		mtime  int64
		etag   string
		size   int64
		mode   int32
	)

	if err := row.Scan(&uuid, &level, &mpath1, &mpath2, &mpath3, &mpath4, &name, &leaf, &mtime, &etag, &size, &mode); err != nil {
		return nil, err
	}
	nodeType := tree.NodeType_LEAF
	if leaf == 0 {
		nodeType = tree.NodeType_COLLECTION
	}

	node := mtree.NewTreeNode()

	var mpath []uint64
	for _, m := range strings.Split(mpath1+mpath2+mpath3+mpath4, ".") {
		i, _ := strconv.ParseUint(m, 10, 64)
		mpath = append(mpath, i)
	}
	node.SetMPath(mpath...)
	// node.SetBytes(rat)

	metaName, _ := json.Marshal(name)
	node.Node = &tree.Node{
		Uuid:      uuid,
		Type:      nodeType,
		MTime:     mtime,
		Etag:      etag,
		Size:      size,
		Mode:      mode,
		MetaStore: map[string]string{"name": string(metaName)},
	}

	return node, nil
}

func (dao *IndexSQL) Path(strpath string, create bool, reqNode ...*tree.Node) (mtree.MPath, []*mtree.TreeNode, error) {

	var path mtree.MPath
	var err error

	created := []*mtree.TreeNode{}

	if len(strpath) == 0 || strpath == "/" {
		return []uint64{1}, created, nil
	}

	names := strings.Split(fmt.Sprintf("/%s", strings.TrimLeft(strpath, "/")), "/")

	path = make([]uint64, len(names))
	path[0] = 1
	parents := make([]*mtree.TreeNode, len(names))

	// Reading root path
	node, err := dao.GetNode(path[0:1])
	if err != nil || node == nil {
		// Making sure we have a node in the database
		rootNodeId := "ROOT"
		if dao.rootNodeId != "" {
			rootNodeId = dao.rootNodeId
		}
		node = NewNode(&tree.Node{
			Uuid: rootNodeId,
			Type: tree.NodeType_COLLECTION,
		}, []uint64{1}, []string{""})

		if err = dao.AddNode(node); err != nil {
			// Has it been created elsewhere ?
			node, err = dao.GetNode(path[0:1])
			if err != nil || node == nil {
				return path, created, err
			}
		} else {
			created = append(created, node)
		}
	}

	parents[0] = node

	maxLevel := len(names) - 1

	for level := 1; level <= maxLevel; level++ {

		p := node

		if create {
			// Making sure we lock the parent node
			cond.L.Lock()
			for {
				current := inserting.Load().(map[string]bool)

				if _, ok := current[p.Uuid]; !ok {
					current[p.Uuid] = true
					inserting.Store(current)
					break
				}

				cond.Wait()
			}
			cond.L.Unlock()
		}

		node, _ = dao.GetNodeChild(path[0:level], names[level])

		if nil != node {
			path[level] = node.MPath[len(node.MPath)-1]
			parents[level] = node

			node.Path = strings.Trim(strings.Join(names[0:level], "/"), "/")
		} else {
			if create {
				if path[level], err = dao.GetNodeFirstAvailableChildIndex(path[0:level]); err != nil {
					return nil, created, err
				}

				if level == len(names)-1 && len(reqNode) > 0 {
					node = NewNode(reqNode[0], path[0:level+1], names[0:level+1])
				} else {
					node = NewNode(&tree.Node{
						Type:  tree.NodeType_COLLECTION,
						Mode:  0777,
						MTime: time.Now().Unix(),
					}, path[0:level+1], names[0:level+1])
				}

				if node.Uuid == "" {
					node.Uuid = uuid.New()
				}

				if node.Etag == "" {
					// Should only happen for folders - generate first Etag from uuid+mtime
					node.Etag = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%d", node.Uuid, node.MTime))))
				}

				err = dao.AddNode(node)

				cond.L.Lock()
				current := inserting.Load().(map[string]bool)
				delete(current, p.Uuid)
				inserting.Store(current)
				cond.L.Unlock()

				cond.Signal()

				if err != nil {
					return nil, created, err
				}

				created = append(created, node)

				parents[level] = node
			} else {
				return nil, created, nil
			}

		}

		if create {
			cond.L.Lock()
			current := inserting.Load().(map[string]bool)
			delete(current, p.Uuid)
			inserting.Store(current)
			cond.L.Unlock()

			cond.Signal()
		}
	}

	return path, created, err
}

func (dao *IndexSQL) GetSQLDAO() sql.DAO {
	return dao
}

func (dao *IndexSQL) lock() {
	if current, ok := mu.Load().(*sync.Mutex); ok {
		current.Lock()
	}
}

func (dao *IndexSQL) unlock() {
	if current, ok := mu.Load().(*sync.Mutex); ok {
		current.Unlock()
	}
}

// NewBatchSend Creation of the channels
func NewBatchSend() *BatchSend {
	b := new(BatchSend)
	b.in = make(chan *mtree.TreeNode)
	b.out = make(chan error, 1)

	return b
}

// Send a node to the batch
func (b *BatchSend) Send(arg interface{}) {
	if node, ok := arg.(*mtree.TreeNode); ok {
		b.in <- node
	}
}

// Close the Batch
func (b *BatchSend) Close() error {
	close(b.in)

	err := <-b.out

	return err
}

// Split node.MPath into 4 strings for storing in DB
func prepareMPathParts(node *mtree.TreeNode) (string, string, string, string) {
	mPath := make([]byte, indexLen*4)
	copy(mPath, []byte(node.MPath.String()))
	mPath1 := string(bytes.Trim(mPath[(indexLen*0):(indexLen*1-1)], "\x00"))
	mPath2 := string(bytes.Trim(mPath[(indexLen*1):(indexLen*2-1)], "\x00"))
	mPath3 := string(bytes.Trim(mPath[(indexLen*2):(indexLen*3-1)], "\x00"))
	mPath4 := string(bytes.Trim(mPath[(indexLen*3):(indexLen*4-1)], "\x00"))
	return mPath1, mPath2, mPath3, mPath4
}

// where t.mpath = ?
func getMPathEquals(mpath []byte) (string, []interface{}) {
	var res []string
	var args []interface{}

	for {
		var cnt int
		cnt = (len(mpath) - 1) / indexLen
		res = append(res, fmt.Sprintf(`mpath%d LIKE ?`, cnt+1))
		args = append(args, mpath[(cnt*indexLen):])

		if idx := cnt * indexLen; idx == 0 {
			break
		}

		mpath = mpath[0 : cnt*indexLen]
	}

	return strings.Join(res, " and "), args
}

// t.mpath LIKE ?
func getMPathLike(mpath []byte) (string, []interface{}) {
	var res []string
	var args []interface{}

	mpath = append(mpath, []byte(".%")...)

	done := false
	for {
		var cnt int
		cnt = (len(mpath) - 1) / indexLen

		if !done {
			res = append(res, fmt.Sprintf(`mpath%d LIKE ?`, cnt+1))
			args = append(args, mpath[(cnt*indexLen):])
			done = true
		} else {
			res = append(res, fmt.Sprintf(`mpath%d LIKE ?`, cnt+1))
			args = append(args, mpath[(cnt*indexLen):])
		}

		if idx := cnt * indexLen; idx == 0 {
			break
		}

		mpath = mpath[0 : cnt*indexLen]
	}

	return strings.Join(res, " and "), args
}

// and (t.mpath = ? OR t.mpath LIKE ?)
func getMPathEqualsOrLike(mpath []byte) (string, []interface{}) {
	var res []string
	var args []interface{}

	mpath = append(mpath, []byte(".%")...)

	done := false
	for {
		var cnt int
		cnt = (len(mpath) - 1) / indexLen

		if !done {
			res = append(res, fmt.Sprintf(`mpath%d LIKE ?`, cnt+1))
			res = append(res, fmt.Sprintf(`mpath%d LIKE ?`, cnt+1))

			args = append(args, mpath[(cnt*indexLen):len(mpath)-2], mpath[(cnt*indexLen):])

			done = true
		} else {
			res = append(res, fmt.Sprintf(`mpath%d LIKE "%s"`, cnt+1, mpath[(cnt*indexLen):]))
		}

		if idx := cnt * indexLen; idx == 0 {
			break
		}

		mpath = mpath[0 : cnt*indexLen]
	}

	return strings.Join(res, " or "), args
}

// where t.mpath in (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
func getMPathesIn(mpathes ...string) (string, []interface{}) {
	var res []string
	var args []interface{}

	for _, mpath := range mpathes {
		r, a := getMPathEquals([]byte(mpath))
		res = append(res, fmt.Sprintf(`(%s)`, r))
		args = append(args, a...)
	}

	return strings.Join(res, " or "), args
}
