import React, {Component} from 'react';
import api from './api';

export default class ArchiveInput extends Component {
    constructor(props) {
        super(props);
        this.state = {requestUrl: '', processing: false, err: null};

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        this.setState({requestUrl: event.target.value});
    }

    handleSubmit(e) {
        this.setState({processing: true});
        e.preventDefault();

        const {requestUrl} = this.state;
        api.CreateArchiveAndRedirect(requestUrl).catch((err) => {
            this.setState({err: true, processing: false});
        });
    }

    render() {
        const {requestUrl, processing, err} = this.state;
        return (
            <form onSubmit={this.handleSubmit}>
                <div className={`form-group ${err ? "has-error has-feedback" : ""}`}>
                    <input className={`flipkart-navbar-input col-xs-10 fk-input`} type="url"
                           placeholder="Kaydetmek istediginiz adres..." name="request_url"
                           value={requestUrl} onChange={this.handleChange}
                           required={true} disabled={processing}
                    />
                    <button className="flipkart-navbar-button col-xs-2 fk-button"
                            disabled={this.state.processing} type="submit">
                        {processing ? "Arşivleniyor..." : "Arşivle"}
                    </button>

                    {err ? <div className="alert alert-danger archive-input-error" role="alert">
                        Bir hata oluştu, lütfen tekrar deneyin.
                    </div> : ""}
                </div>
            </form>
        )
    }
}
