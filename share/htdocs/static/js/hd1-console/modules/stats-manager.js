/**
 * HD1 Stats Manager - Real-time Performance Monitoring with Multiple Graphs
 * 
 * Manages FPS, Memory, GPU, and WebSocket traffic graphs with 30-second history.
 * All graphs feature smooth rendering and real-time data collection.
 */
class HD1StatsManager {
    constructor(domManager) {
        this.dom = domManager;
        this.ready = false;
        this.graphsInitialized = false;
        
        // Graph configurations - unified for all graphs
        this.graphConfig = {
            width: 200, // Will be dynamically updated by resizeCanvas
            height: 100, // Will be dynamically updated by resizeCanvas
            historyLength: 30, // 30 data points (1 per second for 30 seconds)
            updateInterval: 250, // Update every 250ms for smooth real-time feel
            fpsUpdateInterval: 100 // FPS sampled more frequently but averaged per second
        };
        
        // FPS tracking
        this.fpsHistory = [];
        this.fpsCanvas = null;
        this.fpsContext = null;
        
        // Memory tracking
        this.memoryHistory = [];
        this.memoryCanvas = null;
        this.memoryContext = null;
        
        // WebSocket latency tracking
        this.latencyHistory = [];
        this.latencyCanvas = null;
        this.latencyContext = null;
        
        // WebSocket tracking
        this.wsInHistory = [];
        this.wsOutHistory = [];
        this.wsCanvas = null;
        this.wsContext = null;
        this.wsInBytes = 0;
        this.wsOutBytes = 0;
        this.wsLastSample = { in: 0, out: 0 };
        
        // Update intervals
        this.updateInterval = null;
        this.fpsUpdateInterval = null;
        
        // Performance monitoring
        this.lastFrameTime = performance.now();
        this.frameCount = 0;
        this.fpsAccumulator = [];
        this.lastSecond = Math.floor(Date.now() / 1000);
    }

    /**
     * Initialize stats manager
     */
    async initialize() {
        console.log('[HD1-Stats] Initializing Enhanced Stats Manager');
        
        try {
            // Check if console is collapsed on startup
            const debugPanel = this.dom.get('debug-panel');
            const isCollapsed = debugPanel && debugPanel.classList.contains('collapsed');
            
            if (isCollapsed) {
                console.log('[HD1-Stats] Console is collapsed, deferring graph initialization');
                this.graphsInitialized = false;
            } else {
                // Initialize all graphs immediately
                this.initializeGraphs();
                this.graphsInitialized = true;
            }
            
            // Start monitoring (data collection works regardless of visibility)
            this.startMonitoring();
            
            this.ready = true;
            console.log('[HD1-Stats] Enhanced Stats Manager ready');
        } catch (error) {
            console.error('[HD1-Stats] Initialization failed:', error);
        }
    }

    /**
     * Initialize all graphs
     */
    initializeGraphs() {
        this.initializeFpsGraph();
        this.initializeMemoryGraph();
        this.initializeLatencyGraph();
        this.initializeWebSocketGraph();
        console.log('[HD1-Stats] All graphs initialized');
    }

    /**
     * Ensure graphs are initialized (called when expanding)
     */
    ensureGraphsInitialized() {
        if (!this.graphsInitialized) {
            console.log('[HD1-Stats] Initializing graphs on first expand');
            this.initializeGraphs();
            this.graphsInitialized = true;
            
            // Trigger immediate render to show current data
            this.renderAllGraphs();
        }
    }

    /**
     * Initialize FPS graph
     */
    initializeFpsGraph() {
        this.fpsCanvas = this.dom.get('fps-graph-canvas');
        if (this.fpsCanvas) {
            this.fpsContext = this.fpsCanvas.getContext('2d');
            this.resizeCanvas(this.fpsCanvas);
            console.log('[HD1-Stats] FPS graph initialized');
        }
    }

    /**
     * Initialize memory utilization graph
     */
    initializeMemoryGraph() {
        this.memoryCanvas = this.dom.get('memory-graph-canvas');
        if (this.memoryCanvas) {
            this.memoryContext = this.memoryCanvas.getContext('2d');
            this.resizeCanvas(this.memoryCanvas);
            console.log('[HD1-Stats] Memory graph initialized');
        }
    }

