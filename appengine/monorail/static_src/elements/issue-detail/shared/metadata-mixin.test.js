// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import {MetadataMixin, HARDCODED_FIELD_GROUPS} from './metadata-mixin.js';
import {PolymerElement} from '@polymer/polymer';

let element;

class MockMetadataElement extends MetadataMixin(PolymerElement) {
  static get is() {
    return 'mock-metadata-element';
  }
}

customElements.define(MockMetadataElement.is, MockMetadataElement);

suite('mr-metadata', () => {
  setup(() => {
    element = document.createElement('mock-metadata-element');
    document.body.appendChild(element);
  });

  teardown(() => {
    document.body.removeChild(element);
  });

  test('initializes', () => {
    assert.instanceOf(element, MockMetadataElement);
  });

  test('groups only applicable to type=FLT-Launch', () => {
    element.issueType = 'Not-FLT-Launch';

    assert.deepEqual(element._filteredGroups, []);

    element.issueType = 'flt-launch';

    assert.deepEqual(element._filteredGroups, HARDCODED_FIELD_GROUPS);

    element.issueType = '';

    assert.deepEqual(element._filteredGroups, []);

    element.issueType = 'FLT-LAUNCH';

    assert.deepEqual(element._filteredGroups, HARDCODED_FIELD_GROUPS);
  });
});