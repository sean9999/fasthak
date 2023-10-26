
import hak from "./hak";

type registerSSE = () => Promise<EventSource>

export { hak, registerSSE };