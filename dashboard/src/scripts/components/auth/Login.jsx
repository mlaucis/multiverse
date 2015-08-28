import React, { Component, PropTypes, findDOMNode } from 'react'

import { Link } from 'react-router'

import AccountStore from '../../stores/AccountStore'
import { requestAccount } from '../../actions/ConsoleActionCreator'
import { requestLogin } from '../../actions/ConsoleActionCreator'

export class LoginForm extends Component {
  static propTypes = {
    errors: PropTypes.array
  }

  handleSubmit = (event) => {
    event.preventDefault()

    let email = findDOMNode(this.refs.email).value
    let password = findDOMNode(this.refs.password).value

    requestLogin(email, password).then( user => {
      requestAccount(user)
    })
  }

  render() {
    let errors = this.props.errors.map( error => {
      return <p key={error.code} className='error'>{error.message}</p>
    })

    return (
      <form onSubmit={this.handleSubmit}>
        {errors}
        <div className='form-group'>
          <div className='group'>
            <input id='email'
                   name='email'
                   placeholder='Email'
                   ref='email'
                   type='email'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='email'>Email</label>
          </div>
          <div className='group'>
            <input id='password'
                   name='password'
                   placeholder='Password'
                   ref='password'
                   type='password'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='password'>Password</label>
          </div>
          <p>No account yet? <Link to='AUTH_SIGNUP'>Sign Up</Link></p>
        </div>
        <div className='actions'>
          <input
            className='btn-default block'
            ref='submit'
            type='submit'
            value='Log in'/>
        </div>
      </form>
    )
  }
}

export default class Login extends Component {
  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  static propTypes = {
    query: PropTypes.object
  }

  constructor() {
    super()

    this.state = this.getState()
  }

  getState() {
    return {
      errors: AccountStore.errors
    }
  }

  componentDidMount() {
    AccountStore.addChangeListener(this._onChange)
  }

  componentWillUnmount() {
    AccountStore.removeChangeListener(this._onChange)
  }

  _onChange = () => {
    this.setState(this.getState())

    if (AccountStore.isAuthenticated) {
      let next = this.props.query.nextPath || 'CONSOLE'

      this.context.router.transitionTo(next)
    }
  }

  render() {
    return (
      <section className='login'>
        <LoginForm errors={this.state.errors}/>
      </section>
    )
  }
}
