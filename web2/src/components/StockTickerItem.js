import React from 'react';

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
export function StockTickerItem(props) {
    return(
        <div className="StockTicker-item">
            <span className="mono">
                <div className="StockTicker-item-code">{props.team.code}</div>
                <div className="StockTicker-item-delta">{getDelta(props.team)}</div>
                <div className="StockTicker-item-value">{latestValue(props.team)}</div>
            </span>
        </div>
    );
}