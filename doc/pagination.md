# Pagination
Requests that return multiple items will be paginated to 25 items by default.

## Query Parameters
There are two optional query parameters that may be specified to control
pagination:

- `before`: specifies the ID which will be used as a pivot in order to get the
items before it. The item matching the ID specified under `before` will not be
returned in the response. This is mutually exclusive with `after`.
- `after`: specifies the ID which will be used as a pivot in order to get the
items after it. The item matching the ID specified under `after` will not be
returned in the response. This is mutually exclusive with `before`.
- `limit`: The maximum number of collection items to return for a single request.
Minimum value is 1. Maximum value is 100. Default is 25.

## Constraints

If both `before` and `after` are specified then the request will produce a
`409 Conflict` return code for the request.

If neither `before` and `after` are not specified then the request will
return the latest items in the response, according to the `limit` parameter.

If `limit` parameter is not specified for a request then the default value
will be used. Currently the default value is 25 but this value might change in
the future. Clients should be advised to handle this situation on their end
either by explicitly specifying the `limit` parameter or nor relying on the
`limit` default value. In many cases requesting a certain number of items
might not completely satisfiable due to lack of new data present in the system
or a total number that's not perfectly divisible by the number specified by `limit`.

If a response is composed of multiple lists then it may contain one or more
empty lists if the other lists still have items that can be retrieved. Once
no list in the response contains items anymore then the response status will be
changed from `200` to `204`.

## Cursors

In order to provide a consistent pagination cursor for any number of lists that
the response might contain a single, opaque, cursor will be exposed to the
clients.

As such the cursor will have the following properties:

- internally the cursor will be a combined cursor from each list the response
returns
- externally the cursor will be a single, opaque, form representation

### Internal cursors

An internal cursor is a cursor which is not exposed to the API clients.

This cursor will be created from the ID of the elements of the list.

### External cursors

An external cursor is a cursor which is an aggregation of one or multiple
internal cursors. The aggregation is defined as a JSON object with members
named after the list that they are holding the value of the cursor for.

To represent an external cursor the following method will be used:

```
base64Encode({"cursor1": value, "cursor2": value, ....})
```

For example, for cursors: events at 1000 and and posts at 2005

```
base64.StdEncoding.EncodeToString("{\"events\":1000,\"posts\":2005}") = eyJldmVudHMiOjEwMDAsInBvc3RzIjoyMDA1fQ==
```

## Response

Once a successful response is generated, the following header will be present:

- `Link`: contains the links to the previous and next result set. These links
might not always return data as either hard limits might be reached in the system
or no new data might be present. The response code will be accordingly to the
situation either a 200 for a successful response with content or 204 for a
successful response with no content

The response body will also contain the following section:

```json
{
  "pagination:": {
    "next": "<link-to-next-page>",
    "prev": "<link-to-previous-page>"
  }
}
```

## Example

- Sending nothing

Request:
```bash
curl -i \
    -X "GET" "https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D" \
    -H "Accept: application/json" \
    -H "Authorization: Basic [[AUTH_TOKEN]]" \
    -H "User-Agent: Tapglue Sample"
```

Response:
```
HTTP/1.1 200 OK
[...]
Link: <https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&after=eyJldmVudHMiOjEwMDB9&limit=25>; rel="next",
        <https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&before=eyJldmVudHMiOjEwMjR9&limit=25>; rel="prev"
```

```json
{
   "events":[
      {
         "id_string":"1024",
         "user_id_string":"1",
         "id":1025,
         "user_id":1,
         "type":"love",
         "visibility":30,
         "object":{
            "id":"picture1",
            "type":"picture",
            "display_names":{
               "de":"Bild 1",
               "en":"Picture 1"
            }
         },
         "object_id":0,
         "owned":false,
         "created_at":"2015-12-14T10:56:02.314314969Z",
         "updated_at":"2015-12-14T10:56:02.314314969Z",
         "enabled":true
      },
      [...]
      {
        "id_string":"1000",
        "user_id_string":"1",
        "id":1000,
        "user_id":1,
        "type":"love",
        "visibility":30,
        "object":{
            "id":"picture2",
            "type":"picture",
            "display_names":{
              "de":"Bild 2",
              "en":"Picture 2"
            }
        },
        "object_id":0,
        "owned":false,
        "created_at":"2015-12-14T10:56:02.314314969Z",
        "updated_at":"2015-12-14T10:56:02.314314969Z",
        "enabled":true
    }
   ],
   "users":{
      "1":{
         "id_string":"1",
         "id":1,
         "is_friend":false,
         "is_follower":false,
         "is_followed":false,
         "user_name":"Username",
         "first_name":"First",
         "last_name":"Last",
         "email":"e@mail.com",
         "metadata":{
            "key":"valueChanged"
         },
         "enabled":true
      }
   },

   "pagination:": {
     "next": "<link-to-next-page>",
     "prev": "<link-to-previous-page>"
   },
   "events_count":25,
   "users_count":1
}
```

- Specifying only `limit`

Request:
```bash
curl -i \
    -X "GET" "https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&limit=10" \
    -H "Accept: application/json" \
    -H "Authorization: Basic [[AUTH_TOKEN]]" \
    -H "User-Agent: Tapglue Sample"
```

Response:
```
HTTP/1.1 200 OK
[...]
Link: <https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&after=eyJldmVudHMiOjEwMDB9&limit=10>; rel="next",
        <https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&before=eyJldmVudHMiOjEwMDl9&limit=10>; rel="prev"
```

- Specifying only the since parameter

Request:
```bash
curl -i \
    -X "GET" "https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&since=eyJldmVudHMiOjEwMDB9" \
    -H "Accept: application/json" \
    -H "Authorization: Basic [[AUTH_TOKEN]]" \
    -H "User-Agent: Tapglue Sample"
```

Response:
```
HTTP/1.1 200 OK
[...]
Link: <https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&after=eyJldmVudHMiOjEwMDB9&limit=10>; rel="next",
        <https://api.tapglue.com/0.4/events?where=%7B%22type%22%3A%20%7B%22in%22%3A%5B%22love%22%5D%7D%7D&before=eyJldmVudHMiOjEwMDl9&limit=10>; rel="prev"
```

