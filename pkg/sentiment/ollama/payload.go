package ollama

import "time"

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Format   chatFormat    `json:"format"`
	Stream   bool          `json:"stream"`
}

type chatResponse struct {
	Model     string      `json:"model"`
	CreatedAt time.Time   `json:"created_at"`
	Message   chatMessage `json:"message"`
}

type chatMessage struct {
	Role    string `json:"role"` // "system", "user"
	Content string `json:"content"`
}

type chatFormat struct {
	Type       string                          `json:"type"` // "object"
	Properties map[string]chatFormatProperties `json:"properties"`
	Required   []string                        `json:"required"`
}

type chatFormatProperties struct {
	Type string `json:"type"` // "number"
}

type scoreResult struct {
	Score float64 `json:"score"`
}
