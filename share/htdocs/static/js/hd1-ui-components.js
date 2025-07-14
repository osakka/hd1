// ===================================================================
// AUTO-GENERATED UI COMPONENTS - DO NOT MODIFY
// ===================================================================
//
// Standard UI components auto-generated from OpenAPI specification
// Updates automatically when API specification changes
//
// ===================================================================

class HD1UIComponents {
    constructor(apiClient) {
        this.api = apiClient;
        this.components = new Map();
        this.initializeComponents();
    }

    initializeComponents() {
        console.log('🎨 Initializing auto-generated UI components...');

        this.components.set('updateentity', this.createUpdateEntityComponent());

        this.components.set('deleteentity', this.createDeleteEntityComponent());

        this.components.set('moveavatar', this.createMoveAvatarComponent());

        this.components.set('updatescene', this.createUpdateSceneComponent());

        this.components.set('getscene', this.createGetSceneComponent());

        this.components.set('getversion', this.createGetVersionComponent());

        this.components.set('submitoperation', this.createSubmitOperationComponent());

        this.components.set('getmissingoperations', this.createGetMissingOperationsComponent());

        this.components.set('getfullsync', this.createGetFullSyncComponent());

        this.components.set('getsyncstats', this.createGetSyncStatsComponent());

        this.components.set('createentity', this.createCreateEntityComponent());

        console.log('✅ UI components initialized');
    }


    // Component for PUT /threejs/entities/{entityId}
    createUpdateEntityComponent() {
        return {
            name: 'updateentity',
            endpoint: '/threejs/entities/{entityId}',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateEntity</h4><form id="updateEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('updateentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateEntity(formData.param1, formData);
                    this.showResult('updateentity', result);
                    return result;
                } catch (error) {
                    this.showError('updateentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for DELETE /threejs/entities/{entityId}
    createDeleteEntityComponent() {
        return {
            name: 'deleteentity',
            endpoint: '/threejs/entities/{entityId}',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>DeleteEntity</h4><form id="deleteEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="deleteEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('deleteentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.deleteEntity(formData.param1);
                    this.showResult('deleteentity', result);
                    return result;
                } catch (error) {
                    this.showError('deleteentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /threejs/avatars/{sessionId}/move
    createMoveAvatarComponent() {
        return {
            name: 'moveavatar',
            endpoint: '/threejs/avatars/{sessionId}/move',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>MoveAvatar</h4><form id="moveAvatar-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="moveAvatar-result" class="result-area"></div></div>';
                this.attachEventListeners('moveavatar', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.moveAvatar(formData.param1, formData);
                    this.showResult('moveavatar', result);
                    return result;
                } catch (error) {
                    this.showError('moveavatar', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /threejs/scene
    createUpdateSceneComponent() {
        return {
            name: 'updatescene',
            endpoint: '/threejs/scene',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateScene</h4><form id="updateScene-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateScene-result" class="result-area"></div></div>';
                this.attachEventListeners('updatescene', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateScene(formData);
                    this.showResult('updatescene', result);
                    return result;
                } catch (error) {
                    this.showError('updatescene', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /threejs/scene
    createGetSceneComponent() {
        return {
            name: 'getscene',
            endpoint: '/threejs/scene',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetScene</h4><form id="getScene-form"><button type="submit">Execute</button></form><div id="getScene-result" class="result-area"></div></div>';
                this.attachEventListeners('getscene', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getScene();
                    this.showResult('getscene', result);
                    return result;
                } catch (error) {
                    this.showError('getscene', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /system/version
    createGetVersionComponent() {
        return {
            name: 'getversion',
            endpoint: '/system/version',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetVersion</h4><form id="getVersion-form"><button type="submit">Execute</button></form><div id="getVersion-result" class="result-area"></div></div>';
                this.attachEventListeners('getversion', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getVersion();
                    this.showResult('getversion', result);
                    return result;
                } catch (error) {
                    this.showError('getversion', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sync/operations
    createSubmitOperationComponent() {
        return {
            name: 'submitoperation',
            endpoint: '/sync/operations',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>SubmitOperation</h4><form id="submitOperation-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="submitOperation-result" class="result-area"></div></div>';
                this.attachEventListeners('submitoperation', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.submitOperation(formData);
                    this.showResult('submitoperation', result);
                    return result;
                } catch (error) {
                    this.showError('submitoperation', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sync/missing/{from}/{to}
    createGetMissingOperationsComponent() {
        return {
            name: 'getmissingoperations',
            endpoint: '/sync/missing/{from}/{to}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetMissingOperations</h4><form id="getMissingOperations-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getMissingOperations-result" class="result-area"></div></div>';
                this.attachEventListeners('getmissingoperations', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getMissingOperations(formData.param1, formData.param2);
                    this.showResult('getmissingoperations', result);
                    return result;
                } catch (error) {
                    this.showError('getmissingoperations', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sync/full
    createGetFullSyncComponent() {
        return {
            name: 'getfullsync',
            endpoint: '/sync/full',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetFullSync</h4><form id="getFullSync-form"><button type="submit">Execute</button></form><div id="getFullSync-result" class="result-area"></div></div>';
                this.attachEventListeners('getfullsync', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getFullSync();
                    this.showResult('getfullsync', result);
                    return result;
                } catch (error) {
                    this.showError('getfullsync', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sync/stats
    createGetSyncStatsComponent() {
        return {
            name: 'getsyncstats',
            endpoint: '/sync/stats',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSyncStats</h4><form id="getSyncStats-form"><button type="submit">Execute</button></form><div id="getSyncStats-result" class="result-area"></div></div>';
                this.attachEventListeners('getsyncstats', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSyncStats();
                    this.showResult('getsyncstats', result);
                    return result;
                } catch (error) {
                    this.showError('getsyncstats', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /threejs/entities
    createCreateEntityComponent() {
        return {
            name: 'createentity',
            endpoint: '/threejs/entities',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateEntity</h4><form id="createEntity-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('createentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createEntity(formData);
                    this.showResult('createentity', result);
                    return result;
                } catch (error) {
                    this.showError('createentity', error);
                    throw error;
                }
            }
        };
    }


    attachEventListeners(componentName, container) {
        const component = this.components.get(componentName);
        if (!component) return;
        
        const form = container.querySelector('form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const formData = new FormData(form);
                const data = Object.fromEntries(formData.entries());
                
                try {
                    await component.execute(data);
                } catch (error) {
                    console.error('Component execution failed:', error);
                }
            });
        }
    }

    showResult(componentName, result) {
        const resultDiv = document.getElementById(componentName + '-result');
        if (resultDiv) {
            resultDiv.innerHTML = '<pre class="success">' + JSON.stringify(result, null, 2) + '</pre>';
        }
    }

    showError(componentName, error) {
        const resultDiv = document.getElementById(componentName + '-result');
        if (resultDiv) {
            resultDiv.innerHTML = '<pre class="error">Error: ' + error.message + '</pre>';
        }
    }

    getComponent(name) {
        return this.components.get(name);
    }

    renderAll(containerId) {
        const container = document.getElementById(containerId);
        if (!container) {
            console.error('Container not found:', containerId);
            return;
        }
        
        container.innerHTML = '<h2>HD1 API Interface</h2>';
        this.components.forEach((component, name) => {
            const div = document.createElement('div');
            div.className = 'component-container';
            div.id = name + '-container';
            container.appendChild(div);
            component.render(name + '-container');
        });
    }
}

// Export for global use
window.HD1UIComponents = HD1UIComponents;

console.log('🎨 HD1 UI Components loaded - Auto-generated from specification');