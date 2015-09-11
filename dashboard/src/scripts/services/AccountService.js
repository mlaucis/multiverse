import { request } from '../utils/APIUtils'

export function accountGet(user) {
  return request('get', `/organizations/${user.accountId}`)
}

export function accountCreate(name, description, plan, originalReferrer) {
  return request('post', '/organizations', {
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
  return request('post', `/organizations/${accountID}/members`, {
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
  return request('post', '/organizations/members/login', {
    email: email,
    password: password
  })
}

export function logout(user) {
  let url = `/organizations/${user.accountId}/members/${user.id}/logout`

  return request('delete', url)
}
