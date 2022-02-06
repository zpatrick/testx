package docker_test

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/zpatrick/testx/assert"
	"github.com/zpatrick/testx/suite"
)

type MysqlSuite struct {
	containerID string
	db          *sql.DB
}

// var mysqlPort = cfg.Setting[int]{
// 	Default: func() int { return 3306 },
// 	Providers: []cfg.Provider[int]{
// 		cfg.EnvVar("APP_DB_PORT", strconv.Atoi),
// 	},
// }

func (m *MysqlSuite) Setup() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "failed to start docker client")
	}

	// TODO: pull image if not exists
	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        "mysql",
			ExposedPorts: nat.PortSet{"3306": struct{}{}},
			Env: []string{
				"MYSQL_ROOT_PASSWORD=pswd123",
				"MYSQL_DATABASE=users",
			},
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				nat.Port("3306"): {{HostIP: "127.0.0.1", HostPort: "3306"}},
			},
		},
		nil,
		nil,
		"mysql-test",
	)
	if err != nil {
		return errors.Wrap(err, "failed to create mysql docker container")
	}

	m.containerID = resp.ID
	if err := dockerClient.ContainerStart(ctx, m.containerID, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "failed to start mysql docker container")
	}

	for d := time.Millisecond * 500; ; d += time.Millisecond * 500 {
		conn, err := net.DialTimeout("tcp", "0.0.0.0:3306", time.Millisecond*50)
		if err != nil {
			if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "connection refused") {
				log.Printf("waiting for port 3306 to become available (error: %s)", err.Error())
				time.Sleep(d)
				continue
			}

			return errors.Wrap(err, "failed to dial 3306")
		}
		defer conn.Close()

		if _, _, err := bufio.NewReader(conn).ReadLine(); err != nil {
			if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "connection refused") {
				log.Printf("waiting for port 3306 to become available (error: %s)", err.Error())
				time.Sleep(d)
				continue
			}

			return errors.Wrap(err, "failed to read 3306")
		}

		break
	}

	cfg := mysql.Config{
		User:   "root",
		Passwd: "pswd123",
		Net:    "tcp",
		Addr:   "0.0.0.0:3306",
		DBName: "users",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return errors.Wrap(err, "failed to open db")
	}
	m.db = db

	for d := time.Second; ; d += time.Second {
		err := db.PingContext(ctx)
		if err == nil {
			break
		} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		} else {
			log.Printf("waiting for mysql ping to succeed (error: %s)", err.Error())
			time.Sleep(d)
		}
	}

	return nil
}

func (m *MysqlSuite) Teardown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if m.db != nil {
		if err := m.db.Close(); err != nil {
			return errors.Wrap(err, "failed to close db")
		}
	}

	if m.containerID == "" {
		return nil
	}

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "failed to start docker client")
	}

	log.Println("stopping mysql container")
	if err := dockerClient.ContainerStop(ctx, m.containerID, nil); err != nil {
		return errors.Wrap(err, "failed to stop mysql docker container")
	}

	log.Println("removing mysql container")
	if err := dockerClient.ContainerRemove(ctx, m.containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return errors.Wrap(err, "failed to remove mysql docker container")
	}

	return nil
}

func TestMain(m *testing.M) {
	if v := flag.Lookup("test.v"); v == nil || v.Value.String() != "true" {
		log.SetOutput(ioutil.Discard)
	}

	suite.Register(&MysqlSuite{})
	os.Exit(suite.Run(m))
}

func TestMysql_doSomeQuery(t *testing.T) {
	m := suite.Get[*MysqlSuite](t)

	ctx := context.Background()
	row := m.db.QueryRowContext(ctx, "SELECT DATABASE() FROM DUAL")
	assert.NilError(t, row.Err())

	var dbName string
	assert.NilError(t, row.Scan(&dbName))
	assert.Equal(t, dbName, "users")
}
