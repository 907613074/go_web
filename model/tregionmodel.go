package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TRegionModel = (*customTRegionModel)(nil)

type (
	// TRegionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTRegionModel.
	TRegionModel interface {
		tRegionModel
		withSession(session sqlx.Session) TRegionModel
	}

	customTRegionModel struct {
		*defaultTRegionModel
	}
)

// NewTRegionModel returns a model for the database table.
func NewTRegionModel(conn sqlx.SqlConn) TRegionModel {
	return &customTRegionModel{
		defaultTRegionModel: newTRegionModel(conn),
	}
}

func (m *customTRegionModel) withSession(session sqlx.Session) TRegionModel {
	return NewTRegionModel(sqlx.NewSqlConnFromSession(session))
}
