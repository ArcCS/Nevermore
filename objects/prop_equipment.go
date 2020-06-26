package objects

type Equipment struct{
	// Status
	Armor int
	Weight int
	DamageIgnore int
	// TODO: Eventually create an effect system for equipment

	Head *Item
	Chest *Item
	Legs *Item
	Feet *Item
	Arms *Item
	Hands *Item
	Ring1 *Item
	Ring2 *Item

	// Hands, can hold shield or weapon
	Main *Item
	Off *Item
}

func RestoreEquipment(i ...*Item) Equipment{
	return Equipment{}
}

func ( *Equipment) Equip(item *Item) (ok bool){
	ok = false
	// TODO: Remove from inventory
	// Set to attached
	// Update armor values
	// Update weight
	return
}

// Attempt to unequip by name, or type
func (e *Equipment) Unequip(strName string) (ok bool){
	ok = false
	e.Hands = nil
	// TODO: Put into inventory
	// Update armor values
	// Update weight
	return
}