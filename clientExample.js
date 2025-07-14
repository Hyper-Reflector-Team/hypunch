// This expects you are running something like node js for your clients.
var dgram = require("dgram");
let socket = null;
let message = "";

// These variables need to be obtained from somewhere, maybe your front end code some other id generation
// They are used to tell the server that you want information that matches the UID
let myUID = "someUID";
let OpponentUID = "someOtherUID";

// Open an new socket server locally
if (!socket) {
  try {
    socket = dgram.createSocket("udp4");
    socket.bind(() => {
      console.log("Socket bound to random port:", socket.address());
    });
  } catch (error) {
    console.error("Error opening socket:", error);
  }
}

// Start server connection
function sendToServer(myUID, OpponentUID, kill) {
  const serverPort = keys.PUNCH_PORT; // revert this after killing all of the new services
  const serverHost = keys.COTURN_IP;
  const message = Buffer.from(
    JSON.stringify({
      uid: myUID,
      peerUid: OpponentUID,
      kill,
    })
  );
  try {
    socket.send(
      message,
      0,
      message.length,
      serverPort,
      serverHost,
      function (err) {
        if (err) return console.log("Handle Error", err);
        // UDP message was sent to the server
        console.log("UDP message sent Server " + serverHost + ":" + serverPort);
      }
    );
  } catch (error) {
    console.log("could not send message to server");
  }
}

async function sendToOtherClient(address, port, msg = "") {
  if (msg.length >= 1) {
    message = Buffer.from(msg);
  } else {
    message = Buffer.from("ping");
  }
  try {
    if (!socket) return;
    socket.send(message, 0, message.length, port, address, function (err) {
      if (err) return console.log("Handle Error", err);
    });
  } catch (error) {
    // Error sending a seocket message
    if (error) return console.log("Handle Error", error);
  }
}

// KEEP ALIVE MESSAGES
// Make sure our socket exists
if (socket) {
  try {
    socket.on("message", function (message) {
      const messageContent = message.toString();
      // Make sure we are only sending correct messages to the the other client.
      if (messageContent === "ping" || message.includes('"port"')) {
        if (message.includes('"port"') && !keepAliveInterval) {
          keepAliveInterval = setInterval(() => {
            sendToOtherClient(
              opponentEndpoint.peer.address,
              opponentEndpoint.peer.port,
              "ping"
            );
          }, 1000); // Delay between keep alive messages
        }
      }
      try {
        opponentEndpoint = JSON.parse(message);
        currentMatchId = opponentEndpoint.matchId || null;
        sendToOtherClient(
          opponentEndpoint.peer.address,
          opponentEndpoint.peer.port
        );
      } catch (err) {}
    });
  } catch (error) {
    console.log("error in socket", error);
  }
}

// Start the hole punch
sendToServer(myUID, OpponentUID, false);

// To kill a client connection and delete the entries from the server mapping we fire off
function endHolePunch() {
  sendToServer(myUID, OpponentUID, true);
  clearInterval(keepAliveInterval);
  keepAliveInterval = null;
}
