jest.dontMock('../TrackingStore')

import AccountConstants from '../../constants/AccountConstants'
import ApplicationConstants from '../../constants/ApplicationConstants'
import MemberConstants from '../../constants/MemberConstants'

describe('TrackingStore', () => {
  let ConsoleDispatcher
  let callback
  let trackMock = jest.genMockFunction()
  let app = {
    description: 'instant',
    id: 54321,
    name: 'Yo'
  }
  let user = {
    accountId: 123,
    email: 'nyn@cat.biz',
    id: 1234,
    firstName: 'cat',
    lastName: 'nyn',
    password: '4321'
  }

  Object.defineProperty(window, 'analytics', { value: { track: trackMock } })

  beforeEach( () => {
    require('../TrackingStore')
    ConsoleDispatcher = require('../../dispatcher/ConsoleDispatcher')
    callback = ConsoleDispatcher.register.mock.calls[0][0]
  })

  it('registers a callback with the dispatcher', () => {
    expect(ConsoleDispatcher.register).toBeCalled()
  })

  it('tracks ACCOUNTUSER_CREATE_FAILURE', () => {
    callback({
      type: AccountConstants.ACCOUNTUSER_CREATE_FAILURE,
      ...user
    })

    expect(trackMock).lastCalledWith('Member signed up', {
      eventId: 1,
      email: user.email,
      firstName: user.firstName,
      lastName: user.lastName,
      organizationId: user.accountId,
      success: false
    })
  })

  it('tracks ACCOUNTUSER_CREATE_SUCCESS', () => {
    callback({
      type: AccountConstants.ACCOUNTUSER_CREATE_SUCCESS,
      response: { id: user.id },
      ...user
    })

    expect(trackMock).lastCalledWith('Member signed up', {
      eventId: 1,
      email: user.email,
      firstName: user.firstName,
      lastName: user.lastName,
      memberId: user.id,
      organizationId: user.accountId,
      success: true
    })
  })

  it('tracks LOGIN_FAILURE', () => {
    callback({
      type: AccountConstants.LOGIN_FAILURE
    })

    expect(trackMock).lastCalledWith('Member logged in', {
      eventId: 2,
      success: false
    })
  })

  it('tracks LOGIN_SUCCESS', () => {
    callback({
      type: AccountConstants.LOGIN_SUCCESS,
      response: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member logged in', {
      eventId: 2,
      memberId: user.id,
      success: true
    })
  })

  it('tracks LOGOUT_FAILURE', () => {
    callback({
      type: AccountConstants.LOGOUT_FAILURE,
      user: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member logged out', {
      eventId: 3,
      memberId: user.id,
      success: false
    })
  })

  it('tracks LOGOUT_SUCCESS', () => {
    callback({
      type: AccountConstants.LOGOUT_SUCCESS,
      user: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member logged out', {
      eventId: 3,
      memberId: user.id,
      success: true
    })
  })

  it('tracks APP_CREATE_FAILURE', () => {
    callback({
      type: ApplicationConstants.APP_CREATE_FAILURE,
      name: app.name,
      description: app.description
    })

    expect(trackMock).lastCalledWith('Application created', {
      eventId: 7,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: false
    })
  })

  it('tracks APP_CREATE_SUCCESS', () => {
    callback({
      type: ApplicationConstants.APP_CREATE_SUCCESS,
      response: { id: app.id },
      name: app.name,
      description: app.description
    })

    expect(trackMock).lastCalledWith('Application created', {
      eventId: 7,
      appId: app.id,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: true
    })
  })

  it('tracks APP_EDIT_FAILURE', () => {
    callback({
      type: ApplicationConstants.APP_EDIT_FAILURE,
      description: app.description,
      id: app.id,
      name: app.name
    })

    expect(trackMock).lastCalledWith('Application edited', {
      eventId: 8,
      appId: app.id,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: false
    })
  })

  it('tracks APP_EDIT_SUCCESS', () => {
    callback({
      type: ApplicationConstants.APP_EDIT_SUCCESS,
      description: app.description,
      id: app.id,
      name: app.name
    })

    expect(trackMock).lastCalledWith('Application edited', {
      eventId: 8,
      appId: app.id,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: true
    })
  })

  it('tracks APP_DELETE_FAILURE', () => {
    callback({
      type: ApplicationConstants.APP_DELETE_FAILURE,
      id: app.id
    })

    expect(trackMock).lastCalledWith('Application deleted', {
      eventId: 9,
      appId: app.id,
      success: false
    })
  })

  it('tracks APP_DELETE_SUCCESS', () => {
    callback({
      type: ApplicationConstants.APP_DELETE_SUCCESS,
      id: app.id
    })

    expect(trackMock).lastCalledWith('Application deleted', {
      eventId: 9,
      appId: app.id,
      success: true
    })
  })

  it('tracks MEMBER_CREATE_FAILURE', () => {
    callback({
      type: MemberConstants.MEMBER_CREATE_FAILURE
    })

    expect(trackMock).lastCalledWith('Member created', {
      eventId: 10,
      success: false
    })
  })

  it('tracks MEMBER_CREATE_SUCCESS', () => {
    callback({
      type: MemberConstants.MEMBER_CREATE_SUCCESS,
      response: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member created', {
      eventId: 10,
      memberId: user.id,
      success: true
    })
  })

  it('tracks MEMBER_DELETE_FAILURE', () => {
    callback({
      type: MemberConstants.MEMBER_DELETE_FAILURE,
      id: user.id
    })

    expect(trackMock).lastCalledWith('Member deleted', {
      eventId: 11,
      memberId: user.id,
      success: false
    })
  })

  it('tracks MEMBER_DELETE_SUCCESS', () => {
    callback({
      type: MemberConstants.MEMBER_DELETE_SUCCESS,
      id: user.id
    })

    expect(trackMock).lastCalledWith('Member deleted', {
      eventId: 11,
      memberId: user.id,
      success: true
    })
  })

  it('tracks MEMBER_INVITE_FAILURE', () => {
    callback({
      type: MemberConstants.MEMBER_INVITE_FAILURE,
      email: user.email
    })

    expect(trackMock).lastCalledWith('Member invited', {
      eventId: 12,
      email: user.email,
      success: false
    })
  })

  it('tracks MEMBER_INVITE_SUCCESS', () => {
    callback({
      type: MemberConstants.MEMBER_INVITE_SUCCESS,
      email: user.email
    })

    expect(trackMock).lastCalledWith('Member invited', {
      eventId: 12,
      email: user.email,
      success: true
    })
  })
})
