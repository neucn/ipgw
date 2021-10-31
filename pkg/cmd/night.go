package cmd

import (
	"errors"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/handler"
	"github.com/urfave/cli/v2"
	"time"
)

var (
	NightCommand = &cli.Command{
		Name:  "night",
		Usage: "auto login and logout at night",
		Action: func(ctx *cli.Context) error {
			var hasLogin bool = false
			const retryNum int = 5

			const freeStartHour = 0
			const loginDelayMinutes = 5

			const freeEndHour = 6
			const logoutAdvanceMinutes = 5

			console.InfoL("Please check if the time on the current computer is correct!!!")

			for {
				console.Info("\rThe current datetime is:", time.Now())
				//fmt.Println("\ntest")

				if (time.Now().Hour() >= freeStartHour && time.Now().Minute() >= loginDelayMinutes) && (time.Now().Hour() < freeEndHour && time.Now().Minute() < 60-logoutAdvanceMinutes) {
					for i := 0; !hasLogin && i < retryNum; i++ {
						err := loginUseDefaultAccount(ctx)
						if err == nil {
							hasLogin = true
							console.InfoL("\nLogin successfully at", time.Now().String())
						}
					}

				}

				if time.Now().Hour() >= freeEndHour-1 && time.Now().Minute() > 60-logoutAdvanceMinutes {
					for i := 0; hasLogin && i < retryNum; i++ {
						h := handler.NewIpgwHandler()
						connected, loggedIn := h.IsConnectedAndLoggedIn()
						if !connected {
							return errors.New("not in campus network")
						}
						if !loggedIn {
							console.InfoL("not logged in yet")
						}
						info := h.GetInfo()
						if err := h.Logout(); err != nil {
							console.InfoL("fail to logout account '%s':\n\t%v", info.Username, err)
						} else {
							hasLogin = false
							console.InfoL("\nLogout successfully at", time.Now().String())
						}

					}

				}
				time.Sleep(1 * time.Second)

			}
			return nil
		},
	}
)
