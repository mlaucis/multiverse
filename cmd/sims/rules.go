package main

import (
	"fmt"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
)

type conRuleFunc func(*connection.StateChange) (*message, error)
type eventRuleFunc func(*event.StateChange) ([]*message, error)
type objectRuleFunc func(*object.StateChange) ([]*message, error)

func conRuleFollower(fetchUser fetchUserFunc) conRuleFunc {
	return func(change *connection.StateChange) (*message, error) {
		if change.Old != nil ||
			change.New.State != connection.StateConfirmed ||
			change.New.Type != connection.TypeFollow {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.FromID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		target, err := fetchUser(change.Namespace, change.New.ToID)
		if err != nil {
			return nil, fmt.Errorf("target fetch: %s", err)
		}

		return &message{
			message: fmt.Sprintf(
				"%s %s (%s) started following you",
				origin.Firstname,
				origin.Lastname,
				origin.Username,
			),
			recipient: target.ID,
		}, nil
	}
}

func conRuleFriendRequest(fetchUser fetchUserFunc) conRuleFunc {
	return func(change *connection.StateChange) (*message, error) {
		if change.Old != nil ||
			change.New.State != connection.StatePending ||
			change.New.Type != connection.TypeFriend {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.FromID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		target, err := fetchUser(change.Namespace, change.New.ToID)
		if err != nil {
			return nil, fmt.Errorf("target fetch: %s", err)
		}

		return &message{
			message: fmt.Sprintf(
				"%s %s (%s) sent you a friend request.",
				origin.Firstname,
				origin.Lastname,
				origin.Username,
			),
			recipient: target.ID,
		}, nil
	}
}

func eventRuleLikeCreated(
	fetchFriends fetchFriendsFunc,
	fetchUser fetchUserFunc,
) eventRuleFunc {
	return func(change *event.StateChange) ([]*message, error) {
		if change.Old != nil ||
			change.New.Enabled == false ||
			!isLike(change.New) {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.UserID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		fs, err := fetchFriends(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ms := []*message{}

		for _, friend := range fs {
			ms = append(ms, &message{
				message: fmt.Sprintf(
					"Your friend %s %s (%s) liked a Post.",
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: friend.ID,
			})
		}

		return ms, nil
	}
}

func objectRuleCommentCreated(
	fetchFriends fetchFriendsFunc,
	fetchUser fetchUserFunc,
) objectRuleFunc {
	return func(change *object.StateChange) ([]*message, error) {
		if change.Old != nil ||
			change.New.Deleted == true ||
			!isComment(change.New) {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		fs, err := fetchFriends(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ms := []*message{}

		for _, friend := range fs {
			ms = append(ms, &message{
				message: fmt.Sprintf(
					"Your friend %s %s (%s) commented on a Post.",
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: friend.ID,
			})
		}

		return ms, nil
	}
}

func objectRulePostCreated(
	fetchFriends fetchFriendsFunc,
	fetchUser fetchUserFunc,
) objectRuleFunc {
	return func(change *object.StateChange) ([]*message, error) {
		if change.Old != nil ||
			!isPost(change.New) ||
			change.New.Deleted == true {
			return nil, nil
		}

		origin, err := fetchUser(change.Namespace, change.New.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		fs, err := fetchFriends(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ms := []*message{}

		for _, friend := range fs {
			ms = append(ms, &message{
				message: fmt.Sprintf(
					"Your friend %s %s (%s) created a new Post.",
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: friend.ID,
			})
		}

		return ms, nil
	}
}

func isComment(o *object.Object) bool {
	if o.Type != controller.TypeComment {
		return false
	}

	return o.Owned
}

func isLike(e *event.Event) bool {
	if e.Type != controller.TypeLike {
		return false
	}

	return e.Owned
}

func isPost(o *object.Object) bool {
	if o.Type != controller.TypePost {
		return false
	}

	return o.Owned
}
