import * as React from 'react'
import { render } from 'react-dom'
import { RequestDetail } from '~/components/request_detail'

window.catcher = window.catcher || {};

window.catcher.connect = function() {
    if (window.WebSocket) {
      var protocol = "wss://";
      if (window.location.protocol === "http:") {
        protocol = "ws://";
      }
      const conn = new WebSocket(protocol + window.location.host + "/init-client");
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
  $(ts).addClass('timestamp');
  $(ts).text(window.catcher.formatDate(time));

  var remoteAddr = document.createElement('div');
  $(remoteAddr).addClass('remoteaddr');
  $(remoteAddr).text(req.remote_addr);

  $(snippetDiv).append(methodPath, ts, remoteAddr);
  window.catcher.selector.prepend(snippetDiv);

  var selectFn = function() {
    $("#selector .snippet").removeClass('selected');
    $(snippetDiv).addClass('selected');

    render(
        <RequestDetail request={req} />,
        document.getElementById('requests')
    );
  };

  $(snippetDiv).click(selectFn);
  if (window.catcher.selector.children().length == 1) {
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