    /**
     * Initialize WebSocket latency graph
     */
    initializeLatencyGraph() {
        this.latencyCanvas = this.dom.get('latency-graph-canvas');
        if (this.latencyCanvas) {
            this.latencyContext = this.latencyCanvas.getContext('2d');
            this.resizeCanvas(this.latencyCanvas);
            console.log('[HD1-Stats] Latency graph initialized');
        }
    }

    /**
     * Initialize WebSocket traffic graph (bidirectional)
     */
    initializeWebSocketGraph() {
        this.wsCanvas = this.dom.get('websocket-graph-canvas');
        if (this.wsCanvas) {
            this.wsContext = this.wsCanvas.getContext('2d');
            this.resizeCanvas(this.wsCanvas);
            console.log('[HD1-Stats] WebSocket graph initialized');
        }
    }

    /**
     * Resize canvas to match its container dimensions
     */
    resizeCanvas(canvas) {
        if (!canvas) return;
        
        // Get the parent container size
        const container = canvas.parentElement;
        const containerRect = container.getBoundingClientRect();
        const displayWidth = Math.floor(containerRect.width - 12); // Account for padding
        const displayHeight = Math.floor(containerRect.height - 43);
        
        // Set canvas internal size to match display size for crisp rendering
        const pixelRatio = window.devicePixelRatio || 1;
        canvas.width = displayWidth * pixelRatio;
        canvas.height = displayHeight * pixelRatio;
        
        // Set CSS size to ensure proper display
        canvas.style.width = displayWidth + 'px';
        canvas.style.height = displayHeight + 'px';
        
        // Scale the context to match device pixel ratio
        const ctx = canvas.getContext('2d');
        ctx.scale(pixelRatio, pixelRatio);
        
        // Update graph config with actual dimensions
        this.graphConfig.width = displayWidth;
        this.graphConfig.height = displayHeight;
        
        console.log(`[HD1-Stats] Canvas resized to ${displayWidth}x${displayHeight} (ratio: ${pixelRatio})`);
    }


    /**
     * Start all monitoring systems
     */
    startMonitoring() {
        // Start frame counting for FPS (high frequency sampling)
        this.startFpsMonitoring();
        
        // Start unified update loop (1 second intervals)
        this.updateInterval = setInterval(() => {
            this.collectMetrics();
            this.renderAllGraphs();
        }, this.graphConfig.updateInterval);
        
        console.log('[HD1-Stats] Monitoring started - smooth real-time 250ms intervals');
    }

    /**
     * Force resize all canvases
     */
    forceResizeAll() {
        console.log('[HD1-Stats] Force resizing all canvases');
        if (this.fpsCanvas) this.resizeCanvas(this.fpsCanvas);
        if (this.memoryCanvas) this.resizeCanvas(this.memoryCanvas);
        if (this.latencyCanvas) this.resizeCanvas(this.latencyCanvas);
        if (this.wsCanvas) this.resizeCanvas(this.wsCanvas);
    }

    /**
     * Start FPS monitoring with requestAnimationFrame
     */
    startFpsMonitoring() {
        const trackFrame = () => {
            const currentTime = performance.now();
            const currentSecond = Math.floor(Date.now() / 1000);
            
            // Count frames
            this.frameCount++;
            
            // If we've moved to a new second, calculate FPS
            if (currentSecond > this.lastSecond) {
                const fps = this.frameCount; // frames in the last second
                this.fpsHistory.push(fps);
                
                // Trim history to last 30 seconds
                if (this.fpsHistory.length > this.graphConfig.historyLength) {
                    this.fpsHistory.shift();
                }
                
                // Reset for next second
                this.frameCount = 0;
                this.lastSecond = currentSecond;
            }
            
            requestAnimationFrame(trackFrame);
        };
        
        requestAnimationFrame(trackFrame);
    }

    /**
     * Collect all metrics (called every second)
     */
    collectMetrics() {
        this.collectMemoryMetrics();
        // FPS is collected in startFpsMonitoring automatically
        // Latency is updated via trackLatency method from WebSocket pongs
        // WebSocket metrics are updated via trackWebSocketTraffic method
    }

