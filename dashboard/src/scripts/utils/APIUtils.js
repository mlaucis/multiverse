import { camelizeKeys, decamelizeKeys } from 'humps'
import 'whatwg-fetch'

import AccountStore from '../stores/AccountStore'

let _base = 'https://api.tapglue.com/0.4'

if (window.location.hostname === 'localhost') {
  _base = 'http://localhost:8083/0.4'
}

export function request(method, url, body) {
  let opts = {
    method: method,
    headers: {
      'Access-Control-Request-Headers': 'content-type',
      'Access-Control-Request-Method': method,
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    }
  }

  if (AccountStore.isAuthenticated) {
    let user = AccountStore.user
    let b64 = btoa(`${user.accountToken}:${user.token}`)

    opts.headers.Authorization = `Basic ${b64}`
  } else if (AccountStore.account) {
    let b64 = btoa(`${AccountStore.account.token}:`)

    opts.headers.Authorization = `Basic ${b64}`
  }

  if (typeof body === 'object') {
    opts.body = JSON.stringify(decamelizeKeys(body))
  }

  let _storedResponse

  return fetch(`${_base}${url}`, opts)
    .then( response => {
      _storedResponse = response

      if (response.status === 204) {
        return {}
      }

      return response.json()
    })
    .then( json => {
      if (_storedResponse.status >= 400) {
        let err = new Error(_storedResponse.statusText)
        err.errors = json.errors
        err.response = _storedResponse

        throw err
      }

      return camelizeKeys(json)
    })
}
