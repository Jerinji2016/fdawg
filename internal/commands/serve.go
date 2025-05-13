package commands

import (
    "github.com/Jerinji2016/fdawg/internal/server"
    "github.com/Jerinji2016/fdawg/pkg/flutter"
    "github.com/Jerinji2016/fdawg/pkg/utils"
    "github.com/urfave/cli/v2"
)

// ServeCommand returns the CLI command for starting the web server
func ServeCommand() *cli.Command {
    return &cli.Command{
        Name:        "serve",
        Usage:       "Start a web server for Flutter project management",
        Description: "Starts a web server for the specified Flutter project. If no directory is provided, the current directory will be used.",
        ArgsUsage:   "[directory]",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "port",
                Aliases: []string{"p"},
                Value:   "8080",
                Usage:   "Port to run the server on",
                EnvVars: []string{"FDAWG_PORT"},
            },
        },
        Action: func(c *cli.Context) error {
            // Get directory from arguments or use current directory
            dir := "."
            if c.Args().Len() > 0 {
                dir = c.Args().Get(0)
            }
            
            // First validate that it's a Flutter project
            utils.Info("Validating Flutter project before starting server...")
            
            result, err := flutter.ValidateProject(dir)
            if err != nil {
                utils.Error("Cannot start server: %v", err)
                return err
            }
            
            if len(result.MissingDirs) > 0 {
                utils.Warning("Some Flutter directories are missing:")
                for _, missingDir := range result.MissingDirs {
                    utils.Log("  - %s", missingDir)
                }
                utils.Warning("This might be a partial or incomplete Flutter project.")
            }
            
            // Get the port from the flag
            port := c.String("port")
            
            utils.Success("Valid Flutter project detected. Starting server for: %s", result.PubspecInfo.Name)
            utils.Info("Project path: %s", result.ProjectPath)
            utils.Info("Server will run on port: %s", port)
            
            // Start the server with the project info
            return server.Start(port, result)
        },
    }
}
