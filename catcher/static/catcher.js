window.catcher = window.catcher || {};

window.catcher.connect = function() {
    if (window.WebSocket) {
        conn = new WebSocket("ws://" + window.location.host + "/init-client");
        conn.onclose = function(evt) {
            console.log("connection closed", evt);
            // Reconnet after a pause.
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

var td = function(tr, value, className) {
  var cell = document.createElement('td');
  $(cell).text(value);
  if (className) {
    $(cell).addClass(className);
  }
  $(tr).append(cell);
};

var makeDateTime = function(time) {
  var dateDiv = document.createElement('div');
  var timeDiv = document.createElement('div');
  dateDiv.className = 'date';
  timeDiv.className = 'time';
  $(dateDiv).text(window.catcher.formatDate(time));
  $(timeDiv).text(window.catcher.formatTime(time));
  var dateTime = $('<div class="timestamp"></div>');
  dateTime.append(dateDiv);
  dateTime.append(timeDiv);
  return dateTime;
};

var makeDiv = function(text, className) {
  var div = document.createElement('div');
  $(div).text(text);
  $(div).addClass(className);
  return div;
};

window.catcher.insertRequest = function(req) {
  console.log(req);
  var time = new Date(req.time);
  var tr = document.createElement('tr');
  $(tr).data('r', req);

  var mainTd = document.createElement('td');
  $(mainTd).addClass('general-info');
  var methodAndPath = makeDiv(req.method + ' ' + req.path, 'method-and-path');
  var dateTime = makeDateTime(time);
  var remoteAddr = makeDiv(req.remote_addr, 'remote-addr');

  var a = document.createElement('a');
  $(a).addClass('show-raw');
  $(a).text('raw http request');
  $(a).attr('href', '#');
  var optionsDiv = document.createElement('div');
  $(optionsDiv).append(a);
  $(optionsDiv).addClass('options');

  $(mainTd).append(methodAndPath, dateTime, remoteAddr, optionsDiv);
  tr.appendChild(mainTd);

  var bodyTd = document.createElement('td');
  $(bodyTd).addClass('body');
  var code = document.createElement('code');
  $(code).text(req.body);
  bodyTd.appendChild(code);
  tr.appendChild(bodyTd);

  if (req.headers['Content-Type'] &&
      req.headers['Content-Type'][0] === "application/json") {
    $(tr).find('.body code').each(function(i, block) {
      hljs.highlightBlock(block);
    });
  }

  window.catcher.addListeners(tr);
  window.catcher.heading.after(tr);
  if (!window.catcher.table.is(':visible')) {
    window.catcher.noRequests.hide();
    window.catcher.table.show();
  }
};

window.catcher.formatDate = function(date) {
  return moment(date).format("dddd, MMMM Do YYYY");
};

window.catcher.formatTime = function(date) {
  return moment(date).format("h:mm:ss a");
};

window.catcher.addListeners = function(row) {
  $(row).find('.show-raw').click(window.catcher.showRaw);
};

window.catcher.showRaw = function(evt) {
  evt.preventDefault();
  var req = $(evt.currentTarget).closest('tr').data('r');
  window.catcher.rawContent.text(req.raw_request);
  window.catcher.rawContent.each(function(i, block) {
    hljs.highlightBlock(block);
  });
  window.catcher.rawPopup.show();
};

$(document).ready(function() {
  window.catcher.table = $('table#caught-requests');
  window.catcher.heading = $('table#caught-requests tr.heading');
  window.catcher.noRequests = $('#no-requests');
  window.catcher.rawPopup = $('#raw-popup');
  window.catcher.rawContent = $('#raw-popup #raw-content');

  $('#raw-popup .close-popup').click(function(e) {
    window.catcher.rawPopup.hide();
  });

  $('#hostname').text(window.location.host);

  window.catcher.connect();
});
