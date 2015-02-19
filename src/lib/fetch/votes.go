package fetch

import (
	"labix.org/v2/mgo/bson"
	"log"
	proto "proto"
)

type EligibleUser struct {
	User     proto.User
	Eligible bool
}

func (f *Fetcher) StoreVote(v proto.RecipeVote) {
	f.SS.Database.Votes.Insert(v)
}

func (f *Fetcher) VotesForMeal(m proto.Meal) []proto.RecipeVote {
	votes := make([]proto.RecipeVote, 0)
	f.SS.Database.Votes.Find(bson.M{"meal": bson.M{"id": *m.Id}}).All(&votes)

	return votes
}

/**
 * Takes an input vote as a parameter and finds out whether the recipe
 * that was voted on has reached a quorum among all of the users in the
 * group. If yes, it's copied into the meal object. If not, do nothing.
 */
func (f *Fetcher) CheckForQuorum(v proto.RecipeVote) bool {
	// Get the group (full list of users)
	group, err := f.GroupById(*v.Group.Id)

	if err != nil {
		log.Println("Group doesn't exist")
		return false
	}

	// Get the meal (whose votes we're blocking on)
	meal, merr := f.GetCurrentMeal(group)

	if merr != nil {
		log.Println("Couldn't fetch current meal: " + merr.Error())
		return false
	}

	// Get the votes for the current meal
	votes := f.VotesForMeal(meal)

	users := make([]EligibleUser, 0)
	for i := 0; i < len(group.Members); i++ {
		users = append(users, EligibleUser{
			User:     *group.Members[i],
			Eligible: true,
		})
	}

	// Remove users who have abstained.
	for i := 0; i < len(users); i++ {
		// Check if this vote was issued by this user and whether it's an abstain vote.
		for j := 0; j < len(meal.Votes); j++ {
			if *users[i].User.Id == *meal.Votes[j].User.Id && *meal.Votes[j].Status == proto.RecipeVote_ABSTAIN {
				users[i].Eligible = false
			}
		}

		// Also check to see whether this user has voted yes to this recipe.
		for j := 0; j < len(votes); j++ {
			if *users[i].User.Id == *votes[j].User.Id && *votes[j].Status == proto.RecipeVote_YES {
				users[i].Eligible = false
			}
		}
	}

	// If vote yes || abstain for everyone, return yes. If not, return no.
	for i := 0; i < len(users); i++ {
		if users[i].Eligible {
			return false
		}
	}

	return true
}
