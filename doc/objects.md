# Objects

Introducing the new entity `objects`. This resource can be used to create persistent objects such as Posts.

## Use-cases

- Create objects
    - text
    - image
    - #tag
    - mention users in image (not in mvp)
    - mentions (not in mvp)
    - visibility
- Like object
- Comment object (not in mvp)
- Rating of object (not in mvp)
- Retrieve list of object
- Search objects (not in mvp)
- Counts of likes & comments for objects (not in mvp)
- Share / Repost + counts (counts not in mvp)
- Persistency

# Requirements

In the following the requirements for the MVP are listed.

## Interface

@xla, @dlsniper as we didn't talk about this please double check and feel free to adjust.

| Name | Method | Route |
| ---- |:------:| -----:|
|Create Object|`POST`|`/objects`|
|Retrieve Object|`GET`|`/objects/{objectID}`|
|Update Object|`PUT`|`/objects/{objectID}`|
|Delete Object|`DELETE`|`/objects/{objectID}`|
|Retrieve all objects|`GET`|`/objects`|
|Retrieve connections objects|`GET`|`/me/objects/connections`|
|Retrieve my objects|`GET`|`/me/objects`|

## Functionality

- [x] CRUD Object
- [x] Retrieve Objects
  - [x] All objects
  - [x] Connections objects
  - [x] My own objects

## Model

- tags
- visibility
- owner
- target
- type
- location
- longitude
- latitude
- attachments
  - text
  - image
  - video
  - audio
  - url
