import {BotEvent} from "../types";
import type {Client} from "discord.js";

const event: BotEvent = {
    name: "ready",
    once: true,
    execute: (client: Client) => {
        console.log(`Logged in as ${client.user?.username}`)
    }
}

export default event;