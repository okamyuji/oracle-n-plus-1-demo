package config

import (
	"testing"

	"github.com/sijms/go-ora/v3/configurations"
)

// go-ora v3.0.1 の FastLogin は、2 本目以降の接続で Oracle 23c に EOF で切られる（sijms/go-ora#732）
func TestOracleDSN_DisablesFastLogin(t *testing.T) {
	c := &Config{DBUsername: "u", DBPassword: "p", DBHost: "localhost", DBPort: 1521, DBServiceName: "FREEPDB1"}

	cfg, err := configurations.ParseConfig(OracleDSN(c))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FastLogin {
		t.Fatal("FastLogin must be disabled for go-ora v3.0.1")
	}
}

// パスワードに URL の区切り文字があっても、ホストやサービス名が変わらない
func TestOracleDSN_EscapesCredentials(t *testing.T) {
	c := &Config{DBUsername: "u", DBPassword: "p@ss/w?rd#1", DBHost: "db.local", DBPort: 1521, DBServiceName: "FREEPDB1"}

	cfg, err := configurations.ParseConfig(OracleDSN(c))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Password != "p@ss/w?rd#1" || cfg.Servers[0].Addr != "db.local" || cfg.ServiceName != "FREEPDB1" {
		t.Fatalf("got password=%q host=%q service=%q", cfg.Password, cfg.Servers[0].Addr, cfg.ServiceName)
	}
}
