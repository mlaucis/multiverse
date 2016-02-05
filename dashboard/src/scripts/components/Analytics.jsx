import d3 from 'd3'
import React, { Component, PropTypes } from 'react'
import { findDOMNode } from 'react-dom'
import ReactFauxDom from 'react-faux-dom'
import { DateRange } from 'react-date-range'
import Dropdown from 'react-dropdown'
import moment from 'moment'

import AnalyticsConstants from '../constants/AnalyticsConstants'

import AccountStore from '../stores/AccountStore'
import AnalyticsStore from '../stores/AnalyticsStore'
import ApplicationStore from '../stores/ApplicationStore'

import { requestApps } from '../actions/ConsoleActionCreator'
import { requestMetrics } from '../actions/ConsoleActionCreator'

import LoaderIcon from '../../icons/Console_Animated_Loader.svg?t=custom'

const FORMAT = AnalyticsConstants.TIME_FORMAT_RANGE
const COLOR_START = '#2d474c'
const COLOR_END = '#04c298'
const TIME_FORMAT_D3 = '%Y-%m-%d'
const iconsDimensions = {
  'connections': 'ui-2_node',
  'events': 'ui-2_favourite-31',
  'objects': 'ui-2_chat',
  'users': 'users_multiple-11',
}
const ranges = {
  'Last 7 days': {
    startDate: (now) => {
      return now.subtract(7, 'days')
    },
    endDate: (now) => {
      return now
    },
  },
  'Last 30 days': {
    startDate: (now) => {
      return now.subtract(30, 'days')
    },
    endDate: (now) => {
      return now
    },
  },
  'Last 90 days': {
    startDate: (now) => {
      return now.subtract(90, 'days')
    },
    endDate: (now) => {
      return now
    },
  },
  'Last Month': {
    startDate: (now) => {
      return now.subtract(1, 'months').startOf('month')
    },
    endDate: (now) => {
      return now.subtract(1, 'months').endOf('month')
    },
  },
  'Month to Date': {
    startDate: (now) => {
      return now.startOf('month')
    },
    endDate: (now) => {
      return now
    },
  }
}
const rangeTheme = {
  Calendar: {
    margin: '0 1rem',
    padding: 0,
  },
  DayHover: {
    background: '#f46c7c',
    color: '#fff',
  },
  DayInRange: {
    background: '#24383c',
    color: '#eaeaea',
  },
  DayPassive: {
    opacity: 0,
  },
  DaySelected: {
    background: '#f25467',
  },
  MonthAndYear: {
    padding: '0rem 0rem 1rem 0rem',
  },
  MonthArrowPrev: {
    borderRightColor: '#24383c',
  },
  MonthArrowNext: {
    borderLeftColor: '#24383c',
  },
  MonthButton: {
    background: '0',
  },
  PredefinedRanges: {
    width: '11rem',
  }
}

export default class Analytics extends Component {
  constructor() {
    super()

    this.app = undefined
    this.deselected = {}
    this.eWidth = 0
    this.highlight = -1
    this.range = {
      endDate: moment(),
      startDate: moment().subtract(30, 'days'),
    }
    this.refreshing = true
    this.showRange = false
    this.state = this.getState()
  }

  componentDidMount() {
    this.eWidth = (findDOMNode(this.refs.chart) || {}).offsetWidth || 0

    AnalyticsStore.addChangeListener(this.handleChange)
    ApplicationStore.addChangeListener(this.handleChange)

    window.addEventListener('resize', this.handleResize)

    this.handleResize()
    requestApps(AccountStore.user)
  }

  componentWillUnmount() {
    AnalyticsStore.removeChangeListener(this.handleChange)
    ApplicationStore.removeChangeListener(this.handleChange)
    window.removeEventListener('resize', this.handleResize)
  }

  color = (idx) => {
    let keys = Object.keys(this.state.sets)

    return d3
      .scale
      .linear()
      .domain([0, keys.length -1])
      .range([COLOR_START, COLOR_END])(idx)
  };

