package commands

import (
    "github.com/Jerinji2016/fdawg/internal/server"
    "github.com/urfave/cli/v2"
)

// ServeCommand returns the CLI command for starting the web server
func ServeCommand() *cli.Command {
    return &cli.Command{
        Name:  "serve",
        Usage: "Start a web server",
        Action: func(c *cli.Context) error {
            return server.Start(c.String("port"))
        },
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:  "port",
                Value: "8080",
                Usage: "Port to run the server on",
            },
        },
    }
}