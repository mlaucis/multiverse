import React, { Component } from 'react'
import Router from 'react-router'
import { Link } from 'react-router'

let { RouteHandler } = Router

import Logo from './Logo'

export default class Auth extends Component {
  render() {
    return (
      <section className='auth'>
        <div className='grid grid--justify-center'>
          <div className='logo grid__col-md-3'>
            <Logo/>
          </div>
          <nav className='grid__col-md-9'>
            <ul className='grid grid--justify-space-around'>
              <li><a href='https://tapglue.com/news-feed/'>Product</a></li>
              <li><a href='https://tapglue.com/pricing/'>Pricing</a></li>
              <li><a href='https://developers.tapglue.com'>Docs</a></li>
              <li><a href='https://tapglue.com/blog'>Blog</a></li>
              <li><a href='https://tapglue.com/about-us/'>About Us</a></li>
              <li className='actions'>
                <Link
                  to='AUTH_LOGIN'
                  className='btn-default outline small white'>
                  Login
                </Link>
                <Link
                  to='AUTH_SIGNUP'
                  className='btn-default outline small white'>
                  Signup
                </Link>
              </li>
            </ul>
          </nav>
          <section className='grid__col-sm-5 grid--align-content-space-between'>
            <RouteHandler/>
          </section>
        </div>
      </section>
    )
  }
}
