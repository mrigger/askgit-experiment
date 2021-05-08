// +build sqlite_vtable vtable
package main

import (
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type numberModule struct {
}

func (m *numberModule) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	err := c.DeclareVTab(fmt.Sprintf(`
		CREATE TABLE %s (
			idx INT,
			val INT
		)`, args[0]))
	if err != nil {
		return nil, err
	}
	return &number{}, nil
}

func (m *numberModule) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m *numberModule) DestroyModule() {}

type number struct {
	val int
}

func (v *number) Open() (sqlite3.VTabCursor, error) {
	return &valCursor{0}, nil
}

func (v *number) BestIndex(cst []sqlite3.InfoConstraint, ob []sqlite3.InfoOrderBy) (*sqlite3.IndexResult, error) {
	return &sqlite3.IndexResult{}, nil
}

func (v *number) Disconnect() error { return nil }
func (v *number) Destroy() error    { return nil }

type valCursor struct {
	index int
}

func (vc *valCursor) Column(c *sqlite3.SQLiteContext, col int) error {
	switch col {
	case 0:
		c.ResultInt(vc.index)
	case 1:
		c.ResultInt(vc.index * vc.index)
	}
	return nil
}

func (vc *valCursor) Filter(idxNum int, idxStr string, vals []interface{}) error {
	vc.index = 0
	return nil
}

func (vc *valCursor) Next() error {
	vc.index++
	return nil
}

func (vc *valCursor) EOF() bool {
	return false
}

func (vc *valCursor) Rowid() (int64, error) {
	return int64(vc.index), nil
}

func (vc *valCursor) Close() error {
	return nil
}
