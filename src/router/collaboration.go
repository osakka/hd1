package router

import (
	"github.com/gorilla/mux"
	"holodeck1/api/assets"
	"holodeck1/api/ot"
	"holodeck1/api/webrtc"
	assetsPkg "holodeck1/assets"
	authPkg "holodeck1/auth"
	otPkg "holodeck1/ot"
	webrtcPkg "holodeck1/webrtc"
)

func SetupCollaborationRoutes(r *mux.Router, webrtcManager *webrtcPkg.Manager, otManager *otPkg.Manager, assetManager *assetsPkg.Manager, authManager *authPkg.Manager) {
	// Initialize handlers
	webrtcHandler := webrtc.NewHandler(webrtcManager)
	otHandler := ot.NewHandler(otManager)
	assetHandler := assets.NewHandler(assetManager)
	
	// Initialize middleware
	authMiddleware := authPkg.NewMiddleware(authManager)
	
	// API base path
	api := r.PathPrefix("/api").Subrouter()
	
	// WebRTC routes (require authentication)
	webrtcRoutes := api.PathPrefix("/webrtc").Subrouter()
	webrtcRoutes.Use(authMiddleware.AuthRequired)
	webrtcRoutes.HandleFunc("/sessions", webrtcHandler.CreateRTCSession).Methods("POST")
	webrtcRoutes.HandleFunc("/sessions/{sessionId}/join", webrtcHandler.JoinRTCSession).Methods("POST")
	webrtcRoutes.HandleFunc("/sessions/{sessionId}/leave", webrtcHandler.LeaveRTCSession).Methods("POST")
	webrtcRoutes.HandleFunc("/sessions/{sessionId}/stats", webrtcHandler.GetRTCStats).Methods("GET")
	webrtcRoutes.HandleFunc("/sessions/{sessionId}/ws", webrtcHandler.HandleWebSocketConnection).Methods("GET")
	
	// Operational Transform routes (require authentication)
	otRoutes := api.PathPrefix("/ot").Subrouter()
	otRoutes.Use(authMiddleware.AuthRequired)
	otRoutes.HandleFunc("/documents", otHandler.CreateDocument).Methods("POST")
	otRoutes.HandleFunc("/documents/{documentId}", otHandler.GetDocument).Methods("GET")
	otRoutes.HandleFunc("/documents/{documentId}/join", otHandler.JoinDocument).Methods("POST")
	otRoutes.HandleFunc("/documents/{documentId}/leave", otHandler.LeaveDocument).Methods("POST")
	otRoutes.HandleFunc("/documents/{documentId}/operations", otHandler.ApplyOperation).Methods("POST")
	otRoutes.HandleFunc("/documents/{documentId}/operations/history", otHandler.GetOperationHistory).Methods("GET")
	otRoutes.HandleFunc("/sessions/{sessionId}/documents", otHandler.GetDocumentsBySession).Methods("GET")
	
	// Asset routes (require authentication)
	assetRoutes := api.PathPrefix("/assets").Subrouter()
	assetRoutes.Use(authMiddleware.AuthRequired)
	assetRoutes.HandleFunc("/upload", assetHandler.UploadAsset).Methods("POST")
	assetRoutes.HandleFunc("/{assetId}", assetHandler.GetAsset).Methods("GET")
	assetRoutes.HandleFunc("/{assetId}", assetHandler.DeleteAsset).Methods("DELETE")
	assetRoutes.HandleFunc("/{assetId}/content", assetHandler.GetAssetContent).Methods("GET")
	assetRoutes.HandleFunc("/{assetId}/usage", assetHandler.TrackAssetUsage).Methods("POST")
	assetRoutes.HandleFunc("/{assetId}/usage", assetHandler.GetAssetUsage).Methods("GET")
	assetRoutes.HandleFunc("/sessions/{sessionId}", assetHandler.GetAssetsBySession).Methods("GET")
}