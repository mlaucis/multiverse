import React, { Component } from 'react'
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

    if (AccountStore.isAuthenticated) {
      requestAccount(AccountStore.user)
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize)
  }

  getState() {
    return {
      windowHeight: Math.max(window.innerHeight, document.body.clientHeight)
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
      }
      // {
      //   disabled: false,
      //   icon: 'bar-chart',
      //   name: 'Analytics',
      //   route: 'APPS'
      // },
      // {
      //   disabled: false,
      //   icon: 'settings',
      //   name: 'Settings',
      //   route: 'APPS'
      // }
    ]

    let style = {
      minHeight: `${this.state.windowHeight}px`
    }

    return (
      <section className='console'>
        <Header/>
        <div className="clearfix"></div>
        <div className="page-container">
          <Sidebar sections={sections}/>
          <div className="page-content-wrapper">
            <div className="page-content" style={style}>
              <RouteHandler/>
              <div className="clearfix"></div>
            </div>
            <div className="scroll-to-top">
              <i className="glyphicon glyphicon-chevron-up"></i>
            </div>
          </div>
        </div>
      </section>
    )
  }
}
