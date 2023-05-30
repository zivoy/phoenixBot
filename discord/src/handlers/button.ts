import {ButtonInteraction} from "discord.js";
import {VerifyButtonID} from "../consts";
import {request} from "http";

export default async function (interaction: ButtonInteraction) {
    switch (interaction.customId) {
        case VerifyButtonID:
            let req = request({
                host: "closure-compiler.appspot.com",
                port: "80",
                path: "/reload-rsi",
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                }
            }, (r) => {
                console.log(r);
            });
            req.write(JSON.stringify({verified: interaction.user.id}))
            req.end();
            await interaction.update({components: []})
            break;
    }
}