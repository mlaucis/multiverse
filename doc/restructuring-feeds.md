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

### Use-cases

Following overview show the new structure from a use-cases perspective.

## Goal

The goal is to restructure the feed materializations to fullfill the requirements specified in the use-cases overview.

## API

The following table shows an overview the required changes:

| Name                          | Description | Endpoint current | Endpoint (0.5) | Action |
| ----------------------------- | ----- | ----------- | ---------------- | -------------- | ------ |
| News Feed                     |
| Posts Feed                    |
| Notifications Feed            |
| Notification Connections Feed |
| Notifications Me Feed         |
| User Activity Feed            |
| User Posts Feed               |
| User Events                   |
| User Posts                    |

| Name | Description | Endpoint current | Endpoint (0.5) | Action |
| ----------- | ----- | ----------- | ---------------- | -------------- | ------ |
| User Posts  |
| User Events |
