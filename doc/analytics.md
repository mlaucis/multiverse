# Analytics

## Overview

We will build the first version of our analytics to display relevant metrics that can easily extracted from our current data model. The next iteration will be extend the data model in order to gather more granular usage information on a daily basis.

## Components

The solution will consists of the following components:

- Jarvis Query API
- Analytics API
- Analytics Visualization in Dashboard

### Jarvis Query API

The Jarvis query API will the fastest way to retrieve aggregation customer metrics based on a provided `customer namespace`. This will internally query the analytics API with the default time ranges:

- Month to date
- Last 30 days
- Overall

### Analytics API

The analytics API will be the core interface to retrieve analytics data. Analytics data should be provided on `application` and not `organization` level. Following is a  draft for the solution:

The design is slightly inspired by: https://keen.io/docs/api/

#### Data

The data retrieved in the first version will be:

- New users
- New connections
- New events
- New objects

for a given timeframe.

#### Request

**Endpoint**
`GET /:version/organizations/:orgId/applications/:appId/analytics?where={query}`

**Example**


```
curl -X "GET" "https://api.tapglue.com/0.4/organizations/73cd5fd8-c295-5f7d-bb5b-30ec6f531ac1/applications/3c0a0484-8b23-5c7d-b5c2-bc39ee237606/analytics?where={"timeframe": {"start": "2015-08-13T19:00:00.000Z","end":"2015-08-15T19:00:00.000Z"}} \
```

The main question is if we should be consistent with our query logic and stick to a proper `GET` request or if it would be ok to have the query parameter in the payload.

#### Response

**Example**

```json
{
  "data": [
    {
      "summary": {
        "unique_users": 3,
        "new_users": 3,
        "new_connections": 3,
        "new_events": 3,
        "new_objects": 3
      },
      "new_users": [
        {
          "date": "2015-08-13",
          "value": 3
        },
        {
          "date": "2015-08-14",
          "value": 3
        },
        {
          "date": "2015-08-15",
          "value": 3
        }
      ],
      "new_connections": [
        {
          "date": "2015-08-13",
          "value": 3
        },
        {
          "date": "2015-08-14",
          "value": 3
        },
        {
          "date": "2015-08-15",
          "value": 3
        }
      ],
      "new_events": [
        {
          "date": "2015-08-13",
          "value": 3
        },
        {
          "date": "2015-08-14",
          "value": 3
        },
        {
          "date": "2015-08-15",
          "value": 3
        }
      ],
      "new_objects": [
        {
          "date": "2015-08-13",
          "value": 3
        },
        {
          "date": "2015-08-14",
          "value": 3
        },
        {
          "date": "2015-08-15",
          "value": 3
        }
      ],
    }
  ],
  "timeframe": {
    "start": "2015-08-13T19:00:00.000Z",
    "end": "2015-08-15T19:00:00.000Z"
  }
}
```

Please challenge everything in the example above. I think this is the most critical part o the design. Many things are missing here such as the representation of a group by in the response. The idea is to have a data summary for the high level metrics in that timeframe and daily data for the selected timeframe.

Everything could be done with a single query (might not be so important for the dashboard). The summary data would appear in top boxes and the granular data in an histogram.

## Phases

The project will be divided into two phases.

### Phase 1

All of the specification above will be done in the first version of the API and the Dashboard.

### Phase 2

Additional counts to determine unique user activity on a daily bases will be added in phase 2. That is to be designed.

## Visualization

Here is a draft for the Information Architecture for the Dashboard. This is not designed in terms of:

- Shapes
- Fonts
- Colors
- Distances

![Dashboard Draft](http://s14.postimg.org/4z84h1i35/UX_draft.png)
