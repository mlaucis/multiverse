import { Dispatcher } from 'flux'

const dispatcher = new Dispatcher

export function dispatch(type, action = {}) {
  dispatcher.dispatch({ type, ...action })
}

export function dispatchAsync(promise, types, action = {}) {
  const { request, success, failure } = types

  console.log(request)

  // FIXME(xla) Understand why the extra dispatch of request can lead to
  //            invariant violations in the Dispatcher.
  // dispatch(request, action)

  return new Promise( (resolve, reject) => {
    promise
      .then( response => {
        console.log(success)
        dispatch(success, { ...action, response })

        resolve(response)

        return response
      })
      .catch( error => {
        console.log(failure)
        console.error(error)
        dispatch(failure, { ...action, error })

        // FIXME(xla): Understand why this leads to uncaught error if catch
        //             handler downstream is not attached.
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
