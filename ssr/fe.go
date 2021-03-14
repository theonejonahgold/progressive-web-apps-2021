package ssr

import (
	"context"
	"os"
	"os/exec"

	u "github.com/theonejonahgold/pwa/utils"
)

func runSnowpackDevBuilds(ctx context.Context) error {
	file, err := u.RetrieveSnowpackFilePath()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, file, "build")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
