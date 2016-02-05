import React, { Component } from 'react'

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
  };

  handleClick = (event) => {
    event.preventDefault()

    this.setState(this.getState())
  };

  shouldComponentUpdate(props, state) {
    return !!state.user
  }

  render() {
    return (
      <div className='profile'>
        <span>
          {this.state.user.firstName}
        </span>
      </div>
    )
  }
}

export default class Header extends Component {
  render() {
    return (
      <header className='grid'>
        <div className='grid__col-md-2'>
          <Logo/>
        </div>
        <div className='grid__col-md-10 grid__col--bleed'>
          <nav className='grid grid--justify-end'>
            <Profile/>
          </nav>
        </div>
      </header>
    )
  }
}