  getState() {
    let width = Math.max(window.innerWidth, document.body.clientWidth) - 320

    if (this.eWidth > 0) {
      width = this.eWidth
    }

    let id = +new Date()

    return {
      app: this.app,
      apps: ApplicationStore.apps.map((app, idx) => {
        return { label: app.name, value: app.id, id: id }
      }),
      deselected: this.deselected,
      highlight: this.highlight,
      labels: AnalyticsStore.labels,
      range: this.range,
      refreshing: this.refreshing,
      sets: AnalyticsStore.dimensions,
      showRange: this.showRange,
      width: width,
    }
  }

  handleAppSelect = (app) => {
    this.app = app

    this.handleRefresh()

    this.setState(this.getState())
  };

  handleChange = () => {
    this.setState(this.getState())

    if (!this.state.app && this.state.apps.length > 0) {
      let app = this.state.apps[0]
      let end = this.range.endDate.format(FORMAT)
      let start = this.range.startDate.format(FORMAT)

      this.app = app

      requestMetrics(app.value, start, end).catch(err => {
        console.log(err)
        console.log(err.errors)
      })

      this.setState(this.getState())

      return
    }

    if (this.state.refreshing) {
      this.refreshing = false

      this.setState(this.getState())
    }
  };

  handleClick = (event) => {
    if (this.showRange) {
      if (!findDOMNode(this.refs['daterange']).contains(event.target) &&
          !findDOMNode(this.refs['date']).contains(event.target)) {
        this.showRange = false

        this.setState(this.getState())
      }
    }
  };

  handleDimensionDeselect = (d) => {
    let parent = this

    return func => {
      if (parent.deselected[d]) {
        delete parent.deselected[d]
      } else {
        parent.deselected[d] = true
        this.highlight = -1
      }
      parent.setState(parent.getState())
    }
  };

  handleHighlight = (idx) => {
    return () => {
      this.highlight = idx
      this.setState(this.getState())
    }
  };

  handleRangeChange = (range) => {
    this.range = range

    this.setState(this.getState())
  };

  handleRangeToggle = () => {
    this.showRange = !this.showRange

    this.setState(this.getState())
  };

  handleRefresh = () => {
    this.refreshing = true
    this.showRange = false

    this.setState(this.getState())

    let end = this.range.endDate.format(FORMAT)
    let start = this.range.startDate.format(FORMAT)

    requestMetrics(this.state.app.value, start, end).catch(err => {
      console.log(err)
      console.log(err.errors)
    })
  };

  handleResize = () => {
    this.setState(this.getState())
  };

  viewDimension = (d, idx) => {
    let c = 'dimension'
    let col = this.color(idx)
    let deselected = this.state.deselected
    let iconClass = `nc-icon-glyph ${iconsDimensions[d]}`
    let sum = this.state.sets[d].reduce((p, c) => { return p +c })
    let setHighlight = deselected.hasOwnProperty(d) ? (
      func => {}
    ) : (
      this.handleHighlight
    )
    let dimensionName = d

    if (dimensionName == 'objects') {
      dimensionName = 'posts & comments'
    }

    if (dimensionName == 'users') {
      dimensionName = `new ${dimensionName}`
    }

    if (idx == this.state.highlight) {
      c = `${c} highlight`
      col = '#f46c7c'
    }

    let style = {}

    if (deselected.hasOwnProperty(d)) {
      style['backgroundColor'] = '#ddd'
      style['color'] = col
    } else {
      style['backgroundColor'] = col
    }

    c = `${c} grid__col-3`

    return (
      <div
        className={c}
        key={d}
        onClick={this.handleDimensionDeselect(d)}
        onMouseEnter={setHighlight(idx)}
        onMouseLeave={setHighlight(-1)}
        style={style}>
        <div className='grid grid--bleed'>
          <div className='grid__col-2 grid--align-self-center grid--justify-center'>
            <span className={iconClass}/>
          </div>
          <div className='grid__col-10'>
            <p className='sum'>{numberWithCommas(sum)}</p>
            <p>{dimensionName}</p>
          </div>
        </div>
      </div>
    )
  };

