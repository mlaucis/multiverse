import React, { Component } from 'react'
import { findDOMNode } from 'react-dom'
import Router from 'react-router'

let { RouteHandler } = Router

import AccountStore from '../stores/AccountStore'
import { requestAccount } from '../actions/ConsoleActionCreator'

import Header from './Header'
import Sidebar from './Sidebar'

export default class Console extends Component {
  constructor() {
    super()

    this.state = this.getState()
  }

  componentDidMount() {
    window.addEventListener('resize', this.handleResize)

    if (this.state.windowHeight != this.getHeight) {
      this.setState(this.getState())
    }

    if (AccountStore.isAuthenticated) {
      requestAccount(AccountStore.user)
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize)
  }

  getHeight() {
    return Math.max(
      document.body.clientHeight,
      document.body.offsetHeight,
      document.body.scrollHeight,
      window.innerHeight,
    )
  }

  getState() {
    return {
      // FIXME(xla): Find a reliable way to determine the real height.
      windowHeight: this.getHeight()
    }
  }

  handlResize() {
    this.setState(this.getState())
  }

  render() {
    let sections = [
      {
        disalbed: false,
        icon: 'Home',
        name: 'Home',
        route: 'DASHBOARD'
      },
      {
        disabled: false,
        icon: 'Apps',
        name: 'Apps',
        route: 'APPS'
      },
      {
        disabled: false,
        icon: 'Members',
        name: 'Members',
        route: 'MEMBERS'
      },
      {
        disabled: false,
        icon: 'Analytics',
        name: 'Analytics',
        route: 'ANALYTICS'
      },
      {
        disabled: false,
        icon: 'LogOut',
        name: 'Logout',
        route: 'AUTH_LOGOUT'
      }
    ]

    let style = {
      minHeight: `${this.state.windowHeight}px`
    }

    return (
      <section className='console grid'>
        <Header/>
        <div className='sidebar-container grid__col-md-2 grid__col--bleed'>
          <Sidebar sections={sections}/>
        </div>
        <div className='content grid__col-md-10 grid__col--bleed' style={style}>
          <RouteHandler ref='route'/>
        </div>
      </section>
    )
  }
}
