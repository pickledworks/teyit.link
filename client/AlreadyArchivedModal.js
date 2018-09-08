import React, {Component} from "react";
import moment from "moment";
import Modal from "react-modal";
import Api from './api';

import "moment/locale/tr";

export default class AlreadyArchivedModal extends Component {
    constructor () {
        super();
        this.handleCloseModal = this.handleCloseModal.bind(this);
        this.force = this.force.bind(this);
    }


    handleCloseModal () {
        this.setState({ showModal: false });
    }

    componentWillMount() {
        const archivedAt = moment(this.props.archived_at).fromNow();
        this.state = {archivedAt: archivedAt, showModal: true, createInProgress: false};

        if (history.replaceState) {
            history.replaceState({}, document.title, location.href.replace("?fresh=false", ""));
        }
    }

    force() {
        this.setState({createInProgress: true});
        Api.CreateArchiveAndRedirect(this.props.request_url, true);
    }

    render() {
        const {archivedAt, showModal, createInProgress} = this.state;
        return (
            <div>
                <Modal isOpen={showModal} style={{
                    content: {
                        top: "initial", left: "40px",
                        right: "40px", bottom: "100px",
                        background: "rgb(255, 255, 255, 0.9)",
                        "border-radius": "10px",
                    }
                }}>
                    <h2>Bu sayfa en son <i>{archivedAt}</i> arşivlenmiş.</h2>
                    <p>isterseniz tekrar arşivleyebilirsiniz.</p>
                    <button onClick={this.force} disabled={createInProgress}>
                        {createInProgress ? "yeniden arşivleniyor..." : "tekrar arşivle"}</button>
                    <button onClick={this.handleCloseModal} disabled={createInProgress}>devam et</button>
                </Modal>
            </div>
        )
    }
}
