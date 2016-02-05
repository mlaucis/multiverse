# Feed restructuring

This PR describes the changes required for another iteration on our feeds structure.

## Initial situation

One of our customers required to implement an Instagram like Notifications feed where they wanted to have two different tabs:

- Following Feed
- Me Feed

![Image Instagram](http://cdn2.expertreviews.co.uk/sites/expertreviews/files/2015/09/instagram_activity_you.png?itok=AQf8tpp9)

That lead to the fact that we started rethinking the structure of our feeds from a hierarchy and use-case perspective.

The concept for the new feed structure has been discussed in this [Spreadsheet](https://docs.google.com/spreadsheets/d/1tbWwt30eQDYnEjgjpYxH1HM0rHPa3MPM4A6hNQRN8HA/edit)

### Hierarchy

Following diagram describes the hierarchy of the new feed structure.

![Image Hierarchy](http://s14.postimg.org/glqhzrx1t/Feeds.png)

### Use-cases

Following overview show the new structure from a use-cases perspective.

![Image Use-cases](http://s17.postimg.org/ko7jxayun/Use_cases.png)

## Goal

The goal is to restructure the feed materializations to fulfill the requirements specified in the use-cases overview.

## API

The following table shows an overview the new feed structure:

### Graph Feeds

| Name                          | Image         | Description | Endpoint current | Endpoint (0.5) |
| ----------------------------- | --------------| ----------- | ---------------- | -------------- |
| News Feed                     | ![Newsfeed](http://s24.postimg.org/va9af9gdh/newsfeed.png) | Contain all possible entries from within my network. | `/me/feed` | `/me/feed/news` |
| Posts Feed                    |  ![Posts Feed](http://s29.postimg.org/8xbnsnjlj/IMG_2279.png) | Contain posts created within my network. | `/me/feed/posts` | `/me/feed/posts` |
| Notifications Feed            | ![Notifications](http://s2.postimg.org/bkwsmwap5/notifications.png) | Superset of Graph and Me feed. | `/me/feed/events` | `/me/feed/notifications` |
| Notification Connections Feed | ![Notificatons Connections](http://s15.postimg.org/q5q03is3v/graphfeed.png) | Contains connection and event entries from within my network. | ` ` | `me/feed/social` |
| Notifications Me Feed         | ![Notificaitons Me](http://s14.postimg.org/e525flju9/mefeed.png) | Contains entries which target me or my content. | ` ` | `me/feed/self` |

### User Feeds

| Name                          |  Image         | Description | Endpoint current | Endpoint (0.5) |
| ----------------------------- | -------------- | ----------- | ---------------- | -------------- |
| User Activity Feed            | ![User Activity](http://s24.postimg.org/538k2of1h/userevents.png) | Contains entries which originate from me. | ` ` | `/me/activity` |
| User Activity Posts Feed      | ![User Activity Posts](http://s16.postimg.org/89eh6eys5/userposts.png) | All posts of a user. | ` ` | `/me/activity/posts` |
| User Activity Events Feed     | ![User Activity Events](http://s24.postimg.org/538k2of1h/userevents.png) | All activites of a user without users. | ` ` | `/me/activity/events` |

### User Collections

| Name                          |   Image         | Description | Endpoint current | Endpoint (0.5) |
| ----------------------------- | --------------- | ----------- | ---------------- | -------------- |
| User Posts  | ![User Activity](http://s24.postimg.org/538k2of1h/userevents.png) | Contains posts of a single user.  | `/me/posts` | `/me/posts` |
| User Events | ![User Activity Events](http://s24.postimg.org/538k2of1h/userevents.png) | Contains events of a single user. | `/me/events` | `/me/events` |
