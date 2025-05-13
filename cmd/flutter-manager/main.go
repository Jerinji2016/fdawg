package main

import (
    "log"
    "os"

    "github.com/Jerinji2016/fdawg/internal/commands"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "fdawg",
        Usage: "CLI tool to manage Flutter projects",
        Commands: []*cli.Command{
            commands.ServeCommand(),
            // More commands will be added here
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}