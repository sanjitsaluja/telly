package telly

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"time"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	templates := GetBaseTemplates()
	templates = append(templates, "./views/channels.html")
	err := RenderTemplate(w, templates, "base", map[string]interface{}{
		"Title":    "Home",
		"Channels": DefaultSessionManager.channels,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveChannelDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	templates := GetBaseTemplates()
	templates = append(templates, "./views/channel.html")
	channel, found := DefaultSessionManager.FindChannel(id)
	if !found {
		http.Error(w, "Cannot find channel", http.StatusNotFound)
		return
	}
	err := RenderTemplate(w, templates, "base", map[string]interface{}{
		"Title":   "Channel",
		"Channel": channel,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveVideoMpg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-mpegURL")
	vars := mux.Vars(r)
	id := vars["id"]
	session, found := DefaultSessionManager.FindSession(id)
	if !found {
		http.Error(w, "Cannot find session", http.StatusNotFound)
		return
	}
	session.LastRead = time.Now()
	if !session.Ready {
		http.Error(w, "Session is not ready", 420)
		return
	}
	http.ServeFile(w, r, session.TranscodedFilePath())
}

func serveVideoTS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "video/MP2T")
	vars := mux.Vars(r)
	id := vars["id"]
	file := filepath.Join(TranscodingPath(), fmt.Sprintf("%s.ts", id))
	http.ServeFile(w, r, file)
}

func serveTuneChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, found := DefaultSessionManager.FindSession(id)
	if found {
		RenderJSON(w, 200, map[string]bool{"ready": false})
		return
	}
	if _, err := DefaultSessionManager.NewSession(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		RenderJSON(w, 200, map[string]bool{"ready": false})
	}
}

func serveChannelStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if session, found := DefaultSessionManager.FindSession(id); !found {
		http.Error(w, "Cannot find session", http.StatusNotFound)
	} else {
		RenderJSON(w, 200, map[string]bool{"ready": session.Ready})
	}

}
