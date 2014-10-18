window.catcher = window.catcher || {};

window.catcher.connect = function() {
    if (window.WebSocket) {
        conn = new WebSocket("ws://localhost:4000/init-client");
        conn.onclose = function(evt) {
            console.log("connection closed", evt);
        };
        conn.onmessage = function(evt) {
            var req = JSON.parse(evt.data);
            window.catcher.insertRequest(req);
        };
    } else {
      console.log("Your browser doesn't support websockets.");
    }
};

var td = function(tr, value, className) {
  var cell = document.createElement('td');
  $(cell).text(value);
  if (className) {
    $(cell).addClass(className);
  }
  $(tr).append(cell);
};

window.catcher.insertRequest = function(req) {
  console.log(req);
  var time = new Date(req.time);
  var tr = document.createElement('tr');
  td(tr, req.method, 'method');
  td(tr, req.path, 'path');
  td(tr, time.toString(), 'time');
  td(tr, req.body, 'body');

  var optionsTd = document.createElement('td');
  $(optionsTd).addClass('options');
  var a = document.createElement('a');
  $(a).addClass('show-raw');
  $(a).text('raw');
  $(a).attr('href', '#');
  $(optionsTd).append(a);
  $(tr).append(optionsTd);

  window.catcher.heading.after(tr);
  if (!window.catcher.table.is(':visible')) {
    window.catcher.noRequests.hide();
    window.catcher.table.show();
  }
};

$(document).ready(function() {
  window.catcher.table = $('table#caught-requests');
  window.catcher.heading = $('table#caught-requests tr.heading');
  window.catcher.noRequests = $('#no-requests');
  window.catcher.connect();
});
