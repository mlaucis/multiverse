# Marketing Automation

## Initial situation

In the last 10 months, we've created our network layer which enables developer to create user profiles, social graphs and feeds. We've added Posts, Comment and Likes to complete the basic functionality of each social network. This allows to create the fundamental signals within a network and make sense out of it.

## Goal

The goal is to expose the social experience through marketing channels to increase reach of Tapglue's social features. Concretely that will allow app developer to target their users outside of the apps to drive retention and engagement.

## Plan

### Logic Layer

Follow signals can be used:

#### Content

- Events/Posts
- System Events could be better to start with, as we - understand them better
- Trending content in my community
- Trending content globally
- Stats

#### Trigger

- Inactivity
- Custom Logic
- System Event

#### System Events

- `friend` joined app from other `network`
- `friend` confirmed `connection` with `me`
- `friend` confirmed `connection` with other `user`
- `me` has incoming pending `connection`
- `friend` created a `post`
- `friend` created a `comment` on `post`
- `friend` created a `like` on `post`
- `friend` created a `share` on `post`
- `friend` updated `avatar`
- `friend` has their `birthday`
- `me` has new `mention`

### Delivery Layer

Following channels can be used for delivery:

- Push
- Email
- In-app message
- SMS
- Webhooks

### Tracking System

We need to install a tracking system that allows developers to track the performance of the messages that have been sent.

## Architecture

### Discussion

- Design of the Logic Layer
- Design of Deliver Layer
- Design of Tracking System

### Challenges

-  How can MA System work independent from Network layer?
- System events as content/triggers: don't happen too often
- System events as content/triggers: too spammy
- Level of customization
- Localization
- Flexibility vs. Simplicity
- User control
- Dashboard functionality
- How do customers avoid conflicts between push notification systems
