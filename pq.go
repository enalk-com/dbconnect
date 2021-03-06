package dbconnect

import (
	// "database/sql"
	"context"
	"fmt"
	"os"
	"sync"

	// pq is imported to allow sql connection using driver 'postgres'
	// _ "github.com/lib/pq"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PQConfig defines all the parameters to be used for establishing
// a postgresql connection
type PQConfig struct {
	ID              string `json:"id,omitempty" toml:"id,omitempty"`
	Host            string `json:"host,omitempty" toml:"host,omitempty"`
	Port            int    `json:"port,omitempty" toml:"port,omitempty"`
	User            string `json:"user,omitempty" toml:"user,omitempty"`
	Pwd             string `json:"pwd,omitempty" toml:"pwd,omitempty"`
	DB              string `json:"db,omitempty" toml:"db,omitempty"`
	SSLMode         string `json:"sslmode,omitempty" toml:"sslmode,omitempty"` // disable | require | verify-ca | verify-full
	FallbackAppName string `json:"fallback_application_name,omitempty" toml:"fallback_application_name,omitempty"`
	ConnectTimeout  int    `json:"connect_timeout,omitempty" toml:"connect_timeout,omitempty"` // in seconds
	SSLCert         string `json:"sslcert,omitempty" toml:"sslcert,omitempty"`                 // location if PEM encoded cert file
	SSLKey          string `json:"sslkey,omitempty" toml:"sslkey,omitempty"`                   // location of PEM encoded key file
	SSLRootCert     string `json:"sslrootcert,omitempty" toml:"sslrootcert,omitempty"`         // location of PEM encoded root certificate file
	_db             *pgxpool.Pool
	once            sync.Once
}

func (pc *PQConfig) assert() error {
	if pc.Host == "" || pc.User == "" || pc.DB == "" {
		return fmt.Errorf("invalid host, user or database name")
	}

	if pc.SSLMode != "" && pc.SSLMode != "disable" && pc.SSLMode != "verify-ca" && pc.SSLMode != "verify-full" {
		return fmt.Errorf("invalid sslmode: %s", pc.SSLMode)
	}

	return nil
}

func (pc *PQConfig) expandEnv() {
	pc.ID = os.ExpandEnv(pc.ID)
	pc.Host = os.ExpandEnv(pc.Host)
	pc.User = os.ExpandEnv(pc.User)
	pc.Pwd = os.ExpandEnv(pc.Pwd)
	pc.DB = os.ExpandEnv(pc.DB)
	pc.SSLMode = os.ExpandEnv(pc.SSLMode)
	pc.FallbackAppName = os.ExpandEnv(pc.FallbackAppName)
	pc.SSLCert = os.ExpandEnv(pc.SSLCert)
	pc.SSLKey = os.ExpandEnv(pc.SSLKey)
	pc.SSLRootCert = os.ExpandEnv(pc.SSLRootCert)
}

func (pc *PQConfig) defaults() {
	if pc.Port == 0 {
		pc.Port = 5432
	}

	if pc.SSLMode == "" {
		pc.SSLMode = "disable"
	}
}

func (pc *PQConfig) connect() error {
	var gerr error
	pc.once.Do(func() {
		pc.expandEnv()
		if err := pc.assert(); err != nil {
			gerr = err
			return
		}

		pc.defaults()

		connStr := fmt.Sprintf(
			"user=%s dbname=%s host=%s port=%d sslmode=%s",
			pc.User, pc.DB, pc.Host, pc.Port, pc.SSLMode,
		)
		if pc.Pwd != "" {
			connStr = fmt.Sprintf("%s password=%s", connStr, pc.Pwd)
		}

		if pc.FallbackAppName != "" {
			connStr = fmt.Sprintf("%s fallback_application_name=%s", connStr,
				pc.FallbackAppName)
		}

		if pc.ConnectTimeout > 0 {
			connStr = fmt.Sprintf("%s connect_timeout=%d", connStr,
				pc.ConnectTimeout)
		}

		if pc.SSLCert != "" {
			connStr = fmt.Sprintf("%s sslcert=%s", connStr,
				pc.SSLCert)
		}

		if pc.SSLKey != "" {
			connStr = fmt.Sprintf("%s sslkey=%s", connStr,
				pc.SSLKey)
		}

		if pc.SSLRootCert != "" {
			connStr = fmt.Sprintf("%s sslrootcert=%s", connStr,
				pc.SSLRootCert)
		}

		p, err := pgxpool.Connect(context.Background(), connStr)
		// db, err := sql.Open("postgres", connStr)
		if err != nil {
			gerr = err
			return
		}
		pc._db = p
	})
	if gerr != nil {
		return gerr
	}
	return nil
}

func (pc *PQConfig) db() (*pgxpool.Pool, error) {
	if err := pc.connect(); err != nil {
		return nil, fmt.Errorf("[PQConfig.db] -> connect err: %s", err.Error())
	}
	if pc._db == nil {
		return nil, fmt.Errorf("invalid pq configuration. connection could not be made")
	}
	return pc._db, nil
}
