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

const lambdaClient = new AWS.Lambda({
  region: 'eu-central-1'
});

const createArchive = (archive_id, slug, payload = {}, request_url, req, res) => {
  const params = {
    archive_id,
    slug,
    meta_title: payload.title,
    meta_description: payload.description,
    request_url
  };

  if (req.query.request_id) {
    params.request_id = Number(req.query.request_id);
  }

  Archive.create(params).then(() => {
    if (req.query.client) {
      res.status(200).json({
        data: params,
        status: true
      });
    }
    else {
      res.status(200).redirect(`/${slug}`);
    }
  }).catch(() => {
    if (req.query.client) {
      res.status(503).json({
        status: false
      });
    }
    else {
      res.status(503).redirect('/?fail-create');
    }
  });
};

const lambda = (req, res) => {
  const archive_id = uuid();
  const request_url = req.request_url;
  const slug = generate(7);

  const params = {
    FunctionName: 'teyitlink-archive',
    Payload: JSON.stringify({
      request_url,
      archive_id
    })
  };
  lambdaClient.invoke(params, (err, data) => {
    const payload = JSON.parse(data.Payload);
    if (err) {
      if (req.query.client) {
        res.status(503).json({
          status: false
        });
      }
      else {
        res.status(503).redirect('/?fail-create');
      }
    }
    else {
      createArchive(archive_id, slug, payload, request_url, req, res);
    }
  });
};

// exports
export default lambda;
