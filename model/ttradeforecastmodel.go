package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TTradeForecastModel = (*customTTradeForecastModel)(nil)

type (
	// TTradeForecastModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTTradeForecastModel.
	TTradeForecastModel interface {
		tTradeForecastModel
		withSession(session sqlx.Session) TTradeForecastModel
	}

	customTTradeForecastModel struct {
		*defaultTTradeForecastModel
	}
)

// NewTTradeForecastModel returns a model for the database table.
func NewTTradeForecastModel(conn sqlx.SqlConn) TTradeForecastModel {
	return &customTTradeForecastModel{
		defaultTTradeForecastModel: newTTradeForecastModel(conn),
	}
}

func (m *customTTradeForecastModel) withSession(session sqlx.Session) TTradeForecastModel {
	return NewTTradeForecastModel(sqlx.NewSqlConnFromSession(session))
}
