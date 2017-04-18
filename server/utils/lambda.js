// teyit.link

'use strict';

// node modules
import AWS, {Lambda} from 'aws-sdk';
import {generate} from 'randomstring';
import uuid from 'uuid';

// local modules
import {
  AWS_ACCESS_KEY_ID,
  AWS_SECRET_ACCESS_KEY,
  AWS_REGION,
  production
} from '../../config';
import Archive from '../models/archive';

// configuration
AWS.config.update({
  accessKeyId: AWS_ACCESS_KEY_ID,
  secretAccessKey: AWS_SECRET_ACCESS_KEY,
  region: AWS_REGION
});

const createArchive = (archive_id, slug, payload = {}, request_url, res) => {
  Archive.create({
    archive_id,
    slug,
    meta_title: payload.title,
    meta_description: payload.description,
    request_url
  }).then(() => {
    res.status(200).redirect(`/${slug}`);
  }).catch(() => {
    res.status(503).redirect('/?fail-create');
  });
};

const lambda = (req, res) => {
  const archive_id = uuid();
  const request_url = req.request_url;
  const slug = generate(7);

  if (production) {
    const params = {
      FunctionName: 'teyitlink-archive',
      Payload: JSON.stringify({
        request_url,
        archive_id
      })
    };

    Lambda.invoke(params, (err, data) => {
      const payload = JSON.parse(data.Payload);

      if (err) {
        res.status(503).redirect('/?fail-archive');
      }
      else {
        createArchive(archive_id, slug, payload, request_url, res);
      }
    });
  }
  else {
    // poor man's hack to bypass AWS Lambda and to run application on local machine.
    createArchive(archive_id, slug, {title: slug, description: `${slug} - ${archive_id}`}, request_url, res);
  }
};

// exports
export default lambda;
