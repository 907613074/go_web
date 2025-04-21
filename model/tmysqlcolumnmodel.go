package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TMysqlColumnModel = (*customTMysqlColumnModel)(nil)

type (
	// TMysqlColumnModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTMysqlColumnModel.
	TMysqlColumnModel interface {
		tMysqlColumnModel
		withSession(session sqlx.Session) TMysqlColumnModel
	}

	customTMysqlColumnModel struct {
		*defaultTMysqlColumnModel
	}
)

// NewTMysqlColumnModel returns a model for the database table.
func NewTMysqlColumnModel(conn sqlx.SqlConn) TMysqlColumnModel {
	return &customTMysqlColumnModel{
		defaultTMysqlColumnModel: newTMysqlColumnModel(conn),
	}
}

func (m *customTMysqlColumnModel) withSession(session sqlx.Session) TMysqlColumnModel {
	return NewTMysqlColumnModel(sqlx.NewSqlConnFromSession(session))
}
