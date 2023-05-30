export interface Response {
    error?: string
}

export interface DiscordVerifyRequest {
    discord_name?: string
    discord_id?: string
    code: string
}

export interface DiscordVerifyResponse extends Response {
    discord_id: string
}

export type nRPCFunction = (request: any) => Promise<Response>