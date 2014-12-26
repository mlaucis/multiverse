
# To Do's

1. ORM + DB Connection
2. Tests
3. Secure API (Authentication) (check: http://jwt.io/)
4. Add message queue
5. Add cache (Redis)
6. Refine version control (v1, v2 etc.)

# Cases

| Entity | Case | Method | URL | Implementation | Test | Docs |
| ------ | ---- | ------ | --- |:--------------:|:----:|:----:|
|Account|Read Account|GET|`/account/:AccountID`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Gluee/backend/wiki/1.-Account#read-account)|
|Account|Create account|POST|`/account`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Gluee/backend/wiki/1.-Account#create-account)|
|AccountUser|Read account user|GET|`/account/:AccountID/user/:UserID`|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/2.-AccountUser#read-account-user)|
|AccountUser|Create account user|POST|`/account/:AccountID/user`|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/2.-AccountUser#create-account-user)|
|AccountUser|List account users|GET|`/account/:AccountID/users`|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/2.-AccountUser#list-account-users)|
|Application|Read application|GET|`/app/:AppID`|[:x:](https://github.com/Gluee/backend/blob/master/server/application.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Gluee/backend/wiki/3.-Application#read-application)|
|Application|Create application|POST|`/account/:AccountID/app`|[:x:](https://github.com/Gluee/backend/blob/master/server/application.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Gluee/backend/wiki/3.-Application#create-application)|
|Application|List applications|GET|`/account/:AccountID/applications`|[:x:](https://github.com/Gluee/backend/blob/master/server/application.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Gluee/backend/wiki/3.-Application#list-applications)|
|User|Read app user|GET|`/app/:AppID/user/:Token`|[:x:](https://github.com/Gluee/backend/blob/master/server/user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/4.-User#read-app-user)|
|User|Create app user|POST|`/app/:AppID/user/:userToken`|[:x:](https://github.com/Gluee/backend/blob/master/server/user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/4.-User#create-app-user)|
|User|List app users|GET|`/app/:AppID/users`|[:x:](https://github.com/Gluee/backend/blob/master/server/user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/4.-User#list-app-users)|
|UserConnection|Create user connection|POST|`/app/:AppID/user/:Token/connections`|[:x:](https://github.com/Gluee/backend/blob/master/server/connections.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Gluee/backend/wiki/5.-UserConnection#create-user-connection)|
|UserConnection|List user connections|GET|`/app/:AppID/connection`|[:x:](https://github.com/Gluee/backend/blob/master/server/connections.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Gluee/backend/wiki/5.-UserConnection#list-user-connections)|
|Session|Read session|GET|`/app/:AppID/user/:Token/session/:SessionID`|[:x:](https://github.com/Gluee/backend/blob/master/server/session.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Gluee/backend/wiki/6.-Session#read-session)|
|Session|Create session|POST|`/app/:AppID/user/:userToken/session`|[:x:](https://github.com/Gluee/backend/blob/master/server/session.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Gluee/backend/wiki/6.-Session#create-session)|
|Session|List user sessions|GET|`/app/:AppID/user/:userToken/sessions`|[:x:](https://github.com/Gluee/backend/blob/master/server/session.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Gluee/backend/wiki/6.-Session#list-user-sessions)|
|Event|Read event|GET|`/app/:AppID/event/:EventID`|[:x:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#read-event)|
|Event|Create event|POST|`/app/:AppID/user/:userToken/session/:SessionID/event`|[:x:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#create-event)|
|Event|List user events|GET|`/app/:AppID/user/:Token/events`|[:x:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#list-user-events)|
|Event|List session events|GET|`/app/:AppID/user/:userToken/session/:SessionID/events`|[:x:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#list-session-events)|
|Event|List connections events|GET|`/app/:AppID/user/:Token/connections/events`|[:x:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#list-connections-events)|