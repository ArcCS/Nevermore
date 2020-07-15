package spells

type MobNPCSpells interface {
	ChangePlacement(place int) bool
	ApplyEffect(effect string)
	RemoveEffect(effect string)
	ReceiveDamage(damage int) (int, int)
	ReceiveVitalDamage(damage int) int
	Heal(damage int) (int, int)
	HealVital(damage int)
	HealStam(damage int)
	RestoreMana(damage int)
}

var Effects = map[string]func(caller MobNPCSpells, target MobNPCSpells, modifier int){
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

func HealStam(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func HealVit(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func HealBoth(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func HealAll(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func FireDamage(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func EarthDamage(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func AirDamage(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func WaterDamage(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Light(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func CurePoison(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Bless(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Protection(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Invisibility(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func DetectInvisible(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Teleport(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Stun(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Enchant(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Recall(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Summon(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func WizardWalk(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Levitate(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistFire(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistMagic(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func RemoveCurse(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistAir(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistWater(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistEarth(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Clairvoyance(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func RemoveDisease(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func CureBlindness(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Polymorph(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Attraction(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func InertialBarrier(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Surge(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistPoison(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResilientAura(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistDisease(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func DisruptMagic(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Reflection(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Dodge(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func ResistAcid(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }

func Embolden(caller MobNPCSpells, target MobNPCSpells, modifier int){ return }


