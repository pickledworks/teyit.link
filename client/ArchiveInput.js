import {h, Component} from 'preact';
import api from './api';

export default class ArchiveInput extends Component {
    constructor(props) {
        super(props);
        this.state = {requestUrl: '', count: 0, lastArchivedAt: null};

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        this.setState({requestUrl: event.target.value});
    }

    handleSubmit(e) {
        e.preventDefault();

        const { requestUrl } = this.state;
        api.CountPreviousArchives(requestUrl, "24h").then(({count, last_archived_at}) => {
            let ask = true;
            if (count > 0) {
                ask = confirm(`Son 24 saatte bu adres ${count} kez linklenmis. Linklemek istediginize emin misiniz?`);
            }

            if (ask) {
                api.CreateArchive(this.state.requestUrl).then((resp) => window.location = `/${resp.slug}`);
            }
        });
    }

    render() {
        return (
            <form onSubmit={this.handleSubmit}>
                <div class="form-group">
                    <input className={`flipkart-navbar-input col-xs-10 fk-input`} type="url"
                           placeholder="Kaydetmek istediginiz adres..." name="request_url"
                           value={this.state.requestUrl} onChange={this.handleChange}
                           required={true}
                    />
                    <button className="flipkart-navbar-button col-xs-2 fk-button" type="submit">Kaydet</button>
                </div>
            </form>
        )
    }
}
