package snowpack

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

func RunBuild() error {
	file, err := retrieveSnowpackFilePath()
	if err != nil {
		return err
	}
	cmd := exec.Command(file, "build")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "NODE_ENV=production")
	return cmd.Run()
}

func RunDev(ctx context.Context) error {
	file, err := retrieveSnowpackFilePath()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, file, "build")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "NODE_ENV=development")
	return cmd.Start()
}

func retrieveSnowpackFilePath() (string, error) {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "node_modules", ".bin", "snowpack")
	file, err := exec.LookPath(fp)
	return file, err
}
