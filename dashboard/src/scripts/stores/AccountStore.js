import { EventEmitter } from 'events'

import AccountConstants from '../constants/AccountConstants'
import { register } from '../dispatcher/ConsoleDispatcher'

const CHANGE_EVENT = 'change'

class AccountStore extends EventEmitter {
  constructor() {
    super()

    this.dispatchToken = register(this._handleActions)
    this._account = undefined
    this._errors = []
    this._initialized = false
    this._user = undefined
  }

  get account() {
    return this._account
  }

  get accountName() {
    let name = ''

    if (this._account && this._account.name) {
      name = this._account.name
    }

    return name
  }

  get errors() {
    return this._errors
  }

  get isAuthenticated() {
    return !!this._user && !!this._user.token
  }

  get user() {
    return this._user
  }

  init() {
    if (!this._account && localStorage.getItem(AccountConstants.ACCOUNT_KEY)) {
      let accountItem = localStorage.getItem(AccountConstants.ACCOUNT_KEY)

      this._account = JSON.parse(accountItem)
    }

    if (!this._user && localStorage.getItem(AccountConstants.USER_KEY)) {
      this._user = JSON.parse(localStorage.getItem(AccountConstants.USER_KEY))
    }
  }

  logout() {
    localStorage.removeItem(AccountConstants.ACCOUNT_KEY)
    localStorage.removeItem(AccountConstants.USER_KEY)

    this._account = undefined
    this._errors = []
    this._user = undefined
  }

  setAccount(account) {
    localStorage.setItem(
      AccountConstants.ACCOUNT_KEY,
      JSON.stringify(account)
    )

    this._account = account
  }

  setUser(user) {
    localStorage.setItem(AccountConstants.USER_KEY, JSON.stringify(user))
    this._user = user
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
      case AccountConstants.ACCOUNT_SUCCESS:
        this.setAccount(action.response)

        this.emitChange()
        break
      case AccountConstants.ACCOUNT_FAILURE:
        this.logout()
        this._errors = action.error.errors

        this.emitChange()
        break
      case AccountConstants.ACCOUNT_CREATE_SUCCESS:
        this.setAccount(action.response)

        this.emitChange()
        break
      case AccountConstants.ACCOUNT_CREATE_FAILURE:
        this._errors = action.error.errors

        this.emitChange()
        break
      case AccountConstants.ACCOUNTUSER_CREATE_SUCCESS:
        this.setUser(action.response)

        this.emitChange()
        break
      case AccountConstants.ACCOUNTUSER_CREATE_FAILURE:
        this._errors = action.error.errors

        this.emitChange()
        break
      case AccountConstants.LOGIN_SUCCESS:
        let user = action.response

        this.setAccount({ id: user.accountId, token: user.accountToken })
        this.setUser(action.response)

        this.emitChange()
        break
      case AccountConstants.LOGIN_FAILURE:
        this._errors = action.error.errors

        this.emitChange()
        break
      case AccountConstants.LOGOUT_SUCCESS:
        this.logout()

        this.emitChange()
        break
      default:
      // nothing to do
    }
  }
}

export default new AccountStore
