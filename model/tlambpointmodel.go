package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TLambPointModel = (*customTLambPointModel)(nil)

type (
	// TLambPointModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLambPointModel.
	TLambPointModel interface {
		tLambPointModel
		withSession(session sqlx.Session) TLambPointModel
	}

	customTLambPointModel struct {
		*defaultTLambPointModel
	}
)

// NewTLambPointModel returns a model for the database table.
func NewTLambPointModel(conn sqlx.SqlConn) TLambPointModel {
	return &customTLambPointModel{
		defaultTLambPointModel: newTLambPointModel(conn),
	}
}

func (m *customTLambPointModel) withSession(session sqlx.Session) TLambPointModel {
	return NewTLambPointModel(sqlx.NewSqlConnFromSession(session))
}
