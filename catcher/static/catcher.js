window.catcher = window.catcher || {};

window.catcher.connect = function() {
    if (window.WebSocket) {
        conn = new WebSocket("wss://" + window.location.host + "/init-client");
        conn.onclose = function(evt) {
            console.log("connection closed", evt);
            // Reconnect after a pause.
            setTimeout(window.catcher.connect, 1000);
        };
        conn.onmessage = function(evt) {
            var req = JSON.parse(evt.data);
            window.catcher.insertRequest(req);
        };
    } else {
      console.log("Your browser doesn't support websockets.");
    }
};

window.catcher.insertRequest = function(req) {
  console.log(req);
  var time = new Date(req.time);

  var mainDiv = document.createElement('div');
  $(mainDiv).addClass('request');
  $(mainDiv).data('r', req);

  var pre = document.createElement('pre');
  $(pre).text(req.raw_request);

  var metadata = document.createElement('div');
  $(metadata).addClass('metadata');
  $(metadata).text("Received at " + window.catcher.formatDate(time));

  $(mainDiv).append(pre, metadata);

  window.catcher.requests.prepend(mainDiv);
  if (!window.catcher.requests.is(':visible')) {
    window.catcher.noRequests.hide();
    window.catcher.requests.show();
  }
};

window.catcher.formatDate = function(date) {
  return moment(date).format();
};

$(document).ready(function() {
  window.catcher.requests = $('#requests');
  window.catcher.noRequests = $('#no-requests');
  $('#hostname').text(window.location.host);

  window.catcher.connect();
});
