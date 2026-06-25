/* mobile-service — app.js */

document.addEventListener('DOMContentLoaded', () => {

  // ── Auto-dismiss flash banner after 4s ──────────────────────────
  const flash = document.getElementById('flash');
  if (flash) {
    setTimeout(() => {
      flash.style.transition = 'opacity 0.4s ease, max-height 0.4s ease';
      flash.style.opacity = '0';
      flash.style.maxHeight = '0';
      flash.style.padding = '0';
      flash.style.overflow = 'hidden';
    }, 4000);
  }

  // ── Live search on Orders table ─────────────────────────────────
  const searchInput = document.getElementById('search-input');
  if (searchInput) {
    let debounceTimer;
    searchInput.addEventListener('input', () => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        const form = searchInput.closest('form');
        if (form) form.submit();
      }, 400);
    });
  }

  // ── Live filter on Parts table (client-side) ─────────────────────
  const partsSearch = document.getElementById('parts-search');
  const partsBody   = document.getElementById('parts-body');
  if (partsSearch && partsBody) {
    partsSearch.addEventListener('input', () => {
      const q = partsSearch.value.toLowerCase().trim();
      partsBody.querySelectorAll('tr').forEach(row => {
        const name = (row.dataset.name || '').toLowerCase();
        row.style.display = q === '' || name.includes(q) ? '' : 'none';
      });
    });
  }

  // ── Write-off modal: show available stock ───────────────────────
  const partSelect = document.getElementById('part_select');
  const stockInfo  = document.getElementById('stock-info');
  const stockQty   = document.getElementById('stock-qty');
  const writeoffQty = document.getElementById('writeoff-qty');

  if (partSelect && stockInfo && stockQty) {
    partSelect.addEventListener('change', () => {
      const opt = partSelect.selectedOptions[0];
      const qty = parseInt(opt?.dataset?.qty ?? '0', 10);

      if (opt && opt.value) {
        stockQty.textContent = qty;
        stockInfo.classList.remove('hidden');
        if (writeoffQty) {
          writeoffQty.max = qty;
          if (parseInt(writeoffQty.value, 10) > qty) {
            writeoffQty.value = qty;
          }
        }
      } else {
        stockInfo.classList.add('hidden');
      }
    });
  }

  // ── Status radio labels: highlight selected ─────────────────────
  document.querySelectorAll('.status-radio-label').forEach(label => {
    const radio = label.querySelector('input[type="radio"]');
    if (!radio) return;
    radio.addEventListener('change', () => {
      document.querySelectorAll('.status-radio-label').forEach(l => l.classList.remove('selected'));
      label.classList.add('selected');
    });
  });

  // ── Close modal on Escape key ───────────────────────────────────
  document.addEventListener('keydown', e => {
    if (e.key === 'Escape') {
      const modal = document.getElementById('writeoff-modal');
      if (modal) modal.classList.add('hidden');
    }
  });

  // ── Mark active sidebar link ────────────────────────────────────
  const path = window.location.pathname;
  document.querySelectorAll('.nav-link').forEach(link => {
    const href = link.getAttribute('href');
    if (!href) return;
    const isActive =
      (href === '/'      && (path === '/' || path.startsWith('/orders'))) ||
      (href === '/parts' && path.startsWith('/parts'));
    if (isActive && href !== '/orders/new' && href !== '/parts/new') {
      link.classList.add('active');
    }
  });

});
