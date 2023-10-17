document.body.addEventListener('htmx:beforeOnLoad', function (evt) {
  if (evt.detail.xhr.status === 422 || evt.detail.xhr.status === 400) {
    evt.detail.shouldSwap = true;
    evt.detail.isError = false;
  }
});
