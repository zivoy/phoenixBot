import {ButtonInteraction, Interaction} from "discord.js";
import { BotEvent } from "../types";
import button from "../handlers/button";

const event : BotEvent = {
    name: "interactionCreate",
    execute: (interaction: Interaction) => {
        if (interaction.isButton()){
            button(interaction as ButtonInteraction)
        }
    }
}

export default event;