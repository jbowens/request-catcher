import * as React from 'react'

export class RequestDetail extends React.Component {
  public render() {
    let dollars = ~~(this.props.value / 100);
    let cents = this.props.value % 100;

    let dollars_with_commas = dollars.toLocaleString();

    return (
      <div className="request">
        <pre>
          { this.props.request.raw_request }
        </pre>
      </div>
    )
  }
}
