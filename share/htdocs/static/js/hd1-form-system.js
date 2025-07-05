// ===================================================================
// AUTO-GENERATED FORM SYSTEM - DO NOT MODIFY
// ===================================================================
//
// Dynamic form generation system auto-generated from OpenAPI schemas
//
// ===================================================================

class HD1FormSystem {
    constructor() {
        this.formSchemas = new Map();
        this.initializeSchemas();
    }

    initializeSchemas() {
        console.log('📝 Initializing auto-generated form schemas...');

        this.formSchemas.set('joinSessionChannelForm', {
        "title": "JoinSessionChannel",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setEntityTransformsForm', {
        "title": "SetEntityTransforms",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('createAudioSourceForm', {
        "title": "CreateAudioSource",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setLoggingConfigForm', {
        "title": "SetLoggingConfig",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('saveSceneStateForm', {
        "title": "SaveSceneState",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updatePhysicsWorldForm', {
        "title": "UpdatePhysicsWorld",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('createSessionForm', {
        "title": "CreateSession",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('bulkEntityLifecycleOperationForm', {
        "title": "BulkEntityLifecycleOperation",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('applyForceForm', {
        "title": "ApplyForce",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setLogLevelForm', {
        "title": "SetLogLevel",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setCameraPositionForm', {
        "title": "SetCameraPosition",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('loadSceneStateForm', {
        "title": "LoadSceneState",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('createAnimationForm', {
        "title": "CreateAnimation",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('playAnimationForm', {
        "title": "PlayAnimation",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('stopAudioForm', {
        "title": "StopAudio",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('startCameraOrbitForm', {
        "title": "StartCameraOrbit",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('postClientJoinSyncForm', {
        "title": "PostClientJoinSync",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updateEntityForm', {
        "title": "UpdateEntity",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('activateEntityForm', {
        "title": "ActivateEntity",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('deactivateEntityForm', {
        "title": "DeactivateEntity",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('startRecordingForm', {
        "title": "StartRecording",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setSessionAvatarForm', {
        "title": "SetSessionAvatar",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('leaveSessionChannelForm', {
        "title": "LeaveSessionChannel",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updateComponentForm', {
        "title": "UpdateComponent",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setTraceModulesForm', {
        "title": "SetTraceModules",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('createChannelForm', {
        "title": "CreateChannel",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('enableEntityForm', {
        "title": "EnableEntity",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('createSessionSceneForm', {
        "title": "CreateSessionScene",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('activateSessionSceneForm', {
        "title": "ActivateSessionScene",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('playRecordingForm', {
        "title": "PlayRecording",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updateSceneHierarchyForm', {
        "title": "UpdateSceneHierarchy",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('resetSceneStateForm', {
        "title": "ResetSceneState",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('disableEntityForm', {
        "title": "DisableEntity",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updateSceneStateForm', {
        "title": "UpdateSceneState",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('importSceneDefinitionForm', {
        "title": "ImportSceneDefinition",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('stopAnimationForm', {
        "title": "StopAnimation",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setCanvasForm', {
        "title": "SetCanvas",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updateChannelForm', {
        "title": "UpdateChannel",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('addComponentForm', {
        "title": "AddComponent",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('playAudioForm', {
        "title": "PlayAudio",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('bulkComponentOperationForm', {
        "title": "BulkComponentOperation",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('forceRefreshForm', {
        "title": "ForceRefresh",
        "submitText": "Execute POST",
        "fields": {"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('createEntityForm', {
        "title": "CreateEntity",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('setEntityParentForm', {
        "title": "SetEntityParent",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "param2": {"type": "string", "title": "Parameter 2", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('updateSessionGraphForm', {
        "title": "UpdateSessionGraph",
        "submitText": "Execute PUT",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('syncSessionStateForm', {
        "title": "SyncSessionState",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        this.formSchemas.set('stopRecordingForm', {
        "title": "StopRecording",
        "submitText": "Execute POST",
        "fields": {"param1": {"type": "string", "title": "Parameter 1", "required": true}, "data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}}});

        console.log('✅ Form schemas initialized');
    }

    generateForm(schemaName, containerId) {
        const schema = this.formSchemas.get(schemaName);
        if (!schema) {
            console.error('Schema not found:', schemaName);
            return;
        }
        
        const container = document.getElementById(containerId);
        if (!container) {
            console.error('Container not found:', containerId);
            return;
        }
        
        const form = this.createFormFromSchema(schema);
        container.appendChild(form);
    }

    createFormFromSchema(schema) {
        const form = document.createElement('form');
        form.className = 'hd1-auto-form';
        
        if (schema.title) {
            const title = document.createElement('h3');
            title.textContent = schema.title;
            form.appendChild(title);
        }
        
        Object.entries(schema.fields || {}).forEach(([fieldName, fieldSchema]) => {
            const fieldDiv = document.createElement('div');
            fieldDiv.className = 'form-field';
            
            const label = document.createElement('label');
            label.textContent = fieldSchema.title || fieldName;
            label.setAttribute('for', fieldName);
            
            const input = this.createInputForField(fieldName, fieldSchema);
            
            fieldDiv.appendChild(label);
            fieldDiv.appendChild(input);
            form.appendChild(fieldDiv);
        });
        
        const submitBtn = document.createElement('button');
        submitBtn.type = 'submit';
        submitBtn.textContent = schema.submitText || 'Submit';
        submitBtn.className = 'hd1-submit-btn';
        form.appendChild(submitBtn);
        
        return form;
    }

    createInputForField(fieldName, fieldSchema) {
        const input = document.createElement('input');
        input.name = fieldName;
        input.id = fieldName;
        
        switch (fieldSchema.type) {
            case 'string':
                input.type = 'text';
                break;
            case 'number':
                input.type = 'number';
                break;
            case 'boolean':
                input.type = 'checkbox';
                break;
            default:
                input.type = 'text';
        }
        
        if (fieldSchema.required) {
            input.required = true;
        }
        
        if (fieldSchema.placeholder) {
            input.placeholder = fieldSchema.placeholder;
        }
        
        return input;
    }

    getFormData(formElement) {
        const formData = new FormData(formElement);
        const data = {};
        
        for (const [key, value] of formData.entries()) {
            data[key] = value;
        }
        
        return data;
    }
}

// Export for global use
window.HD1FormSystem = HD1FormSystem;

console.log('📝 HD1 Form System loaded - Auto-generated from specification');