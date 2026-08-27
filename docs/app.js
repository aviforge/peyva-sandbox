(function () {
  function safeGet(key) {
    try {
      return window.localStorage.getItem(key);
    } catch (e) {
      return null;
    }
  }

  function safeSet(key, value) {
    try {
      window.localStorage.setItem(key, value);
    } catch (e) {
      // localStorage unavailable (e.g. some browsers under file://) - ignore.
    }
  }

  function initThemeToggle() {
    var toggle = document.querySelector('[data-theme-toggle]');
    if (!toggle) return;

    var saved = safeGet('peyva:theme');
    if (saved === 'dark') {
      document.documentElement.setAttribute('data-theme', 'dark');
    }

    toggle.addEventListener('click', function () {
      var isDark = document.documentElement.getAttribute('data-theme') === 'dark';
      if (isDark) {
        document.documentElement.removeAttribute('data-theme');
        safeSet('peyva:theme', 'light');
      } else {
        document.documentElement.setAttribute('data-theme', 'dark');
        safeSet('peyva:theme', 'dark');
      }
    });
  }

  function initSidebarToggle() {
    var toggle = document.querySelector('[data-sidebar-toggle]');
    var sidebar = document.querySelector('.sidebar');
    var backdrop = document.querySelector('[data-sidebar-backdrop]');
    if (!toggle || !sidebar || !backdrop) return;

    function close() {
      sidebar.classList.remove('is-open');
      backdrop.classList.remove('is-open');
    }

    toggle.addEventListener('click', function () {
      sidebar.classList.toggle('is-open');
      backdrop.classList.toggle('is-open');
    });

    backdrop.addEventListener('click', close);

    sidebar.addEventListener('click', function (e) {
      if (e.target.closest('a')) close();
    });
  }

  function initMarkComplete() {
    var btn = document.querySelector('[data-mark-complete]');
    if (!btn) return;
    var slug = document.body.getAttribute('data-chapter-slug');
    if (!slug) return;

    var key = 'peyva:complete:' + slug;

    function render() {
      var done = safeGet(key) === '1';
      btn.textContent = done ? '✓ Chapter completed' : 'Mark chapter as complete';
      btn.classList.toggle('is-complete', done);
    }

    btn.addEventListener('click', function () {
      safeSet(key, safeGet(key) === '1' ? '0' : '1');
      render();
      updateProgress();
    });

    render();
  }

  function updateProgress() {
    var fill = document.querySelector('[data-progress-fill]');
    var label = document.querySelector('[data-progress-label]');

    var chapters = document.querySelectorAll('.chapter-list a, .chapter-list li.active');
    var total = chapters.length;
    var completed = 0;

    chapters.forEach(function (el) {
      var slug = el.getAttribute('data-slug');
      var done = !!slug && safeGet('peyva:complete:' + slug) === '1';
      el.classList.toggle('is-complete', done);
      if (done) completed++;
    });

    var pct = total ? Math.round((completed / total) * 100) : 0;
    if (fill) fill.style.width = pct + '%';
    if (label) label.textContent = pct + '% complete (' + completed + ' / ' + total + ')';
  }

  function initCopyPrompt() {
    var buttons = document.querySelectorAll('[data-copy-prompt]');
    if (!buttons.length) return;

    // navigator.clipboard needs a secure context. The README tells readers they
    // can open docs/index.html straight from a clone, and a file:// page is not
    // one in most browsers, so fall back rather than silently doing nothing.
    function copy(text) {
      if (navigator.clipboard && window.isSecureContext) {
        return navigator.clipboard.writeText(text);
      }
      return new Promise(function (resolve, reject) {
        var scratch = document.createElement('textarea');
        scratch.value = text;
        scratch.setAttribute('readonly', '');
        scratch.style.position = 'fixed';
        scratch.style.top = '-1000px';
        document.body.appendChild(scratch);
        scratch.select();
        var ok = false;
        try {
          ok = document.execCommand('copy');
        } catch (e) {
          ok = false;
        }
        document.body.removeChild(scratch);
        ok ? resolve() : reject();
      });
    }

    Array.prototype.forEach.call(buttons, function (button) {
      var timer = null;
      // The buttons label different things: a prompt, the spec, the rules file.
      // Restoring what the page shipped keeps each one describing itself.
      var label = button.getAttribute('aria-label');
      button.addEventListener('click', function () {
        var pre = button.parentNode.querySelector('.prompt');
        if (!pre) return;

        // The standing rules are printed once a page and prepended here, so a
        // copied prompt still carries them without the page repeating itself
        // above every block.
        var text = pre.textContent;
        var rulesID = button.parentNode.getAttribute('data-prompt-rules');
        if (rulesID) {
          var rules = document.getElementById(rulesID);
          if (rules) text = rules.textContent.replace(/\s+$/, '') + '\n\n' + text;
        }

        copy(text).then(function () {
          button.classList.add('is-copied');
          button.setAttribute('aria-label', 'Copied');
          window.clearTimeout(timer);
          timer = window.setTimeout(function () {
            button.classList.remove('is-copied');
            button.setAttribute('aria-label', label);
          }, 2000);
        }, function () {
          // Copying is blocked. Select the prompt so it can be copied by hand.
          // The rules stay above; a hand copy of the body is the best that can
          // be offered without moving the selection off what was clicked.
          var range = document.createRange();
          range.selectNodeContents(pre);
          var sel = window.getSelection();
          sel.removeAllRanges();
          sel.addRange(range);
        });
      });
    });
  }

  // The operating system changes the commands a few chapters ask for and the
  // runner script, rather than every prompt the way the language does. Same
  // shape otherwise: the generator writes the text and wraps the choice in a
  // span, and this only ever replaces what is inside that span.
  function initSystem() {
    var names = document.querySelectorAll('[data-os-name]');
    var select = document.querySelector('[data-system-select]');
    if (!names.length && !select) return;

    // Every system's script is already in the page, so switching shows one and
    // hides the rest. Nothing is fetched, which is what lets the site work from
    // a clone with no network.
    var blocks = document.querySelectorAll('[data-runner-for]');

    function apply(name, id) {
      Array.prototype.forEach.call(names, function (el) {
        el.textContent = name;
      });
      if (!blocks.length || !id) return;
      Array.prototype.forEach.call(blocks, function (el) {
        el.hidden = el.getAttribute('data-runner-for') !== id;
      });
    }

    var saved = safeGet('peyva:systemName');
    var savedID = safeGet('peyva:system');
    if (saved && saved.length < 40) apply(saved, savedID);

    if (!select) return;

    if (savedID && select.querySelector('option[value="' + savedID + '"]')) {
      select.value = savedID;
    }
    // Nothing saved yet, so the page is showing whichever script the generator
    // put first. Hide the others so the reader is not handed three.
    apply(select.options[select.selectedIndex].getAttribute('data-prompt') ||
      select.options[select.selectedIndex].textContent, select.value);

    // The option's own text is the short name for the picker. What goes into a
    // prompt names the shell too, so it is carried on the option rather than
    // rebuilt here, where it would drift from the Go that writes the pages.
    select.addEventListener('change', function () {
      var option = select.options[select.selectedIndex];
      var text = option.getAttribute('data-prompt') || option.textContent;
      safeSet('peyva:system', select.value);
      safeSet('peyva:systemName', text);
      apply(text, select.value);
    });
  }

  function initLanguage() {
    var names = document.querySelectorAll('[data-language-name]');
    if (!names.length) return;

    // The generator writes the preamble and wraps the language in a span, so
    // this only ever replaces that one word. Rebuilding the sentence here would
    // put the rules in two places, and they drifted the first time it did.
    function apply(name) {
      Array.prototype.forEach.call(names, function (el) {
        el.textContent = name;
      });
    }

    var saved = safeGet('peyva:languageName');
    if (saved && saved.length < 40) apply(saved);

    // Only the chapter that owns the picker can change the choice. Every other
    // chapter reads it, so a reader cannot switch language halfway through and
    // leave the earlier chapters written in something else.
    var select = document.querySelector('[data-language-select]');
    if (!select) return;

    var savedID = safeGet('peyva:language');
    if (savedID && select.querySelector('option[value="' + savedID + '"]')) {
      select.value = savedID;
    }

    select.addEventListener('change', function () {
      var chosen = select.options[select.selectedIndex].textContent;
      safeSet('peyva:language', select.value);
      safeSet('peyva:languageName', chosen);
      apply(chosen);
    });
  }

  document.addEventListener('DOMContentLoaded', function () {
    initThemeToggle();
    initSidebarToggle();
    initMarkComplete();
    initCopyPrompt();
    initLanguage();
    initSystem();
    updateProgress();
  });
})();
