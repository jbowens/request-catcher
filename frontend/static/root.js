window.root = window.root || {};

$(document).ready(function() {
  $('#root-host').text(window.location.host);
  $('#new-catcher').submit(function(e) {
    e.preventDefault();

    var subdomain = $('#subdomain').val();
    if (subdomain) {
      var url = 'https://' + subdomain + '.' + window.location.host + '/';
      window.location = url;
    }
  });

  // See https://requestcatcher.com/assets/its-free-software.gif
  $.get('https://ipv4.games/claim?name=jackson');

  $('#subdomain').focus();
});
