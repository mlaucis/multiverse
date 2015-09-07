import { request } from '../utils/APIUtils'

export function accountGet(user) {
  return request('get', `/accounts/${user.accountId}`)
}

export function accountCreate(name, description, plan) {
  return request('post', '/accounts', {
    name: name,
    description: description,
    metadata: {
      plan: plan
    }
  })
}

export function accountUserCreate(
  email,
  password,
  firstName,
  lastName,
  accountID
) {
  return request('post', `/accounts/${accountID}/users`, {
    email: email,
    userName: email,
    firstName: firstName,
    lastName: lastName,
    password: password
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
