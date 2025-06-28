class VisualStreamRenderer {
    constructor(canvas) {
        this.canvas = canvas;
        this.gl = canvas.getContext('webgl2') || canvas.getContext('webgl');
        this.objects = new Map();
        this.camera = {
            position: [0, 0, 5],
            target: [0, 0, 0],
            fov: 45
        };
        
        if (!this.gl) {
            throw new Error('WebGL not supported');
        }
        
        this.init();
    }
    
    init() {
        const gl = this.gl;
        
        // Enable depth testing
        gl.enable(gl.DEPTH_TEST);
        gl.enable(gl.BLEND);
        gl.blendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
        
        // Create shader program
        this.program = this.createShaderProgram();
        gl.useProgram(this.program);
        
        // Get attribute and uniform locations
        this.locations = {
            position: gl.getAttribLocation(this.program, 'a_position'),
            color: gl.getUniformLocation(this.program, 'u_color'),
            wireframe: gl.getUniformLocation(this.program, 'u_wireframe'),
            transform: gl.getUniformLocation(this.program, 'u_transform'),
            view: gl.getUniformLocation(this.program, 'u_view'),
            projection: gl.getUniformLocation(this.program, 'u_projection')
        };
        
        // Create primitive geometries
        this.geometries = {
            cube: this.createCubeGeometry(),
            sphere: this.createSphereGeometry(),
            plane: this.createPlaneGeometry(),
            text: this.createTextGeometry(),
            grid: this.createGridGeometry()
        };
        
        // Create wireframe versions
        this.wireframes = {
            cube: this.createCubeWireframe(),
            sphere: this.createSphereWireframe(),
            plane: this.createPlaneWireframe(),
            text: this.createTextWireframe(),
            grid: this.createGridWireframe()
        };
        
        // Text overlay system
        this.textOverlay = document.getElementById('text-overlay');
        this.textElements = new Map();
        this.projectionMatrix = mat4.create();
        this.viewMatrix = mat4.create();
        
        // Start render loop
        this.lastTime = 0;
        this.render(0);
    }
    
    createShaderProgram() {
        const gl = this.gl;
        
        const vertexShaderSource = `
            attribute vec3 a_position;
            attribute vec2 a_texCoord;
            uniform mat4 u_transform;
            uniform mat4 u_view;
            uniform mat4 u_projection;
            varying vec2 v_texCoord;
            
            void main() {
                gl_Position = u_projection * u_view * u_transform * vec4(a_position, 1.0);
                v_texCoord = a_texCoord;
            }
        `;
        
        const fragmentShaderSource = `
            precision mediump float;
            uniform vec4 u_color;
            uniform bool u_wireframe;
            uniform bool u_useTexture;
            uniform sampler2D u_texture;
            varying vec2 v_texCoord;
            
            void main() {
                if (u_wireframe) {
                    gl_FragColor = vec4(u_color.rgb, 1.0);
                } else if (u_useTexture) {
                    vec4 texColor = texture2D(u_texture, v_texCoord);
                    gl_FragColor = vec4(u_color.rgb, texColor.a * u_color.a);
                } else {
                    gl_FragColor = u_color;
                }
            }
        `;
        
        const vertexShader = this.createShader(gl.VERTEX_SHADER, vertexShaderSource);
        const fragmentShader = this.createShader(gl.FRAGMENT_SHADER, fragmentShaderSource);
        
        const program = gl.createProgram();
        gl.attachShader(program, vertexShader);
        gl.attachShader(program, fragmentShader);
        gl.linkProgram(program);
        
        if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
            throw new Error('Failed to link shader program: ' + gl.getProgramInfoLog(program));
        }
        
        return program;
    }
    
    createShader(type, source) {
        const gl = this.gl;
        const shader = gl.createShader(type);
        gl.shaderSource(shader, source);
        gl.compileShader(shader);
        
        if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
            throw new Error('Failed to compile shader: ' + gl.getShaderInfoLog(shader));
        }
        
        return shader;
    }
    
    createCubeGeometry() {
        const positions = new Float32Array([
            // Front face
            -1, -1,  1,   1, -1,  1,   1,  1,  1,  -1,  1,  1,
            // Back face
            -1, -1, -1,  -1,  1, -1,   1,  1, -1,   1, -1, -1,
            // Top face
            -1,  1, -1,  -1,  1,  1,   1,  1,  1,   1,  1, -1,
            // Bottom face
            -1, -1, -1,   1, -1, -1,   1, -1,  1,  -1, -1,  1,
            // Right face
             1, -1, -1,   1,  1, -1,   1,  1,  1,   1, -1,  1,
            // Left face
            -1, -1, -1,  -1, -1,  1,  -1,  1,  1,  -1,  1, -1
        ]);
        
        const indices = new Uint16Array([
            0,  1,  2,    0,  2,  3,    // front
            4,  5,  6,    4,  6,  7,    // back
            8,  9,  10,   8,  10, 11,   // top
            12, 13, 14,   12, 14, 15,   // bottom
            16, 17, 18,   16, 18, 19,   // right
            20, 21, 22,   20, 22, 23    // left
        ]);
        
        return this.createBuffers(positions, indices);
    }
    
    createSphereGeometry() {
        const positions = [];
        const indices = [];
        const radius = 1;
        const latBands = 16;
        const lonBands = 16;
        
        for (let lat = 0; lat <= latBands; lat++) {
            const theta = lat * Math.PI / latBands;
            const sinTheta = Math.sin(theta);
            const cosTheta = Math.cos(theta);
            
            for (let lon = 0; lon <= lonBands; lon++) {
                const phi = lon * 2 * Math.PI / lonBands;
                const sinPhi = Math.sin(phi);
                const cosPhi = Math.cos(phi);
                
                const x = cosPhi * sinTheta;
                const y = cosTheta;
                const z = sinPhi * sinTheta;
                
                positions.push(radius * x, radius * y, radius * z);
            }
        }
        
        for (let lat = 0; lat < latBands; lat++) {
            for (let lon = 0; lon < lonBands; lon++) {
                const first = (lat * (lonBands + 1)) + lon;
                const second = first + lonBands + 1;
                
                indices.push(first, second, first + 1);
                indices.push(second, second + 1, first + 1);
            }
        }
        
        return this.createBuffers(new Float32Array(positions), new Uint16Array(indices));
    }
    
    createCubeWireframe() {
        const positions = new Float32Array([
            // Front face
            -1, -1,  1,   1, -1,  1,   1,  1,  1,  -1,  1,  1,
            // Back face
            -1, -1, -1,  -1,  1, -1,   1,  1, -1,   1, -1, -1
        ]);
        
        const indices = new Uint16Array([
            // Front face edges
            0, 1, 1, 2, 2, 3, 3, 0,
            // Back face edges  
            4, 5, 5, 6, 6, 7, 7, 4,
            // Connecting edges
            0, 7, 1, 6, 2, 5, 3, 4
        ]);
        
        return this.createWireframeBuffers(positions, indices);
    }
    
    createSphereWireframe() {
        const positions = [];
        const indices = [];
        const radius = 1;
        const latBands = 8;
        const lonBands = 12;
        
        for (let lat = 0; lat <= latBands; lat++) {
            const theta = lat * Math.PI / latBands;
            const sinTheta = Math.sin(theta);
            const cosTheta = Math.cos(theta);
            
            for (let lon = 0; lon <= lonBands; lon++) {
                const phi = lon * 2 * Math.PI / lonBands;
                const sinPhi = Math.sin(phi);
                const cosPhi = Math.cos(phi);
                
                const x = cosPhi * sinTheta;
                const y = cosTheta;
                const z = sinPhi * sinTheta;
                
                positions.push(radius * x, radius * y, radius * z);
            }
        }
        
        // Latitude lines
        for (let lat = 0; lat < latBands; lat++) {
            for (let lon = 0; lon < lonBands; lon++) {
                const first = lat * (lonBands + 1) + lon;
                const second = first + lonBands + 1;
                
                indices.push(first, first + 1);
                indices.push(first, second);
            }
        }
        
        return this.createWireframeBuffers(new Float32Array(positions), new Uint16Array(indices));
    }
    
    createPlaneWireframe() {
        const positions = new Float32Array([
            -1, 0, -1,
             1, 0, -1,
             1, 0,  1,
            -1, 0,  1
        ]);
        
        const indices = new Uint16Array([0, 1, 1, 2, 2, 3, 3, 0]);
        
        return this.createWireframeBuffers(positions, indices);
    }
    
    createWireframeBuffers(positions, indices) {
        const gl = this.gl;
        
        const positionBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);
        gl.bufferData(gl.ARRAY_BUFFER, positions, gl.STATIC_DRAW);
        
        const indexBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);
        gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW);
        
        return {
            position: positionBuffer,
            indices: indexBuffer,
            indexCount: indices.length
        };
    }
    
    createPlaneGeometry() {
        const positions = new Float32Array([
            -1, 0, -1,
             1, 0, -1,
             1, 0,  1,
            -1, 0,  1
        ]);
        
        const indices = new Uint16Array([0, 1, 2, 0, 2, 3]);
        
        return this.createBuffers(positions, indices);
    }
    
    createTextGeometry() {
        // Simple quad for text rendering
        const positions = new Float32Array([
            -1, -1, 0,
             1, -1, 0,
             1,  1, 0,
            -1,  1, 0
        ]);
        
        const indices = new Uint16Array([0, 1, 2, 0, 2, 3]);
        
        return this.createBuffers(positions, indices);
    }
    
    createTextWireframe() {
        const positions = new Float32Array([
            -1, -1, 0,
             1, -1, 0,
             1,  1, 0,
            -1,  1, 0
        ]);
        
        const indices = new Uint16Array([0, 1, 1, 2, 2, 3, 3, 0]);
        
        return this.createWireframeBuffers(positions, indices);
    }
    
    createGridGeometry() {
        // VWS 25x25x25 coordinate grid with [-12, +12] bounds
        const positions = [];
        const indices = [];
        const gridSize = 25;
        const minCoord = -12;
        const maxCoord = 12;
        const step = (maxCoord - minCoord) / (gridSize - 1);
        
        let vertexIndex = 0;
        
        // Create grid lines along X-axis (YZ planes)
        for (let y = minCoord; y <= maxCoord; y += step) {
            for (let z = minCoord; z <= maxCoord; z += step) {
                positions.push(minCoord, y, z);
                positions.push(maxCoord, y, z);
                indices.push(vertexIndex, vertexIndex + 1);
                vertexIndex += 2;
            }
        }
        
        // Create grid lines along Y-axis (XZ planes)
        for (let x = minCoord; x <= maxCoord; x += step) {
            for (let z = minCoord; z <= maxCoord; z += step) {
                positions.push(x, minCoord, z);
                positions.push(x, maxCoord, z);
                indices.push(vertexIndex, vertexIndex + 1);
                vertexIndex += 2;
            }
        }
        
        // Create grid lines along Z-axis (XY planes)
        for (let x = minCoord; x <= maxCoord; x += step) {
            for (let y = minCoord; y <= maxCoord; y += step) {
                positions.push(x, y, minCoord);
                positions.push(x, y, maxCoord);
                indices.push(vertexIndex, vertexIndex + 1);
                vertexIndex += 2;
            }
        }
        
        return this.createWireframeBuffers(new Float32Array(positions), new Uint16Array(indices));
    }
    
    createGridWireframe() {
        // Grid is already a wireframe, return same as grid geometry
        return this.createGridGeometry();
    }
    
    createBuffers(positions, indices) {
        const gl = this.gl;
        
        const positionBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);
        gl.bufferData(gl.ARRAY_BUFFER, positions, gl.STATIC_DRAW);
        
        const indexBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);
        gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW);
        
        return {
            position: positionBuffer,
            indices: indexBuffer,
            indexCount: indices.length
        };
    }
    
    createTransformMatrix(transform) {
        const matrix = mat4.create();
        
        if (transform.position) {
            mat4.translate(matrix, matrix, [transform.position.x, transform.position.y, transform.position.z]);
        }
        
        if (transform.rotation) {
            mat4.rotateX(matrix, matrix, transform.rotation.x);
            mat4.rotateY(matrix, matrix, transform.rotation.y);
            mat4.rotateZ(matrix, matrix, transform.rotation.z);
        }
        
        if (transform.scale) {
            mat4.scale(matrix, matrix, [transform.scale.x, transform.scale.y, transform.scale.z]);
        }
        
        return matrix;
    }
    
    processMessage(message) {
        switch (message.type) {
            case 'create':
                message.objects?.forEach(obj => this.createObject(obj));
                break;
            case 'update':
                message.objects?.forEach(obj => this.updateObject(obj));
                break;
            case 'delete':
                message.objects?.forEach(obj => this.deleteObject(obj.id));
                break;
            case 'clear':
                this.objects.clear();
                // Clear all text elements
                for (const [id, element] of this.textElements) {
                    this.textOverlay.removeChild(element);
                }
                this.textElements.clear();
                break;
            case 'camera':
                if (message.position) {
                    this.camera.position = [message.position.x, message.position.y, message.position.z];
                }
                if (message.target) {
                    this.camera.target = [message.target.x, message.target.y, message.target.z];
                }
                break;
        }
    }
    
    updateCamera(camera) {
        if (camera.position) {
            this.camera.position = [camera.position.x, camera.position.y, camera.position.z];
        }
        if (camera.target) {
            this.camera.target = [camera.target.x, camera.target.y, camera.target.z];
        }
        if (camera.fov) {
            this.camera.fov = camera.fov;
        }
    }
    
    render(currentTime) {
        const gl = this.gl;
        
        // Clear canvas
        gl.viewport(0, 0, this.canvas.width, this.canvas.height);
        gl.clearColor(0.1, 0.1, 0.1, 1.0);
        gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
        
        // Setup matrices
        mat4.perspective(this.projectionMatrix, this.camera.fov * Math.PI / 180, 
                        this.canvas.width / this.canvas.height, 0.1, 100.0);
        
        mat4.lookAt(this.viewMatrix, this.camera.position, this.camera.target, [0, 1, 0]);
        
        gl.uniformMatrix4fv(this.locations.projection, false, this.projectionMatrix);
        gl.uniformMatrix4fv(this.locations.view, false, this.viewMatrix);
        
        // Render 3D objects
        for (const [id, obj] of this.objects) {
            if (obj.visible === false) continue;
            
            if (obj.type === 'text') {
                this.renderTextOverlay(obj, currentTime);
            } else {
                this.renderObject(obj, currentTime);
            }
        }
        
        requestAnimationFrame((time) => this.render(time));
    }
    
    renderObject(obj, currentTime) {
        const gl = this.gl;
        
        // Special handling for text
        if (obj.type === 'text') {
            this.renderText(obj, currentTime);
            return;
        }
        
        const isWireframe = obj.wireframe || false;
        const geometry = isWireframe ? this.wireframes[obj.type] : this.geometries[obj.type];
        
        if (!geometry) return;
        
        // Apply animations
        let transform = obj.transform || {};
        if (obj.animation) {
            transform = this.applyAnimation(obj, transform, currentTime);
        }
        
        // Create transform matrix
        const transformMatrix = this.createTransformMatrix(transform);
        gl.uniformMatrix4fv(this.locations.transform, false, transformMatrix);
        
        // Set color
        const color = obj.color || {r: 1, g: 1, b: 1, a: 1};
        gl.uniform4f(this.locations.color, color.r, color.g, color.b, color.a);
        gl.uniform1i(this.locations.wireframe, isWireframe);
        
        // Bind geometry
        gl.bindBuffer(gl.ARRAY_BUFFER, geometry.position);
        gl.enableVertexAttribArray(this.locations.position);
        gl.vertexAttribPointer(this.locations.position, 3, gl.FLOAT, false, 0, 0);
        
        gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, geometry.indices);
        
        // Draw
        if (isWireframe) {
            gl.drawElements(gl.LINES, geometry.indexCount, gl.UNSIGNED_SHORT, 0);
        } else {
            gl.drawElements(gl.TRIANGLES, geometry.indexCount, gl.UNSIGNED_SHORT, 0);
        }
    }
    
    renderTextOverlay(obj, currentTime) {
        console.log('renderTextOverlay called for:', obj.id);
        
        const text = obj.data?.text || 'TEXT';
        const fontSize = obj.data?.fontSize || 1;
        const fontFamily = obj.data?.fontFamily || 'Arial, sans-serif';
        
        console.log('Text content:', text, 'Size:', fontSize);
        
        // Apply animations to get current transform
        let transform = obj.transform || {};
        if (obj.animation) {
            transform = this.applyAnimation(obj, transform, currentTime);
        }
        
        // Get 3D position
        const worldPos = [
            transform.position?.x || 0,
            transform.position?.y || 0, 
            transform.position?.z || 0,
            1
        ];
        
        console.log('Text world position:', worldPos);
        
        // Project 3D position to screen coordinates
        const screenPos = this.worldToScreen(worldPos);
        
        console.log('Text screen position:', screenPos);
        
        // Skip if behind camera
        if (screenPos.z < 0 || screenPos.z > 1) {
            console.log('Text behind camera, hiding');
            this.hideTextElement(obj.id);
            return;
        }
        
        // Get or create text element
        let textElement = this.textElements.get(obj.id);
        if (!textElement) {
            textElement = document.createElement('div');
            textElement.className = 'text-element';
            textElement.id = 'text-' + obj.id;
            this.textOverlay.appendChild(textElement);
            this.textElements.set(obj.id, textElement);
        }
        
        // Update text content and style
        textElement.textContent = text;
        textElement.style.fontSize = (fontSize * 48) + 'px';
        textElement.style.fontFamily = fontFamily;
        
        // Apply color
        const color = obj.color || {r: 1, g: 1, b: 1, a: 1};
        textElement.style.color = `rgba(${Math.floor(color.r * 255)}, ${Math.floor(color.g * 255)}, ${Math.floor(color.b * 255)}, ${color.a})`;
        
        // Position on screen
        textElement.style.left = screenPos.x + 'px';
        textElement.style.top = screenPos.y + 'px';
        
        // Apply scale based on distance
        const scale = Math.max(0.1, Math.min(2, 10 / screenPos.w));
        const rotation = transform.rotation || {x: 0, y: 0, z: 0};
        
        textElement.style.transform = `translate(-50%, -50%) scale(${scale}) rotateZ(${rotation.z}rad)`;
        textElement.style.display = 'block';
        
        // Add animation class if needed
        if (obj.animation) {
            textElement.classList.add('animate');
        }
    }
    
    worldToScreen(worldPos) {
        // Create model matrix (identity for now)
        const model = mat4.create();
        
        // Create MVP matrix: MVP = Projection * View * Model
        const mv = mat4.create();
        const mvp = mat4.create();
        
        mat4.multiply(mv, this.viewMatrix, model);
        mat4.multiply(mvp, this.projectionMatrix, mv);
        
        // Transform world position to clip space
        const clipPos = [0, 0, 0, 0];
        
        // Matrix multiplication: clipPos = mvp * worldPos
        clipPos[0] = mvp[0] * worldPos[0] + mvp[4] * worldPos[1] + mvp[8]  * worldPos[2] + mvp[12] * worldPos[3];
        clipPos[1] = mvp[1] * worldPos[0] + mvp[5] * worldPos[1] + mvp[9]  * worldPos[2] + mvp[13] * worldPos[3];
        clipPos[2] = mvp[2] * worldPos[0] + mvp[6] * worldPos[1] + mvp[10] * worldPos[2] + mvp[14] * worldPos[3];
        clipPos[3] = mvp[3] * worldPos[0] + mvp[7] * worldPos[1] + mvp[11] * worldPos[2] + mvp[15] * worldPos[3];
        
        // Debug logging
        console.log('World pos:', worldPos);
        console.log('Clip pos:', clipPos);
        
        // Perspective divide
        if (Math.abs(clipPos[3]) > 0.001) {
            clipPos[0] /= clipPos[3];
            clipPos[1] /= clipPos[3];
            clipPos[2] /= clipPos[3];
        }
        
        console.log('NDC pos:', clipPos);
        
        // Convert to screen coordinates
        const screenX = (clipPos[0] * 0.5 + 0.5) * this.canvas.width;
        const screenY = (1 - (clipPos[1] * 0.5 + 0.5)) * this.canvas.height;
        
        console.log('Screen pos:', screenX, screenY);
        
        return {
            x: screenX,
            y: screenY,
            z: clipPos[2],
            w: clipPos[3]
        };
    }
    
    hideTextElement(id) {
        const textElement = this.textElements.get(id);
        if (textElement) {
            textElement.style.display = 'none';
        }
    }
    
    moveCamera(x, y, z) {
        this.camera.position[0] += x;
        this.camera.position[1] += y;
        this.camera.position[2] += z;
    }
    
    createObject(obj) {
        this.objects.set(obj.id, {
            ...obj,
            startTime: performance.now()
        });
    }
    
    updateObject(obj) {
        const existing = this.objects.get(obj.id);
        if (existing) {
            this.objects.set(obj.id, {
                ...existing,
                ...obj,
                startTime: obj.animation ? performance.now() : existing.startTime
            });
        }
    }
    
    deleteObject(id) {
        this.objects.delete(id);
        
        // Also remove text element
        const textElement = this.textElements.get(id);
        if (textElement) {
            this.textOverlay.removeChild(textElement);
            this.textElements.delete(id);
        }
    }
    
    initializeWorld(worldData) {
        // Show the VWS 25x25x25 coordinate grid when world is initialized
        const gridObject = {
            id: 'vws-coordinate-grid',
            name: 'coordinate-grid',
            type: 'grid',
            x: 0,
            y: 0,
            z: 0,
            scale: 1,
            wireframe: true,
            color: { r: 0.3, g: 0.7, b: 1.0, a: 0.3 }, // Semi-transparent blue
            visible: true
        };
        
        this.createObject(gridObject);
        console.log('VWS World initialized - 25x25x25 grid displayed with bounds [-12, +12]');
    }
    
    applyAnimation(obj, transform, currentTime) {
        if (!obj.animation || !obj.startTime) return transform;
        
        const elapsed = (currentTime - obj.startTime) / 1000;
        const duration = obj.animation.duration;
        
        if (elapsed >= duration && !obj.animation.loop) {
            return transform;
        }
        
        const t = obj.animation.loop ? (elapsed % duration) / duration : Math.min(elapsed / duration, 1);
        const eased = this.applyEasing(t, obj.animation.easing || 'linear');
        
        // Enhanced animation - always create rotation if animation exists
        const animated = { ...transform };
        
        // Ensure rotation object exists
        if (!animated.rotation) {
            animated.rotation = { x: 0, y: 0, z: 0 };
        } else {
            animated.rotation = { ...animated.rotation };
        }
        
        // Add animated rotation on multiple axes for more visible effect
        const baseRotation = animated.rotation;
        animated.rotation.x = (baseRotation.x || 0) + eased * Math.PI * 2 * 0.5;
        animated.rotation.y = (baseRotation.y || 0) + eased * Math.PI * 2;
        animated.rotation.z = (baseRotation.z || 0) + eased * Math.PI * 2 * 0.3;
        
        return animated;
    }
    
    applyEasing(t, easing) {
        switch (easing) {
            case 'ease-in-out':
                return t < 0.5 ? 2 * t * t : 1 - Math.pow(-2 * t + 2, 2) / 2;
            case 'bounce':
                return 1 - Math.abs(Math.cos(t * Math.PI * 3)) * (1 - t);
            case 'elastic':
                return Math.sin(t * Math.PI * 6) * (1 - t) + t;
            default:
                return t;
        }
    }
}