# Initial situation

Due to the rising amount of customers interested in using Tapglue as their main user management, the need for a password reset mechanism arises.

Currently it is possible to update user passwords with the `PUT` method on the resource. The customer can use the `backend_token` to reset the passwords of their users by manually specifying a new one.

# Goal

The goal is to design a solution that enables developers to implement a `Reset password` feature in their app.

# Approach

The general to process of the password-reset flow should look like this:

- Trigger password-reset
- Send email with link to change password
- User click link
- User enters new password
- Changed password will be stored

(send email –> user clicks a link –> sets new password –> account password encrypted and stored securely)

Changing the user password with an active `SESSION_TOKEN` is possible by performing a `PUT` on the user resource as mentioned above. The only use-case that is not covered is if there if the session is deleted/expired or not known.

This means triggering the password reset need to happen on the collection level `/users/`

# API

## Trigger Password Reset

**Request**

`POST /users/resetPassword`


```
{
  "email": "foo@bar.com"
}
```

Alternatively the username can be used as well:

```
{
  "username": "user123"
}
```

## Internal Logic

If the user with that given `email` or `username` was found Tapglue create a `resetToken`. The resetToken should only be valid for a couple of minutes (15min) and contain information about the app and user.

### Option 1: Set password via API

With a valid resetToken we can then allow a request:

`POST /users/setPassword`


```
{
  "reset_token": "token123",
  "new_password" : "newPW123"
}
```

The logic for the request above would need to be implemented by the customer.

### Option 2: Create secure time-bound link

`https://users.tapglue.com/reset?token=token123`

There has to be a website where the user can type the new password and submit. This website can perform the `setpassword` action mentioned above.

# Questions

There are following open-questions:

- Do we want to provide customization for:
  - a. email template
  - b. custom website
- Do we design the possibility for secret questions?
- Given that we use the email as the source, should we design email verification for signup to enable this?

# Sources

Following sources where considered in the design:

- http://devcenter.kinvey.com/rest/guides/users
- http://docs.stormpath.com/rest/product-guide/#application-password-reset
- https://parse.com/docs/rest/guide#users-requesting-a-password-reset
- http://docs.apigee.com/api-baas/content/resetting-user-password-3
