package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TLinkAinModel = (*customTLinkAinModel)(nil)

type (
	// TLinkAinModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLinkAinModel.
	TLinkAinModel interface {
		tLinkAinModel
		withSession(session sqlx.Session) TLinkAinModel
	}

	customTLinkAinModel struct {
		*defaultTLinkAinModel
	}
)

// NewTLinkAinModel returns a model for the database table.
func NewTLinkAinModel(conn sqlx.SqlConn) TLinkAinModel {
	return &customTLinkAinModel{
		defaultTLinkAinModel: newTLinkAinModel(conn),
	}
}

func (m *customTLinkAinModel) withSession(session sqlx.Session) TLinkAinModel {
	return NewTLinkAinModel(sqlx.NewSqlConnFromSession(session))
}
