import React, { Component } from 'react'

import IconLogo from '../../icons/logo.svg?t=custom'

export default class Logo extends Component {
  render() {
    return (
      <div className="page-logo">
        {IconLogo}
      </div>
    )
  }
}
