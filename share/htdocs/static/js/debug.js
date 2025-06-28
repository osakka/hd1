// VWS Debug Helper - Run in browser console to test renderer directly
// Usage: loadScript('/static/js/debug.js'); then call debugVWS();

function debugVWS() {
    console.log('🔍 VWS Debug Session Started');
    
    // 1. Check if renderer exists and is initialized
    if (typeof renderer === 'undefined') {
        console.error('❌ Renderer not found! Renderer may not be initialized');
        return;
    }
    
    console.log('✅ Renderer found:', renderer);
    console.log('📊 Current objects in renderer:', renderer.objects.size);
    console.log('📷 Current camera position:', renderer.camera.position);
    console.log('🎯 Current camera target:', renderer.camera.target);
    
    // 2. Check WebGL context
    if (!renderer.gl) {
        console.error('❌ WebGL context not available');
        return;
    }
    
    console.log('✅ WebGL context:', renderer.gl);
    console.log('📐 Canvas size:', renderer.canvas.width, 'x', renderer.canvas.height);
    
    // 3. Test direct object creation bypassing API
    console.log('🧪 Testing direct object creation...');
    
    // Create a test cube directly in renderer
    const testCube = {
        id: 'debug-test-cube',
        name: 'test-cube',
        type: 'cube',
        transform: {
            position: { x: 0, y: 0, z: 0 },
            scale: { x: 2, y: 2, z: 2 },
            rotation: { x: 0, y: 0, z: 0 }
        },
        color: { r: 1, g: 0, b: 0, a: 1 },
        wireframe: false,
        visible: true
    };
    
    renderer.createObject(testCube);
    console.log('✅ Test cube created. Objects now:', renderer.objects.size);
    
    // Create a test wireframe sphere
    const testSphere = {
        id: 'debug-test-sphere',
        name: 'test-sphere',
        type: 'sphere',
        transform: {
            position: { x: 3, y: 0, z: 0 },
            scale: { x: 1, y: 1, z: 1 },
            rotation: { x: 0, y: 0, z: 0 }
        },
        color: { r: 0, g: 1, b: 0, a: 1 },
        wireframe: true,
        visible: true
    };
    
    renderer.createObject(testSphere);
    console.log('✅ Test sphere created. Objects now:', renderer.objects.size);
    
    // 4. Force initialize the VWS grid
    console.log('🌐 Initializing VWS coordinate grid...');
    renderer.initializeWorld({
        size: 25,
        transparency: 0.3,
        bounds: { min: -12, max: 12 }
    });
    
    console.log('✅ Grid initialized. Objects now:', renderer.objects.size);
    
    // 5. Check camera positioning
    console.log('📷 Adjusting camera for better view...');
    renderer.camera.position = [8, 8, 8];
    renderer.camera.target = [0, 0, 0];
    renderer.camera.fov = 45;
    
    // 6. Check if objects are being rendered
    console.log('🔍 Checking object rendering...');
    for (const [id, obj] of renderer.objects) {
        console.log(`📦 Object ${id}:`, obj);
        console.log(`   Type: ${obj.type}`);
        console.log(`   Visible: ${obj.visible}`);
        console.log(`   Position: ${obj.transform?.position?.x || 0}, ${obj.transform?.position?.y || 0}, ${obj.transform?.position?.z || 0}`);
    }
    
    // 7. Check geometries
    console.log('🔧 Checking geometries...');
    Object.keys(renderer.geometries).forEach(type => {
        const geom = renderer.geometries[type];
        console.log(`   ${type}: ${geom ? '✅' : '❌'} (indexCount: ${geom?.indexCount})`);
    });
    
    // 8. Manual render test
    console.log('🎨 Forcing manual render...');
    renderer.render(performance.now());
    
    console.log('🔍 VWS Debug Complete! You should now see objects on screen.');
    console.log('📋 Summary:');
    console.log(`   - Total objects: ${renderer.objects.size}`);
    console.log(`   - Camera: [${renderer.camera.position.join(', ')}] → [${renderer.camera.target.join(', ')}]`);
    console.log(`   - Canvas: ${renderer.canvas.width}x${renderer.canvas.height}`);
    
    return {
        renderer: renderer,
        objects: renderer.objects,
        camera: renderer.camera,
        canvas: renderer.canvas
    };
}

// Helper function to load scripts dynamically
function loadScript(src) {
    const script = document.createElement('script');
    script.src = src;
    document.head.appendChild(script);
    return script;
}

console.log('🔧 VWS Debug Helper loaded. Run debugVWS() to start debugging.');