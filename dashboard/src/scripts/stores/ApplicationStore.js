import { EventEmitter } from 'events'

import ApplicationConstants from '../constants/ApplicationConstants'
import { register } from '../dispatcher/ConsoleDispatcher'

const CHANGE_EVENT = 'change'

class ApplicationStore extends EventEmitter {
  constructor() {
    super()

    this.dispatchToken = register(this._handleActions)
    this._apps = {}
    this._errors = []
  }

  get apps() {
    let apps = []

    for (let id in this._apps) {
      apps.push(this._apps[id])
    }

    return apps.sort( (left, right) => {
      let a = Date.parse(left.createdAt)
      let b = Date.parse(right.createdAt)

      if (a < b) {
        return -1
      }

      if (a > b) {
        return 1
      }

      return 0
    })
  }

  get errors() {
    return this._errors
  }

  getAppByID(id) {
    return this._apps[id]
  }

  emitChange() {
    this.emit(CHANGE_EVENT)
  }

  addChangeListener(cb) {
    this.on(CHANGE_EVENT, cb)
  }

  removeChangeListener(cb) {
    this.removeListener(CHANGE_EVENT, cb)
  }

  _handleActions = (action) => {
    switch (action.type) {
      case ApplicationConstants.APP_SUCCESS:
        let app = action.response
        this._apps[app.id] = app
        this.emitChange()

        break
      case ApplicationConstants.APPS_SUCCESS:
        this._apps = {}

        action.response.applications.forEach( (a) => {
          this._apps[a.id] = a
        })

        this.emitChange()

        break
      case ApplicationConstants.APP_CREATE_SUCCESS:
        this._errors = []
        let _app = action.response
        this._apps[_app.id] = _app
        this.emitChange()

        break
      case ApplicationConstants.APP_CREATE_FAILURE:
        this._errors = action.error.errors
        this.emitChange()

        break
      case ApplicationConstants.APP_DELETE_SUCCESS:
        delete this._apps[action.id]
        this._errors = []
        this.emitChange()

        break
      case ApplicationConstants.APP_DELETE_FAILURE:
        this._errors = action.error.errors
        this.emitChange()

        break
      case ApplicationConstants.APP_EDIT_SUCCESS:
        this._apps[action.id] = action.response
        this._errors = []
        this.emitChange()

        break
      case ApplicationConstants.APP_EDIT_FAILURE:
        this._errors = action.error.errors
        this.emitChange()

        break
      default:
      // nothing to do
    }
  }
}

export default new ApplicationStore
