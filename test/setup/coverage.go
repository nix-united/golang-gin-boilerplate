package setup

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	// CoverDirEnv names the host directory that receives coverage data from the
	// application container. When it is empty, coverage collection is disabled
	// and the application image is built exactly like a release image.
	CoverDirEnv = "INTEGRATION_COVERDIR"

	containerCoverDir = "/coverage"
	coverBuildFlags   = "-cover -coverpkg=./internal/...,./cmd/..."

	// The binary flushes coverage counters only on a clean exit, so it must get
	// enough time to finish its graceful shutdown after SIGTERM.
	coverStopTimeout = 30 * time.Second

	coverDirPerm = 0o755
)

// coverage collects Go coverage data (GOCOVERDIR format) written by the
// application binary inside its container.
type coverage struct {
	hostDir string
}

// newCoverageFromEnv returns nil when coverage collection is disabled.
func newCoverageFromEnv() (*coverage, error) {
	dir := os.Getenv(CoverDirEnv)
	if dir == "" {
		return nil, nil
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve %s: %w", CoverDirEnv, err)
	}

	if err := os.MkdirAll(absDir, coverDirPerm); err != nil {
		return nil, fmt.Errorf("create coverage dir %s: %w", absDir, err)
	}

	return &coverage{hostDir: absDir}, nil
}

// apply builds the binary with coverage instrumentation and points GOCOVERDIR
// at a directory that exists inside the container.
func (c *coverage) apply(req *testcontainers.ContainerRequest) {
	flags := coverBuildFlags
	req.BuildArgs["GO_BUILD_FLAGS"] = &flags
	req.Env["GOCOVERDIR"] = containerCoverDir
	req.Files = append(req.Files, testcontainers.ContainerFile{
		Reader:            bytes.NewReader(nil),
		ContainerFilePath: containerCoverDir + "/.keep",
		FileMode:          0o644,
	})
}

// collect stops the container gracefully so the binary writes its counters,
// then copies the coverage files to the host directory. It works with remote
// Docker daemons (e.g. Docker-in-Docker) because it copies through the Docker
// API instead of a bind mount.
func (c *coverage) collect(ctx context.Context, container testcontainers.Container) error {
	timeout := coverStopTimeout
	if err := container.Stop(ctx, &timeout); err != nil {
		return fmt.Errorf("stop app container: %w", err)
	}

	state, err := container.State(ctx)
	if err != nil {
		return fmt.Errorf("get app container state: %w", err)
	}

	if state.ExitCode != 0 {
		return fmt.Errorf("app container exited with code %d; coverage data may be incomplete", state.ExitCode)
	}

	dockerClient, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("new docker client: %w", err)
	}
	defer dockerClient.Close()

	reader, _, err := dockerClient.CopyFromContainer(ctx, container.GetContainerID(), containerCoverDir)
	if err != nil {
		return fmt.Errorf("copy %s from app container: %w", containerCoverDir, err)
	}
	defer reader.Close()

	copied, err := c.extract(reader)
	if err != nil {
		return err
	}

	if copied == 0 {
		return errors.New("app binary wrote no coverage files")
	}

	return nil
}

// extract writes the regular files from a tar stream into the host directory,
// flattening paths. Coverage files are named covmeta.* and covcounters.*, and
// the counters file names are unique per process, so runs never overwrite each other.
func (c *coverage) extract(r io.Reader) (int, error) {
	tarReader := tar.NewReader(r)
	copied := 0

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return copied, nil
		}
		if err != nil {
			return copied, fmt.Errorf("read coverage archive: %w", err)
		}

		name := path.Base(header.Name)
		if header.Typeflag != tar.TypeReg || name == ".keep" {
			continue
		}

		if err := c.writeFile(name, tarReader); err != nil {
			return copied, err
		}

		copied++
	}
}

func (c *coverage) writeFile(name string, r io.Reader) error {
	file, err := os.Create(filepath.Join(c.hostDir, name))
	if err != nil {
		return fmt.Errorf("create coverage file %s: %w", name, err)
	}

	//nolint:gosec // coverage files are small and produced by our own binary
	if _, err := io.Copy(file, r); err != nil {
		return errors.Join(fmt.Errorf("write coverage file %s: %w", name, err), file.Close())
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close coverage file %s: %w", name, err)
	}

	return nil
}
