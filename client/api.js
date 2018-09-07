import axios from "axios";

const Api = {};
Api.CountPreviousArchives = (requestUrl, period = "24h") =>
    axios.get(`/api/count-previous-archives?request_url=${requestUrl}&period=${period}`)
        .then((response) => response.data);

Api.CreateArchive = (requestUrl) => {
    const data = new FormData();
    data.append('request_url', requestUrl);

    return axios.post("/api/archive", data)
        .then((response) => response.data);
};

Api.GetArchive = (slug) => axios.get(`/api/archives/${slug}`).then((resp) => resp.data);

Api.RefreshWhenArchived = (slug) => new Promise((resolve, reject) => {
    Api.GetArchive(slug).then((data) => {
        if (data.archived_at) {
            window.location.reload();
        } else {
            setTimeout(() => Api.RefreshWhenArchived(slug), 500);
        }
    })
});

export default Api;
