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
      button.addEventListener('click', function () {
        var pre = button.parentNode.querySelector('.prompt');
        if (!pre) return;
        copy(pre.textContent).then(function () {
          button.classList.add('is-copied');
          button.setAttribute('aria-label', 'Prompt copied');
          window.clearTimeout(timer);
          timer = window.setTimeout(function () {
            button.classList.remove('is-copied');
            button.setAttribute('aria-label', 'Copy prompt to clipboard');
          }, 2000);
        }, function () {
          // Copying is blocked. Select the prompt so it can be copied by hand.
          var range = document.createRange();
          range.selectNodeContents(pre);
          var sel = window.getSelection();
          sel.removeAllRanges();
          sel.addRange(range);
        });
      });
    });
  }

  function initLanguage() {
    var lines = document.querySelectorAll('[data-prompt-lang]');
    if (!lines.length) return;

    // The picker lives on one chapter only, so the names come from a list the
    // generator writes into every page instead of from the <option> elements.
    var names = {};
    var data = document.querySelector('[data-languages]');
    if (data) {
      try {
        JSON.parse(data.textContent).forEach(function (l) {
          names[l.id] = l.name;
        });
      } catch (e) {
        return;
      }
    }

    function apply(id) {
      var name = names[id];
      if (!name) return;
      Array.prototype.forEach.call(lines, function (line) {
        var text = 'Build this in ' + name + ', using its standard library.';
        // Chapter 0 has nothing before it, so it never asks the assistant to
        // continue an existing codebase.
        if (line.getAttribute('data-continues')) {
          text += '\nContinue the same codebase you built in earlier chapters.';
        }
        line.textContent = text;
      });
      Array.prototype.forEach.call(document.querySelectorAll('[data-language-name]'), function (el) {
        el.textContent = name;
      });
    }

    var saved = safeGet('peyva:language');
    if (saved && names[saved]) apply(saved);

    // Only the chapter that owns the picker can change the choice. Every other
    // chapter reads it, so a reader cannot switch language halfway through and
    // end up with earlier chapters written in something else.
    var select = document.querySelector('[data-language-select]');
    if (!select) return;
    if (saved && select.querySelector('option[value="' + saved + '"]')) {
      select.value = saved;
    }
    select.addEventListener('change', function () {
      safeSet('peyva:language', select.value);
      apply(select.value);
    });
  }

  document.addEventListener('DOMContentLoaded', function () {
    initThemeToggle();
    initSidebarToggle();
    initMarkComplete();
    initCopyPrompt();
    initLanguage();
    updateProgress();
  });
})();
