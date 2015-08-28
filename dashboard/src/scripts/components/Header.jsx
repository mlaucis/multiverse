import React, { Component } from 'react'
import { Link } from 'react-router'

import AccountStore from '../stores/AccountStore'

import Logo from './Logo'

class Profile extends Component {
  constructor() {
    super()

    this.state = this.getState()
  }

  componentDidMount() {
    AccountStore.addChangeListener(this.handleChange)
  }

  componentWillUnmount() {
    AccountStore.removeChangeListener(this.handleChange)
  }

  getState() {
    return {
      user: AccountStore.user
    }
  }

  handleChange = () => {
    this.setState(this.getState())
  }

  handleClick = (event) => {
    event.preventDefault()

    this.setState(this.getState())
  }

  shouldComponentUpdate(props, state) {
    return !!state.user
  }

  render() {
    return (
      <div className='profile'>
        <Link to='AUTH_LOGOUT'>
          <span>
            {this.state.user.firstName}
          </span>
          <span className="glyphicon glyphicon-log-out"></span>
        </Link>
      </div>
    )
  }
}

class HeaderNav extends Component {
  render() {
    return (
      <nav className='header-nav'>
        <Profile/>
      </nav>
    )
  }
}

export default class Header extends Component {
  render() {
    return (
      <div className="page-header navbar navbar-fixed-top">
        <div className="page-header-inner">
          <Logo/>
          <a
            href='#'
            className="menu-toggler responsive-toggler"
            data-toggle="collapse"
            data-target=".navbar-collapse">
          </a>
          <div className="page-top">
            <HeaderNav/>
          </div>
        </div>
      </div>
    )
  }
}
