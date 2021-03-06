/**
 * @module modules/task-repeater-sk
 * @description A custom element for allowing the user to choose whether
 * their task should be repeated either never, daily, every other day,
 * or weekly.
 */

import { $$ } from 'common-sk/modules/dom';
import { define } from 'elements-sk/define';
import { html } from 'lit-html';

import { ElementSk } from '../../../infra-sk/modules/ElementSk';
import 'elements-sk/select-sk';

const frequencies = [
  { num: '0', desc: 'Never repeat' },
  { num: '1', desc: 'Repeat daily' },
  { num: '2', desc: 'Repeat every other day' },
  { num: '7', desc: 'Repeat weekly' }];

const template = () => html`
<div class=tr-container>
  <select-sk>
    ${frequencies.map((f) => html`
    <div>
      <span class=num>${f.num}</span><span>${f.desc}</span>
    </div>`)}
  </select-sk>
</div>
`;

define('task-repeater-sk', class extends ElementSk {
  constructor() {
    super(template);
  }

  connectedCallback() {
    super.connectedCallback();
    this._render();
    this._selector = $$('select-sk', this);
    this._selector.selection = 0;
  }

  /**
   * @prop {string} frequency - Number string representing user selected
   * frequency of their task.
   */
  get frequency() {
    return frequencies[this._selector.selection].num;
  }

  set frequency(val) {
    this._selector.selection = frequencies.findIndex((f) => f.num === val);
  }
});
