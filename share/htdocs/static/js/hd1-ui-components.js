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

        this.components.set('listaudiosources', this.createListAudioSourcesComponent());

        this.components.set('createaudiosource', this.createCreateAudioSourceComponent());

        this.components.set('forcerefresh', this.createForceRefreshComponent());

        this.components.set('getsessionavatar', this.createGetSessionAvatarComponent());

        this.components.set('setsessionavatar', this.createSetSessionAvatarComponent());

        this.components.set('getentitychildren', this.createGetEntityChildrenComponent());

        this.components.set('playaudio', this.createPlayAudioComponent());

        this.components.set('playrecording', this.createPlayRecordingComponent());

        this.components.set('getrecordingstatus', this.createGetRecordingStatusComponent());

        this.components.set('listchannels', this.createListChannelsComponent());

        this.components.set('createchannel', this.createCreateChannelComponent());

        this.components.set('listsessions', this.createListSessionsComponent());

        this.components.set('createsession', this.createCreateSessionComponent());

        this.components.set('startrecording', this.createStartRecordingComponent());

        this.components.set('setloglevel', this.createSetLogLevelComponent());

        this.components.set('deletesession', this.createDeleteSessionComponent());

        this.components.set('getsession', this.createGetSessionComponent());

        this.components.set('syncsessionstate', this.createSyncSessionStateComponent());

        this.components.set('gethierarchytree', this.createGetHierarchyTreeComponent());

        this.components.set('deactivateentity', this.createDeactivateEntityComponent());

        this.components.set('activatesessionscene', this.createActivateSessionSceneComponent());

        this.components.set('listanimations', this.createListAnimationsComponent());

        this.components.set('createanimation', this.createCreateAnimationComponent());

        this.components.set('setloggingconfig', this.createSetLoggingConfigComponent());

        this.components.set('getloggingconfig', this.createGetLoggingConfigComponent());

        this.components.set('getsessionchannelstatus', this.createGetSessionChannelStatusComponent());

        this.components.set('getchannel', this.createGetChannelComponent());

        this.components.set('updatechannel', this.createUpdateChannelComponent());

        this.components.set('deletechannel', this.createDeleteChannelComponent());

        this.components.set('listavatars', this.createListAvatarsComponent());

        this.components.set('getavatarasset', this.createGetAvatarAssetComponent());

        this.components.set('listentities', this.createListEntitiesComponent());

        this.components.set('createentity', this.createCreateEntityComponent());

        this.components.set('getcameraposition', this.createGetCameraPositionComponent());

        this.components.set('setcameraposition', this.createSetCameraPositionComponent());

        this.components.set('disableentity', this.createDisableEntityComponent());

        this.components.set('resetscenestate', this.createResetSceneStateComponent());

        this.components.set('getentitylifecyclestatus', this.createGetEntityLifecycleStatusComponent());

        this.components.set('getscenestate', this.createGetSceneStateComponent());

        this.components.set('updatescenestate', this.createUpdateSceneStateComponent());

        this.components.set('savescenestate', this.createSaveSceneStateComponent());

        this.components.set('applyforce', this.createApplyForceComponent());

        this.components.set('getversion', this.createGetVersionComponent());

        this.components.set('destroyentity', this.createDestroyEntityComponent());

        this.components.set('settracemodules', this.createSetTraceModulesComponent());

        this.components.set('setcanvas', this.createSetCanvasComponent());

        this.components.set('getavatarspecification', this.createGetAvatarSpecificationComponent());

        this.components.set('getcomponent', this.createGetComponentComponent());

        this.components.set('updatecomponent', this.createUpdateComponentComponent());

        this.components.set('removecomponent', this.createRemoveComponentComponent());

        this.components.set('exportscenedefinition', this.createExportSceneDefinitionComponent());

        this.components.set('startcameraorbit', this.createStartCameraOrbitComponent());

        this.components.set('bulkcomponentoperation', this.createBulkComponentOperationComponent());

        this.components.set('getentityparent', this.createGetEntityParentComponent());

        this.components.set('setentityparent', this.createSetEntityParentComponent());

        this.components.set('getscenehierarchy', this.createGetSceneHierarchyComponent());

        this.components.set('updatescenehierarchy', this.createUpdateSceneHierarchyComponent());

        this.components.set('stopanimation', this.createStopAnimationComponent());

        this.components.set('stoprecording', this.createStopRecordingComponent());

        this.components.set('listentitycomponents', this.createListEntityComponentsComponent());

        this.components.set('addcomponent', this.createAddComponentComponent());

        this.components.set('listsessionscenes', this.createListSessionScenesComponent());

        this.components.set('createsessionscene', this.createCreateSessionSceneComponent());

        this.components.set('leavesessionchannel', this.createLeaveSessionChannelComponent());

        this.components.set('getsessiongraph', this.createGetSessionGraphComponent());

        this.components.set('updatesessiongraph', this.createUpdateSessionGraphComponent());

        this.components.set('getlogs', this.createGetLogsComponent());

        this.components.set('importscenedefinition', this.createImportSceneDefinitionComponent());

        this.components.set('listrigidbodies', this.createListRigidBodiesComponent());

        this.components.set('enableentity', this.createEnableEntityComponent());

        this.components.set('bulkentitylifecycleoperation', this.createBulkEntityLifecycleOperationComponent());

        this.components.set('loadscenestate', this.createLoadSceneStateComponent());

        this.components.set('stopaudio', this.createStopAudioComponent());

        this.components.set('joinsessionchannel', this.createJoinSessionChannelComponent());

        this.components.set('getentity', this.createGetEntityComponent());

        this.components.set('updateentity', this.createUpdateEntityComponent());

        this.components.set('deleteentity', this.createDeleteEntityComponent());

        this.components.set('getentitytransforms', this.createGetEntityTransformsComponent());

        this.components.set('setentitytransforms', this.createSetEntityTransformsComponent());

        this.components.set('activateentity', this.createActivateEntityComponent());

        this.components.set('playanimation', this.createPlayAnimationComponent());

        this.components.set('updatephysicsworld', this.createUpdatePhysicsWorldComponent());

        this.components.set('getphysicsworld', this.createGetPhysicsWorldComponent());

        console.log('✅ UI components initialized');
    }


    // Component for GET /sessions/{sessionId}/audio/sources
    createListAudioSourcesComponent() {
        return {
            name: 'listaudiosources',
            endpoint: '/sessions/{sessionId}/audio/sources',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListAudioSources</h4><form id="listAudioSources-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="listAudioSources-result" class="result-area"></div></div>';
                this.attachEventListeners('listaudiosources', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listAudioSources(formData.param1);
                    this.showResult('listaudiosources', result);
                    return result;
                } catch (error) {
                    this.showError('listaudiosources', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/audio/sources
    createCreateAudioSourceComponent() {
        return {
            name: 'createaudiosource',
            endpoint: '/sessions/{sessionId}/audio/sources',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateAudioSource</h4><form id="createAudioSource-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createAudioSource-result" class="result-area"></div></div>';
                this.attachEventListeners('createaudiosource', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createAudioSource(formData.param1, formData);
                    this.showResult('createaudiosource', result);
                    return result;
                } catch (error) {
                    this.showError('createaudiosource', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>ForceRefresh</h4><form id="forceRefresh-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="forceRefresh-result" class="result-area"></div></div>';
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

    // Component for GET /sessions/{sessionId}/avatar
    createGetSessionAvatarComponent() {
        return {
            name: 'getsessionavatar',
            endpoint: '/sessions/{sessionId}/avatar',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSessionAvatar</h4><form id="getSessionAvatar-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSessionAvatar-result" class="result-area"></div></div>';
                this.attachEventListeners('getsessionavatar', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSessionAvatar(formData.param1);
                    this.showResult('getsessionavatar', result);
                    return result;
                } catch (error) {
                    this.showError('getsessionavatar', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/avatar
    createSetSessionAvatarComponent() {
        return {
            name: 'setsessionavatar',
            endpoint: '/sessions/{sessionId}/avatar',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>SetSessionAvatar</h4><form id="setSessionAvatar-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setSessionAvatar-result" class="result-area"></div></div>';
                this.attachEventListeners('setsessionavatar', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setSessionAvatar(formData.param1, formData);
                    this.showResult('setsessionavatar', result);
                    return result;
                } catch (error) {
                    this.showError('setsessionavatar', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children
    createGetEntityChildrenComponent() {
        return {
            name: 'getentitychildren',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/hierarchy/children',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetEntityChildren</h4><form id="getEntityChildren-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getEntityChildren-result" class="result-area"></div></div>';
                this.attachEventListeners('getentitychildren', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getEntityChildren(formData.param1, formData.param2);
                    this.showResult('getentitychildren', result);
                    return result;
                } catch (error) {
                    this.showError('getentitychildren', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/audio/sources/{audioId}/play
    createPlayAudioComponent() {
        return {
            name: 'playaudio',
            endpoint: '/sessions/{sessionId}/audio/sources/{audioId}/play',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>PlayAudio</h4><form id="playAudio-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="playAudio-result" class="result-area"></div></div>';
                this.attachEventListeners('playaudio', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.playAudio(formData.param1, formData.param2, formData);
                    this.showResult('playaudio', result);
                    return result;
                } catch (error) {
                    this.showError('playaudio', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>PlayRecording</h4><form id="playRecording-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="playRecording-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>GetRecordingStatus</h4><form id="getRecordingStatus-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getRecordingStatus-result" class="result-area"></div></div>';
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

    // Component for GET /channels
    createListChannelsComponent() {
        return {
            name: 'listchannels',
            endpoint: '/channels',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListChannels</h4><form id="listChannels-form"><button type="submit">Execute</button></form><div id="listChannels-result" class="result-area"></div></div>';
                this.attachEventListeners('listchannels', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listChannels();
                    this.showResult('listchannels', result);
                    return result;
                } catch (error) {
                    this.showError('listchannels', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /channels
    createCreateChannelComponent() {
        return {
            name: 'createchannel',
            endpoint: '/channels',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateChannel</h4><form id="createChannel-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createChannel-result" class="result-area"></div></div>';
                this.attachEventListeners('createchannel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createChannel(formData);
                    this.showResult('createchannel', result);
                    return result;
                } catch (error) {
                    this.showError('createchannel', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>ListSessions</h4><form id="listSessions-form"><button type="submit">Execute</button></form><div id="listSessions-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateSession</h4><form id="createSession-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createSession-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>StartRecording</h4><form id="startRecording-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="startRecording-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>SetLogLevel</h4><form id="setLogLevel-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setLogLevel-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>DeleteSession</h4><form id="deleteSession-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="deleteSession-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSession</h4><form id="getSession-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSession-result" class="result-area"></div></div>';
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

    // Component for POST /sessions/{sessionId}/channel/sync
    createSyncSessionStateComponent() {
        return {
            name: 'syncsessionstate',
            endpoint: '/sessions/{sessionId}/channel/sync',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>SyncSessionState</h4><form id="syncSessionState-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="syncSessionState-result" class="result-area"></div></div>';
                this.attachEventListeners('syncsessionstate', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.syncSessionState(formData.param1, formData);
                    this.showResult('syncsessionstate', result);
                    return result;
                } catch (error) {
                    this.showError('syncsessionstate', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/hierarchy/tree
    createGetHierarchyTreeComponent() {
        return {
            name: 'gethierarchytree',
            endpoint: '/sessions/{sessionId}/entities/hierarchy/tree',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetHierarchyTree</h4><form id="getHierarchyTree-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getHierarchyTree-result" class="result-area"></div></div>';
                this.attachEventListeners('gethierarchytree', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getHierarchyTree(formData.param1);
                    this.showResult('gethierarchytree', result);
                    return result;
                } catch (error) {
                    this.showError('gethierarchytree', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate
    createDeactivateEntityComponent() {
        return {
            name: 'deactivateentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>DeactivateEntity</h4><form id="deactivateEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="deactivateEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('deactivateentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.deactivateEntity(formData.param1, formData.param2, formData);
                    this.showResult('deactivateentity', result);
                    return result;
                } catch (error) {
                    this.showError('deactivateentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/scenes/{sceneId}/activate
    createActivateSessionSceneComponent() {
        return {
            name: 'activatesessionscene',
            endpoint: '/sessions/{sessionId}/scenes/{sceneId}/activate',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ActivateSessionScene</h4><form id="activateSessionScene-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="activateSessionScene-result" class="result-area"></div></div>';
                this.attachEventListeners('activatesessionscene', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.activateSessionScene(formData.param1, formData.param2, formData);
                    this.showResult('activatesessionscene', result);
                    return result;
                } catch (error) {
                    this.showError('activatesessionscene', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/animations
    createListAnimationsComponent() {
        return {
            name: 'listanimations',
            endpoint: '/sessions/{sessionId}/animations',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListAnimations</h4><form id="listAnimations-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="listAnimations-result" class="result-area"></div></div>';
                this.attachEventListeners('listanimations', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listAnimations(formData.param1);
                    this.showResult('listanimations', result);
                    return result;
                } catch (error) {
                    this.showError('listanimations', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/animations
    createCreateAnimationComponent() {
        return {
            name: 'createanimation',
            endpoint: '/sessions/{sessionId}/animations',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateAnimation</h4><form id="createAnimation-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createAnimation-result" class="result-area"></div></div>';
                this.attachEventListeners('createanimation', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createAnimation(formData.param1, formData);
                    this.showResult('createanimation', result);
                    return result;
                } catch (error) {
                    this.showError('createanimation', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>SetLoggingConfig</h4><form id="setLoggingConfig-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setLoggingConfig-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>GetLoggingConfig</h4><form id="getLoggingConfig-form"><button type="submit">Execute</button></form><div id="getLoggingConfig-result" class="result-area"></div></div>';
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

    // Component for GET /sessions/{sessionId}/channel/status
    createGetSessionChannelStatusComponent() {
        return {
            name: 'getsessionchannelstatus',
            endpoint: '/sessions/{sessionId}/channel/status',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSessionChannelStatus</h4><form id="getSessionChannelStatus-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSessionChannelStatus-result" class="result-area"></div></div>';
                this.attachEventListeners('getsessionchannelstatus', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSessionChannelStatus(formData.param1);
                    this.showResult('getsessionchannelstatus', result);
                    return result;
                } catch (error) {
                    this.showError('getsessionchannelstatus', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /channels/{channelId}
    createGetChannelComponent() {
        return {
            name: 'getchannel',
            endpoint: '/channels/{channelId}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetChannel</h4><form id="getChannel-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getChannel-result" class="result-area"></div></div>';
                this.attachEventListeners('getchannel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getChannel(formData.param1);
                    this.showResult('getchannel', result);
                    return result;
                } catch (error) {
                    this.showError('getchannel', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /channels/{channelId}
    createUpdateChannelComponent() {
        return {
            name: 'updatechannel',
            endpoint: '/channels/{channelId}',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateChannel</h4><form id="updateChannel-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateChannel-result" class="result-area"></div></div>';
                this.attachEventListeners('updatechannel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateChannel(formData.param1, formData);
                    this.showResult('updatechannel', result);
                    return result;
                } catch (error) {
                    this.showError('updatechannel', error);
                    throw error;
                }
            }
        };
    }

    // Component for DELETE /channels/{channelId}
    createDeleteChannelComponent() {
        return {
            name: 'deletechannel',
            endpoint: '/channels/{channelId}',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>DeleteChannel</h4><form id="deleteChannel-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="deleteChannel-result" class="result-area"></div></div>';
                this.attachEventListeners('deletechannel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.deleteChannel(formData.param1);
                    this.showResult('deletechannel', result);
                    return result;
                } catch (error) {
                    this.showError('deletechannel', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /avatars
    createListAvatarsComponent() {
        return {
            name: 'listavatars',
            endpoint: '/avatars',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListAvatars</h4><form id="listAvatars-form"><button type="submit">Execute</button></form><div id="listAvatars-result" class="result-area"></div></div>';
                this.attachEventListeners('listavatars', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listAvatars();
                    this.showResult('listavatars', result);
                    return result;
                } catch (error) {
                    this.showError('listavatars', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /avatars/{avatarType}/asset
    createGetAvatarAssetComponent() {
        return {
            name: 'getavatarasset',
            endpoint: '/avatars/{avatarType}/asset',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetAvatarAsset</h4><form id="getAvatarAsset-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getAvatarAsset-result" class="result-area"></div></div>';
                this.attachEventListeners('getavatarasset', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getAvatarAsset(formData.param1);
                    this.showResult('getavatarasset', result);
                    return result;
                } catch (error) {
                    this.showError('getavatarasset', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities
    createListEntitiesComponent() {
        return {
            name: 'listentities',
            endpoint: '/sessions/{sessionId}/entities',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListEntities</h4><form id="listEntities-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="listEntities-result" class="result-area"></div></div>';
                this.attachEventListeners('listentities', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listEntities(formData.param1);
                    this.showResult('listentities', result);
                    return result;
                } catch (error) {
                    this.showError('listentities', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/entities
    createCreateEntityComponent() {
        return {
            name: 'createentity',
            endpoint: '/sessions/{sessionId}/entities',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateEntity</h4><form id="createEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('createentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createEntity(formData.param1, formData);
                    this.showResult('createentity', result);
                    return result;
                } catch (error) {
                    this.showError('createentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/camera/position
    createGetCameraPositionComponent() {
        return {
            name: 'getcameraposition',
            endpoint: '/sessions/{sessionId}/camera/position',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetCameraPosition</h4><form id="getCameraPosition-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getCameraPosition-result" class="result-area"></div></div>';
                this.attachEventListeners('getcameraposition', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getCameraPosition(formData.param1);
                    this.showResult('getcameraposition', result);
                    return result;
                } catch (error) {
                    this.showError('getcameraposition', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>SetCameraPosition</h4><form id="setCameraPosition-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setCameraPosition-result" class="result-area"></div></div>';
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

    // Component for PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable
    createDisableEntityComponent() {
        return {
            name: 'disableentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/lifecycle/disable',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>DisableEntity</h4><form id="disableEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="disableEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('disableentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.disableEntity(formData.param1, formData.param2, formData);
                    this.showResult('disableentity', result);
                    return result;
                } catch (error) {
                    this.showError('disableentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/scene/state/reset
    createResetSceneStateComponent() {
        return {
            name: 'resetscenestate',
            endpoint: '/sessions/{sessionId}/scene/state/reset',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ResetSceneState</h4><form id="resetSceneState-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="resetSceneState-result" class="result-area"></div></div>';
                this.attachEventListeners('resetscenestate', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.resetSceneState(formData.param1, formData);
                    this.showResult('resetscenestate', result);
                    return result;
                } catch (error) {
                    this.showError('resetscenestate', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status
    createGetEntityLifecycleStatusComponent() {
        return {
            name: 'getentitylifecyclestatus',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/lifecycle/status',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetEntityLifecycleStatus</h4><form id="getEntityLifecycleStatus-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getEntityLifecycleStatus-result" class="result-area"></div></div>';
                this.attachEventListeners('getentitylifecyclestatus', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getEntityLifecycleStatus(formData.param1, formData.param2);
                    this.showResult('getentitylifecyclestatus', result);
                    return result;
                } catch (error) {
                    this.showError('getentitylifecyclestatus', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/scene/state
    createGetSceneStateComponent() {
        return {
            name: 'getscenestate',
            endpoint: '/sessions/{sessionId}/scene/state',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSceneState</h4><form id="getSceneState-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSceneState-result" class="result-area"></div></div>';
                this.attachEventListeners('getscenestate', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSceneState(formData.param1);
                    this.showResult('getscenestate', result);
                    return result;
                } catch (error) {
                    this.showError('getscenestate', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/scene/state
    createUpdateSceneStateComponent() {
        return {
            name: 'updatescenestate',
            endpoint: '/sessions/{sessionId}/scene/state',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateSceneState</h4><form id="updateSceneState-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateSceneState-result" class="result-area"></div></div>';
                this.attachEventListeners('updatescenestate', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateSceneState(formData.param1, formData);
                    this.showResult('updatescenestate', result);
                    return result;
                } catch (error) {
                    this.showError('updatescenestate', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/scene/state/save
    createSaveSceneStateComponent() {
        return {
            name: 'savescenestate',
            endpoint: '/sessions/{sessionId}/scene/state/save',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>SaveSceneState</h4><form id="saveSceneState-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="saveSceneState-result" class="result-area"></div></div>';
                this.attachEventListeners('savescenestate', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.saveSceneState(formData.param1, formData);
                    this.showResult('savescenestate', result);
                    return result;
                } catch (error) {
                    this.showError('savescenestate', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force
    createApplyForceComponent() {
        return {
            name: 'applyforce',
            endpoint: '/sessions/{sessionId}/physics/rigidbodies/{entityId}/force',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ApplyForce</h4><form id="applyForce-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="applyForce-result" class="result-area"></div></div>';
                this.attachEventListeners('applyforce', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.applyForce(formData.param1, formData.param2, formData);
                    this.showResult('applyforce', result);
                    return result;
                } catch (error) {
                    this.showError('applyforce', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /version
    createGetVersionComponent() {
        return {
            name: 'getversion',
            endpoint: '/version',
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

    // Component for DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy
    createDestroyEntityComponent() {
        return {
            name: 'destroyentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/lifecycle/destroy',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>DestroyEntity</h4><form id="destroyEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="destroyEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('destroyentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.destroyEntity(formData.param1, formData.param2);
                    this.showResult('destroyentity', result);
                    return result;
                } catch (error) {
                    this.showError('destroyentity', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>SetTraceModules</h4><form id="setTraceModules-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setTraceModules-result" class="result-area"></div></div>';
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
                
                container.innerHTML = '<div class="hd1-component"><h4>SetCanvas</h4><form id="setCanvas-form"><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setCanvas-result" class="result-area"></div></div>';
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

    // Component for GET /avatars/{avatarType}
    createGetAvatarSpecificationComponent() {
        return {
            name: 'getavatarspecification',
            endpoint: '/avatars/{avatarType}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetAvatarSpecification</h4><form id="getAvatarSpecification-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getAvatarSpecification-result" class="result-area"></div></div>';
                this.attachEventListeners('getavatarspecification', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getAvatarSpecification(formData.param1);
                    this.showResult('getavatarspecification', result);
                    return result;
                } catch (error) {
                    this.showError('getavatarspecification', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}
    createGetComponentComponent() {
        return {
            name: 'getcomponent',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/components/{componentType}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetComponent</h4><form id="getComponent-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getComponent-result" class="result-area"></div></div>';
                this.attachEventListeners('getcomponent', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getComponent(formData.param1, formData.param2);
                    this.showResult('getcomponent', result);
                    return result;
                } catch (error) {
                    this.showError('getcomponent', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}
    createUpdateComponentComponent() {
        return {
            name: 'updatecomponent',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/components/{componentType}',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateComponent</h4><form id="updateComponent-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateComponent-result" class="result-area"></div></div>';
                this.attachEventListeners('updatecomponent', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateComponent(formData.param1, formData.param2, formData);
                    this.showResult('updatecomponent', result);
                    return result;
                } catch (error) {
                    this.showError('updatecomponent', error);
                    throw error;
                }
            }
        };
    }

    // Component for DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}
    createRemoveComponentComponent() {
        return {
            name: 'removecomponent',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/components/{componentType}',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>RemoveComponent</h4><form id="removeComponent-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="removeComponent-result" class="result-area"></div></div>';
                this.attachEventListeners('removecomponent', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.removeComponent(formData.param1, formData.param2);
                    this.showResult('removecomponent', result);
                    return result;
                } catch (error) {
                    this.showError('removecomponent', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/scene/export
    createExportSceneDefinitionComponent() {
        return {
            name: 'exportscenedefinition',
            endpoint: '/sessions/{sessionId}/scene/export',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ExportSceneDefinition</h4><form id="exportSceneDefinition-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="exportSceneDefinition-result" class="result-area"></div></div>';
                this.attachEventListeners('exportscenedefinition', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.exportSceneDefinition(formData.param1);
                    this.showResult('exportscenedefinition', result);
                    return result;
                } catch (error) {
                    this.showError('exportscenedefinition', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>StartCameraOrbit</h4><form id="startCameraOrbit-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="startCameraOrbit-result" class="result-area"></div></div>';
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

    // Component for POST /sessions/{sessionId}/entities/{entityId}/components/bulk
    createBulkComponentOperationComponent() {
        return {
            name: 'bulkcomponentoperation',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/components/bulk',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>BulkComponentOperation</h4><form id="bulkComponentOperation-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="bulkComponentOperation-result" class="result-area"></div></div>';
                this.attachEventListeners('bulkcomponentoperation', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.bulkComponentOperation(formData.param1, formData.param2, formData);
                    this.showResult('bulkcomponentoperation', result);
                    return result;
                } catch (error) {
                    this.showError('bulkcomponentoperation', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
    createGetEntityParentComponent() {
        return {
            name: 'getentityparent',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/hierarchy/parent',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetEntityParent</h4><form id="getEntityParent-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getEntityParent-result" class="result-area"></div></div>';
                this.attachEventListeners('getentityparent', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getEntityParent(formData.param1, formData.param2);
                    this.showResult('getentityparent', result);
                    return result;
                } catch (error) {
                    this.showError('getentityparent', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
    createSetEntityParentComponent() {
        return {
            name: 'setentityparent',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/hierarchy/parent',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>SetEntityParent</h4><form id="setEntityParent-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setEntityParent-result" class="result-area"></div></div>';
                this.attachEventListeners('setentityparent', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setEntityParent(formData.param1, formData.param2, formData);
                    this.showResult('setentityparent', result);
                    return result;
                } catch (error) {
                    this.showError('setentityparent', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/scene/hierarchy
    createGetSceneHierarchyComponent() {
        return {
            name: 'getscenehierarchy',
            endpoint: '/sessions/{sessionId}/scene/hierarchy',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSceneHierarchy</h4><form id="getSceneHierarchy-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSceneHierarchy-result" class="result-area"></div></div>';
                this.attachEventListeners('getscenehierarchy', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSceneHierarchy(formData.param1);
                    this.showResult('getscenehierarchy', result);
                    return result;
                } catch (error) {
                    this.showError('getscenehierarchy', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/scene/hierarchy
    createUpdateSceneHierarchyComponent() {
        return {
            name: 'updatescenehierarchy',
            endpoint: '/sessions/{sessionId}/scene/hierarchy',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateSceneHierarchy</h4><form id="updateSceneHierarchy-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateSceneHierarchy-result" class="result-area"></div></div>';
                this.attachEventListeners('updatescenehierarchy', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateSceneHierarchy(formData.param1, formData);
                    this.showResult('updatescenehierarchy', result);
                    return result;
                } catch (error) {
                    this.showError('updatescenehierarchy', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/animations/{animationId}/stop
    createStopAnimationComponent() {
        return {
            name: 'stopanimation',
            endpoint: '/sessions/{sessionId}/animations/{animationId}/stop',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>StopAnimation</h4><form id="stopAnimation-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="stopAnimation-result" class="result-area"></div></div>';
                this.attachEventListeners('stopanimation', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.stopAnimation(formData.param1, formData.param2, formData);
                    this.showResult('stopanimation', result);
                    return result;
                } catch (error) {
                    this.showError('stopanimation', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>StopRecording</h4><form id="stopRecording-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="stopRecording-result" class="result-area"></div></div>';
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

    // Component for GET /sessions/{sessionId}/entities/{entityId}/components
    createListEntityComponentsComponent() {
        return {
            name: 'listentitycomponents',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/components',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListEntityComponents</h4><form id="listEntityComponents-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="listEntityComponents-result" class="result-area"></div></div>';
                this.attachEventListeners('listentitycomponents', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listEntityComponents(formData.param1, formData.param2);
                    this.showResult('listentitycomponents', result);
                    return result;
                } catch (error) {
                    this.showError('listentitycomponents', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/entities/{entityId}/components
    createAddComponentComponent() {
        return {
            name: 'addcomponent',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/components',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>AddComponent</h4><form id="addComponent-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="addComponent-result" class="result-area"></div></div>';
                this.attachEventListeners('addcomponent', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.addComponent(formData.param1, formData.param2, formData);
                    this.showResult('addcomponent', result);
                    return result;
                } catch (error) {
                    this.showError('addcomponent', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/scenes
    createListSessionScenesComponent() {
        return {
            name: 'listsessionscenes',
            endpoint: '/sessions/{sessionId}/scenes',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListSessionScenes</h4><form id="listSessionScenes-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="listSessionScenes-result" class="result-area"></div></div>';
                this.attachEventListeners('listsessionscenes', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listSessionScenes(formData.param1);
                    this.showResult('listsessionscenes', result);
                    return result;
                } catch (error) {
                    this.showError('listsessionscenes', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/scenes
    createCreateSessionSceneComponent() {
        return {
            name: 'createsessionscene',
            endpoint: '/sessions/{sessionId}/scenes',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>CreateSessionScene</h4><form id="createSessionScene-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="createSessionScene-result" class="result-area"></div></div>';
                this.attachEventListeners('createsessionscene', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.createSessionScene(formData.param1, formData);
                    this.showResult('createsessionscene', result);
                    return result;
                } catch (error) {
                    this.showError('createsessionscene', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/channel/leave
    createLeaveSessionChannelComponent() {
        return {
            name: 'leavesessionchannel',
            endpoint: '/sessions/{sessionId}/channel/leave',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>LeaveSessionChannel</h4><form id="leaveSessionChannel-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="leaveSessionChannel-result" class="result-area"></div></div>';
                this.attachEventListeners('leavesessionchannel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.leaveSessionChannel(formData.param1, formData);
                    this.showResult('leavesessionchannel', result);
                    return result;
                } catch (error) {
                    this.showError('leavesessionchannel', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/channel/graph
    createGetSessionGraphComponent() {
        return {
            name: 'getsessiongraph',
            endpoint: '/sessions/{sessionId}/channel/graph',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetSessionGraph</h4><form id="getSessionGraph-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getSessionGraph-result" class="result-area"></div></div>';
                this.attachEventListeners('getsessiongraph', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getSessionGraph(formData.param1);
                    this.showResult('getsessiongraph', result);
                    return result;
                } catch (error) {
                    this.showError('getsessiongraph', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/channel/graph
    createUpdateSessionGraphComponent() {
        return {
            name: 'updatesessiongraph',
            endpoint: '/sessions/{sessionId}/channel/graph',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateSessionGraph</h4><form id="updateSessionGraph-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateSessionGraph-result" class="result-area"></div></div>';
                this.attachEventListeners('updatesessiongraph', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateSessionGraph(formData.param1, formData);
                    this.showResult('updatesessiongraph', result);
                    return result;
                } catch (error) {
                    this.showError('updatesessiongraph', error);
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
                
                container.innerHTML = '<div class="hd1-component"><h4>GetLogs</h4><form id="getLogs-form"><button type="submit">Execute</button></form><div id="getLogs-result" class="result-area"></div></div>';
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

    // Component for POST /sessions/{sessionId}/scene/import
    createImportSceneDefinitionComponent() {
        return {
            name: 'importscenedefinition',
            endpoint: '/sessions/{sessionId}/scene/import',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ImportSceneDefinition</h4><form id="importSceneDefinition-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="importSceneDefinition-result" class="result-area"></div></div>';
                this.attachEventListeners('importscenedefinition', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.importSceneDefinition(formData.param1, formData);
                    this.showResult('importscenedefinition', result);
                    return result;
                } catch (error) {
                    this.showError('importscenedefinition', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/physics/rigidbodies
    createListRigidBodiesComponent() {
        return {
            name: 'listrigidbodies',
            endpoint: '/sessions/{sessionId}/physics/rigidbodies',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ListRigidBodies</h4><form id="listRigidBodies-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="listRigidBodies-result" class="result-area"></div></div>';
                this.attachEventListeners('listrigidbodies', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.listRigidBodies(formData.param1);
                    this.showResult('listrigidbodies', result);
                    return result;
                } catch (error) {
                    this.showError('listrigidbodies', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable
    createEnableEntityComponent() {
        return {
            name: 'enableentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/lifecycle/enable',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>EnableEntity</h4><form id="enableEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="enableEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('enableentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.enableEntity(formData.param1, formData.param2, formData);
                    this.showResult('enableentity', result);
                    return result;
                } catch (error) {
                    this.showError('enableentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/entities/lifecycle/bulk
    createBulkEntityLifecycleOperationComponent() {
        return {
            name: 'bulkentitylifecycleoperation',
            endpoint: '/sessions/{sessionId}/entities/lifecycle/bulk',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>BulkEntityLifecycleOperation</h4><form id="bulkEntityLifecycleOperation-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="bulkEntityLifecycleOperation-result" class="result-area"></div></div>';
                this.attachEventListeners('bulkentitylifecycleoperation', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.bulkEntityLifecycleOperation(formData.param1, formData);
                    this.showResult('bulkentitylifecycleoperation', result);
                    return result;
                } catch (error) {
                    this.showError('bulkentitylifecycleoperation', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/scene/state/load
    createLoadSceneStateComponent() {
        return {
            name: 'loadscenestate',
            endpoint: '/sessions/{sessionId}/scene/state/load',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>LoadSceneState</h4><form id="loadSceneState-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="loadSceneState-result" class="result-area"></div></div>';
                this.attachEventListeners('loadscenestate', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.loadSceneState(formData.param1, formData);
                    this.showResult('loadscenestate', result);
                    return result;
                } catch (error) {
                    this.showError('loadscenestate', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/audio/sources/{audioId}/stop
    createStopAudioComponent() {
        return {
            name: 'stopaudio',
            endpoint: '/sessions/{sessionId}/audio/sources/{audioId}/stop',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>StopAudio</h4><form id="stopAudio-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="stopAudio-result" class="result-area"></div></div>';
                this.attachEventListeners('stopaudio', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.stopAudio(formData.param1, formData.param2, formData);
                    this.showResult('stopaudio', result);
                    return result;
                } catch (error) {
                    this.showError('stopaudio', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/channel/join
    createJoinSessionChannelComponent() {
        return {
            name: 'joinsessionchannel',
            endpoint: '/sessions/{sessionId}/channel/join',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>JoinSessionChannel</h4><form id="joinSessionChannel-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="joinSessionChannel-result" class="result-area"></div></div>';
                this.attachEventListeners('joinsessionchannel', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.joinSessionChannel(formData.param1, formData);
                    this.showResult('joinsessionchannel', result);
                    return result;
                } catch (error) {
                    this.showError('joinsessionchannel', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/{entityId}
    createGetEntityComponent() {
        return {
            name: 'getentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetEntity</h4><form id="getEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('getentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getEntity(formData.param1, formData.param2);
                    this.showResult('getentity', result);
                    return result;
                } catch (error) {
                    this.showError('getentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/entities/{entityId}
    createUpdateEntityComponent() {
        return {
            name: 'updateentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdateEntity</h4><form id="updateEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updateEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('updateentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updateEntity(formData.param1, formData.param2, formData);
                    this.showResult('updateentity', result);
                    return result;
                } catch (error) {
                    this.showError('updateentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for DELETE /sessions/{sessionId}/entities/{entityId}
    createDeleteEntityComponent() {
        return {
            name: 'deleteentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}',
            method: 'DELETE',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>DeleteEntity</h4><form id="deleteEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="deleteEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('deleteentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.deleteEntity(formData.param1, formData.param2);
                    this.showResult('deleteentity', result);
                    return result;
                } catch (error) {
                    this.showError('deleteentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms
    createGetEntityTransformsComponent() {
        return {
            name: 'getentitytransforms',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/hierarchy/transforms',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetEntityTransforms</h4><form id="getEntityTransforms-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><button type="submit">Execute</button></form><div id="getEntityTransforms-result" class="result-area"></div></div>';
                this.attachEventListeners('getentitytransforms', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getEntityTransforms(formData.param1, formData.param2);
                    this.showResult('getentitytransforms', result);
                    return result;
                } catch (error) {
                    this.showError('getentitytransforms', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms
    createSetEntityTransformsComponent() {
        return {
            name: 'setentitytransforms',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/hierarchy/transforms',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>SetEntityTransforms</h4><form id="setEntityTransforms-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="setEntityTransforms-result" class="result-area"></div></div>';
                this.attachEventListeners('setentitytransforms', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.setEntityTransforms(formData.param1, formData.param2, formData);
                    this.showResult('setentitytransforms', result);
                    return result;
                } catch (error) {
                    this.showError('setentitytransforms', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate
    createActivateEntityComponent() {
        return {
            name: 'activateentity',
            endpoint: '/sessions/{sessionId}/entities/{entityId}/lifecycle/activate',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>ActivateEntity</h4><form id="activateEntity-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="activateEntity-result" class="result-area"></div></div>';
                this.attachEventListeners('activateentity', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.activateEntity(formData.param1, formData.param2, formData);
                    this.showResult('activateentity', result);
                    return result;
                } catch (error) {
                    this.showError('activateentity', error);
                    throw error;
                }
            }
        };
    }

    // Component for POST /sessions/{sessionId}/animations/{animationId}/play
    createPlayAnimationComponent() {
        return {
            name: 'playanimation',
            endpoint: '/sessions/{sessionId}/animations/{animationId}/play',
            method: 'POST',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>PlayAnimation</h4><form id="playAnimation-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="playAnimation-result" class="result-area"></div></div>';
                this.attachEventListeners('playanimation', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.playAnimation(formData.param1, formData.param2, formData);
                    this.showResult('playanimation', result);
                    return result;
                } catch (error) {
                    this.showError('playanimation', error);
                    throw error;
                }
            }
        };
    }

    // Component for PUT /sessions/{sessionId}/physics/world
    createUpdatePhysicsWorldComponent() {
        return {
            name: 'updatephysicsworld',
            endpoint: '/sessions/{sessionId}/physics/world',
            method: 'PUT',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>UpdatePhysicsWorld</h4><form id="updatePhysicsWorld-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div><button type="submit">Execute</button></form><div id="updatePhysicsWorld-result" class="result-area"></div></div>';
                this.attachEventListeners('updatephysicsworld', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.updatePhysicsWorld(formData.param1, formData);
                    this.showResult('updatephysicsworld', result);
                    return result;
                } catch (error) {
                    this.showError('updatephysicsworld', error);
                    throw error;
                }
            }
        };
    }

    // Component for GET /sessions/{sessionId}/physics/world
    createGetPhysicsWorldComponent() {
        return {
            name: 'getphysicsworld',
            endpoint: '/sessions/{sessionId}/physics/world',
            method: 'GET',
            
            render: (containerId) => {
                const container = document.getElementById(containerId);
                if (!container) {
                    console.error('Container not found:', containerId);
                    return;
                }
                
                container.innerHTML = '<div class="hd1-component"><h4>GetPhysicsWorld</h4><form id="getPhysicsWorld-form"><div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div><button type="submit">Execute</button></form><div id="getPhysicsWorld-result" class="result-area"></div></div>';
                this.attachEventListeners('getphysicsworld', container);
            },
            
            execute: async (formData) => {
                try {
                    const result = await this.api.getPhysicsWorld(formData.param1);
                    this.showResult('getphysicsworld', result);
                    return result;
                } catch (error) {
                    this.showError('getphysicsworld', error);
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