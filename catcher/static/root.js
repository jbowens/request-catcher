window.root = window.root || {};

$(document).ready(function() {
  $('#root-host').text(window.location.host);

  $('#new-catcher').submit(function(e) {
    e.preventDefault();

    var subdomain = $('#subdomain').val();
    if (subdomain) {
      var url = 'http://' + subdomain + '.' + window.location.host + '/';
      window.location = url;
    }
  });
});
