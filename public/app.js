const TYPE_WS_JOIN = "TYPE_WS_JOIN";
const TYPE_WS_LEAVE = "TYPE_WS_LEAVE";
const TYPE_WS_MSG = "TYPE_WS_MSG";
const TYPE_WS_CHECK_ONLINE_JOIN = "TYPE_WS_CHECK_ONLINE_JOIN";

new Vue({
  el: "#app",

  data: {
    ws: null,
    newMsg: "",
    chatContent: "",
    username: null,
    joined: false,
    amountOnline: 0,
    amountJoin: 0,
  },

  created: function () {
    var self = this;
    this.ws = new WebSocket("ws://" + window.location.host + "/ws");
    this.ws.addEventListener("message", function (e) {
      var msg = JSON.parse(e.data);

      if (msg.type === TYPE_WS_CHECK_ONLINE_JOIN) {
        self.amountOnline = msg.amount_online;
        self.amountJoin = msg.amount_join;
        return;
      }

      if (msg.type === TYPE_WS_JOIN) {
        self.amountOnline = msg.amount_online;
        self.amountJoin = msg.amount_join;
        self.chatContent +=
          '<div class="chip green lighten-4">' + emojione.toImage(msg.username + " join chat" + " 🎉") + "</div>" + "<br/>";
        return;
      }

      if (msg.type === TYPE_WS_LEAVE) {
        self.amountOnline = msg.amount_online;
        self.amountJoin = msg.amount_join;
        self.chatContent +=
          '<div class="chip red lighten-4">' + emojione.toImage(msg.username + " leave chat" + " 😭") + "</div>" + "<br/>";
        return;
      }

      if (msg.type === TYPE_WS_MSG) {
        var temp = '<div class="chip">'
        if (self.username == msg.username) {
          temp = '<div class="chip blue lighten-4">'
        }
        self.chatContent +=
          temp +
          '<img src="' +
          self.gravatarURL(msg.username) +
          '">' +
          msg.username +
          "</div>" +
          emojione.toImage(msg.message) +
          "<br/>";
        return;
      }
    });
  },

  watch: {
    chatContent: function(val) {
      if (val) {
        window.requestAnimationFrame(() => {
          var element = document.getElementById("chat-messages");
          element.scrollTop = element.scrollHeight
        })
      }
    }
  },

  methods: {
    send: function () {
      if (this.newMsg != "") {
        this.ws.send(
          JSON.stringify({
            type: TYPE_WS_MSG,
            username: this.username,
            message: $("<p>").html(this.newMsg).text(),
          })
        );
        this.newMsg = "";
      }
    },
    join: function () {
      var self = this
      if (!this.username) {
        Materialize.toast("You must choose a username", 2000);
        return;
      }
      this.username = $("<p>").html(this.username).text();
      $.ajax({
        url: "http://" + window.location.host + "/check-username?username=" + this.username,
        type: "GET",
        dataType: "json",
        success: function(res) {
          if (res.status_code == 400) {
            Materialize.toast(res.reason, 2000);
            return
          }
          if (res.status_code == 200 && !res.is_allow) {
            Materialize.toast(res.reason, 2000);
            return
          }
          if (res.status_code == 200 && res.is_allow) {
            self.username = res.username
            self.ws.send(
              JSON.stringify({
                type: TYPE_WS_JOIN,
                username: self.username,
              })
            );
            self.joined = true;
            return
          }
        }
      });
    },
    gravatarURL: function (username) {
      return "http://www.gravatar.com/avatar/" + CryptoJS.MD5(username) + "?d=monsterid";
    },
  },
});
