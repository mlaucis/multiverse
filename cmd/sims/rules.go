package main

import (
	"fmt"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
)

const (
	fmtCommentPost     = "%s %s (%s) commented on a Post."
	fmtCommentPostOwn  = "%s %s (%s) commented on your Post."
	fmtFollow          = "%s %s (%s) started following you"
	fmtFriendConfirmed = "%s %s (%s) accepted your friend request."
	fmtFriendRequest   = "%s %s (%s) sent you a friend request."
	fmtLikePost        = "%s %s (%s) liked a Post."
	fmtLikePostOwn     = "%s %s (%s) liked your Post."
	fmtPostCreated     = "%s %s (%s) created a new Post."

	urnComment = "tapglue/posts/%d/comments/%d"
	urnPost    = "tapglue/posts/%d"
	urnUser    = "tapglue/users/%d"
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
				fmtFollow,
				origin.Firstname,
				origin.Lastname,
				origin.Username,
			),
			recipient: target.ID,
			urn:       fmt.Sprintf(urnUser, origin.ID),
		}, nil
	}
}

func conRuleFriendConfirmed(fetchUser fetchUserFunc) conRuleFunc {
	return func(change *connection.StateChange) (*message, error) {
		if change.Old == nil ||
			change.Old.Type != connection.TypeFriend ||
			change.Old.State != connection.StatePending ||
			change.New.State != connection.StateConfirmed {
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
				fmtFriendConfirmed,
				target.Firstname,
				target.Lastname,
				target.Username,
			),
			recipient: origin.ID,
			urn:       fmt.Sprintf(urnUser, origin.ID),
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
				fmtFriendRequest,
				origin.Firstname,
				origin.Lastname,
				origin.Username,
			),
			recipient: target.ID,
			urn:       fmt.Sprintf(urnUser, origin.ID),
		}, nil
	}
}

func eventRuleLikeCreated(
	fetchFollowerIDs fetchFollowerIDsFunc,
	fetchFriendIDs fetchFriendIDsFunc,
	fetchObject fetchObjectFunc,
	fetchUser fetchUserFunc,
	fetchUsers fetchUsersFunc,
) eventRuleFunc {
	return func(change *event.StateChange) ([]*message, error) {
		if change.Old != nil ||
			change.New.Enabled == false ||
			!isLike(change.New) {
			return nil, nil
		}

		post, err := fetchObject(change.Namespace, change.New.ObjectID)
		if err != nil {
			return nil, fmt.Errorf("post fetch: %s", err)
		}

		origin, err := fetchUser(change.Namespace, change.New.UserID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		owner, err := fetchUser(change.Namespace, post.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("owner fetch: %s", err)
		}

		followIDs, err := fetchFollowerIDs(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		friendIDs, err := fetchFriendIDs(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ids := filterIDs(append(followIDs, friendIDs...), owner.ID)

		rs, err := fetchUsers(change.Namespace, ids...)
		if err != nil {
			return nil, err
		}

		rs = append(rs, owner)

		ms := []*message{}

		for _, recipient := range rs {
			f := fmtLikePost

			if post.OwnerID == recipient.ID {
				f = fmtLikePostOwn
			}

			ms = append(ms, &message{
				message: fmt.Sprintf(
					f,
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: recipient.ID,
				urn:       fmt.Sprintf(urnPost, post.ID),
			})
		}

		return ms, nil
	}
}

func objectRuleCommentCreated(
	fetchFollowerIDs fetchFollowerIDsFunc,
	fetchFriendIDs fetchFriendIDsFunc,
	fetchObject fetchObjectFunc,
	fetchUser fetchUserFunc,
	fetchUsers fetchUsersFunc,
) objectRuleFunc {
	return func(change *object.StateChange) ([]*message, error) {
		if change.Old != nil ||
			change.New.Deleted == true ||
			!isComment(change.New) {
			return nil, nil
		}

		post, err := fetchObject(change.Namespace, change.New.ObjectID)
		if err != nil {
			return nil, fmt.Errorf("post fetch: %s", err)
		}

		origin, err := fetchUser(change.Namespace, change.New.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("origin fetch: %s", err)
		}

		owner, err := fetchUser(change.Namespace, post.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("owner fetch: %s", err)
		}

		followIDs, err := fetchFollowerIDs(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		friendIDs, err := fetchFriendIDs(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		ids := filterIDs(append(followIDs, friendIDs...), owner.ID)

		rs, err := fetchUsers(change.Namespace, ids...)
		if err != nil {
			return nil, err
		}

		rs = append(rs, owner)

		ms := []*message{}

		for _, recipient := range rs {
			f := fmtCommentPost

			if post.OwnerID == recipient.ID {
				f = fmtCommentPostOwn
			}

			ms = append(ms, &message{
				message: fmt.Sprintf(
					f,
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: recipient.ID,
				urn:       fmt.Sprintf(urnComment, post.ID, change.New.ID),
			})
		}

		return ms, nil
	}
}

func objectRulePostCreated(
	fetchFollowerIDs fetchFollowerIDsFunc,
	fetchFriendIDs fetchFriendIDsFunc,
	fetchUser fetchUserFunc,
	fetchUsers fetchUsersFunc,
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

		followIDs, err := fetchFollowerIDs(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		friendIDs, err := fetchFriendIDs(change.Namespace, origin.ID)
		if err != nil {
			return nil, err
		}

		rs, err := fetchUsers(change.Namespace, append(followIDs, friendIDs...)...)
		if err != nil {
			return nil, err
		}

		ms := []*message{}

		for _, recipient := range rs {
			ms = append(ms, &message{
				message: fmt.Sprintf(
					fmtPostCreated,
					origin.Firstname,
					origin.Lastname,
					origin.Username,
				),
				recipient: recipient.ID,
				urn:       fmt.Sprintf(urnPost, change.New.ID),
			})
		}

		return ms, nil
	}
}

func filterIDs(ids []uint64, fs ...uint64) []uint64 {
	var (
		is   = []uint64{}
		seen = map[uint64]struct{}{}
	)

	for _, id := range fs {
		seen[id] = struct{}{}
	}

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}

		is = append(is, id)
	}

	return is
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
