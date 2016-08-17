## Goal

Add ability for users to block other users to remove them from the community experience and prevent "stalking".

## Description

This PR specifies user blocking abilities, requested by a potential customer Zya (Ditty).

## Methodology

- Requirements
- Action Items
- Deliverables
- Open Questions

## Requirements

- Block users

## Action items

### Block users

- New connection of type "blocked"
- Add is_blocked flag
- bi-directional block (anti-stalk)

#### API

- New block endpoints
 - block
 - unblock
 - receive
- Materialization of Lists & feeds
  - Business logic (remove users/content or just return with flag)
  - Product decision to be made regarding behaviour

#### Estimation

- 1 iteration

## Deliverables

- [ ] Block Users API

## Open Questions

- Receiving blocked users
- Appearance of blocked users in search
