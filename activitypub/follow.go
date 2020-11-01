package activitypub

import (
	"fmt"
	"litepub1/models"
	"net/url"
	"reflect"
)

func ProcessFollow(actor *Actor, activity *Activity) {
	instanceUrl, err := url.Parse(activity.Actor)
	if err != nil {
		logger.Debugf("could not parse actor url (%s): %s", activity.Actor, err.Error())
		return
	}

	// Check blocklist
	if models.BlockedInstanceExists(instanceUrl.Host) {
		logger.Warningf("blocked instanceUrl (%s) tried to follow relay", activity.Actor)
		return
	}

	logger.Tracef("Actor ID: %v", actor.ID)
	if !models.FollowedInstanceExists(instanceUrl.Host) {
		instance := models.FollowedInstance{
			Hostname: instanceUrl.Host,
		}

		err := models.CreateFollowedInstance(&instance)
		if err != nil {
			logger.Errorf("could not add followed instanceUrl (%s): %s", instance.Hostname, err.Error())
		}

		// check if object is actor
		logger.Tracef("Follow object (%s) is: %v", reflect.TypeOf(activity.Object), activity.Object)

		message := Activity{
			Context: "https://www.w3.org/ns/activitystreams",
			Type: "Accept",
			To: []string{actor.ID},
			Actor: fmt.Sprintf("https://%s/actor", myAPHost),
			Object: Activity{
				Type: "Follow",
				ID: activity.ID,
				Object: fmt.Sprintf("https://%s/actor", myAPHost),
				Actor: actor.ID,
			},
		}

		err = actor.PushActivity(&message)
		if err != nil {
			logger.Errorf("could not push follow Accept to actor (%s): %s", activity.Actor, err.Error())
		}
	}

}