package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ondrovic/common/types"
	"github.com/ondrovic/common/utils"
	"github.com/pterm/pterm"
)

var (
	application types.Application
	version     string = "dev"
)

func init() {

	name := "Common"
	desc := "Share common types and utils"
	usage := fmt.Sprintf("%s [root-directory]", name)

	application = types.Application{
		Name:        name,
		Description: desc,
		Style: types.Styles{
			Color: types.Colors{
				Background: pterm.BgDarkGray,
				Foreground: pterm.FgWhite,
			},
		},
		Usage:   usage,
		Version: version,
	}

}

func display(item interface{}) {
	fmt.Println(item)
}

func main() {
	err := utils.ClearTerminalScreen(runtime.GOOS)
	if err != nil {
		fmt.Print(err)
	}
	display(fmt.Sprintf("Name: %s", application.Name))
	display(fmt.Sprintf("Desc: %s", application.Description))
	display(fmt.Sprintf("Usage: %s [root-directory]", strings.ToLower(application.Name)))
	display(fmt.Sprintf("Version: %s", application.Version))

	err = utils.ApplicationBanner(&application, utils.ClearTerminalScreen)
	if err != nil {
		fmt.Print(err)
	}
}
