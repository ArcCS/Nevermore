package spells

type MobNPCSpells interface {
	ChangePlacement(place int) bool
	ApplyEffect(effectName string, length string, interval string,  effect func(), effectOff func())
	RemoveEffect(effectName string)
	ReceiveDamage(damage int) (int, int)
	ReceiveVitalDamage(damage int) int
	Heal(damage int) (int, int)
	HealVital(damage int)
	HealStam(damage int)
	RestoreMana(damage int)
}

var Effects = map[string]func(caller MobNPCSpells, target MobNPCSpells, modifier int) string {
	"HealStam": HealStam,
	"HealVit": HealVit,
	"HealBoth": HealBoth,
	"HealAll": HealAll,
	"FireDamage": FireDamage,
	"EarthDamage": EarthDamage,
	"AirDamage": AirDamage,
	"WaterDamage": WaterDamage,
	"Light": Light,
	"CurePoison": CurePoison,
	"Bless": Bless,
	"Protection": Protection,
	"Invisibility": Invisibility,
	"DetectInvisible": DetectInvisible,
	"Teleport": Teleport,
	"Stun": Stun,
	"Enchant": Enchant,
	"Recall": Recall,
	"Summon": Summon,
	"WizardWalk": WizardWalk,
	"Levitate": Levitate,
	"ResistFire": ResistFire,
	"ResistMagic": ResistMagic,
	"RemoveCurse": RemoveCurse,
	"ResistAir": ResistAir,
	"ResistWater": ResistWater,
	"ResistEarth": ResistEarth,
	"Clairvoyance": Clairvoyance,
	"RemoveDisease": RemoveDisease,
	"CureBlindness": CureBlindness,
	"Polymorph": Polymorph,
	"Attraction": Attraction,
	"InertialBarrier": InertialBarrier,
	"Surge": Surge,
	"ResistPoison": ResistPoison,
	"ResilientAura": ResilientAura,
	"ResistDisease": ResistDisease,
	"DisruptMagic": DisruptMagic,
	"Reflection": Reflection,
	"Dodge": Dodge,
	"ResistAcid": ResistAcid,
	"Embolden": Embolden,
}

func HealStam(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func HealVit(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func HealBoth(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func HealAll(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func FireDamage(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func EarthDamage(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func AirDamage(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func WaterDamage(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Light(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func CurePoison(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Bless(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Protection(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Invisibility(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func DetectInvisible(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Teleport(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Stun(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Enchant(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Recall(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Summon(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func WizardWalk(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Levitate(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistFire(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistMagic(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func RemoveCurse(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistAir(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistWater(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistEarth(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Clairvoyance(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func RemoveDisease(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func CureBlindness(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Polymorph(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Attraction(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func InertialBarrier(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Surge(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistPoison(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResilientAura(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistDisease(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func DisruptMagic(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Reflection(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Dodge(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func ResistAcid(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }

func Embolden(caller MobNPCSpells, target MobNPCSpells, modifier int) string{ return "" }


