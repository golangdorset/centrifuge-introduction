var Centrifuge = require("centrifuge");
const WebSocket = require('ws');
const util = require("util");
const args = require('minimist')(process.argv.slice(2))

// Start the Centrifuge client.
var centrifuge = new Centrifuge(
    util.format('ws://%s:%s/v1/connection/websocket', args['host'], args['port'].toString()),
    {websocket: WebSocket},
);

// Subscribe to channels.
centrifuge.subscribe('broadcast', function(ctx) {
    console.log(util.format("Someone says via broadcast channel: %s", ctx.data))
})
centrifuge.subscribe(args['user'].toString(), function(ctx) {
    console.log(util.format("Someone says via user channel: %s", ctx.data))
})

centrifuge.connect();