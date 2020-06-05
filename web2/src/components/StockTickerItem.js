import React from 'react';
import axios from 'axios';

// getDelta takes a team object and returns the difference betwen 
// the two most recent values as a span object. If there are not 
// two values for the team, 0 is returned.
function getDelta(team) {
    if(team.values.length <= 1) {
        return(
            <span>-</span>
        );
    }
    let delta = team.values[0] - team.values[1];
    return (delta < 0)
        ? <span className="dMinus">{Math.abs(delta)}</span>
        : <span className="dPlus">{delta}</span>
}

// latestValue takes a team object and returns the most recent stock
// value the team has had.
function latestValue(team) {
    if(team.values.length === 0) {
        return 0;
    }
    return team.values[0];
}

// StockTickerItem is a react component that displays current team 
// information. Meant to be composed in the StockTicker component.
export class StockTickerItem extends React.Component {
    constructor(props) {
        super(props);

        this.state = { team: null };
    }

    // make api call after component mounts to grab team info
    componentDidMount() {
        const id = this.props.team;
        const endpoint = 'http://localhost:8080/api/v2/teams/' + id;
        axios.get(endpoint).then(res => {
            this.setState({ team: res.data });
        }).catch(err => {
            console.log(err);
        });
    }

    render() {
        let content = (
            <span className="mono">
                <div className="StockTicker-item-code" 
                    style={{visibility: `hidden`}}>
                    GG
                </div>
            </span>
        );

        if(this.state.team != null) {
            content = (
                <span className="mono">
                    <div className="StockTicker-item-code">
                        {this.state.team.code}
                    </div>
                    <div className="StockTicker-item-delta">
                        {getDelta(this.state.team)}
                    </div>
                    <div className="StockTicker-item-value">
                        {latestValue(this.state.team)}
                    </div>
                </span>
            );
        }

        return(
            <div className="StockTicker-item">
                {content}
            </div>
        );
    }
}