package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

func init() {
	addHandler(find{},
		"Usage:  find (room|mob|item) name|desc|maxdam|range (text) (page #) \n \n Use this command to search the database and find a list of matching items \n Items can also be replace name or desc with maxdam example: \n find item maxdam 25 \n To find a weapon with a maximum damage associated with it. \n Mobs/Items can also be replace name or desc with range example: \n find item range 225-250 \n To find all mobs in a specific id range. ",
		permissions.Builder,
		"find")
}

type find cmd

func (find) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Search what?  Missing parameters")
		return
	}

	objectType := strings.ToLower(s.words[0])
	searchType := strings.ToLower(s.words[1])
	searchText := strings.ToLower(s.words[2])
	searchPage := 0
	if len(s.words) == 4 {
		page, _ := strconv.Atoi(s.words[3])
		searchPage = page
	}
	switch objectType {
	case "room":
		if searchType == "name" {
			results := data.SearchRoomName(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["room_id"].(int64))) + ")(" + itemData["creator"].(string) + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find room name "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else if searchType == "desc" {
			results := data.SearchRoomDesc(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["room_id"].(int64))) + ")(" + itemData["creator"].(string) + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find room desc "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else {
			s.msg.Actor.SendBad("Search which field?")
		}
	case "mob":
		if searchType == "name" {
			results := data.SearchMobName(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["mob_id"].(int64))) + ")(" + strconv.Itoa(int(itemData["level"].(int64))) + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find mob name "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else if searchType == "desc" {
			results := data.SearchMobDesc(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["mob_id"].(int64))) + ")(" + strconv.Itoa(int(itemData["level"].(int64))) + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find mob desc "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else if searchType == "range" {
			if !strings.Contains(searchText, "-") {
				s.msg.Actor.SendBad("Dash in range not optional, examples: 250-280 3000-3500")
				return
			}
			idRange := strings.Split(searchText, "-")
			loId, _ := strconv.Atoi(idRange[0])
			hiId, _ := strconv.Atoi(idRange[1])
			if loId > hiId {
				s.msg.Actor.SendBad("Bad ID Range")
				return
			}
			results := data.SearchMobRange(loId, hiId, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["mob_id"].(int64))) + ")(" + strconv.Itoa(int(itemData["level"].(int64))) + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find mob range "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else {
			s.msg.Actor.SendBad("Search which field?")
		}
	case "item":
		if searchType == "name" {
			results := data.SearchItemName(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["item_id"].(int64))) + ")(" + config.ItemTypes[int(itemData["type"].(int64))] + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find item name "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else if searchType == "desc" {
			results := data.SearchItemDesc(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["item_id"].(int64))) + ")(" + config.ItemTypes[int(itemData["type"].(int64))] + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find item desc "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else if searchType == "maxdam" {
			results := data.SearchItemMaxDamage(searchText, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["item_id"].(int64))) + ")(" + config.ItemTypes[int(itemData["type"].(int64))] + ") " + itemData["name"].(string) + " Max Damage: " + strconv.Itoa(int(itemData["max_damage"].(int64))))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find item range "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else if searchType == "range" {
			if !strings.Contains(searchText, "-") {
				s.msg.Actor.SendBad("Dash in range not optional, examples: 250-280 3000-3500")
				return
			}
			idRange := strings.Split(searchText, "-")
			loId, _ := strconv.Atoi(idRange[0])
			hiId, _ := strconv.Atoi(idRange[1])
			if loId > hiId {
				s.msg.Actor.SendBad("Bad ID Range")
				return
			}
			results := data.SearchItemRange(loId, hiId, config.Server.SearchResults*searchPage)
			s.msg.Actor.SendGood("===== Search Results =====")
			for _, item := range results {
				if item != nil {
					itemData := item.(map[string]interface{})
					s.msg.Actor.SendGood("(" + strconv.Itoa(int(itemData["item_id"].(int64))) + ")(" + config.ItemTypes[int(itemData["type"].(int64))] + ") " + itemData["name"].(string))
				}
			}
			s.msg.Actor.SendGood("===== Type 'more' for another page of results =====")
			s.actor.AddCommands("more", "find item range "+searchText+" "+strconv.Itoa(searchPage+1))
			return
		} else {
			s.msg.Actor.SendBad("Search which field?")
		}

	}
	s.ok = true
	return
}
