import { Dispatcher } from 'flux'

const dispatcher = new Dispatcher

export function dispatch(type, action = {}) {
  dispatcher.dispatch({ type, ...action })
}

export function dispatchAsync(promise, types, action = {}) {
  const { request, success, failure } = types

  // FIXME(xla) Understand why the extra dispatch of request can lead to
  //            invariant violations in the Dispatcher.
  // dispatch(request, action)

  return new Promise( (resolve, reject) => {
    promise
      .then( response => {
        dispatch(success, { ...action, response })

        resolve(response)

        return response
      })
      .catch( error => {
        console.error(error)
        dispatch(failure, { ...action, error })

        reject(error)
      })
  })
}

export function register(callback) {
  return dispatcher.register(callback)
}

export function waitFor(tokens) {
  return dispatcher.waitFor(tokens)
}
