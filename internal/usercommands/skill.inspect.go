package usercommands

import (
	"errors"
	"fmt"

	"github.com/GoMudEngine/GoMud/internal/events"
	"github.com/GoMudEngine/GoMud/internal/items"
	"github.com/GoMudEngine/GoMud/internal/rooms"
	"github.com/GoMudEngine/GoMud/internal/skills"
	"github.com/GoMudEngine/GoMud/internal/templates"
	"github.com/GoMudEngine/GoMud/internal/users"
)

/*
Peep Skill
Level 1 - Reveals the type and value of items.
Level 2 - Reveals weapon damage or uses an item has left.
Level 3 - Reveals any stat modifiers an item has.
Level 4 - Reveals special magical properties like elemental effects.
*/
func Inspect(rest string, user *users.UserRecord, room *rooms.Room, flags events.EventFlag) (bool, error) {

	if user.Character.GetSkillLevel(skills.Inspect) == 0 {
		user.SendText("You don't know how to inspect.")
		return true, fmt.Errorf("you don't know how to inspect")
	}

	if len(rest) == 0 {
		user.SendText("Type `help inspect` for more information on the inspect skill.")
		return true, nil
	}

	skillLevel := user.Character.GetSkillLevel(skills.Inspect)

	// Check whether the user has an item in their inventory that matches
	matchItem, found := user.Character.FindInBackpack(rest)

	if !found {
		user.SendText(fmt.Sprintf("You don't have a %s to inspect. Is it still worn, perhaps?", rest))
	} else {

		if !user.Character.TryCooldown(skills.Inspect.String(), "3 rounds") {
			user.SendText(
				fmt.Sprintf("You need to wait %d more rounds to use that skill again.", user.Character.GetCooldown(skills.Inspect.String())),
			)
			return true, errors.New(`you're doing that too often`)
		}

		user.SendText(
			fmt.Sprintf(`You inspect the <ansi fg="item">%s</ansi>.`, matchItem.DisplayName()),
		)
		room.SendText(
			fmt.Sprintf(`<ansi fg="username">%s</ansi> inspects their <ansi fg="item">%s</ansi>...`, user.Character.Name, matchItem.DisplayName()),
			user.UserId,
		)

		type inspectDetails struct {
			InspectLevel int
			Item         *items.Item
			ItemSpec     *items.ItemSpec
		}

		iSpec := matchItem.GetSpec()

		details := inspectDetails{
			InspectLevel: skillLevel,
			Item:         &matchItem,
			ItemSpec:     &iSpec,
		}

		inspectTxt, _ := templates.Process("descriptions/inspect", details, user.UserId)
		user.SendText(inspectTxt)

	}

	return true, nil
}
