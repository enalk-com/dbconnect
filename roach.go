package dbconnect

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type roachOps struct {
	ClusterName string `json:"cluster_name,omitempty" toml:"cluster_name,omitempty"`
	C           string `json:"c,omitempty" toml:"c,omitempty"`
}

type RoachConfig struct {
	ID              string   `json:"id,omitempty" toml:"id,omitempty"`
	Host            string   `json:"host,omitempty" toml:"host,omitempty"`
	Port            int      `json:"port,omitempty" toml:"port,omitempty"`
	User            string   `json:"user,omitempty" toml:"user,omitempty"`
	Pwd             string   `json:"pwd,omitempty" toml:"pwd,omitempty"`
	DB              string   `json:"db,omitempty" toml:"db,omitempty"`
	SSLMode         string   `json:"sslmode,omitempty" toml:"sslmode,omitempty"`                   // disable | require | verify-ca | verify-full
	ApplicationName string   `json:"application_name,omitempty" toml:"application_name,omitempty"` // in seconds
	SSLCert         string   `json:"sslcert,omitempty" toml:"sslcert,omitempty"`                   // location if PEM encoded cert file
	SSLKey          string   `json:"sslkey,omitempty" toml:"sslkey,omitempty"`                     // location of PEM encoded key file
	SSLRootCert     string   `json:"sslrootcert,omitempty" toml:"sslrootcert,omitempty"`           // location of PEM encoded root certificate file
	Options         roachOps `json:"options,omitempty" toml:"options,omitempty"`
	_db             *pgxpool.Pool
	once            sync.Once
}

func (rc *RoachConfig) assert() error {
	if rc.Host == "" || rc.User == "" || rc.DB == "" {
		return fmt.Errorf("invalid host, user or database name")
	}

	if rc.SSLMode != "require" &&
		rc.SSLMode != "disable" &&
		rc.SSLMode != "verify-ca" &&
		rc.SSLMode != "verify-full" &&
		rc.SSLMode != "allow" &&
		rc.SSLMode != "prefer" {
		return fmt.Errorf("invalid sslmode: %s", rc.SSLMode)
	}

	return nil
}

func (rc *RoachConfig) defaults() {
	if rc.Port == 0 {
		rc.Port = 26257
	}

	if rc.SSLMode == "" {
		rc.SSLMode = "verify-full"
	}
}

func (rc RoachConfig) connString() string {
	var auth string
	if rc.User != "" {
		auth = rc.User
		if rc.Pwd != "" {
			auth += ":" + rc.Pwd
		}
		auth += "@"
	}

	s := fmt.Sprintf("postgresql://%s%s:%d", auth, rc.Host, rc.Port)
	if rc.DB != "" {
		s += "/" + rc.DB
	}

	qps := []string{
		fmt.Sprintf("sslmode=%s", rc.SSLMode),
	}

	if rc.SSLCert != "" {
		qps = append(qps, fmt.Sprintf("sslcert=%s", rc.SSLCert))
	}

	if rc.SSLKey != "" {
		qps = append(qps, fmt.Sprintf("sslkey=%s", rc.SSLKey))
	}

	if rc.SSLRootCert != "" {
		qps = append(qps, fmt.Sprintf("sslrootcert=%s", rc.SSLRootCert))
	}

	if rc.ApplicationName != "" {
		qps = append(qps, fmt.Sprintf("application_name=%s", rc.ApplicationName))
	}

	var ops []string
	if rc.Options.ClusterName != "" {
		ops = append(ops, "--cluster_name="+rc.Options.ClusterName)
	}

	if rc.Options.C != "" {
		ops = append(ops, "-c "+rc.Options.C)
	}
	if len(ops) > 0 {
		qps = append(qps, fmt.Sprintf("options=%s", strings.Join(ops, " ")))
	}

	var qs string
	if len(qps) > 0 {
		qs = "?" + strings.Join(qps, "&")
	}

	cs := s + qs
	return cs
}

func (rc *RoachConfig) connect() error {
	var gerr error
	rc.once.Do(func() {
		if err := rc.assert(); err != nil {
			gerr = err
			return
		}
		rc.defaults()

		p, err := pgxpool.Connect(context.Background(), rc.connString())
		if err != nil {
			gerr = err
			return
		}
		rc._db = p
	})
	if gerr != nil {
		return gerr
	}
	return nil
}

func (rc *RoachConfig) db() (*pgxpool.Pool, error) {
	if err := rc.connect(); err != nil {
		return nil, fmt.Errorf("[RoachConfig.db] -> connect err: %s", err.Error())
	}
	if rc._db == nil {
		return nil, fmt.Errorf("invalid pq configuration. connection could not be made")
	}
	return rc._db, nil
}
