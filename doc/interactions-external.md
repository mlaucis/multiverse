
# External interactions

Wherein we describe how to offer the core interactions we offer for Posts: Comments and Likes (in the future: Votes, Shares, etc.). The UGC package evolves around the idea of a fully user created corpus of entities, which ignores a variety of use-cases where customers already own a large set of entities which they own and are not easily mapped to Posts. We assume that the reasons are unnegiotable, while they still want the core interactions.

### API

We going to provide a new set of routes which at their roots share `/externals` and offer the same functionality that can be found for `/posts`. Instead of expecting a post id we offer a that a set of alphanumeric characters can be passed which is used to associate the entity created with it for later reference and list retrieval.

```
/externals/{externalID:[a-zA-Z0-9]+}/comments
/externals/{externalID:[a-zA-Z0-9]+}/likes
...
```

To offer the same convenience that we have for Posts (e.g. Counts) for external objects we going to serve a limited response for `/externals/:id`.

```
"counts": {
  "comments": 3,
  "likes": 12,
  "shares": 1
},
"id": "1a2b7c4d3e6f3g1h3i0j3k6l4m7n8o3p9q"
```

Another assumption which doesn't hold true is the knowledge around `Visibility` which means that we have to force the caller to provide that information.

```
curl -X POST /posts/1a2b7c4d3e6f3g1h3i0j3k6l4m7n8o3p9q/comments \
    -d $'
    {
        "content": "Do like.",
				"visibility": 30
    }
    '
```

```
201 Created
{
  "content": "Do like.",
  "external_id": "1a2b7c4d3e6f3g1h3i0j3k6l4m7n8o3p9q",
  "id": "12743631303647840",
  "user_id": "11286878691041888",
	"visibility": 30,
  "created_at": "2015-11-27T16:03:36.171840385Z",
  "updated_at": "2015-11-27T16:03:36.171840478Z"
}
```

```
curl -X POST /posts/1a2b7c4d3e6f3g1h3i0j3k6l4m7n8o3p9q/like \
    -d $'
    {
				"visibility": 30
    }
    '
```

```
201 Created
{
  "external_id": "1a2b7c4d3e6f3g1h3i0j3k6l4m7n8o3p9q",
  "id": "37363583381716140",
  "user_id": "11286878691041888",
  "created_at":"2015-03-21T14:28:02.4+01:00",
  "updated_at":"2015-03-21T14:28:02.4+01:00"
}
```

### Internals

Symmetrically we are going to store Comments as Objects and Likes as Events. Instead of setting the `ObjectID` to the corresponding post we going to set `TargetID` and `Target.ID` respectively which laster is used to retrieve per external object.

``` go
&object.Object{
	Attachments: []object.Attachment{
	  object.NewTextAttachment(attachmentContent, content),
	},
	OwnerID:    owner.ID,
	Owned:      true,
	TargetID:   externalID,
	Type:       typeComment,
	Visibility: ps[0].Visibility,
}
```
