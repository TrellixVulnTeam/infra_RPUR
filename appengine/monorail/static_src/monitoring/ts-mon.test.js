import {assert} from 'chai';
import sinon from 'sinon';

import MonorailTSMon, {PAGE_TYPES} from './ts-mon.js';

suite('MonorailTSMon', () => {
  let mts;

  setup(() => {
    window.CS_env = {
      token: 'rutabaga-token',
      tokenExpiresSec: 1234,
      app_version: 'rutabaga-version',
    };
    window.chops = {rpc: {PrpcClient: sinon.spy()}};
    MonorailTSMon.prototype.disableAfterNextFlush = sinon.spy();
    mts = new MonorailTSMon();
  });

  teardown(() => {
    delete window.CS_env;
  });

  suite('constructor', () => {
    test('initializes a prpcClient', () => {
      assert.equal(mts.prpcClient.constructor.name, 'AutoRefreshPrpcClient');
    });

    test('sets a client ID', () => {
      assert.isNotNull(mts.clientId);
    });

    test('disables sending after next flush', () => {
      sinon.assert.calledOnce(mts.disableAfterNextFlush);
    });
  });

  test('generateClientId', () => {
    const clientID = MonorailTSMon.generateClientId();
    assert.isNaN(clientID);
    const clientIDNum = parseInt(clientID, 32);
    assert.isNotNaN(clientIDNum);
    assert.isAtLeast(clientIDNum, 0);
    assert.isAtMost(clientIDNum, Math.pow(2, 32));
  });

  suite('recordUserTiming', () => {
    test('records a timing metric only if matches', () => {
      const metric = {add: sinon.spy() };
      mts._userTimingMetrics = [{
        category: 'rutabaga',
        eventName: 'rutabaga-name',
        eventLabel: 'rutabaga-label',
        metric: metric,
      }];

      mts.recordUserTiming('kohlrabi', 'rutabaga-name', 'rutabaga-label', 1);
      sinon.assert.notCalled(metric.add);
      metric.add.resetHistory();

      mts.recordUserTiming('rutabaga', 'is-a-tuber', 'rutabaga-label', 1);
      sinon.assert.notCalled(metric.add);
      metric.add.resetHistory();

      mts.recordUserTiming('rutabaga', 'rutabaga-name', 'went bad', 1);
      sinon.assert.notCalled(metric.add);
      metric.add.resetHistory();

      mts.recordUserTiming('rutabaga', 'rutabaga-name', 'rutabaga-label', 1);
      sinon.assert.calledOnce(metric.add);
      assert.equal(metric.add.args[0][0], 1);
      const argsKeys = Array.from(metric.add.args[0][1].keys());
      assert.deepEqual(argsKeys, ['client_id', 'host_name', 'document_visible']);
    });
  });

  suite('recordPageLoadTiming', () => {
    let getPageTypeFromPathStub;
    setup(() => {
      getPageTypeFromPathStub = sinon.stub(MonorailTSMon, 'getPageTypeFromPath');
      mts.pageLoadMetric = {add: sinon.spy()};
      sinon.stub(MonorailTSMon, 'isPageVisible').callsFake(() => (true));
    });

    teardown(() => {
      MonorailTSMon.getPageTypeFromPath.restore();
      MonorailTSMon.isPageVisible.restore();
    });

    test('records page load on issue list page', () => {
      getPageTypeFromPathStub.returns(PAGE_TYPES.ISSUE_LIST);
      mts.recordPageLoadTiming();
      sinon.assert.calledOnce(mts.pageLoadMetric.add);
      assert.isNumber(mts.pageLoadMetric.add.getCall(0).args[0]);
      assert.isString(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'client_id'));
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'host_name'), 'rutabaga-version');
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'template_name'), 'issue_list');
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'document_visible'), true);
    });

    test('records page load on issue detail page', () => {
      getPageTypeFromPathStub.returns(PAGE_TYPES.ISSUE_DETAIL);
      mts.recordPageLoadTiming();
      sinon.assert.calledOnce(mts.pageLoadMetric.add);
      assert.isNumber(mts.pageLoadMetric.add.getCall(0).args[0]);
      assert.isString(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'client_id'));
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'host_name'), 'rutabaga-version');
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'template_name'), 'issue_detail');
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'document_visible'), true);
    });

    test('does not record page load timing on other pages', () => {
      getPageTypeFromPathStub.returns('rutabaga');
      mts.recordPageLoadTiming();
      sinon.assert.notCalled(mts.pageLoadMetric.add);
    });

    test('does not record page load timing if no page type found', () => {
      getPageTypeFromPathStub.returns(null);
      mts.recordPageLoadTiming();
      sinon.assert.notCalled(mts.pageLoadMetric.add);
    });

    test('does not record page load timing if over max threshold', () => {
      getPageTypeFromPathStub.returns(PAGE_TYPES.ISSUE_DETAIL);
      window.performance = {
        timing: {
          navigationStart: 1000,
          domContentLoadedEventEnd: 2001,
        },
      };
      mts.recordPageLoadTiming(1000);
      sinon.assert.notCalled(mts.pageLoadMetric.add);
    });

    test('records page load on issue detail page if under threshold', () => {
      getPageTypeFromPathStub.returns(PAGE_TYPES.ISSUE_DETAIL);
      MonorailTSMon.isPageVisible.restore();
      sinon.stub(MonorailTSMon, 'isPageVisible').callsFake(() => (false));
      window.performance = {
        timing: {
          navigationStart: 1000,
          domContentLoadedEventEnd: 1999,
        },
      };
      mts.recordPageLoadTiming(1000);
      sinon.assert.calledOnce(mts.pageLoadMetric.add);
      assert.isNumber(mts.pageLoadMetric.add.getCall(0).args[0]);
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[0], 999);
      assert.isString(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'client_id'));
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'host_name'), 'rutabaga-version');
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'template_name'), 'issue_detail');
      assert.equal(mts.pageLoadMetric.add.getCall(0).args[1].get(
        'document_visible'), false);
    });
  });

  suite('getPageTypeFromPath', () => {
    test('returns null by default', () => {
      const actual = MonorailTSMon.getPageTypeFromPath('/rutabgaga');
      assert.equal(actual, null);
    });

    test('picks up issue detail path', () => {
      const path = '/p/rutabaga/issues/detail?id=1#c3';
      const actual = MonorailTSMon.getPageTypeFromPath(path);
      assert.equal(actual, PAGE_TYPES.ISSUE_DETAIL);
      assert.equal(actual, 'issue_detail');
    });

    test('picks up issue list path', () => {
      const path = '/p/ruta-0baga/issues/list?mode=list#things';
      const actual = MonorailTSMon.getPageTypeFromPath(path);
      assert.equal(actual, PAGE_TYPES.ISSUE_LIST);
      assert.equal(actual, 'issue_list');
    });
  });

  suite('getGlobalClient', () => {
    test('only creates one global client', () => {
      delete window.__tsMonClient;
      const client1 = MonorailTSMon.getGlobalClient();
      assert.equal(client1, window.__tsMonClient);

      const client2 = MonorailTSMon.getGlobalClient();
      assert.equal(client2, window.__tsMonClient);
      assert.equal(client2, client1);
    });
  });
});