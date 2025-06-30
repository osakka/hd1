// ===================================================================
// AUTO-GENERATED UI COMPONENTS - DO NOT MODIFY
// ===================================================================
//
// Standard UI components auto-generated from OpenAPI specification
// Updates automatically when API specification changes
//
// ===================================================================

class THDUIComponents {
    constructor(apiClient) {
        this.api = apiClient;
        this.components = new Map();
        this.initializeComponents();
    }

    initializeComponents() {
        console.log('🎨 Initializing auto-generated UI components...');

        this.components.set('initializeworld', this.createInitializeWorldComponent());

        this.components.set('getworldspec', this.createGetWorldSpecComponent());

        this.components.set('listscenes', this.createListScenesComponent());

        this.components.set('startrecording', this.createStartRecordingComponent());

        this.components.set('getobject', this.createGetObjectComponent());

        this.components.set('updateobject', this.createUpdateObjectComponent());

        this.components.set('deleteobject', this.createDeleteObjectComponent());

        this.components.set('getloggingconfig', this.createGetLoggingConfigComponent());

        this.components.set('setloggingconfig', this.createSetLoggingConfigComponent());

        this.components.set('forcerefresh', this.createForceRefreshComponent());

        this.components.set('loadscene', this.createLoadSceneComponent());

        this.components.set('playrecording', this.createPlayRecordingComponent());

        this.components.set('setloglevel', this.createSetLogLevelComponent());

        this.components.set('settracemodules', this.createSetTraceModulesComponent());

        this.components.set('listsessions', this.createListSessionsComponent());

        this.components.set('createsession', this.createCreateSessionComponent());

        this.components.set('stoprecording', this.createStopRecordingComponent());

        this.components.set('listobjects', this.createListObjectsComponent());

        this.components.set('createobject', this.createCreateObjectComponent());

        this.components.set('getlogs', this.createGetLogsComponent());

        this.components.set('startcameraorbit', this.createStartCameraOrbitComponent());

        this.components.set('deletesession', this.createDeleteSessionComponent());

        this.components.set('getsession', this.createGetSessionComponent());

        this.components.set('savescenefromsession', this.createSaveSceneFromSessionComponent());

        this.components.set('forkscene', this.createForkSceneComponent());

        this.components.set('getrecordingstatus', this.createGetRecordingStatusComponent());

        this.components.set('setcanvas', this.createSetCanvasComponent());

        this.components.set('setcameraposition', this.createSetCameraPositionComponent());

        console.log('✅ UI components initialized');
    }


