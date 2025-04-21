package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TLambUserModel = (*customTLambUserModel)(nil)

type (
	// TLambUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLambUserModel.
	TLambUserModel interface {
		tLambUserModel
		withSession(session sqlx.Session) TLambUserModel
	}

	customTLambUserModel struct {
		*defaultTLambUserModel
	}
)

// NewTLambUserModel returns a model for the database table.
func NewTLambUserModel(conn sqlx.SqlConn) TLambUserModel {
	return &customTLambUserModel{
		defaultTLambUserModel: newTLambUserModel(conn),
	}
}

func (m *customTLambUserModel) withSession(session sqlx.Session) TLambUserModel {
	return NewTLambUserModel(sqlx.NewSqlConnFromSession(session))
}
