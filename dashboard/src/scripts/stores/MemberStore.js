import { EventEmitter } from 'events'

import { register } from '../dispatcher/ConsoleDispatcher'

import MemberConstants from '../constants/MemberConstants'

const CHANGE_EVENT = 'change'

class MemberStore extends EventEmitter {
  constructor() {
    super()

    this.dispatchToken = register(this._handleActions)
    this._errors = []
    this._members = {}
  }

  get errors() {
    return this._errors
  }

  get members() {
    let members = []

    for (let id in this._members) {
      members.push(this._members[id])
    }

    return members.filter( member => {
      // return member.enabled
      return true
    }).sort( (left, right) => {
      let a = new Date(left.createdAt)
      let b = new Date(right.createdAt)

      if (a > b) {
        return -1
      }

      if (a < b) {
        return 1
      }

      return 0
    })
  }

  getMemberByID(id) {
    return this._members[id]
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
      case MemberConstants.MEMBER_SUCCESS:
        let member = action.response

        this._members[member.id] = member
        this._errors = []

        this.emitChange()
        break
      case MemberConstants.MEMBERS_SUCCESS:
        this._errors = []
        this._members = {}

        action.response.accountUsers.forEach( (m) => {
          this._members[m.id] = m
        })

        this.emitChange()
        break
      case MemberConstants.MEMBER_CREATE_SUCCESS:
        let m = action.response

        this._errors = []
        this._members[m.id] = m

        this.emitChange()
        break
      case MemberConstants.MEMBER_CREATE_FAILURE:
        this._errors = action.error.errors

        this.emitChange()
        break
      case MemberConstants.MEMBER_DELETE_SUCCESS:
        delete this._members[action.id]
        this._errors = []

        this.emitChange()
        break
      case MemberConstants.MEMBER_DELETE_FAILURE:
        this._errors = action.error.errors

        this.emitChange()
        break
      default:
      // nothing to do
    }
  }
}

export default new MemberStore
