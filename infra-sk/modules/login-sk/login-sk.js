/**
 * @module infra-sk/modules/login-sk
 * @description <h2><code>login-sk</code></h2>
 *
 * <p>
 * The <login-sk> custom element. Uses the Login promise to display the
 * current login status and provides login/logout links. Reports errors via
 * {@linkcode module:elements-sk/error-toast-sk}.
 * </p>
 *
 * @attr {boolean} testing_offline - If we should really poll for Login status
 *  or use mock data.
 * @attr {string} login_host - If it needs to be something other than skia.org.
 *  If login_host is not specified then skia.org will be used.
 *
 * <p>
 */
import { define } from 'elements-sk/define'
import { errorMessage } from 'elements-sk/errorMessage';
import { LoginTo } from '../login';

define('login-sk', class extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `<span class=email>Loading...</span><a class=logInOut></a>`;
    const host = this.loginHost ? this.loginHost : 'skia.org';
    if (this.testingOffline) {
      this.querySelector('.email').textContent = "test@example.com";
      const logInOut = this.querySelector('.logInOut');
      logInOut.href = `https://${host}/logout/?redirect=` + encodeURIComponent(document.location);
      logInOut.textContent = 'Logout';
    } else {
      LoginTo(`https://${host}/loginstatus/`).then((status) => {
        this.querySelector('.email').textContent = status.Email;
        let logInOut = this.querySelector('.logInOut');
        if (!status.Email) {
            logInOut.href = status.LoginURL;
            logInOut.textContent = 'Login';
        } else {
            logInOut.href = `https://${host}/logout/?redirect=` + encodeURIComponent(document.location);
            logInOut.textContent = 'Logout';
        }
      }).catch(errorMessage);
    }
  }

  /** @prop testingOffline {boolean} Reflects the testing_offline attribute for ease of use.
   */
  get testingOffline() { return this.hasAttribute('testing_offline'); }
  set testingOffline(val) {
    if (val) {
      this.setAttribute('testing_offline', '');
    } else {
      this.removeAttribute('testing_offline');
    }
  }

  /** @prop loginHost {string} Which host should be used for login URLs.
   */
  get loginHost() { return this.getAttribute('login_host'); }
  set loginHost(val) {
    if (val) {
      this.setAttribute('login_host', val);
    }
  }
});