  viewRangePicker() {
    let c = ''
    if (this.state.showRange) {
      c = 'show'
    }

    return (
      <div className={`daterange ${c}`} ref='daterange'>
        <DateRange
          endDate={this.range.endDate}
          firstDayOfWeek={1}
          format={FORMAT}
          maxDate={moment()}
          onChange={this.handleRangeChange}
          startDate={this.range.startDate}
          ranges={ranges}
          theme={rangeTheme} />

        <div className='actions'>
          <buton
            className='btn-default'
            onClick={this.handleRefresh}>Apply</buton>
        </div>
      </div>
    )
  }

  render() {
    let highlight = this.state.highlight
    let color = this.color
    let deselected = this.state.deselected
    let dimensions = Object.keys(this.state.sets)
    let sets = dimensions.filter(
      d => !deselected[d]
    ).map(
      d => this.state.sets[d]
    ).filter(
      set => typeof set !== 'undefined'
    )
    let colors = dimensions.map((d, i) => {
      if (i == highlight) {
        return '#f46c7c'
      }
      if (deselected[d]) {
        return undefined
      }
      return color(i)
    }).filter(c => {
      return typeof c !== 'undefined'
    })
    let buttons = dimensions.map(this.viewDimension)
    let start = this.range.startDate.format(FORMAT)
    let end = this.range.endDate.format(FORMAT)

    return (
      <div
        id='analytics'
        className='grid grid--align-content-start'
        onClick={this.handleClick}>
        <div className='grid__col-1 label grid--align-self-center'>App:</div>
        <div className='grid__col-3 grid__col--bleed grid--align-self-center'>
          <Dropdown
            onChange={this.handleAppSelect}
            options={this.state.apps}
            value={this.state.app} />
        </div>
        <div className='grid__col-5 grid__col--bleed'>
          <div className='grid grid--align-center datepicker'>
            <div className='grid__col-2 label'>Dates:</div>
            <div className='grid__col-6'>
              <div className='date' onClick={this.handleRangeToggle} ref='date'>
                {start} - {end}
              </div>
            </div>
          </div>
        </div>
        <div className='grid__col-8 grid__col--bleed'>
          {this.viewRangePicker()}
        </div>
        <div className='dimensions grid__col-12 grid__col--bleed'>
          <div className='grid'>
            {buttons}
          </div>
        </div>
        <div className='chart grid__col-12' ref='chart'>
          <StackChart
            colors={colors}
            dimensions={dimensions}
            labels={this.state.labels}
            sets={sets}
            width={this.state.width}/>
        </div>
      </div>
    )
  }
}

class StackChart extends Component {
  static propTypes = {
    colors: PropTypes.array.isRequired,
    dimensions: PropTypes.array.isRequired,
    labels: PropTypes.array.isRequired,
    sets: PropTypes.array.isRequired,
    width: PropTypes.number.isRequired,
  };

  constructor() {
    super()

    this.showTooltip = false
    this.hovered = null
    this.state = this.getState()
  }

  getState() {
    return {
      hovered: this.hovered,
      showTooltip: this.showTooltip,
    }
  }

  handleOut = (rect, vals, bucket, d) => {
    this.showTooltip = false
    this.setState(this.getState())
  };

  handleOver = (rect, vals, bucket, d) => {
    this.showTooltip = true
    this.hovered = {dimension: d, element: rect, x: vals.x, y: vals.y}
    this.setState(this.getState())
  };