    /**
     * Collect memory utilization metrics
     */
    collectMemoryMetrics() {
        try {
            let memoryUsage = 0;
            
            // Try to get actual memory info
            if (performance.memory) {
                // Chrome/Edge performance.memory (actual heap usage)
                const used = performance.memory.usedJSHeapSize;
                const total = performance.memory.totalJSHeapSize;
                const limit = performance.memory.jsHeapSizeLimit;
                
                // Calculate percentage against the limit (not total)
                memoryUsage = (used / limit) * 100;
            } else {
                // No memory API available - set to 0
                memoryUsage = 0;
            }
            
            // Add to history
            this.memoryHistory.push(Math.min(memoryUsage, 100));
            
            // Trim history to last 30 seconds
            if (this.memoryHistory.length > this.graphConfig.historyLength) {
                this.memoryHistory.shift();
            }
        } catch (error) {
            console.error('[HD1-Stats] Memory collection failed:', error);
            this.memoryHistory.push(0);
            
            // Trim history to last 30 seconds
            if (this.memoryHistory.length > this.graphConfig.historyLength) {
                this.memoryHistory.shift();
            }
        }
    }

    /**
     * Track WebSocket latency (called from WebSocket manager)
     */
    trackLatency(latencyMs) {
        this.latencyHistory.push(latencyMs);
        
        // Trim history to last 30 seconds
        if (this.latencyHistory.length > this.graphConfig.historyLength) {
            this.latencyHistory.shift();
        }
        
        console.log(`[HD1-Stats] Latency tracked: ${latencyMs}ms (history: ${this.latencyHistory.length})`);
    }



    /**
     * Track WebSocket traffic (called from WebSocket manager)
     */
    trackWebSocketTraffic(inBytes, outBytes) {
        const now = Date.now();
        
        // Calculate bytes per second
        const timeDelta = this.wsLastSample.time ? (now - this.wsLastSample.time) / 1000 : 1;
        const inRate = (inBytes - this.wsLastSample.in) / timeDelta;
        const outRate = (outBytes - this.wsLastSample.out) / timeDelta;
        
        // Add to history
        this.wsInHistory.push(Math.max(inRate, 0));
        this.wsOutHistory.push(Math.max(outRate, 0));
        
        // Trim history
        if (this.wsInHistory.length > this.graphConfig.historyLength) {
            this.wsInHistory.shift();
            this.wsOutHistory.shift();
        }
        
        // Update last sample
        this.wsLastSample = { in: inBytes, out: outBytes, time: now };
        
        // Update totals
        this.wsInBytes = inBytes;
        this.wsOutBytes = outBytes;
    }

    /**
     * Render all graphs
     */
    renderAllGraphs() {
        this.renderFpsGraph();
        this.renderMemoryGraph();
        this.renderLatencyGraph();
        this.renderWebSocketGraph();
    }

    /**
     * Render FPS graph
     */
    renderFpsGraph() {
        if (!this.fpsContext || !this.fpsCanvas) return;
        
        const ctx = this.fpsContext;
        const width = this.graphConfig.width;
        const height = this.graphConfig.height;
        
        // Clear canvas
        ctx.clearRect(0, 0, width, height);
        
        // Background
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.fillRect(0, 0, width, height);
        
        // Grid lines (30, 60, 120 FPS)
        ctx.strokeStyle = 'rgba(0, 255, 255, 0.1)';
        ctx.lineWidth = 1;
        [30, 60, 120].forEach(fps => {
            const y = height - (fps / 120) * height;
            if (y >= 0 && y <= height) {
                ctx.beginPath();
                ctx.moveTo(0, y);
                ctx.lineTo(width, y);
                ctx.stroke();
            }
        });
        
        // FPS line (right-to-left scrolling)
        if (this.fpsHistory.length > 0) {
            ctx.strokeStyle = '#00ffff'; // HD1 theme cyan
            ctx.lineWidth = 2;
            ctx.beginPath();
            
            this.fpsHistory.forEach((fps, index) => {
                // Right-to-left: newest data on the right
                const x = width - ((this.fpsHistory.length - 1 - index) / this.graphConfig.historyLength) * width;
                const y = height - Math.min(fps / 120, 1) * height;
                
                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });
            
            ctx.stroke();
            
            // Current FPS text
            const currentFps = this.fpsHistory[this.fpsHistory.length - 1];
            ctx.fillStyle = '#00ffff';
            ctx.font = '9px monospace';
            ctx.fillText(`${currentFps}fps`, 5, 12);
            
            // Update collapsed value
            this.dom.setText('collapsed-fps-value', `${currentFps}`);
        }
    }

