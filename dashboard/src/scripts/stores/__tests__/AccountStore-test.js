jest.dontMock('../AccountStore')

import AccountConstants from '../../constants/AccountConstants'

describe('AccountStore', () => {
  let AccountStore
  let ConsoleDispatcher
  let callback
  let account = {
    id: 4321,
    name: 'CatApp',
    description: 'Make your cats social',
    token: 10001,
    enabled: true,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    receivedAt: Date.now()
  }
  let user = {
    id: 1234,
    accountId: 4321,
    accountToken: 'tokenString',
    token: 3214,
    firstName: 'Nyan',
    lastName: 'Cat'
  }

  beforeEach(() => {
    AccountStore = require('../AccountStore')
    ConsoleDispatcher = require('../../dispatcher/ConsoleDispatcher')
    callback = ConsoleDispatcher.register.mock.calls[0][0]

    localStorage.clear()
    AccountStore.init()
  })

  it('registers a callback with the dispatcher', () => {
    expect(ConsoleDispatcher.register).toBeCalled()
  })

  it('initializes without data', () => {
    expect(AccountStore.account).toBeUndefined()
    expect(AccountStore.accountName).toBe('')
    expect(AccountStore.errors.length).toBe(0)
    expect(AccountStore.isAuthenticated).toBeFalsy()
    expect(AccountStore.user).toBeUndefined()
  })

  it('restores state from localStorage', () => {
    localStorage.setItem(AccountConstants.ACCOUNT_KEY, JSON.stringify(account))
    localStorage.setItem(AccountConstants.USER_KEY, JSON.stringify(user))

    AccountStore.init()

    expect(AccountStore.account).toEqual(account)
    expect(AccountStore.accountName).toEqual(account.name)
    expect(AccountStore.errors.length).toBe(0)
    expect(AccountStore.isAuthenticated).toBeTruthy()
    expect(AccountStore.user).toEqual(user)
  })

  it('reacts to ACCOUNT_SUCCESS', () => {
    callback({
      type: AccountConstants.ACCOUNT_SUCCESS,
      response: account
    })

    expect(AccountStore.account).toEqual(account)
    expect(AccountStore.accountName).toEqual(account.name)
    expect(AccountStore.errors.length).toBe(0)
  })

  it('reacts to ACCOUNT_FAILURE', () => {
    let err = new Error('Bad Request')
    err.errors = [
      { code: 4004, message: 'Invalid Account Credentials' }
    ]

    callback({
      type: AccountConstants.ACCOUNT_FAILURE,
      error: err
    })

    expect(AccountStore.account).toBeUndefined()
    expect(AccountStore.errors).toEqual(err.errors)
    expect(AccountStore.isAuthenticated).toBeFalsy()
  })

  it('reacts to ACCOUNT_CREATE_SUCCESS', () => {
    callback({
      type: AccountConstants.ACCOUNT_CREATE_SUCCESS,
      response: account
    })

    expect(AccountStore.account).toEqual(account)
    expect(AccountStore.accountName).toEqual(account.name)
    expect(AccountStore.errors.length).toBe(0)
  })

  it('reacts to ACCOUNT_CREATE_FAILURE', () => {
    let err = new Error('Bad Request')
    err.errors = [
      { code: 4001, message: 'Generic Auth error' }
    ]

    callback({
      type: AccountConstants.ACCOUNT_CREATE_FAILURE,
      error: err
    })

    expect(AccountStore.account).toBeUndefined()
    expect(AccountStore.errors).toEqual(err.errors)
    expect(AccountStore.isAuthenticated).toBeFalsy()
  })

  it('reacts to LOGIN_SUCCESS', () => {
    callback({
      type: AccountConstants.LOGIN_SUCCESS,
      response: user
    })

    expect(AccountStore.errors.length).toBe(0)
    expect(AccountStore.isAuthenticated).toBeTruthy()
    expect(AccountStore.user).toEqual(user)

    let storageUser = localStorage.getItem(AccountConstants.USER_KEY)
    expect(JSON.parse(storageUser)).toEqual(user)
  })

  it('reacts to LOGIN_FAILURE', () => {
    let err = new Error('Bad request')
    err.errors = [
      { code: 4011, message: 'username and password do not match' }
    ]

    callback({
      type: AccountConstants.LOGIN_FAILURE,
      error: err
    })

    expect(AccountStore.errors).toEqual(err.errors)
    expect(AccountStore.isAuthenticated).toBeFalsy()
    expect(AccountStore.user).toBeUndefined()
  })

  it('reacts to LOGOUT_SUCCESS', () => {
    localStorage.setItem(AccountConstants.USER_KEY, JSON.stringify(user))

    AccountStore.init()

    callback({
      type: AccountConstants.LOGOUT_SUCCESS
    })

    expect(AccountStore.errors.length).toBe(0)
    expect(AccountStore.isAuthenticated).toBeFalsy()
    expect(AccountStore.user).toBeUndefined()
    expect(localStorage.getItem(AccountConstants.USER_KEY)).toBeUndefined()
  })
})
