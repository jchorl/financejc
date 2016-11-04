import React from 'react';

class Recurring extends React.Component {
  render () {
    return (
      <div>
        Recurring transactions
      </div>
    )
  }
}

export default class RecurringWrapper extends React.Component {
  render () {
    return (
      <div>
        <Recurring />
      </div>
    )
  }
}