    /**
     * Render memory utilization graph
     */
    renderMemoryGraph() {
        if (!this.memoryContext || !this.memoryCanvas) return;
        
        const ctx = this.memoryContext;
        const width = this.graphConfig.width;
        const height = this.graphConfig.height;
        
        // Clear canvas
        ctx.clearRect(0, 0, width, height);
        
        // Background
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.fillRect(0, 0, width, height);
        
        // Grid lines (25%, 50%, 75%, 100%)
        ctx.strokeStyle = 'rgba(255, 255, 0, 0.1)';
        ctx.lineWidth = 1;
        [25, 50, 75, 100].forEach(percent => {
            const y = height - (percent / 100) * height;
            ctx.beginPath();
            ctx.moveTo(0, y);
            ctx.lineTo(width, y);
            ctx.stroke();
        });
        
        // Memory usage line (right-to-left scrolling)
        if (this.memoryHistory.length > 0) {
            ctx.strokeStyle = '#00ffff'; // HD1 theme cyan
            ctx.lineWidth = 2;
            ctx.beginPath();
            
            this.memoryHistory.forEach((usage, index) => {
                // Right-to-left: newest data on the right
                const x = width - ((this.memoryHistory.length - 1 - index) / this.graphConfig.historyLength) * width;
                const y = height - (usage / 100) * height;
                
                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });
            
            ctx.stroke();
            
            // Current usage text
            const currentUsage = this.memoryHistory[this.memoryHistory.length - 1];
            ctx.fillStyle = '#00ffff';
            ctx.font = '9px monospace';
            ctx.fillText(`${Math.round(currentUsage)}%`, 5, 12);
            
            // Update collapsed value
            this.dom.setText('collapsed-memory-value', `${Math.round(currentUsage)}%`);
        }
    }

    /**
     * Render WebSocket latency graph
     */
    renderLatencyGraph() {
        if (!this.latencyContext || !this.latencyCanvas) return;
        
        const ctx = this.latencyContext;
        const width = this.graphConfig.width;
        const height = this.graphConfig.height;
        
        // Clear canvas
        ctx.clearRect(0, 0, width, height);
        
        // Background
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.fillRect(0, 0, width, height);
        
        // Grid lines (50ms, 100ms, 200ms, 500ms)
        ctx.strokeStyle = 'rgba(255, 165, 0, 0.1)';
        ctx.lineWidth = 1;
        [50, 100, 200, 500].forEach(latency => {
            const y = height - (latency / 500) * height;
            if (y >= 0 && y <= height) {
                ctx.beginPath();
                ctx.moveTo(0, y);
                ctx.lineTo(width, y);
                ctx.stroke();
            }
        });
        
        // Latency line (right-to-left scrolling)
        if (this.latencyHistory.length > 0) {
            ctx.strokeStyle = '#00ffff'; // HD1 theme cyan
            ctx.lineWidth = 2;
            ctx.beginPath();
            
            this.latencyHistory.forEach((latency, index) => {
                // Right-to-left: newest data on the right
                const x = width - ((this.latencyHistory.length - 1 - index) / this.graphConfig.historyLength) * width;
                const y = height - Math.min(latency / 500, 1) * height; // Scale to 500ms max
                
                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });
            
            ctx.stroke();
            
            // Current latency text
            const currentLatency = this.latencyHistory[this.latencyHistory.length - 1];
            ctx.fillStyle = '#00ffff';
            ctx.font = '9px monospace';
            ctx.fillText(`${Math.round(currentLatency)}ms`, 5, 12);
            
            // Update collapsed value
            this.dom.setText('collapsed-latency-value', `${Math.round(currentLatency)}ms`);
        }
    }

