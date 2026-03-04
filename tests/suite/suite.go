package suite

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"

	ssov1 "github.com/Kefir4c/protos_sso/gen/go/sso"
	"github.com/Kefir4c/sso-service/internal/config"
	"github.com/Kefir4c/sso-service/tests/testdata"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/lib/pq"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
	db         *sql.DB
}

var (
	once = sync.Once{}
)

func configPath() string {
	dir, _ := os.Getwd()
	fmt.Printf("DEBUG: current directory: %s\n", dir)

	if v := os.Getenv("CONFIG_PATH"); v != "" {
		fmt.Printf("DEBUG: CONFIG_PATH=%s\n", v)
		if _, err := os.Stat(v); err == nil {
			fmt.Println("DEBUG: file exists")
			return v
		} else {
			fmt.Printf("DEBUG: file NOT exists: %v\n", err)
		}
	}

	defaultPath := "./config/local.yaml"
	fmt.Printf("DEBUG: trying default path: %s\n", defaultPath)
	if _, err := os.Stat(defaultPath); err == nil {
		fmt.Println("DEBUG: default file exists")
		return defaultPath
	} else {
		fmt.Printf("DEBUG: default file NOT exists: %v\n", err)
	}

	panic("config file not found")
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.Storage.Username,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "failed to connet to db")

	err = db.Ping()
	require.NoError(t, err, "db is not reachable")

	once.Do(func() {
		err = testdata.Seed(db)
		require.NoError(t, err, "failed to seed test db")
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	host := "app"

	clientConn, err := grpc.NewClient(
		net.JoinHostPort(host, strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "grpc server connection failed")

	t.Cleanup(func() {
		clientConn.Close()
	})

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(clientConn),
	}

}
