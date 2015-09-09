import { request } from '../utils/APIUtils'

export function accountGet(user) {
  return request('get', `/accounts/${user.accountId}`)
}

export function accountCreate(name, description, plan, originalReferrer) {
  return request('post', '/accounts', {
    name: name,
    description: description,
    metadata: {
      originalReferrer: originalReferrer,
      plan: plan
    }
  })
}

export function accountUserCreate(
  email,
  password,
  firstName,
  lastName,
  accountID,
  originalReferrer,
  referrer
) {
  return request('post', `/accounts/${accountID}/users`, {
    email: email,
    userName: email,
    firstName: firstName,
    lastName: lastName,
    password: password,
    metadata: {
      originalReferrer: originalReferrer,
      referrer: referrer
    }
  })
}

export function login(email, password) {
  return request('post', '/accounts/users/login', {
    email: email,
    password: password
  })
}

export function logout(user) {
  let url = `/accounts/${user.accountId}/users/${user.id}/logout`

  return request('delete', url)
}
