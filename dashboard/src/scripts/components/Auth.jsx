import React, { Component } from 'react'
import Router from 'react-router'

let { RouteHandler } = Router

import Logo from './Logo'

export default class Auth extends Component {
  render() {
    return (
      <section className='auth'>
        <header>
        <div className='logo'>
          <Logo/>
        </div>
        <nav>
          <ul>
            <li><a href='#'>Product</a></li>
            <li><a href='#'>Docs</a></li>
            <li><a href='#'>Blog</a></li>
            <li><a href='#'>About Us</a></li>
          </ul>
        </nav>
        </header>
        <RouteHandler/>
      </section>
    )
  }
}
