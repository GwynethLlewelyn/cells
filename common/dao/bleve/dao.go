/*
 * Copyright (c) 2019-2021. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

// Package bleve provides tools for using Bolt as a standard persistence layer for services.
package bleve

import (
	"github.com/pydio/cells/v4/common/dao"
	"github.com/pydio/cells/v4/common/utils/configx"
)

// DAO defines the functions specific to the boltdb dao
type DAO interface {
	dao.DAO
	BleveConfig() *dao.BleveConfig
}

// Handler for the main functions of the DAO
type Handler struct {
	dao.DAO
}

// NewDAO creates a new handler for the boltdb dao
func NewDAO(driver string, dsn string, prefix string) *Handler {
	conn, err := dao.NewConn(driver, dsn)
	if err != nil {
		return nil
	}
	return &Handler{
		DAO: dao.NewDAO(conn, driver, prefix),
	}
}

// Init initialises the handler
func (h *Handler) Init(configx.Values) error {
	return nil
}

// BleveConfig returns the folder to lookup for bleve index
func (h *Handler) BleveConfig() *dao.BleveConfig {
	if h == nil {
		return &dao.BleveConfig{}
	}

	if conn := h.GetConn(); conn != nil {
		return conn.(*dao.BleveConfig)
	}
	return &dao.BleveConfig{}
}
