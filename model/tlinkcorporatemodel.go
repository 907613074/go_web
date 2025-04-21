package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TLinkCorporateModel = (*customTLinkCorporateModel)(nil)

type (
	// TLinkCorporateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTLinkCorporateModel.
	TLinkCorporateModel interface {
		tLinkCorporateModel
		withSession(session sqlx.Session) TLinkCorporateModel
	}

	customTLinkCorporateModel struct {
		*defaultTLinkCorporateModel
	}
)

// NewTLinkCorporateModel returns a model for the database table.
func NewTLinkCorporateModel(conn sqlx.SqlConn) TLinkCorporateModel {
	return &customTLinkCorporateModel{
		defaultTLinkCorporateModel: newTLinkCorporateModel(conn),
	}
}

func (m *customTLinkCorporateModel) withSession(session sqlx.Session) TLinkCorporateModel {
	return NewTLinkCorporateModel(sqlx.NewSqlConnFromSession(session))
}
