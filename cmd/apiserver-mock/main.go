// Package main entrypoint for the application.
package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"

	"github.com/christgf/env"
	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/skpr/api/internal/model"
	"github.com/skpr/api/internal/server/mock/backup"
	"github.com/skpr/api/internal/server/mock/compass"
	"github.com/skpr/api/internal/server/mock/config"
	"github.com/skpr/api/internal/server/mock/cron"
	"github.com/skpr/api/internal/server/mock/environment"
	"github.com/skpr/api/internal/server/mock/events"
	"github.com/skpr/api/internal/server/mock/project"
	"github.com/skpr/api/internal/server/mock/purge"
	"github.com/skpr/api/internal/server/mock/restore"
	"github.com/skpr/api/internal/server/mock/version"
	"github.com/skpr/api/pb"
)

const (
	// Orange which aligns with the Skpr theme.
	Orange = "#EE5622"
)

const cmdExample = `
  # Run the mock API server on port 443
  apiserver-mock --port 443`

const cmdLong = `  __  __  ___   ____ _  __     _    ____ ___ ____  _____ ______     _______ ____  
 |  \/  |/ _ \ / ___| |/ /    / \  |  _ \_ _/ ___|| ____|  _ \ \   / / ____|  _ \ 
 | |\/| | | | | |   | ' /    / _ \ | |_) | |\___ \|  _| | |_) \ \ / /|  _| | |_) |
 | |  | | |_| | |___| . \   / ___ \|  __/| | ___) | |___|  _ < \ V / | |___|  _ < 
 |_|  |_|\___/ \____|_|\_\ /_/   \_\_|  |___|____/|_____|_| \_\ \_/  |_____|_| \_\
                                                                                                                                             
                                                                                    
Mock implementation of the Skpr API.`

// Options for the CLI.
type Options struct {
	Port int
}

func main() {
	o := Options{}

	globalModel := model.NewModel()
	globalModel.CreateEnvironment("dev", 1)
	globalModel.CreateEnvironment("stg", 2)
	globalModel.CreateEnvironment("prod", 4)

	cmd := &cobra.Command{
		Use:     "apiserver-mock",
		Short:   "Mock implementation of the Skpr API.",
		Long:    cmdLong,
		Example: cmdExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			server := grpc.NewServer()

			log.Println("Registering service: Backup")
			pb.RegisterBackupServer(server, &backup.Server{
				Model: globalModel,
			})

			log.Println("Registering service: Compass")
			pb.RegisterCompassServer(server, &compass.Server{})

			log.Println("Registering service: Config")
			pb.RegisterConfigServer(server, &config.Server{
				Model: globalModel,
			})

			log.Println("Registering service: Cron")
			pb.RegisterCronServer(server, &cron.Server{
				Model: globalModel,
			})

			log.Println("Registering service: Environments")
			pb.RegisterEnvironmentServer(server, &environment.Server{
				Model: globalModel,
			})

			log.Println("Registering service: Events")
			pb.RegisterEventsServer(server, &events.Server{})

			log.Println("Registering service: Project")
			pb.RegisterProjectServer(server, &project.Server{})

			log.Println("Registering service: Purge")
			pb.RegisterPurgeServer(server, &purge.Server{
				Model: globalModel,
			})

			log.Println("Registering service: Restore")
			pb.RegisterRestoreServer(server, &restore.Server{
				Model: globalModel,
			})

			log.Println("Registering service: Version")
			pb.RegisterVersionServer(server, &version.Server{})

			reflection.Register(server)

			log.Println("Starting server on port:", o.Port)

			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", o.Port))
			if err != nil {
				return fmt.Errorf("failed to setup listener: %w", err)
			}

			return server.Serve(listener)
		},
	}

	cmd.PersistentFlags().IntVar(&o.Port, "port", env.Int("APISERVER_MOCK_PORT", 50051), "Port for the API server")

	cobra.AddTemplateFunc("StyleHeading", func(data string) string {
		return gchalk.WithHex(Orange).Bold(data)
	})

	usageTemplate := cmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
	).Replace(usageTemplate)

	re := regexp.MustCompile(`(?m)^Flags:\s*$`)
	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)
	cmd.SetUsageTemplate(usageTemplate)

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
