import axios from "axios";

const Api = {};

Api.GetLocalData = (key, def = null) => {
    const keyInPage = window["TL_DATA"][key];
    if (keyInPage) {
        return keyInPage;
    }
    return def;
};

Api.CountPreviousArchives = (requestUrl, period = "24h") =>
    axios.get(`/api/count-previous-archives?request_url=${requestUrl}&period=${period}`)
        .then((response) => response.data);

Api.CreateArchive = (requestUrl, force) => {
    const data = new FormData();
    data.append('request_url', requestUrl);
    data.append('force', force);

    return axios.post("/api/archive", data)
        .then((response) => {
            const data = response.data;
            let redirectUrl = `/${data.slug}`;

            // If this wasn't just created, set the URL so we can show the modal
            if (response.status !== 201) {
                redirectUrl = `${redirectUrl}?fresh=false`;
                data.isFresh = false;
            }

            data.redirectUrl = redirectUrl;

            return data;
        });
};

Api.CreateArchiveAndRedirect = (requestUrl, force) =>
    Api.CreateArchive(requestUrl, force).then((data) => window.location = data.redirectUrl);

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
