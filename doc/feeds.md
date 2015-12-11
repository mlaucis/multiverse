# Feeds

With the introduction of the UGC family (namely Posts and Comments), system-level events and more in the pipeline the feed at its core needs to be looked at. This is driven by the idea that with a small set of pre-defined feed the majority of use-cases should be covered while leaving room to expand on feeds as new concepts and engines are introduced to the platform. At this point we identified five distinct use-cases:

|Endpoint|Name|Description|
|---|---|---|
|`/me/feed`|news feed|superset of events and posts driven by social and interest graph|
|`/me/feed/events`|event feed|events from social and interest graph|
|`/me/feed/posts`|post feed|posts from social and interest graph|
|`/users/:id/events`|user events|events of a user as seen by the current user|
|`/users/:id/posts`|user posts|posts of a user as seen by the current user|

*Currently there seems to be no reasonable use-case for a superset of events and posts of a user as seen by the current user*

### API

This section outlines the possible response structure of the feeds.

##### news feed

The important point here is the speration of `events` and `posts` in two lists to avoid the complexities that arise when interleaving two fundamentally different types.

``` json
{
  "events": [
    {
      "id_string": "12570250134426155",
      "user_id_string": "12442021620879682",
      "id": 12570250134426156,
      "user_id": 12442021620879682,
      "type": "love",
      "visibility": 30,
      "created_at": "2015-11-26T11:21:12.945437973Z",
      "updated_at": "2015-11-26T11:21:12.945437973Z"
    }
    ...
  ],
	"posts": [
		{
			"attachments": [
				{
					"content": "http://bit.ly/123gif",
					"name": "teaser",
					"type": "url"
				},
				{
					"content": "Lorem ipsum...",
					"name": "body",
					"type": "text"
				}
			],
			"id": "12743631303647839",
			"tags": [
				"review"
			],
			"visibility": 30,
			"user_id": "11286878691041887",
			"created_at": "2015-11-27T16:03:36.171840385Z",
			"updated_at": "2015-11-27T16:03:36.171840478Z"
		}
		...
	],
  "users": {
    "12442021620879682": {
      "id_string": "12442021620879682",
      "id": 12442021620879682,
      "social_ids": {
        "facebook": "fb54321"
      },
      "is_friend": true,
      "is_follower": false,
      "is_followed": false,
      "user_name": "Quemby",
      "first_name": "Perry",
      "last_name": "Mccoy",
      "email": "hello@yahoo.fr",
      "metadata": {
        "key": "value"
      },
      "enabled": true
    }
    ...
  },
  "events_count": 8,
  "posts_count": 3,
  "users_count": 4,
  "unread_events_count": 0
}

```

##### event feed

Variant of the above example but only returning the subset of events.

##### post feed

Variant of first example but only returning the subset of posts.

### Internals

Ultimately all feeds prefixed with `/me` are split in three phases:

1. traverse the social graph to receive followings and friends
2. receive entities from users returned by step one
3. translate entities into feed response

We should be able to build feeds without it being its own entity, rather a composition of existing services. User specific feeds should already be covered with the current implementations of events and posts. To cover the meta feeds we need a control structure:

``` go
import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type FeedController struct {
	connections connection.Service
	events      event.Service
	posts       object.Service
}

func (c *FeedController) NewsFeed(
	app *v04_entity.Application,
	user *v04_entity.ApplicationUser,
) ([]*v04_entity.Event, []*object.Object, error) {
	// ...
}

func (c *FeedController) EventFeed(
	app *v04_entity.Application,
	user *v04_entity.ApplicationUser,
) ([]*v04_entity.Event, error) {
	// ...
}

func (c *FeedController) PostFeed(
	app *v04_entity.Application,
	user *v04_entity.ApplicationUser,
) ([]*object.Object, error) {
	// ...
}
```

Which in turn will be consumed by the handlers assigned for each of the feed endpoints.

### Entity Constraints

This section describes how feeds or parts of them (entities) are composed by outlining the logic used to retrieve them. We use `user` as the current user for which we compute the feed.

##### events

* events with **global visbility**
* events from **following** connections with **connection** visibility
* events from **following** connections with **public** visibility
* events from **friend** connections with **connection** visibility
* events from **friend** connections with **public** visibility
* events whith **target.id** equals **user.id**

*ordered by creation time*

##### objects

* objects from **following** connections with **connection** visibility and **tg_post** type
* objects from **following** connections with **public** visibility and **tg_post** type
* objects from **friend** connections with **connection** visibility and **tg_post** type
* objects from **friend** connections with **public** visibility and **tg_post** type
