'use strict';
const uploader = (body, client, bucket, archive_id, name, type) =>
    new Promise((resolve, reject) => {
        client.putObject({
            Bucket: bucket,
            Key: archive_id + '/' + name,
            Body: body,
            ACL: 'public-read',
            ContentType: type,
        }, (err, res) => {
            console.log(err,res);
            if(err){
                reject(err);
            }else{
                resolve();
            }
        });

    });

module.exports = uploader;