  render() {
    let colors = this.props.colors
    let dimensions = this.props.dimensions
    let labels = this.props.labels
    let sets = this.props.sets

    let handleOut = this.handleOut
    let handleOver = this.handleOver

    let margin = {top: 20, right: 40, bottom: 70, left: 80}
    let width = this.props.width - margin.left - margin.right
    let height = 400 - margin.bottom - margin.top

    let parseDate = d3.time.format(TIME_FORMAT_D3).parse
    let n = sets.length
    let m = labels.length
    let stack = d3.layout.stack()
    let layers = stack(d3.range(n).map(idx => {
      return sets[idx].map((v, idx) => {
        return {x: idx, y: v}
      })
    }))
    let yGroupMax = d3.max(layers, l => {
      return d3.max(l, d => {
        return d.y
      })
    })
    let yStackMax = d3.max(layers, l => {
      return d3.max(l, d => {
        return d.y0 + d.y
      })
    })


    let x = d3
      .scale.ordinal()
      .domain(labels.map(l => { return parseDate(l) }))
      .rangeRoundBands([0, width], 0.25, 0)
    let y = d3
      .scale.linear()
      .domain([0, yStackMax])
      .range([height, 0])

    let xAxis = d3
      .svg.axis()
      .scale(x)
      .tickFormat(d3.time.format('%b %d'))
      .tickSize(0)
      .tickPadding(25)
      .orient('bottom')

    let yAxis = d3
      .svg.axis()
      .scale(y)
      .orient('left')
      .innerTickSize(4)
      .ticks(6)

    let container = ReactFauxDom.createElement('div')
    let node = ReactFauxDom.createElement('svg')

    container.appendChild(node)

    let svg = d3
      .select(node)
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('transform', `translate(${margin.left}, ${margin.right})`)

    svg
      .append('g')
      .attr('class', 'x axis')
      .attr('transform', `translate(0, ${height})`)
      .call(xAxis)
      .selectAll('text')
      .style('text-anchor', 'end')
      .attr('dx', '-1rem')
      .attr('dy', '-1rem')
      .attr('transform', 'rotate(-45)')

    svg
      .append('g')
      .attr('class', 'y axis')
      .call(yAxis)
      .append('text')
      .attr('transform', 'rotate(-90)')
      .attr('y', 6)
      .attr('dy', '0.71rem')
      .style('text-anchor', 'end')

    let layer = svg
      .selectAll('.layer')
      .data(layers)
      .enter()
      .append('g')
      .attr('class', 'layer')

    layer
      .selectAll('rect')
      .data(d => { return d })
      .enter()
      .append('rect')
      .attr('x', (d, idx) => { return x(parseDate(labels[idx])) })
      .attr('y', d =>  y(d.y0 + d.y) )
      .attr('width', x.rangeBand())
      .attr('height', d => y(d.y0) - y(d.y0 + d.y) )
      .style('fill', (d, bucket, idx) => {
        if (this.state.showTooltip && this.state.hovered.dimension == idx) {
          if (this.hovered.x == d.x && this.hovered.y == d.y) {
            return '#f46c7c'
          }
        }
        return colors[idx]
      })
      .style('stroke', (d, bucket, idx) => {
        if (this.state.showTooltip && this.state.hovered.dimension == idx) {
          if (this.hovered.x == d.x && this.hovered.y == d.y) {
            return '#f46c7c'
          }
        }

        return colors[idx]
      })
      .on('mouseout', function (vals, bucket, d) {
        handleOut(this, vals, bucket, d)
      })
      .on('mouseover', function (vals, bucket, d) {
        handleOver(this, vals, bucket, d)
      })

    let tooltip = ReactFauxDom.createElement('div')
    tooltip.setAttribute('class', 'tooltip')

    let tooltipMetric = ReactFauxDom.createElement('p')
    tooltipMetric.setAttribute('class', 'metric')

    tooltip.appendChild(tooltipMetric)

    let tooltipDimension = ReactFauxDom.createElement('p')
    tooltipDimension.setAttribute('class', 'dimension')

    tooltip.appendChild(tooltipDimension)

    if (this.state.hovered) {
      let hovered = this.state.hovered
      let xPos = parseFloat(d3.select(hovered.element).attr('x'))
			let yPos = parseFloat(d3.select(hovered.element).attr('y'))
			let height = parseFloat(d3.select(hovered.element).attr('height'))

      tooltip.style.setProperty('left', xPos)
      tooltip.style.setProperty('top', yPos)

      tooltipMetric.innerHTML = numberWithCommas(hovered.y)
      tooltipDimension.innerHTML = `${dimensions[hovered.dimension]}`

      if (dimensions[hovered.dimension] == 'objects') {
        tooltipDimension.innerHTML = 'Posts & Comments'
      }

      if (dimensions[hovered.dimension] == 'users') {
        tooltipDimension.innerHTML = 'new users'
      }

      if (this.state.showTooltip) {
        tooltip.setAttribute('class', 'tooltip active')
      } else {
        tooltip.setAttribute('class', 'tooltip')
      }
    }

    container.appendChild(tooltip)

    return container.toReact()
  }
}

function numberWithCommas(x) {
  return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}
