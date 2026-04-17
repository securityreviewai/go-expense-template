document.body.addEventListener('htmx:responseError', (e) => {
  console.error('HTMX error', e.detail);
});
