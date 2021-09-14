package dbconnect

import (
	"os"
	"strconv"
	"testing"
)

func TestRoachConnect(t *testing.T) {
	type tt struct {
		name string
		rc   *RoachConfig
		err  error
	}

	host := os.Getenv("DBC_TEST_ROACH_HOST")
	p := os.Getenv("DBC_TEST_ROACH_PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		port = 26257
	}
	user := os.Getenv("DBC_TEST_ROACH_USER")
	pwd := os.Getenv("DBC_TEST_ROACH_PWD")
	db := os.Getenv("DBC_TEST_ROACH_DB")
	sslmode := os.Getenv("DBC_TEST_ROACH_SSLMODE")
	sslcert := os.Getenv("DBC_TEST_ROACH_SSLCERT")
	sslkey := os.Getenv("DBC_TEST_ROACH_SSLKEY")
	sslrootcert := os.Getenv("DBC_TEST_ROACH_SSLROOTCERT")

	tsts := []tt{
		{
			name: "valid",
			rc: &RoachConfig{
				ID:          "roachtest",
				Host:        host,
				Port:        port,
				User:        user,
				Pwd:         pwd,
				DB:          db,
				SSLMode:     sslmode,
				SSLCert:     sslcert,
				SSLKey:      sslkey,
				SSLRootCert: sslrootcert,
			},
		},
	}

	for _, tst := range tsts {
		t.Run(tst.name, func(t *testing.T) {
			if err := tst.rc.connect(); err != nil {
				if tst.err == nil {
					t.Fatal(err)
				}
			} else if tst.err != nil {
				t.Fatal("was supposed to error")
			}
			if tst.err == nil && tst.rc._db == nil {
				t.Fatal("_db instance is empty")
			}
		})
	}
}
