'use strict';
module.exports = (body, client, bucket, archiveID, name, type) =>
    new Promise((resolve, reject) => {
        client.putObject({
            Bucket: bucket,
            Key: archiveID + '/' + name,
            Body: body,
            ACL: 'public-read',
            ContentType: type,
        }, (err, res) => {
            console.log(err,res);
            if (err) {
                reject(err);
            } else {
                resolve();
            }
        });
    });
