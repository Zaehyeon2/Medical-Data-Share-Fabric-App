'use strict';

var FabricClient = require('fabric-client');
var superagent = require('superagent');
var path =  require('path');
var fs = require('fs');

async function main() {
    try{
        const client = new FabricClient();
        let envelope_bytes = fs.readFileSync(path.join(__dirname, '../../bin/kingod.tx'));
        let config_json = JSON.parse(fs.readFileSync(path.join(__dirname, '../../bin/kingod.json')));
        var config_update = client.extractChannelConfig(envelope_bytes)
        var response = superagent.post('http://127.0.0.1:7059/protolator/encode/common.ConfigUpdate',
            config_json.payload.data.toString())
            .buffer()
            .end((err, res) => {
                if(err) {
                console.error(err);
                return;
                }
                config_proto = res.body;
            });
        var signature = client.signChannelConfig(config_proto);
        config_json.signatures.push(signature);

        console.log(123);
        var orderer = client.newOrderer("0.0.0.0:8888");

        let tx_id = client.newTransactionID();

        request = {
            config: config_proto,
            signatures: signatures,
            name: 'kingod',
            orderer: orderer,
            txId: tx_id
        };

        var channel = client.createChannel(request);


    } catch (error) {
        console.error(`Failed to Create Channel :${error}`);
        process.exit(1);
    }

}

main();