package usercommands

import (
	"fmt"

	"github.com/GoMudEngine/GoMud/internal/buffs"
	"github.com/GoMudEngine/GoMud/internal/characters"
	"github.com/GoMudEngine/GoMud/internal/events"
	"github.com/GoMudEngine/GoMud/internal/rooms"
	"github.com/GoMudEngine/GoMud/internal/scripting"
	"github.com/GoMudEngine/GoMud/internal/skills"
	"github.com/GoMudEngine/GoMud/internal/spells"
	"github.com/GoMudEngine/GoMud/internal/users"
)

/*
Protection Skill
Level 1 - Aid (revive) a player
Level 3 - Aid (revive) a player, even during combat
*/
func Aid(rest string, user *users.UserRecord, room *rooms.Room, flags events.EventFlag) (bool, error) {

	skillLevel := user.Character.GetSkillLevel(skills.Protection)

	if skillLevel == 0 {
		user.SendText("You don't know how to provide aid.")
		return true, fmt.Errorf("you don't know how to provide aid")
	}

	if skillLevel < 3 && !room.IsCalm() {
		user.SendText("You can only do that in calm rooms!")
		return true, nil
	}

	aidPlayerId, _ := room.FindByName(rest, rooms.FindDowned)

	if aidPlayerId == user.UserId {
		aidPlayerId = 0
	}

	if aidPlayerId > 0 {

		p := users.GetByUserId(aidPlayerId)

		if p != nil {

			if p.Character.Health > 0 {
				user.SendText(fmt.Sprintf(`<ansi fg="username">%s</ansi> is not in need of aid!`, p.Character.Name))
				return true, nil
			}

			if user.Character.Aggro != nil {
				user.SendText("You are too busy to aid anyone!")
				return true, nil
			}

			// Set spell Aid
			spellAggro := characters.SpellAggroInfo{
				SpellId:              `aidskill`,
				SpellRest:            ``,
				TargetUserIds:        []int{aidPlayerId},
				TargetMobInstanceIds: []int{},
			}

			continueCasting := true
			if handled, err := scripting.TrySpellScriptEvent(`onCast`, user.UserId, 0, spellAggro); err == nil {
				continueCasting = handled
			}

			if continueCasting {
				spellInfo := spells.GetSpell(`aidskill`)
				user.Character.CancelBuffsWithFlag(buffs.Hidden)
				user.Character.SetCast(spellInfo.WaitRounds, spellAggro)
			}

		}

		return true, nil
	}

	user.SendText("Aid whom?")
	return true, nil
}
