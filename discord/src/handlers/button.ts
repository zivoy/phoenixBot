import {ButtonInteraction} from "discord.js";
import {VerifyButtonID} from "../consts";
import {request} from "http";

export default async function (interaction: ButtonInteraction) {
    switch (interaction.customId) {
        case VerifyButtonID:
            let u = new URL(process.env.VERIFYURL || "http://localhost:80/reload-rsi")
            let req = request({
                host: u.host,
                port: u.port,
                path: u.pathname,
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                }
            }, (r) => {
                console.log(r);
            });
            req.write(JSON.stringify({verified: interaction.user.id})) //todo also return which code so lookup is easier
            req.end();
            await interaction.update({components: []})
            break;
    }
}