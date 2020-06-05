import React from 'react';
import axios from 'axios';
import {CustomSelect} from './CustomSelect.js';
import {StockTickerItem} from './StockTickerItem.js';

// StockTicker displays StockTickerItems for each of the currently
// selected League and Season.
export class StockTicker extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            teams: [],
            leagues: [],
            splits: [],
            selectedLeague: 0,
            selectedSplit: 0,
        };

        this.updateLeague = this.updateLeague.bind(this);
        this.updateSplit = this.updateSplit.bind(this);
    }

    // called when the value for CustomSelect changes
    updateLeague(value) {
        this.setState({ selectedLeague: value });
    }
    updateSplit(value) {
        this.setState({ selectedSplit: value }, () => {
            const splitId = this.currentSplit().id;
            axios.get('http://localhost:8080/api/v2/splits/' + splitId)
                .then((r) => {
                    this.setState({ teams: r.data });
                }).catch((err) => {
                    console.log(err);
                });
        });
    }

    // helpers to fetch the current league/split selected from the
    // list of leagues and splits in state
    currentLeague() {
        if(this.state.leagues.length === 0) {
            return { id: 0, name: "fetching..." };
        } else {
            return this.state.leagues[this.state.selectedLeague];
        }
    }
    currentSplit() {
        if(this.state.splits.length === 0) {
            return { id: 0, name: "fetching..." };
        } else {
            return this.state.splits[this.state.selectedSplit];
        }
    }

    // make API calls after the componenet mounts
    componentDidMount() {
        axios.get('http://localhost:8080/api/v2/leagues').then(res => {
            this.setState({ leagues: res.data });
        }).catch(err => {
            console.log(err);
        });

        axios.get('http://localhost:8080/api/v2/splits').then(res => {
            const splits = res.data.map(s => {
                s.name = s.year + ' ' + s.split;
                return s;
            });
            this.setState({ splits }, () => {
                if(this.state.splits.length > 0) {
                    // if there are splits, pick the first one by default
                    this.updateSplit(0);
                }
            });
        }).catch(err => {
            console.log(err);
        });
    }

    render() {
        let listItems = <p>No teams found in selected split.</p>;
        if(this.state.teams.length > 0) {
            listItems = this.state.teams.map((team) =>
                <StockTickerItem key={team} team={team} />
            );
        }

        return(
            <aside className="StockTicker">
                <div className="StockTicker-CustomSelect-row">
                    <CustomSelect
                        label="LEAGUE"
                        options={this.state.leagues}
                        onValueChange={this.updateLeague} />
                    <CustomSelect
                        label="SEASON"
                        options={this.state.splits} 
                        onValueChange={this.updateSplit} />
                </div>
                {listItems}
            </aside>
        );
    }
}