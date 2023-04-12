import fetch, { Response } from "node-fetch";

export interface Tureng {
    translate(word: string, dictionary: Dictionary): Promise<TranslationResult | null>;
}

export enum Dictionary {
    EnglishTurkish = "entr",
    EnglishGerman = "ende",
    EnglishSpanish = "enes",
    EnglishFrench = "enfr"
}

export class Client implements Tureng {
    private static defaultURL = "http://api.tureng.com/v1/dictionary";
    private static defaultOptions = {
        debug: false,
        timeout: 1 * 1000
    }

    private static headers = {
        "User-Agent": "Tureng/2012061663 CFNetwork/1335.0.3 Darwin/21.6.0",
        "Accept": "*/*",
        "Accept-Language": "en-GB,en;q=0.9"
    }
    private url: string;
    private options: { debug: boolean, timeout: number };

    constructor(url: string = Client.defaultURL, options: { debug?: boolean, timeout?: number } = Client.defaultOptions) {
        this.url = url;
        this.options = { ...Client.defaultOptions, ...options };

    }

    async translate(word: string, dictionary: Dictionary = Dictionary.EnglishTurkish): Promise<TranslationResult | null> {
        const promise = this.doRequest(word, dictionary)

        const response = await promise;
        if (!response.ok) {
            if (response.status == 404) {
                return null;
            }

            throw new Error(`translate: request: ${response.status} ${response.statusText}`)
        }

        return await response.json() as TranslationResult;
    }

    private doRequest(word: string, dictionary: Dictionary): Promise<Response> {
        let addr = `${this.url}/${dictionary}/${word}`;
        return fetch(addr, {
            method: "GET",
            headers: Client.headers,
            signal: AbortSignal.timeout(this.options.timeout)
        })
    }
}

export interface TermResult {
    TermA: string;
    TermB: string;
    CategoryTextA: string;
    CategoryTextB: string;
    TermTypeTextA: string;
    TermTypeTextB: string;
    IsSlang: boolean;
}

export interface TranslationResult {
    SearchedTerm: string;
    IsFound: boolean;
    AResults: TermResult[];
    BResults: TermResult[];
    PrimeATerm: string;
}
