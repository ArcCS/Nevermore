package spells

type MobNPCSpells interface {
	ChangePlacement(place int64)
	ApplyEffect()
	RemoteEffect(effect string)
	ReceiveDamage(damage int)
	ReceiveVitalDamage(damage int)
	Heal(damage int)
	HealVital(damage int)
	RestoreMana(damage int)
	InflictDamage() (damage int)
	CastSpell(spell string) bool
}

var Effects = map[string]func(caller MobNPCSpells, target MobNPCSpells){
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

func HealStam(caller MobNPCSpells, target MobNPCSpells){ return }

func HealVit(caller MobNPCSpells, target MobNPCSpells){ return }

func HealBoth(caller MobNPCSpells, target MobNPCSpells){ return }

func HealAll(caller MobNPCSpells, target MobNPCSpells){ return }

func FireDamage(caller MobNPCSpells, target MobNPCSpells){ return }

func EarthDamage(caller MobNPCSpells, target MobNPCSpells){ return }

func AirDamage(caller MobNPCSpells, target MobNPCSpells){ return }

func WaterDamage(caller MobNPCSpells, target MobNPCSpells){ return }

func Light(caller MobNPCSpells, target MobNPCSpells){ return }

func CurePoison(caller MobNPCSpells, target MobNPCSpells){ return }

func Bless(caller MobNPCSpells, target MobNPCSpells){ return }

func Protection(caller MobNPCSpells, target MobNPCSpells){ return }

func Invisibility(caller MobNPCSpells, target MobNPCSpells){ return }

func DetectInvisible(caller MobNPCSpells, target MobNPCSpells){ return }

func Teleport(caller MobNPCSpells, target MobNPCSpells){ return }

func Stun(caller MobNPCSpells, target MobNPCSpells){ return }

func Enchant(caller MobNPCSpells, target MobNPCSpells){ return }

func Recall(caller MobNPCSpells, target MobNPCSpells){ return }

func Summon(caller MobNPCSpells, target MobNPCSpells){ return }

func WizardWalk(caller MobNPCSpells, target MobNPCSpells){ return }

func Levitate(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistFire(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistMagic(caller MobNPCSpells, target MobNPCSpells){ return }

func RemoveCurse(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistAir(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistWater(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistEarth(caller MobNPCSpells, target MobNPCSpells){ return }

func Clairvoyance(caller MobNPCSpells, target MobNPCSpells){ return }

func RemoveDisease(caller MobNPCSpells, target MobNPCSpells){ return }

func CureBlindness(caller MobNPCSpells, target MobNPCSpells){ return }

func Polymorph(caller MobNPCSpells, target MobNPCSpells){ return }

func Attraction(caller MobNPCSpells, target MobNPCSpells){ return }

func InertialBarrier(caller MobNPCSpells, target MobNPCSpells){ return }

func Surge(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistPoison(caller MobNPCSpells, target MobNPCSpells){ return }

func ResilientAura(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistDisease(caller MobNPCSpells, target MobNPCSpells){ return }

func DisruptMagic(caller MobNPCSpells, target MobNPCSpells){ return }

func Reflection(caller MobNPCSpells, target MobNPCSpells){ return }

func Dodge(caller MobNPCSpells, target MobNPCSpells){ return }

func ResistAcid(caller MobNPCSpells, target MobNPCSpells){ return }

func Embolden(caller MobNPCSpells, target MobNPCSpells){ return }