    /**
     * Render WebSocket traffic graph (bidirectional with zero center)
     */
    renderWebSocketGraph() {
        if (!this.wsContext || !this.wsCanvas) return;
        
        const ctx = this.wsContext;
        const width = this.graphConfig.width;
        const height = this.graphConfig.height;
        const centerY = height / 2;
        
        // Clear canvas
        ctx.clearRect(0, 0, width, height);
        
        // Background
        ctx.fillStyle = 'rgba(0, 0, 0, 0.2)';
        ctx.fillRect(0, 0, width, height);
        
        // Center line (zero point)
        ctx.strokeStyle = 'rgba(0, 255, 255, 0.3)';
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(0, centerY);
        ctx.lineTo(width, centerY);
        ctx.stroke();
        
        // Get max value for scaling (use actual max or minimum of 1 for visibility)
        const maxIn = this.wsInHistory.length > 0 ? Math.max(...this.wsInHistory) : 1;
        const maxOut = this.wsOutHistory.length > 0 ? Math.max(...this.wsOutHistory) : 1;
        const maxValue = Math.max(maxIn, maxOut, 1); // Ensure minimum of 1 for scaling
        
        // Incoming data (above center) - right-to-left scrolling
        if (this.wsInHistory.length > 0) {
            ctx.strokeStyle = '#00ffff'; // HD1 theme cyan
            ctx.lineWidth = 2;
            ctx.beginPath();
            
            this.wsInHistory.forEach((rate, index) => {
                // Right-to-left: newest data on the right
                const x = width - ((this.wsInHistory.length - 1 - index) / this.graphConfig.historyLength) * width;
                const y = centerY - (rate / maxValue) * (centerY - 5);
                
                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });
            
            ctx.stroke();
        }
        
        // Outgoing data (below center) - right-to-left scrolling
        if (this.wsOutHistory.length > 0) {
            ctx.strokeStyle = '#00ffff'; // HD1 theme cyan (matching theme)
            ctx.lineWidth = 1;
            ctx.setLineDash([2, 2]); // Dashed line to differentiate
            ctx.beginPath();
            
            this.wsOutHistory.forEach((rate, index) => {
                // Right-to-left: newest data on the right
                const x = width - ((this.wsOutHistory.length - 1 - index) / this.graphConfig.historyLength) * width;
                const y = centerY + (rate / maxValue) * (centerY - 5);
                
                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });
            
            ctx.stroke();
            ctx.setLineDash([]); // Reset line dash
        }
        
        // Labels
        ctx.fillStyle = '#00ffff';
        ctx.font = '8px monospace';
        ctx.fillText('IN', 5, 10);
        ctx.fillText('OUT', 5, height - 3);
        
        // Current rates
        if (this.wsInHistory.length > 0) {
            const currentIn = this.wsInHistory[this.wsInHistory.length - 1];
            const currentOut = this.wsOutHistory[this.wsOutHistory.length - 1];
            
            ctx.fillStyle = '#00ffff';
            ctx.font = '8px monospace';
            ctx.fillText(`${Math.round(currentIn)}`, width - 30, 10);
            ctx.fillText(`${Math.round(currentOut)}`, width - 30, height - 3);
            
            // Update collapsed values (separate RX and TX)
            this.dom.setText('collapsed-rx-value', `${Math.round(currentIn)}`);
            this.dom.setText('collapsed-tx-value', `${Math.round(currentOut)}`);
        }
    }

    /**
     * Get current FPS
     */
    getCurrentFPS() {
        return this.fpsHistory.length > 0 ? this.fpsHistory[this.fpsHistory.length - 1] : 0;
    }

    /**
     * Get current memory usage
     */
    getCurrentMemoryUsage() {
        return this.memoryHistory.length > 0 ? this.memoryHistory[this.memoryHistory.length - 1] : 0;
    }

    /**
     * Get current latency
     */
    getCurrentLatency() {
        return this.latencyHistory.length > 0 ? this.latencyHistory[this.latencyHistory.length - 1] : 0;
    }

    /**
     * Get WebSocket traffic stats
     */
    getWebSocketStats() {
        return {
            totalIn: this.wsInBytes,
            totalOut: this.wsOutBytes,
            currentInRate: this.wsInHistory.length > 0 ? this.wsInHistory[this.wsInHistory.length - 1] : 0,
            currentOutRate: this.wsOutHistory.length > 0 ? this.wsOutHistory[this.wsOutHistory.length - 1] : 0
        };
    }

    /**
     * Cleanup stats manager
     */
    cleanup() {
        if (this.updateInterval) {
            clearInterval(this.updateInterval);
        }
        
        // Clear all histories
        this.fpsHistory = [];
        this.memoryHistory = [];
        this.latencyHistory = [];
        this.wsInHistory = [];
        this.wsOutHistory = [];
        
        this.ready = false;
    }
}

// Export for use in console manager
window.HD1StatsManager = HD1StatsManager;