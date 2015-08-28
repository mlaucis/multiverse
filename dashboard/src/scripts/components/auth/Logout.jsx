import React, { Component, PropTypes } from 'react'

import AccountStore from '../../stores/AccountStore'
import { requestLogout } from '../../actions/ConsoleActionCreator'

export default class Logout extends Component {
  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  componentDidMount() {
    if (!AccountStore.isAuthenticated) {
      this.context.router.transitionTo('AUTH_LOGIN')
      return
    }

    requestLogout(AccountStore.user).then( () => {
      this.context.router.transitionTo('AUTH_LOGIN')
    })
  }

  render() {
    return <p>Logging out...</p>
  }
}
