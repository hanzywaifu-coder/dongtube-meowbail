export interface SectionRow {
    title: string;
    description?: string;
    id: string;
}

export interface Section {
    title: string;
    highlight_label?: string;
    rows: SectionRow[];
}

export interface InteractiveMenuOptions {
    body?: string;
    footer?: string;
    thumbnail?: Buffer | string;
    fileName?: string;
    sections?: Section[];
    ctaUrl?: string;
    ctaText?: string;
    copyCode?: string;
    copyText?: string;
}

export interface SocketConfig {
    printQRInTerminal?: boolean;
    pairingCode?: string;
    phoneNumber?: string;
    newsletterJid?: string;
    newsletterName?: string;
}

export class DongtubeMeowbail {
    constructor(options?: SocketConfig);
    requestPairingCode(phoneNumber: string, customCode?: string): Promise<string>;
    sendGroupStatus(groupJid: string, content: any): Promise<any>;
    sendInteractiveMenu(jid: string, menuOptions: InteractiveMenuOptions): Promise<any>;
    sendMessage(jid: string, content: any, options?: any): Promise<any>;
}

export function makeWASocket(config?: SocketConfig): DongtubeMeowbail;
export default makeWASocket;
