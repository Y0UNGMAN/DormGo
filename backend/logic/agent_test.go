package logic

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestPendingActionKeyUsesUserID(t *testing.T) {
	got := PendingActionKey(42)
	want := "dormgo:agent:pending_action:user:42"
	if got != want {
		t.Fatalf("PendingActionKey() = %q, want %q", got, want)
	}
}

func TestChatWithAgentPostsUserContextToPython(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		var req agentPythonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.UserID != 7 {
			t.Fatalf("UserID = %d, want 7", req.UserID)
		}
		if req.Message != "帮我找南校区拼车" {
			t.Fatalf("Message = %q", req.Message)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AgentChatResponse{
			Reply: "找到 1 个帖子",
			Type:  "candidates",
		})
	}))
	defer server.Close()

	viper.Set("agent.ServiceURL", server.URL)
	viper.Set("agent.InternalToken", "test-token")

	resp, err := ChatWithAgent(7, "帮我找南校区拼车")
	if err != nil {
		t.Fatalf("ChatWithAgent returned error: %v", err)
	}
	if resp.Reply != "找到 1 个帖子" {
		t.Fatalf("Reply = %q", resp.Reply)
	}
	if resp.Type != "candidates" {
		t.Fatalf("Type = %q", resp.Type)
	}
}
