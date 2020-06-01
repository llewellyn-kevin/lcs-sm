import React from 'react';

function getDelta(team) {
    if(team.values.length <= 1) {
        return 0;
    }
    return team.values[0] - team.values[1];
}

function latestValue(team) {
    if(team.values.length === 0) {
        return 0;
    }
    return team.values[0];
}

export function StockTickerItem(props) {
    return(
        <div className="StockTicker-item">
            {props.team.code} 
            {getDelta(props.team)} 
            {latestValue(props.team)}
        </div>
    );
}