    // Component for POST /sessions/{sessionId}/world
    createInitializeWorldComponent() {
        return {
            name: 'initializeworld',
            endpoint: '/sessions/{sessionId}/world',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>InitializeWorld</h4><form id="initializeWorld-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="initializeWorld-result" class="result-area"></div></div>';
                this.attachEventListeners('initializeworld', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.initializeWorld(formData.param1, formData);
                    this.showResult('initializeworld', result);
                    return result;
                } catch (error) {
                    this.showError('initializeworld', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/world
    createGetWorldSpecComponent() {
        return {
            name: 'getworldspec',
            endpoint: '/sessions/{sessionId}/world',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>GetWorldSpec</h4><form id="getWorldSpec-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getWorldSpec-result" class="result-area"></div></div>';
                this.attachEventListeners('getworldspec', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getWorldSpec(formData.param1);
                    this.showResult('getworldspec', result);
                    return result;
                } catch (error) {
                    this.showError('getworldspec', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /scenes
    createListScenesComponent() {
        return {
            name: 'listscenes',
            endpoint: '/scenes',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>ListScenes</h4><form id="listScenes-form"><button type="submit">Execute</button></form><div id="listScenes-result" class="result-area"></div></div>';
                this.attachEventListeners('listscenes', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listScenes();
                    this.showResult('listscenes', result);
                    return result;
                } catch (error) {
                    this.showError('listscenes', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/recording/start
    createStartRecordingComponent() {
        return {
            name: 'startrecording',
            endpoint: '/sessions/{sessionId}/recording/start',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>StartRecording</h4><form id="startRecording-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="startRecording-result" class="result-area"></div></div>';
                this.attachEventListeners('startrecording', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.startRecording(formData.param1, formData);
                    this.showResult('startrecording', result);
                    return result;
                } catch (error) {
                    this.showError('startrecording', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/objects/{objectName}
    createGetObjectComponent() {
        return {
            name: 'getobject',
            endpoint: '/sessions/{sessionId}/objects/{objectName}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>GetObject</h4><form id="getObject-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getObject-result" class="result-area"></div></div>';
                this.attachEventListeners('getobject', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getObject(formData.param1, formData.param2);
                    this.showResult('getobject', result);
                    return result;
                } catch (error) {
                    this.showError('getobject', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/objects/{objectName}
    createUpdateObjectComponent() {
        return {
            name: 'updateobject',
            endpoint: '/sessions/{sessionId}/objects/{objectName}',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>UpdateObject</h4><form id="updateObject-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateObject-result" class="result-area"></div></div>';
                this.attachEventListeners('updateobject', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateObject(formData.param1, formData.param2, formData);
                    this.showResult('updateobject', result);
                    return result;
                } catch (error) {
                    this.showError('updateobject', error);
                    throw error;
                }
            }
        };
    }

    // Component for DELETE /sessions/{sessionId}/objects/{objectName}
    createDeleteObjectComponent() {
        return {
            name: 'deleteobject',
            endpoint: '/sessions/{sessionId}/objects/{objectName}',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>DeleteObject</h4><form id="deleteObject-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="deleteObject-result" class="result-area"></div></div>';
                this.attachEventListeners('deleteobject', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.deleteObject(formData.param1, formData.param2);
                    this.showResult('deleteobject', result);
                    return result;
                } catch (error) {
                    this.showError('deleteobject', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /admin/logging/config
    createGetLoggingConfigComponent() {
        return {
            name: 'getloggingconfig',
            endpoint: '/admin/logging/config',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>GetLoggingConfig</h4><form id="getLoggingConfig-form"><button type="submit">Execute</button></form><div id="getLoggingConfig-result" class="result-area"></div></div>';
                this.attachEventListeners('getloggingconfig', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getLoggingConfig();
                    this.showResult('getloggingconfig', result);
                    return result;
                } catch (error) {
                    this.showError('getloggingconfig', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /admin/logging/config
    createSetLoggingConfigComponent() {
        return {
            name: 'setloggingconfig',
            endpoint: '/admin/logging/config',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>SetLoggingConfig</h4><form id="setLoggingConfig-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setLoggingConfig-result" class="result-area"></div></div>';
                this.attachEventListeners('setloggingconfig', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setLoggingConfig(formData);
                    this.showResult('setloggingconfig', result);
                    return result;
                } catch (error) {
                    this.showError('setloggingconfig', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /browser/refresh
    createForceRefreshComponent() {
        return {
            name: 'forcerefresh',
            endpoint: '/browser/refresh',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>ForceRefresh</h4><form id="forceRefresh-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="forceRefresh-result" class="result-area"></div></div>';
                this.attachEventListeners('forcerefresh', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.forceRefresh(formData);
                    this.showResult('forcerefresh', result);
                    return result;
                } catch (error) {
                    this.showError('forcerefresh', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /scenes/{sceneId}
    createLoadSceneComponent() {
        return {
            name: 'loadscene',
            endpoint: '/scenes/{sceneId}',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>LoadScene</h4><form id="loadScene-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="loadScene-result" class="result-area"></div></div>';
                this.attachEventListeners('loadscene', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.loadScene(formData.param1, formData);
                    this.showResult('loadscene', result);
                    return result;
                } catch (error) {
                    this.showError('loadscene', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/recording/play
    createPlayRecordingComponent() {
        return {
            name: 'playrecording',
            endpoint: '/sessions/{sessionId}/recording/play',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>PlayRecording</h4><form id="playRecording-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="playRecording-result" class="result-area"></div></div>';
                this.attachEventListeners('playrecording', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.playRecording(formData.param1, formData);
                    this.showResult('playrecording', result);
                    return result;
                } catch (error) {
                    this.showError('playrecording', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /admin/logging/level
    createSetLogLevelComponent() {
        return {
            name: 'setloglevel',
            endpoint: '/admin/logging/level',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>SetLogLevel</h4><form id="setLogLevel-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setLogLevel-result" class="result-area"></div></div>';
                this.attachEventListeners('setloglevel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setLogLevel(formData);
                    this.showResult('setloglevel', result);
                    return result;
                } catch (error) {
                    this.showError('setloglevel', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /admin/logging/trace
    createSetTraceModulesComponent() {
        return {
            name: 'settracemodules',
            endpoint: '/admin/logging/trace',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>SetTraceModules</h4><form id="setTraceModules-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setTraceModules-result" class="result-area"></div></div>';
                this.attachEventListeners('settracemodules', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setTraceModules(formData);
                    this.showResult('settracemodules', result);
                    return result;
                } catch (error) {
                    this.showError('settracemodules', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions
    createListSessionsComponent() {
        return {
            name: 'listsessions',
            endpoint: '/sessions',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>ListSessions</h4><form id="listSessions-form"><button type="submit">Execute</button></form><div id="listSessions-result" class="result-area"></div></div>';
                this.attachEventListeners('listsessions', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listSessions();
                    this.showResult('listsessions', result);
                    return result;
                } catch (error) {
                    this.showError('listsessions', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions
    createCreateSessionComponent() {
        return {
            name: 'createsession',
            endpoint: '/sessions',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>CreateSession</h4><form id="createSession-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createSession-result" class="result-area"></div></div>';
                this.attachEventListeners('createsession', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createSession(formData);
                    this.showResult('createsession', result);
                    return result;
                } catch (error) {
                    this.showError('createsession', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/recording/stop
    createStopRecordingComponent() {
        return {
            name: 'stoprecording',
            endpoint: '/sessions/{sessionId}/recording/stop',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>StopRecording</h4><form id="stopRecording-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="stopRecording-result" class="result-area"></div></div>';
                this.attachEventListeners('stoprecording', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.stopRecording(formData.param1, formData);
                    this.showResult('stoprecording', result);
                    return result;
                } catch (error) {
                    this.showError('stoprecording', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/objects
    createListObjectsComponent() {
        return {
            name: 'listobjects',
            endpoint: '/sessions/{sessionId}/objects',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>ListObjects</h4><form id="listObjects-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="listObjects-result" class="result-area"></div></div>';
                this.attachEventListeners('listobjects', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listObjects(formData.param1);
                    this.showResult('listobjects', result);
                    return result;
                } catch (error) {
                    this.showError('listobjects', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/objects
    createCreateObjectComponent() {
        return {
            name: 'createobject',
            endpoint: '/sessions/{sessionId}/objects',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>CreateObject</h4><form id="createObject-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createObject-result" class="result-area"></div></div>';
                this.attachEventListeners('createobject', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createObject(formData.param1, formData);
                    this.showResult('createobject', result);
                    return result;
                } catch (error) {
                    this.showError('createobject', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /admin/logging/logs
    createGetLogsComponent() {
        return {
            name: 'getlogs',
            endpoint: '/admin/logging/logs',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>GetLogs</h4><form id="getLogs-form"><button type="submit">Execute</button></form><div id="getLogs-result" class="result-area"></div></div>';
                this.attachEventListeners('getlogs', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getLogs();
                    this.showResult('getlogs', result);
                    return result;
                } catch (error) {
                    this.showError('getlogs', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/camera/orbit
    createStartCameraOrbitComponent() {
        return {
            name: 'startcameraorbit',
            endpoint: '/sessions/{sessionId}/camera/orbit',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>StartCameraOrbit</h4><form id="startCameraOrbit-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="startCameraOrbit-result" class="result-area"></div></div>';
                this.attachEventListeners('startcameraorbit', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.startCameraOrbit(formData.param1, formData);
                    this.showResult('startcameraorbit', result);
                    return result;
                } catch (error) {
                    this.showError('startcameraorbit', error);
                    throw error;
                }
            }
        };
    }

    // Component for DELETE /sessions/{sessionId}
    createDeleteSessionComponent() {
        return {
            name: 'deletesession',
            endpoint: '/sessions/{sessionId}',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>DeleteSession</h4><form id="deleteSession-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="deleteSession-result" class="result-area"></div></div>';
                this.attachEventListeners('deletesession', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.deleteSession(formData.param1);
                    this.showResult('deletesession', result);
                    return result;
                } catch (error) {
                    this.showError('deletesession', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}
    createGetSessionComponent() {
        return {
            name: 'getsession',
            endpoint: '/sessions/{sessionId}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>GetSession</h4><form id="getSession-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSession-result" class="result-area"></div></div>';
                this.attachEventListeners('getsession', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSession(formData.param1);
                    this.showResult('getsession', result);
                    return result;
                } catch (error) {
                    this.showError('getsession', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/scenes/save
    createSaveSceneFromSessionComponent() {
        return {
            name: 'savescenefromsession',
            endpoint: '/sessions/{sessionId}/scenes/save',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>SaveSceneFromSession</h4><form id="saveSceneFromSession-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="saveSceneFromSession-result" class="result-area"></div></div>';
                this.attachEventListeners('savescenefromsession', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.saveSceneFromSession(formData.param1, formData);
                    this.showResult('savescenefromsession', result);
                    return result;
                } catch (error) {
                    this.showError('savescenefromsession', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /scenes/{sceneId}/fork
    createForkSceneComponent() {
        return {
            name: 'forkscene',
            endpoint: '/scenes/{sceneId}/fork',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>ForkScene</h4><form id="forkScene-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="forkScene-result" class="result-area"></div></div>';
                this.attachEventListeners('forkscene', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.forkScene(formData.param1, formData);
                    this.showResult('forkscene', result);
                    return result;
                } catch (error) {
                    this.showError('forkscene', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/recording/status
    createGetRecordingStatusComponent() {
        return {
            name: 'getrecordingstatus',
            endpoint: '/sessions/{sessionId}/recording/status',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>GetRecordingStatus</h4><form id="getRecordingStatus-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getRecordingStatus-result" class="result-area"></div></div>';
                this.attachEventListeners('getrecordingstatus', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getRecordingStatus(formData.param1);
                    this.showResult('getrecordingstatus', result);
                    return result;
                } catch (error) {
                    this.showError('getrecordingstatus', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /browser/canvas
    createSetCanvasComponent() {
        return {
            name: 'setcanvas',
            endpoint: '/browser/canvas',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>SetCanvas</h4><form id="setCanvas-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setCanvas-result" class="result-area"></div></div>';
                this.attachEventListeners('setcanvas', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setCanvas(formData);
                    this.showResult('setcanvas', result);
                    return result;
                } catch (error) {
                    this.showError('setcanvas', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/camera/position
    createSetCameraPositionComponent() {
        return {
            name: 'setcameraposition',
            endpoint: '/sessions/{sessionId}/camera/position',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="thd-component"><h4>SetCameraPosition</h4><form id="setCameraPosition-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setCameraPosition-result" class="result-area"></div></div>';
                this.attachEventListeners('setcameraposition', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setCameraPosition(formData.param1, formData);
                    this.showResult('setcameraposition', result);
                    return result;
                } catch (error) {
                    this.showError('setcameraposition', error);
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
        
        container.innerHTML = '<h2>THD API Interface</h2>';
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
window.THDUIComponents = THDUIComponents;

console.log('🎨 THD UI Components loaded - Auto-generated from specification');
