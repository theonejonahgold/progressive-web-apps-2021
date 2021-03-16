package prerender

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
)

// ScheduleBuild schedules the build task 5 minutes to every hour.
func ScheduleBuild() error {
	log.Println("Starting build scheduler...")
	ctx, stop := context.WithCancel(context.Background())
	go gracefullyShutDownOnSignal(stop)

	c := cron.New()
	_, err := c.AddFunc("55 * * * *", func() { Build() })
	if err != nil {
		return err
	}

	c.Start()
	<-ctx.Done()
	if err := ctx.Err(); err != context.Canceled {
		return err
	}

	ctx = c.Stop()
	<-ctx.Done()
	if err := ctx.Err(); err != context.Canceled {
		return err
	}
	return nil
}

func gracefullyShutDownOnSignal(stop context.CancelFunc) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	<-exit
	fmt.Println("Shutting down...")
	stop()
}
