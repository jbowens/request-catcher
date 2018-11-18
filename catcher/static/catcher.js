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
  window.catcher.noRequests.hide();
  console.log(req);
  var time = new Date(req.time);

  var snippetDiv = document.createElement('div');
  $(snippetDiv).addClass('snippet');
  $(snippetDiv).data('r', req);

  var methodPath = document.createElement('h2');
  $(methodPath).text(req.method + " " + req.path);

  var ts = document.createElement('div');
  $(metadata).addClass('timestamp');
  $(metadata).text(window.catcher.formatDate(time));

  var remoteAddr = document.createElement('div');
  $(metadata).addClass('remoteaddr');
  $(metadata).text(req.remote_addr);

  $(snippetDiv).append(methodPath, ts, remoteAddr);
  window.catcher.selector.prepend(snippetDiv);

  var selectFn = function() {
    $("#selector .snippet").removeClass('selected');
    $(snippetDiv).addClass('selected');

    var mainDiv = document.createElement('div');
    $(mainDiv).addClass('request');
    $(mainDiv).data('r', req);

    var pre = document.createElement('pre');
    $(pre).text(req.raw_request);
    $(mainDiv).append(pre);

    window.catcher.requests.clear();
    window.catcher.requests.prepend(mainDiv);
  };

  $(snippetDiv).click(selectFn);
  if (window.catcher.selector.children().length() == 1) {
    selectFn();
  }
};

window.catcher.formatDate = function(date) {
  return moment(date).format();
};

$(document).ready(function() {
  window.catcher.requests = $('#requests');
  window.catcher.selector = $('#selector');
  window.catcher.noRequests = $('#no-requests');
  $('#hostname').text(window.location.host);

  window.catcher.connect();
});
