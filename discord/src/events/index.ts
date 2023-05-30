import ready from "./ready";
import interactionCreate from "./interactionCreate";
import {Client} from "discord.js";

export function ImportEvents(client: Client) {
    [
        ready,
        interactionCreate
    ].forEach((event) => {
        if (event.once)
            client.once(event.name, event.execute)
        else
            client.on(event.name, event.execute)
        console.log(`Loaded event ${event.name}`)
    })
}