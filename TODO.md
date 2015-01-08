
# To Do's

## Done

- [x] Draft API design
- [x] Design entities
- [x] Create db tables
- [x] Implement routes
- [x] Mock API responses
- [x] Document code
- [x] Document api
- [x] Implement basic logging
- [x] Adjust response headers and caching
- [x] Create test files
- [x] ORM + DB Connection
- [x] Implement request validation in routes with regex 
- [x] Implement queries to read/write entities 

## Concept & API Design

### Access

- [ ] Add friend request confirmation
- [ ] Define what APIs to make unaccessible from outside (app,account etc.)
- [ ] Remove API for app/users?
- [ ] Limit data access to only friends?
- [ ] Define API versioning in header and/or URI

### Requests 
- [ ] Define HTTP Redirects if needed
- [ ] Require User-Agent
- [ ] Define & Implement request parameters
- [ ] Define & Implement condition requests
- [ ] Evaluate & Define rate limitation

### Responses
- [ ] Implement proper GET status codes
- [ ] Implement proper POST/PATCH status codes
- [ ] Implement proper DELETE status codes
- [ ] Implement proper Error status codes

### Authentication

- [ ] Implement request authentication with `user_token`
- [ ] Implement user login (`username`+`password`)
- [ ] Implement password recovery
- [ ] Check (check: http://jwt.io/)

### Model

- [ ] Refine Follwer vs. Friend Model
- [ ] Add `IDFA`, `IDFV`, `GPS_ADID`, `Game_Center_ID`, `FB_ID`, etc. and `push_token to users?
- [ ] Attach push certificate to app entity to use for push.

### Data & Processing

- [ ] Evaluate [ffjson](https://github.com/pquerna/ffjson)
- [ ] Implement UTF-8 encoding everywhere
- [ ] Omit or blank or `null` fields
- [ ] Implement PATCH routes and queries
- [ ] Implement DELETE routes and queries
- [ ] Implement Links (`href`) in responses
- [ ] Implement Expansion for links/tokens
- [ ] Define & Implement pagination
- [ ] Define & Implement sorting
- [ ] Evaluate, Define & Implement search

### Security

- [ ] Evaluate caching headers if SSL encrypted
- [ ] Implement SSL only
- [ ] Implement `base-64 encoded` token in requests

### Webhooks

- [ ] Evaluate need for webhooks
- [ ] Define webhook use cases
- [ ] Design and implement webhooks

## Architecture

### Database

- [ ] Evaluate Redis as main DB for users, connections, sessions and events
- [ ] Implement [go-redis](https://github.com/go-redis/redis)

### Message Queue

- [ ] Implement MQ for write requests (i.e. NATS, RabbitMQ, etc.)

### Testing & Benchmarks

- [ ] Implement unit tests db
- [ ] Implement unit tests server
- [ ] Define and implement integration tests
- [ ] Define and implement benchmarks

### Monitoring

- [ ] Evaluate Graphite as monitoring/alerts/metrics solution

## Documentation

- [ ] Evaluate further use cases of godep

# Cases

| Entity | Case | Method | URL | Implementation | Test | Docs |
| ------ | ---- | ------ | --- |:--------------:|:----:|:----:|
|Account|Read Account|`GET`|`/account/:AccountID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Account#read-account-get-accountid)|
|Account|Create account|`POST`|`/account`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Account#create-account-post-account)|
|Account|Update account|`PATCH`|`/account/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/account.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account#update-account-patch-accountid)|
|Account|Delete account|`DELETE`|`/account/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/account.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account#delete-account-delete-accountid)|
|Account users|Read account user|`GET`|`/account/:AccountID/user/:ID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Account-users#read-account-user-get-accountaccountiduserid)|
|Account users|Create account user|`POST`|`/account/:AccountID/user`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Account-users#create-account-user-post-accountaccountiduser)|
|Account users|Update account user|`PATCH`|`/account/:AccountID/user/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account-users#update-account-user-patch-accountaccountiduserid)|
|Account users|Delete account user|`DELETE`|`/account/:AccountID/user/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account-users#delete-account-user-delete-accountaccountiduserid)|
|Account users|List account users|`GET`|`/account/:AccountID/users`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Account-users#list-account-users--get-accountaccountidusers)|
|Applications|Read application|`GET`|`/app/:AppID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Applications#read-application-get-appid)|
|Applications|Create application|`POST`|`/account/:AccountID/app`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Applications#create-application-post-accountaccountidapp)|
|Applications|Update application|`PATCH`|`/app/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Applications#update-application-patch-appid)|
|Applications|Delete application|`DELETE`|`/app/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Applications#delete-application-delete-appid)|
|Applications|List applications|`GET`|`/account/:AccountID/applications`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Applications#list-applications--get-accountaccountidapplications)|
|Users|Read app user|`GET`|`/app/:AppID/user/:Token`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Users#read-app-user-get-appappidusertoken)|
|Users|Create app user|`POST`|`/app/:AppID/user`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Users#create-app-user-post-appappiduser)|
|Users|Update app user|`PATCH`|`/app/:AppID/user/:Token`|[:x:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Users#update-app-user-patch-appappidusertoken)|
|Users|Delete app user|`DELETE`|`/app/:AppID/user/:Token`|[:x:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Users#delete-app-user-delete-appappidusertoken)|
|Users|List app users|`GET`|`/app/:AppID/users`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Users#list-app-users--get-appappidusers)|
|User connections|Create user connection|`POST`|`/app/:AppID/connection`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/connections.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/User-connections#create-user-connection-post-appappidconnection)|
|User connections|List user connections|`GET`|`/app/:AppID/user/:Token/connections`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/connections.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/User-connections#list-user-connections--get-appappidusertokenconnections)|
|User connections|Update user connection|`PATCH`|`/app/:AppID/user/:Token/connection/:ID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/connections.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/User-connections#update-user-connection-patch-appappidusertokenconnectionid)|
|User connections|Delete user connection|`DELETE`|`/app/:AppID/user/:Token/connection/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/User-connections#delete-connection-delete-appappidusertokenconnectionid)|
|Sessions|Read session|`GET`|`/app/:AppID/user/:Token/session/:ID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Sessions#read-user-session-get-appappidusertokensessionid)|
|Sessions|Create session|`POST`|`/app/:AppID/user/:userToken/session`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Sessions#create-user-session-post-appappidusertokensession)|
|Sessions|Update session|`PATCH`|`/app/:AppID/user/:Token/session/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Sessions#update-user-session-patch-appappidusertokensessionid)|
|Sessions|Delete session|`DELETE`|`/app/:AppID/user/:Token/session/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Sessions#delete-user-session-delete-appappidusertokensessionid)|
|Sessions|List user sessions|`GET`|`/app/:AppID/user/:Token/sessions`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Sessions#list-user-sessions-get-appappidusertokensessions)|
|Events|Read event|`GET`|`/app/:AppID/event/:ID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Events#read-event-get-appappideventid)|
|Events|Create event|`POST`|`/app/:AppID/user/:Token/session/:SessionID/event`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Events#create-event-post-appappidusertokensessionsessionidevent)|
|Events|Delete event|`DELETE`|`/app/:AppID/event/:ID`|[:x:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Events#delete-event-delete-appappideventeventid)|
|Events|List user events|`GET`|`/app/:AppID/user/:Token/events`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Events#list-user-events--get-appappidusertokenevents)|
|Events|List session events|`GET`|`/app/:AppID/user/:Token/session/:SessionID/events`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Events#list-session-events--get-appappidusertokensessionsessionidevents)|
|Events|List connections events|`GET`|`/app/:AppID/user/:Token/connections/events`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:white_check_mark:](https://github.com/Tapglue/backend/wiki/Events#list-connection-events--get-appappidusertokenconnectionsevents)|