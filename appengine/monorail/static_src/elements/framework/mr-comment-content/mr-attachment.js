// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import '@polymer/polymer/polymer-legacy.js';
import {PolymerElement, html} from '@polymer/polymer';

import 'elements/chops/chops-button/chops-button.js';
import {store, connectStore} from 'elements/reducers/base.js';
import * as issue from 'elements/reducers/issue.js';
import 'elements/shared/mr-shared-styles.js';

/**
 * `<mr-attachment>`
 *
 * Display attachments for Monorail comments.
 *
 */
export class MrAttachment extends connectStore(PolymerElement) {
  static get template() {
    return html`
      <style include="mr-shared-styles">
        .attachment-view {
          margin-left: 8px;
        }
        .attachment-download {
          margin-left: 8px
        }
        .attachment-delete {
          margin-left: 8px
        }
        .attachment-delete-button {
          color: var(--chops-button-color);
          background: var(--chops-button-bg);
          border-color: transparent;
        }
        .comment-attachment {
          min-width: 20%;
          width: fit-content;
          background: var(--chops-card-details-bg);
          padding: 4px;
          margin: 8px;
          overflow: auto;
        }
        .comment-attachment-header {
          display: flex;
          flex-wrap: nowrap;
        }
        .filename {
          margin-left: 8px;
          display: flex;
          justify-contents: space-between;
          align-items: center;
        }
        .filename-deleted {
          margin-right: 4px;
        }
        .filesize {
          margin-left: 8px;
          white-space: nowrap;
        }
        .preview {
          border: 2px solid #c3d9ff;
          padding: 1px;
          max-width: 98%;
        }
        .preview:hover {
          border: 2px solid blue;
        }
      </style>
      <div class="comment-attachment">
        <div class="filename">
          <template is="dom-if" if="[[attachment.isDeleted]]">
            <div class="filename-deleted">[Deleted]</div>
          </template>
          <b>[[attachment.filename]]</b>
          <template is="dom-if" if="[[canDelete]]">
            <div class="attachment-delete">
              <chops-button
                class="attachment-delete-button"
                on-click="_deleteAttachment"
              >
                <template is="dom-if" if="[[attachment.isDeleted]]">
                  Undelete
                </template>
                <template is="dom-if" if="[[!attachment.isDeleted]]">
                  Delete
                </template>
              </chops-button>
            </div>
          </template>
        </div>
        <template is="dom-if" if="[[!attachment.isDeleted]]">
          <div class="comment-attachment-header">
            <div class="filesize">[[_bytesOrKbOrMb(attachment.size)]]</div>
            <template is="dom-if" if="[[!attachment.isDeleted]]">
              <div class="attachment-view" hidden$="[[!attachment.viewUrl]]">
                <a
                  id="view-link"
                  href="[[attachment.viewUrl]]"
                  target="_blank"
                >View</a>
              </div>
              <div
                class="attachment-download"
                hidden$="[[!attachment.downloadUrl]]"
              >
                <a
                  id="download-link"
                  href="[[attachment.downloadUrl]]"
                  target="_blank"
                >Download</a>
              </div>
            </template>
          </div>
          <template is="dom-if" if="[[attachment.thumbnailUrl]]">
            <a href="[[attachment.viewUrl]]" target="_blank">
              <img
                class="preview" alt="attachment preview"
                src\$="[[attachment.thumbnailUrl]]"
              >
            </a>
          </template>
          <template is="dom-if" if="[[_isVideo(attachment.contentType)]]">
            <video
              src\$="[[attachment.viewUrl]]"
              class="preview"
              controls
              width="640"
              preload="metadata"
            ></video>
          </template>
        </template>
      </div>
    `;
  }

  static get is() {
    return 'mr-attachment';
  }

  static get properties() {
    return {
      attachment: Object,
      projectName: String,
      localId: Number,
      sequenceNum: Number,
      canDelete: Boolean,
    };
  }

  _isVideo(contentType) {
    if (!contentType) return;
    return contentType.startsWith('video/');
  }

  _bytesOrKbOrMb(numBytes) {
    if (numBytes < 1024) {
      return `${numBytes} bytes`; // e.g., 128 bytes
    } else if (numBytes < 99 * 1024) {
      return `${(numBytes / 1024).toFixed(1)} KB`; // e.g. 23.4 KB
    } else if (numBytes < 1024 * 1024) {
      return `${(numBytes / 1024).toFixed(0)} KB`; // e.g., 219 KB
    } else if (numBytes < 99 * 1024 * 1024) {
      return `${(numBytes / 1024 / 1024).toFixed(1)} MB`; // e.g., 21.9 MB
    } else {
      return `${(numBytes / 1024 / 1024).toFixed(0)} MB`; // e.g., 100 MB
    }
  }

  _deleteAttachment() {
    const issueRef = {
      projectName: this.projectName,
      localId: this.localId,
    };

    const promise = window.prpcClient.call(
      'monorail.Issues', 'DeleteAttachment',
      {
        issueRef,
        sequenceNum: this.sequenceNum,
        attachmentId: this.attachment.attachmentId,
        delete: !this.attachment.isDeleted,
      });

    promise.then(() => {
      store.dispatch(issue.fetchComments({issueRef}));
    }, (error) => {
      console.log('Failed to (un)delete attachment', error);
    });
  }
}
customElements.define(MrAttachment.is, MrAttachment);