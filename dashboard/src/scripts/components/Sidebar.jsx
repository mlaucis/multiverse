import React, { Component, PropTypes } from 'react'
import { Link } from 'react-router'

export default class Sidebar extends Component {
  static contextTypes = {
    router: PropTypes.func.isRequired
  };

  static propTypes = {
    activeApp: PropTypes.string,
    sections: PropTypes.array.isRequired
  };

  render() {
    let router = this.context.router
    let sections = this.props.sections.map( (section) => {
      let icon = require(`../../icons/Sidebar_Icon_${section.icon}.svg`)
      let aClass = section.disabled ? 'inactive' : ''
      let lClass = ''

      if (router.isActive(section.route, section.params)) {
        lClass = 'active'
      }

      return (
        <li className={lClass} key={section.name}>
          <Link
            className={aClass}
            to={section.route}
            params={section.params}>
            <img src={icon}/>
            <span>{section.name}</span>
          </Link>
        </li>
      )
    })

    return (
      <ul className='sidebar'>
        {sections}
      </ul>
    )
  }
}
