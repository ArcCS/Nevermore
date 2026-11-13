package config

import "math"

// Bards may cast while singing. Below BardSingCastFreeTier the spell's
// damage or healing is cut by BardSingCastPenalty; from that tier on the
// song costs the spell nothing.
var (
	BardSingCastPenalty  = .2
	BardSingCastFreeTier = 7
)

// BardSingCastMod is the multiplier applied to a cast spell's damage or
// healing. It is below 1 only for a low-tier bard who is singing.
func BardSingCastMod(class int, tier int, singing bool) float64 {
	if class != BARD || !singing || tier >= BardSingCastFreeTier {
		return 1
	}
	return 1 - BardSingCastPenalty
}

// Song timing (2026-11-08: was 45 + 10*con seconds, which let a high-con bard
// outlast the cooldown, and a tick of 8 - tier/5 seconds, which compounded
// with pie and instrument scaling). A song now lasts SongBaseDuration plus
// SongDurationPerTier a tier, always short of SongCooldown, and pulses every
// SongTickSlow seconds until SongFastTickTier, then every SongTickFast.
var (
	SongBaseDuration    = 180
	SongDurationPerTier = 2
	SongCooldown        = 300
	SongTickSlow        = 8
	SongTickFast        = 6
	SongFastTickTier    = 15
	// SongMixCooldown is the wait between switches of a running song with
	// sing mix. It is longer than any song, so a bard can mix once a song.
	SongMixCooldown = 300
	// SongStopRefund is the share of the song cooldown handed back for the
	// share of the song left unsung when a bard stops early.
	SongStopRefund = .9
)

// SongDuration is how many seconds a song sung at the given tier lasts.
func SongDuration(tier int) int {
	return SongBaseDuration + SongDurationPerTier*tier
}

// SongTickSeconds is how often a song sung at the given tier pulses.
func SongTickSeconds(tier int) int {
	if tier >= SongFastTickTier {
		return SongTickFast
	}
	return SongTickSlow
}

// SongStopRefundSeconds is how much of the song cooldown a bard gets back for
// stopping a song of length seconds with remaining seconds still to go:
// SongStopRefund of the cooldown, scaled by the unsung share of the song.
func SongStopRefundSeconds(remaining float64, length float64) float64 {
	if length <= 0 || remaining <= 0 {
		return 0
	}
	if remaining > length {
		remaining = length
	}
	return SongStopRefund * remaining / length * float64(SongCooldown)
}

// Draen's Tale restores DraensRollScale of the instrument roll plus pie over
// DraensPieDivisor in mana to everyone present on each pulse (2026-11-08: was
// the bard's full pie per pulse).
var (
	DraensRollScale  = .1
	DraensPieDivisor = 2.0
)

// DraensTaleMana is the mana one pulse of Draen's Tale restores, for an
// instrument roll (adjustment included) and the bard's piety.
func DraensTaleMana(roll int, pie int) int {
	return max(0, int(math.Round(DraensRollScale*float64(roll)+float64(pie)/DraensPieDivisor)))
}

// Celebration Night heals each player present for CelebrationHealScale of the
// instrument roll, split over the number of players raised to
// CelebrationDivisorExp, with the divisor capped at CelebrationDivisorCap.
// Each share shrinks as the room fills but the whole heal grows (2026-11-08:
// was the roll split evenly, so the total never grew).
var (
	CelebrationHealScale  = .85
	CelebrationDivisorExp = 2.0 / 3
	CelebrationDivisorCap = 3.0
)

// CelebrationHeal is what one pulse of Celebration Night heals each player for,
// given an instrument roll (adjustment included) and the characters present.
func CelebrationHeal(roll int, players int) int {
	if players < 1 {
		players = 1
	}
	divisor := math.Min(CelebrationDivisorCap, math.Pow(float64(players), CelebrationDivisorExp))
	return max(0, int(math.Round(CelebrationHealScale*float64(roll)/divisor)))
}
