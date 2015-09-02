import React, { Component } from 'react'
import Router from 'react-router'

let { RouteHandler } = Router

import Logo from './Logo'

export default class Auth extends Component {
  render() {
    return (
      <section className='auth'>
        <div className='inner'>
          <header>
            <div className='logo'>
              <Logo/>
            </div>
            <nav>
              <ul>
                <li><a href='//developers.tapglue.com'>Docs</a></li>
                <li><a href='//tapglue.com/blog'>Blog</a></li>
                <li><a href='//tapglue.com/about-us'>About Us</a></li>
              </ul>
            </nav>
          </header>
          <RouteHandler/>
        </div>
      </section>
    )
  }
}
