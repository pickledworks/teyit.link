import React, {Component} from 'react';
import api from './api';

export default class ArchiveInput extends Component {
    constructor(props) {
        super(props);
        this.state = {requestUrl: '', processing: false, error: null};

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        this.setState({requestUrl: event.target.value});
    }

    handleSubmit(e) {
        this.setState({processing: true});
        e.preventDefault();

        const { requestUrl } = this.state;
        api.CreateArchiveAndRedirect(requestUrl).catch(() => this.setState({error: true}));
    }

    render() {
        return (
            <form onSubmit={this.handleSubmit}>
                <div class="form-group">
                    <input className={`flipkart-navbar-input col-xs-10 fk-input`} type="url"
                           placeholder="Kaydetmek istediginiz adres..." name="request_url"
                           value={this.state.requestUrl} onChange={this.handleChange}
                           required={true} disabled={this.state.processing}
                    />
                    <button className="flipkart-navbar-button col-xs-2 fk-button"
                            disabled={this.state.processing} type="submit">
                        {this.state.processing ? "Arşivleniyor..." : "Arşivle"}
                    </button>
                </div>
            </form>
        )
    }
}
