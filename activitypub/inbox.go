package activitypub

import "reflect"

func ProcessInbox(actor *Actor, activity *Activity) {
	logger.Tracef("processing %s activity %s from %s [%s](%v)", activity.Type, activity.ID, actor.ID,
		reflect.TypeOf(activity.Object), activity.Object)

	if activity.Type == "Follow" {
		ProcessFollow(actor, activity)
	}


}