package common_controllers

import (
	"bitsflow/common/db"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var lock = sync.RWMutex{}
var PersistenceRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type PersistenceController struct {
	clientConfig *db.DBClientConfig
	serverConfig *db.DBServerConfig

	connData     *db.ConnData
	Error        error
	RowsAffected int64
	Rows         pgx.Rows
	Response     pgx.Row
}

func (c *PersistenceController) Setup(connData *db.ConnData, clientConfig *db.DBClientConfig, serverConfig *db.DBServerConfig) *db.ConnData {
	c.clientConfig = clientConfig
	c.serverConfig = serverConfig
	c.connData = connData
	c.getConnection()

	for ; c.Error != nil && c.Error.Error() == "common_db_max_conn_reached"; c.getConnection() {
		time.Sleep(10 * time.Millisecond)
	}

	return c.connData
}
func (c *PersistenceController) getConnection() {
	//lock.Lock()
	//defer lock.Unlock()
	c.connData, c.Error = db.GetConnection(c.connData, c.clientConfig, c.serverConfig)
	if c.Error != nil {
		fmt.Println("Persistence getConnection Error: ", c.Error)
	}
}

func (c *PersistenceController) ReleaseConnection() {
	//lock.Lock()
	//defer lock.Unlock()
	c.Error = db.ReleaseConnection(c.connData)
	if c.Error != nil {
		fmt.Println("Persistence ReleaseConnection Error: ", c.Error)
	}
}

func (c *PersistenceController) Exec(ctx context.Context, sql string, arguments ...interface{}) {

	c.getConnection()
	var res pgconn.CommandTag = nil
	res, c.Error = c.connData.Conn.Exec(ctx, sql, arguments...)
	c.RowsAffected = res.RowsAffected()
	if c.Error != nil {
		fmt.Println("Persistence Exec Error: ", c.Error)
	}
}

func (c *PersistenceController) QueryRow(ctx context.Context, sql string, arguments ...interface{}) {

	var res pgconn.CommandTag = nil
	c.Response = c.connData.Conn.QueryRow(ctx, sql, arguments...)
	c.RowsAffected = res.RowsAffected()
}

func (c *PersistenceController) Query(ctx context.Context, sql string, arguments ...interface{}) {

	c.Rows, c.Error = c.connData.Conn.Query(ctx, sql, arguments...)
	c.RowsAffected = c.Rows.CommandTag().RowsAffected()
	if c.Error != nil {
		fmt.Println("Persistence Query Error: ", c.Error)
	}
}

func (c *PersistenceController) Scan(entities ...interface{}) {

	if c.Response != nil {
		c.Error = c.Response.Scan(entities...)
	}
	if c.Error != nil {
		fmt.Println("Persistence Scan Error: ", c.Error)
	}
}
func (c *PersistenceController) ScanRow(entities ...interface{}) {

	c.Error = c.Rows.Scan(entities...)
	if c.Error != nil {
		fmt.Println("Persistence ScanRow Error: ", c.Error)
	}
}

func (c *PersistenceController) Next() bool {
	return c.Rows.Next()
}
