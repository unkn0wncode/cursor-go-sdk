package cursor

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

var testClient *Client
var testRepository string

// TestMain initializes a shared client using CURSOR_API_KEY and optional env config.
func TestMain(m *testing.M) {
	loadDotEnv(".env")
	cfg, _ := ConfigFromEnv()
	if cfg.APIKey == "" {
		fmt.Fprintln(os.Stderr, "CURSOR_API_KEY must be set for integration tests")
		os.Exit(2)
	}
	// Identify tests via a dedicated User-Agent.
	cfg.UserAgent = "cursor-go-sdk-tests"
	testClient = NewClientFromConfig(cfg)

	// Determine repository to use for agent tests. Allow override via env.
	testRepository = os.Getenv("CURSOR_TEST_REPOSITORY")
	if testRepository == "" {
		if r, err := discoverThisRepo(); err == nil {
			testRepository = r
		}
	}
	if testRepository == "" {
		fmt.Fprintln(os.Stderr, "Warning: could not determine repo; set CURSOR_TEST_REPOSITORY for agent tests")
	}

	os.Exit(m.Run())
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(line[len("export "):])
		}
		eq := strings.IndexRune(line, '=')
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if len(val) >= 2 {
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) || (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				val = val[1 : len(val)-1]
			}
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}

func discoverThisRepo() (string, error) {
	out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", err
	}
	return parseOwnerRepo(strings.TrimSpace(string(out))), nil
}

func parseOwnerRepo(remote string) string {
	if remote == "" {
		return ""
	}
	remote = strings.TrimSuffix(remote, ".git")
	if strings.HasPrefix(remote, "git@") {
		if i := strings.Index(remote, ":"); i >= 0 {
			remote = remote[i+1:]
		}
		remote = strings.TrimPrefix(remote, "github.com/")
	}
	if strings.HasPrefix(remote, "http://") || strings.HasPrefix(remote, "https://") || strings.HasPrefix(remote, "ssh://") {
		if u, err := url.Parse(remote); err == nil {
			remote = strings.TrimPrefix(u.Path, "/")
		}
	}
	parts := strings.Split(remote, "/")
	if len(parts) >= 2 {
		return fmt.Sprintf("https://github.com/%s/%s", parts[0], parts[1])
	}
	return ""
}

func TestMeEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	start := time.Now()

	resp, err := testClient.Me(ctx)
	require.NoError(t, err)
	t.Logf("Me response: %s", spew.Sdump(resp))
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.APIKeyName)
	require.False(t, resp.CreatedAt.IsZero())
	require.NotEmpty(t, resp.UserEmail)

	// Ensure completed within timeout.
	if deadline, ok := ctx.Deadline(); ok {
		require.LessOrEqual(t, time.Since(start), time.Until(deadline))
	}
}

func TestListRepositoriesEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	start := time.Now()

	resp, err := testClient.ListRepositories(ctx)
	require.NoError(t, err)
	t.Logf("ListRepositories response: %s", spew.Sdump(resp))
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Repositories)

	if deadline, ok := ctx.Deadline(); ok {
		require.LessOrEqual(t, time.Since(start), time.Until(deadline))
	}
}

func TestListModelsEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	start := time.Now()

	models, err := testClient.ListModels(ctx)
	require.NoError(t, err)
	t.Logf("ListModels response: %s", spew.Sdump(models))
	require.NotNil(t, models)
	require.NotEmpty(t, models.Models)

	if deadline, ok := ctx.Deadline(); ok {
		require.LessOrEqual(t, time.Since(start), time.Until(deadline))
	}
}

func TestAgentsStatusEndpoints(t *testing.T) {
	require.NotEmpty(t, testRepository, "repository owner/name required; set CURSOR_TEST_REPOSITORY if autodetect fails")

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	start := time.Now()

	// list agents
	agents, err := testClient.ListAgents(ctx, 10, nil)
	require.NoError(t, err)
	t.Logf("ListAgents response: %s", spew.Sdump(agents))
	require.NotNil(t, agents)

	if len(agents.Agents) > 0 {
		require.NotEmpty(t, agents.Agents)

		// get status of the first agent
		agent, err := testClient.GetAgent(ctx, agents.Agents[0].ID)
		require.NoError(t, err)
		t.Logf("GetAgent response: %s", spew.Sdump(agent))
		require.NotNil(t, agent)
		require.Equal(t, agents.Agents[0].ID, agent.ID)
		require.NotEmpty(t, agent.Status)

		if agent.Status != AgentStatusExpired {
			// get conversation history
			conversation, err := testClient.GetConversation(ctx, agents.Agents[0].ID)
			require.NoError(t, err)
			t.Logf("GetConversation response: %s", spew.Sdump(conversation))
			require.NotNil(t, conversation)
			require.NotEmpty(t, conversation.ID)
			require.NotEmpty(t, conversation.Messages)
		}
	}

	if deadline, ok := ctx.Deadline(); ok {
		require.LessOrEqual(t, time.Since(start), time.Until(deadline))
	}
}

func TestAgentsRunFlowEndpoints(t *testing.T) {
	require.NotEmpty(t, testRepository, "repository owner/name required; set CURSOR_TEST_REPOSITORY if autodetect fails")

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	start := time.Now()

	launchReq := LaunchRequest{
		Prompt: Prompt{Text: "You are a test agent for the Cursor Go SDK. Do not make any repository changes, do not open PRs or branches. Reply 'ok' and exit immediately."},
		Source: Source{Repository: testRepository},
		Target: &LaunchTarget{AutoCreatePR: false},
	}

	agent, err := testClient.LaunchAgent(ctx, launchReq)
	require.NoError(t, err)
	t.Logf("LaunchAgent response: %s", spew.Sdump(agent))
	require.NotNil(t, agent)
	require.NotEmpty(t, agent.ID)
	require.NotEmpty(t, agent.Status)

	// Ensure we clean up the agent regardless of later assertions.
	t.Cleanup(func() {
		_, _ = testClient.DeleteAgent(context.Background(), agent.ID)
	})

	followupID, err := testClient.AddFollowup(ctx, agent.ID, FollowupRequest{Prompt: Prompt{Text: "Confirm no changes were made and exit."}})
	require.NoError(t, err)
	t.Logf("AddFollowup response: %s", followupID)
	require.NotEmpty(t, followupID)

	got, err := testClient.GetAgent(ctx, agent.ID)
	require.NoError(t, err)
	t.Logf("GetAgent response: %s", spew.Sdump(got))
	require.NotNil(t, got)
	require.Equal(t, agent.ID, got.ID)

	listed, err := testClient.ListAgents(ctx, 10, nil)
	require.NoError(t, err)
	t.Logf("ListAgents response: %s", spew.Sdump(listed))
	require.NotNil(t, listed)
	require.NotEmpty(t, listed.Agents)

	time.Sleep(10 * time.Second)

	conv, err := testClient.GetConversation(ctx, agent.ID)
	require.NoError(t, err)
	t.Logf("GetConversation response: %s", spew.Sdump(conv))
	require.NotNil(t, conv)
	require.NotEmpty(t, conv.ID)

	deletedID, err := testClient.DeleteAgent(ctx, agent.ID)
	require.NoError(t, err)
	t.Logf("DeleteAgent response: %s", spew.Sdump(deletedID))
	require.NotEmpty(t, deletedID)

	if deadline, ok := ctx.Deadline(); ok {
		require.LessOrEqual(t, time.Since(start), time.Until(deadline))
	}
}
