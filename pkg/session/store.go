package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Store provides file-based persistence for session state.
// Each session is stored as a JSON file named <id>.json under the base directory.
type Store struct {
	baseDir string
}

// NewStore creates a Store that reads and writes session files under baseDir.
// The directory is created on first write if it does not exist.
func NewStore(baseDir string) *Store {
	return &Store{baseDir: baseDir}
}

// Load reads a session from disk by ID. If the file does not exist,
// a new session with the given ID is returned.
func (s *Store) Load(id string) (*Session, error) {
	path, err := s.path(id)
	if err != nil {
		return nil, fmt.Errorf("load session: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			now := time.Now().UTC()
			return &Session{
				ID:        id,
				StartedAt: now,
			}, nil
		}
		return nil, fmt.Errorf("load session %s: %w", id, err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("parse session %s: %w", id, err)
	}
	return &sess, nil
}

// Save persists a session to disk using atomic write (write to temp, then rename).
func (s *Store) Save(sess *Session) error {
	if err := os.MkdirAll(s.baseDir, 0o755); err != nil {
		return fmt.Errorf("create session dir: %w", err)
	}

	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal session %s: %w", sess.ID, err)
	}

	// Write to a temporary file in the same directory, then rename for atomicity.
	tmp, err := os.CreateTemp(s.baseDir, ".session-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}

	dest, err := s.path(sess.ID)
	if err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("save session: %w", err)
	}
	if err := os.Rename(tmpName, dest); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename session file: %w", err)
	}
	return nil
}

// List returns all sessions stored on disk, sorted by file modification time
// (most recent first is left to the caller).
func (s *Store) List() ([]*Session, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	var sessions []*Session
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")
		sess, err := s.Load(id)
		if err != nil {
			continue // skip corrupt files
		}
		sessions = append(sessions, sess)
	}
	return sessions, nil
}

// Cleanup removes sessions whose LastEventAt (or StartedAt if no events)
// is older than maxAge. Returns the number of sessions deleted.
func (s *Store) Cleanup(maxAge time.Duration) (int, error) {
	sessions, err := s.List()
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().UTC().Add(-maxAge)
	deleted := 0

	for _, sess := range sessions {
		ts := sess.LastEventAt
		if ts.IsZero() {
			ts = sess.StartedAt
		}
		if ts.Before(cutoff) {
			path, err := s.path(sess.ID)
			if err != nil {
				continue // skip invalid IDs
			}
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return deleted, fmt.Errorf("remove session %s: %w", sess.ID, err)
			}
			deleted++
		}
	}
	return deleted, nil
}

// sanitizeID validates that a session ID is safe to use as a filename.
// It rejects path separators, relative components, and empty strings.
func sanitizeID(id string) error {
	if id == "" {
		return fmt.Errorf("empty session ID")
	}
	if id == "." || id == ".." {
		return fmt.Errorf("invalid session ID: %q", id)
	}
	if strings.ContainsAny(id, `/\`) || filepath.Base(id) != id {
		return fmt.Errorf("session ID contains path separators: %q", id)
	}
	return nil
}

// path returns the file path for a session ID.
// Returns an error if the ID contains path traversal components.
func (s *Store) path(id string) (string, error) {
	if err := sanitizeID(id); err != nil {
		return "", err
	}
	return filepath.Join(s.baseDir, id+".json"), nil
}
