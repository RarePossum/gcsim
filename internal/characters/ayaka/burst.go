package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 104

func init() {
	burstFrames = frames.InitAbilSlice(125)
	burstFrames[action.ActionAttack] = 124
	burstFrames[action.ActionDash] = 124
	burstFrames[action.ActionJump] = 114
	burstFrames[action.ActionSwap] = 123
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Soumetsu",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
	}

	//5 second, 20 ticks, so once every 15 frames, bloom after 5 seconds
	ai.Mult = burstBloom[c.TalentLvlBurst()]
	ai.Abil = "Soumetsu (Bloom)"
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), burstHitmark, burstHitmark+300, c.c4)

	// C2 mini-frostflake bloom
	var aiC2 combat.AttackInfo
	if c.Base.Cons >= 2 {
		aiC2 = ai
		aiC2.Mult = burstBloom[c.TalentLvlBurst()] * .2
		aiC2.Abil = "C2 Mini-Frostflake Seki no To (Bloom)"
		// TODO: Not sure about the positioning/size...
		c.Core.QueueAttack(aiC2, combat.NewDefCircHit(2, false, combat.TargettableEnemy), burstHitmark, burstHitmark+300, c.c4)
		c.Core.QueueAttack(aiC2, combat.NewDefCircHit(2, false, combat.TargettableEnemy), burstHitmark, burstHitmark+300, c.c4)
	}

	for i := 0; i < 19; i++ {
		ai.Mult = burstCut[c.TalentLvlBurst()]
		ai.Abil = "Soumetsu (Cutting)"
		c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), burstHitmark, burstHitmark+i*15, c.c4)

		// C2 mini-frostflake cutting
		if c.Base.Cons >= 2 {
			aiC2.Mult = burstCut[c.TalentLvlBurst()] * .2
			aiC2.Abil = "C2 Mini-Frostflake Seki no To (Cutting)"
			// TODO: Not sure about the positioning/size...
			c.Core.QueueAttack(aiC2, combat.NewDefCircHit(2, false, combat.TargettableEnemy), burstHitmark, burstHitmark+i*15, c.c4)
			c.Core.QueueAttack(aiC2, combat.NewDefCircHit(2, false, combat.TargettableEnemy), burstHitmark, burstHitmark+i*15, c.c4)
		}
	}

	c.ConsumeEnergy(8)
	c.SetCD(action.ActionBurst, 20*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,

		State: action.BurstState,
	}
}
