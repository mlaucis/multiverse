import React, { Component } from 'react'
import Router from 'react-router'

let { RouteHandler } = Router

import Logo from './Logo'

export default class Auth extends Component {
  render() {
    return (
      <section className='auth'>
        <div className='grid grid--justify-center'>
          <div className='logo grid__col-md-12 grid grid--justify-center'>
            <a href='https://tapglue.com'>
              <Logo/>
            </a>
          </div>
          <section className='grid__col-sm-5 grid--align-content-space-between'>
            <RouteHandler/>
          </section>
        </div>
      </section>
    )
  }
}
