package cmd

/*
// Syntax: $ACTION item
func init() {
addHandler(action{}, "$action")
}

type action cmd

func (action) process(s *state) {

// Do we have item to perform action specified on command?
if len(s.words) == 0 {
	return
}

// Search for item to perform action.
alias := s.words[0]
what := s.where.Search(alias)

// If item not found all we can do is bail
if what == nil {
	return
}

// Reschedule event and bail early if there are no players here to see the
// action or it's too crowded to see the action.
if !s.where.Players() || s.where.Crowded() {
	Item.FindAction(what).Action()
	s.ok = true
	return
}

// See if item actually has actions. If not, bail without rescheduling. There
// is no point in rescheduling if there are no actions.
oa := Item.FindOnAction(what)
if !oa.Found() {
	return
}

// Display action and schedule next action. Only notify the actor if it's not
// the thing issuing the command.
if s.actor.UID() != what.UID() {
	s.msg.Actor.SendInfo(oa.ActionText())
}
s.msg.Observer.SendInfo(oa.ActionText())
items.FindAction(what).Action()

s.ok = true

}
*/