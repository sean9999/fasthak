/**
 * hak is a singleton
 */
interface Ihak {
    DEBUG: boolean;
    PREFIX: string;
    waitFor: (ms: number) => Promise<void>;
    run: (fn: () => void) => void;
    sse: EventSource | null;
}
declare const hak: Ihak;
export default hak;
