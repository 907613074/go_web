package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TTradeStockModel = (*customTTradeStockModel)(nil)

type (
	// TTradeStockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTTradeStockModel.
	TTradeStockModel interface {
		tTradeStockModel
		withSession(session sqlx.Session) TTradeStockModel
	}

	customTTradeStockModel struct {
		*defaultTTradeStockModel
	}
)

// NewTTradeStockModel returns a model for the database table.
func NewTTradeStockModel(conn sqlx.SqlConn) TTradeStockModel {
	return &customTTradeStockModel{
		defaultTTradeStockModel: newTTradeStockModel(conn),
	}
}

func (m *customTTradeStockModel) withSession(session sqlx.Session) TTradeStockModel {
	return NewTTradeStockModel(sqlx.NewSqlConnFromSession(session))
}
