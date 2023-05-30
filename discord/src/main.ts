import {ImportEvents} from "./events"
import {ConnectNatsListener} from "./nats/nrpc";
import {client} from "./bot";
import process from "process";

const token = process.env.DISCORD_TOKEN;
if (token === undefined) {
    console.error("discord token not provided");
    process.exit(1);
}

ConnectNatsListener(process.env.NATS_URI)

ImportEvents(client)
client.login(token)
