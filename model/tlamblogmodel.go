package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TLambLogModel = (*customTLambLogModel)(nil)

type (
	// TLambLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLambLogModel.
	TLambLogModel interface {
		tLambLogModel
		withSession(session sqlx.Session) TLambLogModel
	}

	customTLambLogModel struct {
		*defaultTLambLogModel
	}
)

// NewTLambLogModel returns a model for the database table.
func NewTLambLogModel(conn sqlx.SqlConn) TLambLogModel {
	return &customTLambLogModel{
		defaultTLambLogModel: newTLambLogModel(conn),
	}
}

func (m *customTLambLogModel) withSession(session sqlx.Session) TLambLogModel {
	return NewTLambLogModel(sqlx.NewSqlConnFromSession(session))
}
