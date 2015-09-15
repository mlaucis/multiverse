import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'

export default class Sidebar extends Component {
  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  static propTypes = {
    activeApp: PropTypes.string,
    sections: PropTypes.array.isRequired
  }

  render() {
    let router = this.context.router
    let sections = this.props.sections.map( (section) => {
      let icon = require(`../../icons/Sidebar_Icon_${section.icon}.svg`)
      let c = ''

      if (router.isActive(section.route, section.params)) {
        c = 'active'
      }

      return (
        <li className={c} key={section.name}>
          <Link
            to={section.route}
            params={section.params}>
            <img src={icon}/>
            <span>{section.name}</span>
          </Link>
        </li>
      )
    })

    return (
      <div className="page-sidebar-wrapper">
        <div className="page-sidebar navbar-collapse collapse">
          <ul className="page-sidebar-menu page-sidebar-menu-hover-submenu">
            {sections}
            <li>
              <a className='inactive' href='#'>
                <img src={require('../../icons/Sidebar_Analytics.png')}/>
                <span className='title'>Analytics</span>
              </a>
            </li>
            <li>
              <a className='inactive' href='#'>
                <img src={require('../../icons/Sidebar_Settings.png')}/>
                <span className='title'>Settings</span>
              </a>
            </li>
          </ul>
        </div>
      </div>
    )
  }
}
