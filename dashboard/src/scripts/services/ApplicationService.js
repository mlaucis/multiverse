import { request } from '../utils/APIUtils'

export function app(id, user) {
  return request('get', `/accounts/${user.accountId}/applications/${id}`)
}

export function apps(user) {
  return request(
    'get',
    `/accounts/${user.accountId}/applications?=${Date.now()}`
  )
}

export function appCreate(name, description, user) {
  return request('post', `/accounts/${user.accountId}/applications`, {
    name: name,
    description: description
  })
}

export function appDelete(id, user) {
  return request('delete', `/accounts/${user.accountId}/applications/${id}`)
}

export function appUpdate(id, name, description, user) {
  return request('put', `/accounts/${user.accountId}/applications/${id}`, {
    name: name,
    description: description
  })
}
