package server

import (
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"holodeck/logging"
)

// ScenesWatcher monitors the scenes directory and broadcasts changes via WebSocket
type ScenesWatcher struct {
	hub        *Hub
	watcher    *fsnotify.Watcher
	scenesPath string
}

// NewScenesWatcher creates a new scenes directory watcher
func NewScenesWatcher(hub *Hub) *ScenesWatcher {
	return &ScenesWatcher{
		hub:        hub,
		scenesPath: "/opt/holo-deck/share/scenes",
	}
}

// Start begins watching the scenes directory for changes
func (sw *ScenesWatcher) Start() error {
	var err error
	sw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Add the scenes directory to the watcher
	err = sw.watcher.Add(sw.scenesPath)
	if err != nil {
		logging.Error("failed to add scenes directory to watcher", map[string]interface{}{
			"error": err.Error(),
			"path":  sw.scenesPath,
		})
		return err
	}

	// Start watching in a goroutine  
	go sw.watchLoop()

	logging.Info("scenes watcher started", map[string]interface{}{
		"scenes_path": sw.scenesPath,
	})
	return nil
}

// Stop stops the file watcher
func (sw *ScenesWatcher) Stop() {
	if sw.watcher != nil {
		sw.watcher.Close()
	}
}

// watchLoop handles file system events
func (sw *ScenesWatcher) watchLoop() {
	logging.Info("scenes watcher loop started", map[string]interface{}{
		"scenes_path": sw.scenesPath,
	})
	
	// Debounce multiple rapid changes
	var lastChange time.Time
	debounceInterval := 500 * time.Millisecond

	for {
		select {
		case event, ok := <-sw.watcher.Events:
			if !ok {
				return
			}

			// Only care about .sh files (scene scripts)
			if filepath.Ext(event.Name) != ".sh" {
				continue
			}

			// Debounce rapid changes
			now := time.Now()
			if now.Sub(lastChange) < debounceInterval {
				continue
			}
			lastChange = now

			// Check if this is a scene-related change
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove {

				logging.Info("scene directory changed", map[string]interface{}{
					"file":      event.Name,
					"operation": event.Op.String(),
				})
				sw.broadcastSceneListChanged()
			}

		case err, ok := <-sw.watcher.Errors:
			if !ok {
				return
			}
			logging.Error("scenes watcher error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

// broadcastSceneListChanged sends a WebSocket message to all clients
func (sw *ScenesWatcher) broadcastSceneListChanged() {
	message := map[string]interface{}{
		"type":      "scene_list_changed",
		"timestamp": time.Now().Unix(),
		"message":   "Scene list has been updated",
	}

	sw.hub.BroadcastUpdate("scene_list_changed", message)
	
	// Get client count safely
	sw.hub.mutex.RLock()
	clientCount := len(sw.hub.clients)
	sw.hub.mutex.RUnlock()
	
	logging.Info("scene list change broadcast", map[string]interface{}{
		"clients_notified": clientCount,
		"message_type":     "scene_list_changed",
	})
}