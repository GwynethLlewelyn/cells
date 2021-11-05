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

package caddy

import (
	"net/url"
	"path/filepath"
	"strings"

	// "github.com/caddyserver/caddy/caddytls"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/pydio/cells/v4/common/caddy/maintenance"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/crypto/providers"
	"github.com/pydio/cells/v4/common/proto/install"
	"github.com/pydio/cells/v4/common/utils/statics"
)

var maintenanceDir string

// SiteConf wraps install.ProxyConfig with caddy-specific configurations
type SiteConf struct {
	*install.ProxyConfig
	// Parsed values from proto oneOf
	TLS     string
	TLSCert string
	TLSKey  string
	// Parsed External host if any
	ExternalHost string
	// Custom Root for this site
	WebRoot string
}

// Redirects returns computed redirects - This function is used inside templates!
func (s SiteConf) Redirects() map[string]string {
	rr := make(map[string]string)
	for _, bind := range s.GetBinds() {
		parts := strings.Split(bind, ":")
		var host, port string
		if len(parts) == 2 {
			host = parts[0]
			port = parts[1]
			if host == "" {
				continue
			}
		} else {
			host = bind
		}
		if port == "" {
			rr["http://"+host] = "https://" + host
		} else if port == "80" {
			continue
		} else {
			rr["http://"+host] = "https://" + host + ":" + port
		}
	}

	return rr
}

// SitesToCaddyConfigs computes all SiteConf from all *install.ProxyConfig by analyzing
// TLSConfig, ReverseProxyURL and Maintenance fields values
func SitesToCaddyConfigs(sites []*install.ProxyConfig) (caddySites []SiteConf, er error) {
	for _, proxyConfig := range sites {
		if bc, er := computeSiteConf(proxyConfig); er == nil {
			caddySites = append(caddySites, bc)
			/*
				// TODO V4 Enable these in caddy generated config
				if proxyConfig.HasTLS() && proxyConfig.GetLetsEncrypt() != nil {
					le := proxyConfig.GetLetsEncrypt()
					if le.AcceptEULA {
						caddytls.Agreed = true
					}
					if le.StagingCA {
						caddytls.DefaultCAUrl = common.DefaultCaStagingUrl
					} else {
						caddytls.DefaultCAUrl = common.DefaultCaUrl
					}
				}
			*/
		} else {
			return caddySites, er
		}
	}
	return caddySites, nil
}

// GetMaintenanceRoot provides a static root folder for serving maintenance page
func GetMaintenanceRoot() (string, error) {
	if maintenanceDir != "" {
		return maintenanceDir, nil
	}
	dir, err := statics.GetAssets("./maintenance/src")
	if err != nil {
		dir = filepath.Join(config.ApplicationWorkingDir(), "static", "maintenance")
		if _, _, err := statics.RestoreAssets(dir, maintenance.PydioMaintenanceBox, nil); err != nil {
			return "", errors.Wrap(err, "could not restore maintenance package")
		}
	}
	maintenanceDir = dir
	return dir, nil
}

func computeSiteConf(pc *install.ProxyConfig) (SiteConf, error) {
	bc := SiteConf{
		ProxyConfig: proto.Clone(pc).(*install.ProxyConfig),
	}
	if pc.ReverseProxyURL != "" {
		if u, e := url.Parse(pc.ReverseProxyURL); e == nil {
			bc.ExternalHost = u.Host
		}
	}
	if bc.TLSConfig == nil {
		for i, b := range bc.Binds {
			bc.Binds[i] = "http://" + b
		}
	} else {
		switch v := bc.TLSConfig.(type) {
		case *install.ProxyConfig_Certificate, *install.ProxyConfig_SelfSigned:
			certFile, keyFile, err := providers.LoadCertificates(pc)
			if err != nil {
				return bc, err
			}
			bc.TLSCert = certFile
			bc.TLSKey = keyFile
		case *install.ProxyConfig_LetsEncrypt:
			bc.TLS = v.LetsEncrypt.Email
		}
	}
	if bc.Maintenance {
		mDir, e := GetMaintenanceRoot()
		if e != nil {
			return bc, e
		}
		bc.WebRoot = mDir
	}
	return bc, nil
}
