const redis = require("redis");
const TronWeb = require("tronweb");
const axios = require("axios");
const { scheduleJob } = require("node-schedule");

// initialize tron web
let tron_network = "https://api.trongrid.io";
let fullNode = tron_network;
let solidityNode = tron_network;
let eventServer = tron_network;
let privateKey = '';
const tronWeb = new TronWeb(fullNode, solidityNode, eventServer, privateKey);

// redis client
const redis_client = redis.createClient({ url: 'redis://localhost:6379' });

// token contract address to monitor
const token_contract_address = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"; // USDT
const redis_address_key = `tron_token_${token_contract_address}_users`;

let current_block, block_number, job;

(async () => {
    await redis_client.connect()
    current_block = await tronWeb.trx.getCurrentBlock();
    block_number = current_block.block_header.raw_data.number;


    job = scheduleJob(`*/3 * * * * *`, main);

})();

async function main() {
    current_block = await tronWeb.trx.getCurrentBlock();

    let gap = current_block.block_header.raw_data.number - block_number;
    if ((gap) <= 0) {
        return;
    }

    try {
        console.log(`start checking ${block_number}...`);
        let resp = await axios.get(`${tron_network}/v1/contracts/${token_contract_address}/events?event_name=Transfer&only_unconfirmed=false&only_confirmed=false&block_number=${block_number}&order_by=block_timestamp,asc&limit=200`);

        resp.data.data.map(async event => {
            let to_addr = tronWeb.address.fromHex(event.result.to);
            let from_addr = tronWeb.address.fromHex(event.result.from);
            let transfer_amount = parseInt(event.result.value);

            // filter target user address
            if (await redis_client.sIsMember(redis_address_key, to_addr) || await redis_client.sIsMember(redis_address_key, from_addr)) {
                console.log(`Monitor a transfer from ${from_addr} to ${to_addr} amount ${transfer_amount} for token ${token_contract_address}.`)
            }
        });
        block_number += 1;
    } catch (e) {
        console.log(`error ${e} skip.`)
    }
}
