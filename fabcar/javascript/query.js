/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
<<<<<<< HEAD
const { createDiffieHellman } = require('crypto');
const BufferList = require('bl/BufferList');
var fs = require("fs");
const CID = require("cids");
const ipfsClient = require('ipfs-http-client');
=======
const fs = require('fs');
>>>>>>> 0031db5de282eb61a57d774a5dd8e12de79aa076


async function main() {
    try {
<<<<<<< HEAD
        if (process.argv.length != 3) {
            console.error("Failed: Expecting argument is 1");
            process.exit(1);
        }
=======
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

>>>>>>> 0031db5de282eb61a57d774a5dd8e12de79aa076
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), '../wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('appUser');
        if (!identity) {
            console.log('An identity for the user "appUser" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'appUser', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('fabcar');

        // Evaluate the specified transaction.
        // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
        // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
        console.time('Tx');
        var result = await contract.evaluateTransaction('GetData', process.argv[2]);
        console.timeEnd('Tx')

        var Data = JSON.parse(result);
        console.log(Data)
        
        const ipfs = ipfsClient('http://localhost:5001')
        console.time('File Write');
        for await (const file of ipfs.get(new CID(Data.meta_data.hash))) {
          
            const content = new BufferList()
            for await (const chunk of file.content) {
              content.append(chunk)
            }
            await fs.writeFileSync("./tmp/" + Data.meta_data.hash, content, 'binary', "wb");
          }


        console.timeEnd('File Write');
        

        gateway.disconnect();

        // Disconnect from the gateway.
        await gateway.disconnect();
        
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

main();
