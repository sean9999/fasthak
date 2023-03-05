
import hak from "./hak.js";

type registerSSE = () => Promise<EventSource>

export { hak, registerSSE };