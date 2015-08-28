import { request } from '../utils/APIUtils'

export function member(id, user) {
  return request('get', `/accounts/${user.accountId}/users/${id}`)
}

export function memberDelete(id, user) {
  return request('delete', `/accounts/${user.accountId}/users/${id}`)
}

export function memberUpdate(
  email,
  password,
  firstName,
  lastName,
  id,
  accountID
) {
  return request('put', `/accounts/${accountID}/users/${id}`, {
    email: email,
    userName: email,
    firstName: firstName,
    lastName: lastName,
    password: password
  })
}

export function members(user) {
  return request('get', `/accounts/${user.accountId}/users`)
}
