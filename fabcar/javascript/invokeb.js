/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const ipfsClient = require('ipfs-http-client');
var fs = require('fs');

const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');

async function main() {
    try {
        if (process.argv.length != 8) {
            console.error("Failed: Expecting argument is 6");
            process.exit(1);
        }
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), '../wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('user1');
        if (!identity) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const Key = process.argv[2];
        const HospitalID = process.argv[3];
        const MedicalInfo = process.argv[4];
        const Hash = process.argv[5];
        const SecurityLv = process.argv[6];
        const Signature = process.argv[7];
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');


        // Get the contract from the network.
        const contract = network.getContract('fabcar');


        console.time('Tx');
        const Now = new Date();
        const result = await contract.submitTransaction('AddData', Key, HospitalID, Now.toISOString(), MedicalInfo, Hash, SecurityLv, Signature);
        console.timeEnd('Tx');
        
        console.log(`Transaction has been evaluated`);
        
        gateway.disconnect();

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')
        

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
