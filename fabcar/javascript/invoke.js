/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function main() {
    try {
<<<<<<< HEAD
        if (process.argv.length != 4) {
            console.error("Failed: Expecting argument is 2");
            process.exit(1);
        }
=======
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        let ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

>>>>>>> 0031db5de282eb61a57d774a5dd8e12de79aa076
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), '../wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('appUser');
        if (!identity) {
            console.log('An identity for the user "appUser" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        console.log('before');
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
<<<<<<< HEAD
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });
        
=======
        await gateway.connect(ccp, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

>>>>>>> 0031db5de282eb61a57d774a5dd8e12de79aa076
        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');
        
        // Get the contract from the network.
        const contract = network.getContract('fabcar');
        
        var fs = require('fs');
        var ip = require('ipfs-http-client');

        console.time('Base64');
        var data = fs.readFileSync(process.argv[3], 'binary');
        var b64data = Buffer.from(data).toString('base64');
        console.timeEnd('Base64');

        var filetype = process.argv[3].split('.');
        var mime = 'image/' + filetype[1];
        //console.log(process.argv[2], b64data, mime);
        console.time('Tx');
        const result20 = await contract.submitTransaction('add', process.argv[2], b64data, mime);
        console.timeEnd('Tx');
        console.log(`Transaction has been evaluated`);
        
        gateway.disconnect();

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
<<<<<<< HEAD
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')
        
=======
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        await contract.submitTransaction('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom');
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();
>>>>>>> 0031db5de282eb61a57d774a5dd8e12de79aa076

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
