## Goal

Add community management features to our platform to enable our customers to manage content in their community.

## Description

This PR specifies community management functionalities, requested by a potential customer Zya (Ditty). It includes capabilities of reporting and removing offensive or inappropriate content.
Zya wants the ability for users to report content and being able to review and delete it on their end.

## Methodology

- Requirements
- Action Items
- Deliverables
- Open Questions

## Requirements

- Report users / content (from other users)
- Delete users / content (from moderator)
- Web interface to consume reports and delete items
- Data browser
- Promote content
- Promote user

## Action items

### Report, Review & Delete offensive Users / Content

#### API

- Create new entity "report" (implement report service)
  - reported_type: user,post,comment
  - state: pending, confirmed, declined
  - member_id
  - user_id
  - reported_id
  - reason
  - type
  - [Optional] summary

- Create report endpoints in core HTTP API
  - users
  - posts
  - comments
- Create report endpoints for dashboard
  - list entities by type order by amount of reports
  - list reports for entity
  - endpoint to resolve report
  - optionally: resolve all reports on same entity

#### Dashboards

- List representation
- Add actions

#### Estimation

- 1 iteration for API
- 2 iteration for Dashboard

#### Open points

- Sending report summaries
- Track issues
- Smarter UX around consumption of reports

### Data Browser

#### API

`TBD`

#### Dashboard

`TBD`

#### Estimation

`TBD`

### Promote Content & Users

#### API

`TBD`

#### Dashboard

`TBD`

#### Estimation

`TBD`

## Deliverables

- [ ] Report User & Content API
- [ ] Review Users & Content Dashboard
- [ ] Data Browser
- [ ] Promote Users & Content API & Actions

## Open Questions

- Reported / Blocked multiple times (threshold) what should happen?
- Should users be notified that have been reported
- Reasons for reporting a user
- How to receive a report?
- Process of accepting / rejecting a report
- Minimum features they require for the web interface (promise less)
- How to promote content they injected in feeds (paid feature)?
- How will the user promotion be consumed?
- How much in the past will historical data be available?
- Will content of a deleted user be automatically deleted?
