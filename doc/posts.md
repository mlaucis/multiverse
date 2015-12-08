# Posts

At the heart of user-generated content (referred to as UCG) the introduction of the Posts API will serve as a convenience to build commonly understood mechanics of social networks. Namly the creation of an object which can carry rich data like text, images or videos and can be interacted with in certain ways (likes, comments, shares). For specific use-cases not covered with this we still offer the underlying Objects API which Posts and their companions are build upon.

Biggest objective is to offer a UCG package to our customers which is rich enough to build typical core features of social networks while being simple to integrate with (SDK integration is paramount here). 

In its first iteration all functionality will be bundled in the Posts API until isolated use-cases for the sub-resources arise.

### API

**Post creation**

```
curl -X POST /posts \
  -d $'
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
        "type": "text
      }
    ],
    "tags": [
      "review"
    ],
    "visibility": 30
  }
  '
```

```
201 Created
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
      "type": "text
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
```

**Comment creation**

```
curl -X POST /posts/12743631303647839/comments \
	-d $'
	{
		"content": "Do like.",
	}
	'
```

```
201 Created
{
  "content": "Do like.",
  "id": "12743631303647840",
  "post_id": "12743631303647839",
  "user_id": "11286878691041888",
  "created_at": "2015-11-27T16:03:36.171840385Z",
  "updated_at": "2015-11-27T16:03:36.171840478Z"
}
```

**Comment retrieval**

```
curl -X GET /posts/12743631303647839/comments
```

```
200 OK
{
  "comments": [
    {
      "content": "Do like.",
      "id": "12743631303647840",
      "post_id": "12743631303647839",
      "user_id": "11286878691041888",
      "created_at": "2015-11-27T16:03:36.171840385Z",
      "updated_at": "2015-11-27T16:03:36.171840478Z"
    }
  ],
  "comments_count": 1,
  "post": {
    "attachments": [
      {
        "content": "http://bit.ly/123gif",
        "name": "teaser",
        "type": "url"
      },
      {
        "content": "Lorem ipsum...",
        "name": "body",
        "type": "text
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
}
```

**Comment deletion**

```
curl -X DELETE /posts/12743631303647839/comments/12743631303647840
```

```
204 No Content
```

**Like creation**

```
curl -X POST /posts/12743631303647839/likes \
  -H 'Content-Length: 0'
```

```
201 Created
{
	"id": "37363583381716140",
	"post_id": "12743631303647839",
  "user_id": "11286878691041888",
  "created_at":"2015-03-21T14:28:02.4+01:00",
  "updated_at":"2015-03-21T14:28:02.4+01:00"
}
```

**Like retrieval**

```
curl -X GET /posts/12743631303647839/likes
```

```
200 OK
{
  "likes": [
    {
      "id": "37363583381716140",
      "post_id": "12743631303647839",
      "user_id": "11286878691041888",
      "created_at":"2015-03-21T14:28:02.4+01:00",
      "updated_at":"2015-03-21T14:28:02.4+01:00"
    }
  ],
  "likes_count": 1,
  "post": {
    "attachments": [
      {
        "content": "http://bit.ly/123gif",
        "name": "teaser",
        "type": "url"
      },
      {
        "content": "Lorem ipsum...",
        "name": "body",
        "type": "text
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
}
```

**Like deletion**

```
curl -X DELETE /posts/12743631303647839/likes/37363583381716140
```

```
204 No Content
```

**Share creation**

```
curl -X POST /posts/12743631303647839/shares \
  -H 'Content-Length: 0'
```

```
{
  "id": "37363583381716141",
  "post_id": "12743631303647839",
  "user_id": "11286878691041888",
  "created_at":"2015-03-21T14:28:02.4+01:00",
  "updated_at":"2015-03-21T14:28:02.4+01:00"
}
```

**Share retrieval**

```
curl -X GET /posts/12743631303647839/likes
```

```
200 OK
{
  "shares": [
    {
      "id": "37363583381716141",
      "post_id": "12743631303647839",
      "user_id": "11286878691041888",
      "created_at":"2015-03-21T14:28:02.4+01:00",
      "updated_at":"2015-03-21T14:28:02.4+01:00"
    }
  ],
  "shares_count": 1,
  "post": {
    "attachments": [
      {
        "content": "http://bit.ly/123gif",
        "name": "teaser",
        "type": "url"
      },
      {
        "content": "Lorem ipsum...",
        "name": "body",
        "type": "text
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
}
```

**Share deletion**

```
curl -X DELETE /posts/12743631303647839/likes/37363583381716141
```

```
204 No Content
```

### Internals

The entire Posts API should be implemented with the introduction of a set of HTTP handlers which call out to a `PostController`. Named controller will re-use existing Objects and Events services to represent Posts and their companions.

#### Types

**Posts**

A post is an object of type `tg_post` and controlled by Tapglue via `owned`. As normal Objects, Posts accept attachments of all supported types and in general are close to vanilla Objects.

**Comments***

A comment is an object of type `tg_comment` and controlled by Tapglue via `owned`. In the beginning they carry one attachment of type `text` which holds the comment body.

**Likes**

A like is an event of type `tg_like` and controlled by Tapglue via `owned`. Their `object_id` is that one of the post.

**Shares**

A share is an object of type `tg_share` and controlled by Tapglue via `owned`. Additionally it should result in a corresponding event. *In the first iteration it should be left out until further specified.*
