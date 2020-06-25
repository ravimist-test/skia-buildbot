/**
 * @module modules/query-count-sk
 * @description <h2><code>query-count-sk</code></h2>
 *
 * Reports the number of matches for a given query.
 *
 * @attr {string} current_query - The current query to count against.
 *
 * @attr {string} url - The URL to POST the query to.
 *
 * @evt paramset-changed - An event with the updated paramset in e.detail
 *   from the fetch response.
 *
 */
import { define } from 'elements-sk/define';
import { html } from 'lit-html';
import { errorMessage } from 'elements-sk/errorMessage';
import { jsonOrThrow } from 'common-sk/modules/jsonOrThrow';
import { ElementSk } from '../../../infra-sk/modules/ElementSk';
import { CountHandlerRequest, CountHandlerResponse, ParamSet } from '../json';
import 'elements-sk/spinner-sk';

export class QueryCountSk extends ElementSk {
  private _last_query = '';
  private _count = '';
  private _requestInProgress = false;

  private static template = (ele: QueryCountSk) => html`
    <div>
      <span>${ele._count}</span>
      <spinner-sk ?active=${ele._requestInProgress}></spinner-sk>
    </div>
  `;

  constructor() {
    super(QueryCountSk.template);
    this._last_query = '';
    this._count = '';
    this._requestInProgress = false;
  }

  connectedCallback() {
    super.connectedCallback();
    this._upgradeProperty('url');
    this._upgradeProperty('current_query');
    this._render();
    this._fetch();
  }

  static get observedAttributes() {
    return ['current_query', 'url'];
  }

  _fetch() {
    if (!this._connected) {
      return;
    }
    if (!this.url) {
      return;
    }
    if (this._requestInProgress) {
      return;
    }
    this._requestInProgress = true;
    this._last_query = this.current_query;
    const now = Math.floor(Date.now() / 1000);
    const body: CountHandlerRequest = {
      q: this.current_query,
      end: now,
      begin: now - 24 * 60 * 60,
    };
    this._render();
    fetch(this.url, {
      method: 'POST',
      body: JSON.stringify(body),
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then(jsonOrThrow)
      .then((json: CountHandlerResponse) => {
        this._count = `${json.count}`;
        this._requestInProgress = false;
        this._render();
        if (this._last_query !== this.current_query) {
          this._fetch();
        }
        this.dispatchEvent(
          new CustomEvent<ParamSet>('paramset-changed', {
            detail: json.paramset,
            bubbles: true,
          })
        );
      })
      .catch((msg) => {
        this._requestInProgress = false;
        this._render();
        errorMessage(msg);
      });
  }

  attributeChangedCallback() {
    this._fetch();
  }

  /** @prop url - The URL to make POST requests to.  */
  get url() {
    return this.getAttribute('url') || '';
  }

  set url(val) {
    this.setAttribute('url', val);
  }

  /** @prop current_query - The current trace query. */
  get current_query() {
    return this.getAttribute('current_query') || '';
  }

  set current_query(val) {
    this.setAttribute('current_query', val);
  }
}

define('query-count-sk', QueryCountSk);