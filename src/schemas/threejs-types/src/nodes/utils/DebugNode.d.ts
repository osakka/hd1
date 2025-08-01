import Node from "../core/Node.js";
import NodeBuilder from "../core/NodeBuilder.js";
import TempNode from "../core/TempNode.js";
import { ShaderNodeObject } from "../tsl/TSLCore.js";

declare class DebugNode extends TempNode {
    constructor(node: Node, callback?: ((code: string) => void) | null);
}

export default DebugNode;

export const debug: (
    node: Node,
    callback?: ((node: NodeBuilder, code: string) => void) | null,
) => ShaderNodeObject<DebugNode>;

declare module "../tsl/TSLCore.js" {
    interface NodeElements {
        debug: typeof debug;
    }
}
