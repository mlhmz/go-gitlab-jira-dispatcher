package auth

import (
	"github.com/gofiber/fiber/v2/log"

	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
	"github.com/urfave/cli/v2"
)

func GetCli(userStore store.UserStore, login Login) *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Manage the users of the application",
		Subcommands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a user with the name and password",
				Action: func(ctx *cli.Context) error {
					user := ctx.Args().Get(0)
					password := ctx.Args().Get(1)
					var hashedPassword string
					if err := login.CreatePassword(&password, &hashedPassword); err != nil {
						return err
					}
					err := userStore.CreateUser(&store.User{
						Username:     user,
						PasswordHash: hashedPassword,
					})
					if err != nil {
						return err
					}
					log.Infof("Created user '%s'", user)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a user by the username",
				Action: func(ctx *cli.Context) error {
					user := ctx.Args().Get(0)
					if err := userStore.DeleteUser(user); err != nil {
						return err
					}
					log.Infof("Deleted user '%s'", user)
					return nil
				},
			},
		},
	}
}
