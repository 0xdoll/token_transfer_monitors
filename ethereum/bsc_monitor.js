const ethers = require("ethers");
const { scheduleJob } = require("node-schedule");
const redis = require("redis");
require('dotenv').config();

const provider = new ethers.providers.JsonRpcProvider(`https://bsc-ws-node.nariox.org`);
const erc20_abi = ["event Transfer(address indexed from, address indexed to, uint256 value)", { "constant": true, "inputs": [], "name": "decimals", "outputs": [{ "name": "", "type": "uint256" }], "payable": false, "stateMutability": "view", "type": "function" }, { "constant": true, "inputs": [], "name": "name", "outputs": [{ "name": "", "type": "string" }], "payable": false, "stateMutability": "view", "type": "function" }, { "constant": true, "inputs": [], "name": "symbol", "outputs": [{ "name": "", "type": "string" }], "payable": false, "stateMutability": "view", "type": "function" }];
const erc20_token_address = "0x55d398326f99059ff775485246999027b3197955"; //  USDT
const erc20_contract = new ethers.Contract(erc20_token_address, erc20_abi, provider);

// redis client
const redis_client = redis.createClient({ url: 'redis://localhost:6379' });
const redis_address_key = `bsc_token_${erc20_token_address}_users`;

let current_block, latest_block, token_decimals, token_symbol, token_name, job;
(async () => {
    current_block = await provider.getBlockNumber();

    token_decimals = await erc20_contract.decimals();
    token_symbol = await erc20_contract.symbol();
    token_name = await erc20_contract.name();
    await redis_client.connect()

    console.log(`Start monitor ERC20 token ${token_name}(${token_symbol}) @ address ${erc20_token_address}...`);

    job = scheduleJob(`*/5 * * * * *`, main);

})();

async function main() {

    latest_block = await provider.getBlockNumber();

    if ((latest_block - current_block) <= 0) {
        return;
    }
    console.log(`\tcheck current_block ${current_block} - latest_block ${latest_block}`);
    let transfer_logs = await provider.getLogs({
        address: erc20_token_address,
        fromBlock: current_block + 1,
        toBlock: latest_block,
        topics: erc20_contract.interface.encodeFilterTopics("Transfer", []),
        // [ ethers.utils.keccak256(ethers.utils.toUtf8Bytes("Transfer(address,address,uint256)")) ],
    });

    for (let transfer_log of transfer_logs) {

        let decode_transfer_log = erc20_contract.interface.decodeEventLog("Transfer", transfer_log.data, transfer_log.topics)

        console.log(`\t Transfer @ block ${transfer_log.blockNumber} from ${decode_transfer_log.from} to ${decode_transfer_log.to} amount ${ethers.utils.formatUnits(decode_transfer_log.value, token_decimals)} for token ${erc20_token_address}.`);

        // filter useful user address
        if (await redis_client.sIsMember(redis_address_key, decode_transfer_log.to) || await redis_client.sIsMember(redis_address_key, decode_transfer_log.from)) {
            console.log(`Monitor a transfer @ block ${transfer_log.blockNumber} from ${decode_transfer_log.from} to ${decode_transfer_log.to} amount ${decode_transfer_log.value} for token ${erc20_token_address}.`);
        }
    }

    current_block = latest_block;

}

