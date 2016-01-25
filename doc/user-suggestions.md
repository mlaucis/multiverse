# User suggestions

Before building a fully-fledged user recommendations API we discussed building a simple heuristic which provides a list of users as a suggestion to another user. The main use-case would be to show list of users to another user during the Onboarding to the community.

# Goal

The goal is to overcome the empty-room-problem by providing features to the customers that help to grow their network fast.

# Approach

Building a prototype for this feature should be approach in four steps:

1. Research heuristics for suggestions
2. Create Query to fetch users
3. Implement API for suggestions
4. Implement in SDK API

# API

One idea is to provide three versions initially:

## New users

`GET users/new`

Get the 10 latest users.

## Trending users

`GET users/trending`

Get the 10 trending users (heuristic to be defined).

## Recommended users

`GET users/recommendations`

Recommend 10 users (heuristic to be defined).
