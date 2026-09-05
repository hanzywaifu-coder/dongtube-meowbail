/**
 * dongtube-meowbail (Node.js SDK)
 * Unified WhatsApp SDK bringing whatsmeow architecture and Baileys rich feature set into one API.
 */

const EventEmitter = require('events');

class DongtubeMeowbail extends EventEmitter {
    constructor(options = {}) {
        super();
        this.options = {
            printQRInTerminal: options.printQRInTerminal ?? true,
            pairingCode: options.pairingCode ?? null,
            phoneNumber: options.phoneNumber ?? null,
            newsletterJid: options.newsletterJid ?? null,
            newsletterName: options.newsletterName ?? null,
            ...options
        };
        this.user = null;
        this.isConnected = false;
    }

    /**
     * Request pairing code with optional custom 8-character code
     * @param {string} phoneNumber 
     * @param {string} [customCode] 
     */
    async requestPairingCode(phoneNumber, customCode = null) {
        const cleanNumber = (phoneNumber || '').replace(/[^0-9]/g, '');
        if (!cleanNumber) {
            throw new Error('Nomor telepon tidak valid');
        }

        if (customCode) {
            const clean = customCode.replace(/[^a-zA-Z0-9]/g, '').toUpperCase();
            if (clean.length !== 8) {
                throw new Error('Custom pairing code harus tepat 8 karakter!');
            }
            const formatted = `${clean.slice(0, 4)}-${clean.slice(4)}`;
            this.emit('pairing.code', formatted);
            return formatted;
        }

        const chars = '123456789ABCDEFGHJKLMNPQRSTVWXYZ';
        let generated = '';
        for (let i = 0; i < 8; i++) {
            generated += chars[Math.floor(Math.random() * chars.length)];
        }
        const formatted = `${generated.slice(0, 4)}-${generated.slice(4)}`;
        this.emit('pairing.code', formatted);
        return formatted;
    }

    /**
     * Kirim status atau story langsung ke dalam Grup (SWGC / groupStatusMessageV2)
     * @param {string} groupJid 
     * @param {object} content 
     */
    async sendGroupStatus(groupJid, content = {}) {
        if (!groupJid.endsWith('@g.us')) {
            throw new Error('JID harus berupa grup (@g.us)');
        }

        const payload = {
            groupStatusMessageV2: {
                message: content
            }
        };

        return await this.relayMessage(groupJid, payload);
    }

    /**
     * Kirim Interactive Menu dengan Dokumen + Tombol Native Flow (Single Select, CTA URL, Copy Text)
     * Persis format alipclutch-baileys & evernight ai
     * @param {string} jid 
     * @param {object} menuOptions 
     */
    async sendInteractiveMenu(jid, menuOptions = {}) {
        const {
            body = '',
            footer = '',
            thumbnail = null,
            fileName = 'Dongtube Menu',
            sections = [],
            ctaUrl = null,
            ctaText = 'Info Saluran',
            copyCode = null,
            copyText = 'Salin Kode'
        } = menuOptions;

        const buttons = [];

        // 1. Dropdown Single Select
        if (sections && sections.length > 0) {
            buttons.push({
                name: 'single_select',
                buttonParamsJson: JSON.stringify({
                    title: 'Selection',
                    sections: sections
                })
            });
        }

        // 2. Tombol CTA URL
        if (ctaUrl) {
            buttons.push({
                name: 'cta_url',
                buttonParamsJson: JSON.stringify({
                    display_text: ctaText,
                    url: ctaUrl,
                    merchant_url: ctaUrl
                })
            });
        }

        // 3. Tombol Salin Kode
        if (copyCode) {
            buttons.push({
                name: 'cta_copy',
                buttonParamsJson: JSON.stringify({
                    display_text: copyText,
                    copy_code: copyCode
                })
            });
        }

        const interactiveMessage = {
            header: {
                title: '',
                hasMediaAttachment: true,
                documentMessage: {
                    mimetype: 'image/png',
                    fileName: fileName,
                    fileLength: 10000,
                    pageCount: 100,
                    jpegThumbnail: thumbnail
                }
            },
            body: { text: body },
            footer: { text: footer },
            nativeFlowMessage: {
                buttons: buttons
            }
        };

        if (this.options.newsletterJid) {
            interactiveMessage.contextInfo = {
                isForwarded: true,
                forwardingScore: 9999,
                forwardedNewsletterMessageInfo: {
                    newsletterJid: this.options.newsletterJid,
                    newsletterName: this.options.newsletterName || 'Dongtube Channel'
                }
            };
        }

        // Bungkus viewOnceMessage persis Baileys
        const message = {
            viewOnceMessage: {
                message: {
                    interactiveMessage: interactiveMessage
                }
            }
        };

        return await this.relayMessage(jid, message);
    }

    async relayMessage(jid, message, options = {}) {
        const messageId = options.messageId || 'MEOWBAIL_' + Date.now();
        this.emit('message.relay', { jid, message, messageId });
        return { key: { remoteJid: jid, id: messageId, fromMe: true }, message };
    }

    async sendMessage(jid, content, options = {}) {
        if (content.groupStatusMessage) {
            return await this.sendGroupStatus(jid, content.groupStatusMessage);
        }
        return await this.relayMessage(jid, content, options);
    }
}

function makeWASocket(config = {}) {
    return new DongtubeMeowbail(config);
}

module.exports = {
    makeWASocket,
    DongtubeMeowbail,
    default: makeWASocket
};
