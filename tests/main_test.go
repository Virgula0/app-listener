package tests

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	amd64Bin string
	arm64Bin string
)

func TestMain(m *testing.M) {
	if err := os.MkdirAll("../build/test", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir build/test: %v\n", err)
		os.Exit(1)
	}

	amd64Bin = absPath("../build/test/app-listener-amd64")

	cmd := exec.Command("go", "build", "-o", amd64Bin, "..")
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build amd64 binary: %v\n", err)
		os.Exit(1)
	}

	arm64Bin = absPath("../build/test/app-listener-arm64")
	cmd2 := exec.Command("go", "build", "-o", arm64Bin, "..")
	cmd2.Env = append(os.Environ(), "GOOS=linux", "GOARCH=arm64", "CGO_ENABLED=0")
	cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "arm64 cross-build: %v (will skip arm64 tests)\n", err)
		arm64Bin = ""
	}

	os.Exit(m.Run())
}

type IntegrationSuite struct {
	suite.Suite
	ctx context.Context
}

func TestIntegrationSuite(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration tests")
	}
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	s.ctx = context.Background()
}

func (s *IntegrationSuite) startContainer(
	image string,
	platform string,
	privileged bool,
	binaryPath string,
) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:         image,
		ImagePlatform: platform,
		Cmd:           []string{"sleep", "3600"},
		WaitingFor:    wait.ForLog("").WithPollInterval(100 * time.Millisecond).WithStartupTimeout(60 * time.Second),
	}

	if binaryPath != "" {
		req.Files = []testcontainers.ContainerFile{
			{
				HostFilePath:      binaryPath,
				ContainerFilePath: "/app-listener",
				FileMode:          0755,
			},
		}
	}

	if privileged {
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.Privileged = true
			hc.CapAdd = []string{"BPF", "SYS_ADMIN"}
		}
		req.Mounts = testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{HostPath: "/sys/kernel/btf"},
				Target: "/sys/kernel/btf",
			},
		}
	}

	c, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err, "start container %s/%s", image, platform)
	return c
}

func (s *IntegrationSuite) exec(c testcontainers.Container, cmd []string) (int, string) {
	exitCode, reader, err := c.Exec(s.ctx, cmd)
	s.Require().NoError(err, "exec: %v", cmd)

	b, err := io.ReadAll(reader)
	s.Require().NoError(err, "read exec output: %v", cmd)

	return exitCode, strings.TrimSpace(string(b))
}

func (s *IntegrationSuite) logs(c testcontainers.Container) string {
	r, err := c.Logs(s.ctx)
	s.Require().NoError(err)
	defer r.Close()

	b, err := io.ReadAll(r)
	s.Require().NoError(err)
	return string(b)
}

func absPath(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	return abs
}
