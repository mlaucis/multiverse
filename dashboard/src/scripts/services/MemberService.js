import { request } from '../utils/APIUtils'

export function member(id, user) {
  return request('get', `/organizations/${user.accountId}/members/${id}`)
}

export function memberDelete(id, user) {
  return request('delete', `/organizations/${user.accountId}/members/${id}`)
}

export function memberUpdate(
  email,
  password,
  firstName,
  lastName,
  id,
  accountID
) {
  return request('put', `/organizations/${accountID}/members/${id}`, {
    email: email,
    userName: email,
    firstName: firstName,
    lastName: lastName,
    password: password
  })
}

export function members(user) {
  return request('get', `/organizations/${user.accountId}/members`)
}
