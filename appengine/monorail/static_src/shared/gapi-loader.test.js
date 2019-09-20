// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import sinon from 'sinon';
import loadGapi from './gapi-loader.js';

describe('loadGapi', () => {
  beforeEach(() => {
    window.CS_env = {gapi_client_id: 'rutabaga'};
    // Pre-load gapi with a fake signin object to prevent loading the
    // real gapi.js.
    loadGapi({
      init: () => {},
      getUserProfileAsync: () => Promise.resolve({}),
    });
  });

  afterEach(() => {
    delete window.CS_env;
  });

  it('errors out if no client_id', () => {
    window.CS_env.gapi_client_id = undefined;
    assert.throws(() => {
      loadGapi();
    });
  });

  it('returns the same promise when called multiple times', () => {
    const callOne = loadGapi();
    const callTwo = loadGapi();

    assert.strictEqual(callOne, callTwo);
    assert.instanceOf(callOne, Promise);
  });
});
