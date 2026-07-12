package collect

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArtifactCategory represents the category of an artifact
type ArtifactCategory string

const (
	CategoryManual    ArtifactCategory = "manual"
	CategoryStandard  ArtifactCategory = "standard"
	CategorySpec     ArtifactCategory = "specification"
	CategoryDatasheet ArtifactCategory = "datasheet"
	CategoryDrawing  ArtifactCategory = "drawing"
	CategoryImage    ArtifactCategory = "image"
	CategoryVideo    ArtifactCategory = "video"
	CategoryArchive  ArtifactCategory = "archive"
	CategoryDocument ArtifactCategory = "document"
	CategoryUnknown  ArtifactCategory = "unknown"
)

// DiscoveredArtifact represents an artifact discovered in the repository
type DiscoveredArtifact struct {
	ID           string           `json:"id"`
	Path         string           `json:"path"`
	RelativePath string           `json:"relative_path"`
	Name         string           `json:"name"`
	Category     ArtifactCategory `json:"category"`
	Size         int64            `json:"size"`
	Hash         string           `json:"hash"`
	Modified     string           `json:"modified"`
	Extension    string           `json:"extension"`
	CollectionID string           `json:"collection_id"`
}

// CollectionResult contains the results of an artifact collection operation
type CollectionResult struct {
	SessionID      string              `json:"session_id"`
	StartedAt      string              `json:"started_at"`
	CompletedAt    string              `json:"completed_at"`
	Repository     string              `json:"repository"`
	ArtifactsFound []DiscoveredArtifact `json:"artifacts_found"`
	TotalSize      int64               `json:"total_size"`
	ProcessingTime float64             `json:"processing_time_seconds"`
}

// NewCollectionResult creates a new collection result
func NewCollectionResult(sessionID, repoPath string) *CollectionResult {
	return &CollectionResult{
		SessionID:      sessionID,
		StartedAt:      time.Now().Format(time.RFC3339),
		Repository:     repoPath,
		ArtifactsFound: []DiscoveredArtifact{},
	}
}

// SetCompleted marks the collection as complete
func (r *CollectionResult) SetCompleted() {
	r.CompletedAt = time.Now().Format(time.RFC3339)
}

// AddArtifact adds a discovered artifact to the result
func (r *CollectionResult) AddArtifact(artifact DiscoveredArtifact) {
	r.ArtifactsFound = append(r.ArtifactsFound, artifact)
	r.TotalSize += artifact.Size
}

// DetectCategory determines the artifact category from file extension and path
func DetectCategory(path string) ArtifactCategory {
	ext := strings.ToLower(filepath.Ext(path))
	lowerPath := strings.ToLower(path)

	switch {
	case strings.Contains(lowerPath, "manual") || strings.Contains(lowerPath, "guide") || strings.Contains(lowerPath, "handbook"):
		return CategoryManual
	case strings.Contains(lowerPath, "standard") || strings.Contains(lowerPath, "iec") || strings.Contains(lowerPath, "ieee") || strings.Contains(lowerPath, "iso") || strings.Contains(lowerPath, "nist"):
		return CategoryStandard
	case strings.Contains(lowerPath, "spec") || strings.Contains(lowerPath, "requirement"):
		return CategorySpec
	case strings.Contains(lowerPath, "datasheet") || strings.Contains(lowerPath, "data-sheet"):
		return CategoryDatasheet
	case strings.Contains(lowerPath, "drawing") || strings.Contains(lowerPath, "diagram") || strings.Contains(lowerPath, "schematic"):
		return CategoryDrawing
	case ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".svg" || ext == ".bmp":
		return CategoryImage
	case ext == ".mp4" || ext == ".avi" || ext == ".mov" || ext == ".mkv" || ext == ".webm":
		return CategoryVideo
	case ext == ".zip" || ext == ".tar" || ext == ".gz" || ext == ".rar" || ext == ".7z":
		return CategoryArchive
	case ext == ".pdf" || ext == ".md" || ext == ".txt" || ext == ".doc" || ext == ".docx" || ext == ".rst" || ext == ".adoc":
		return CategoryDocument
	default:
		return CategoryUnknown
	}
}

// CalculateHash computes SHA-256 hash of a file
func CalculateHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// GenerateArtifactID generates a unique artifact ID
func GenerateArtifactID(repoPath, artifactPath string, index int) string {
	relPath, _ := filepath.Rel(repoPath, artifactPath)
	hash := sha256.Sum256([]byte(relPath + time.Now().Format("20060102")))
	return "ART-" + hex.EncodeToString(hash[:4]) + "-" + string(rune('A'+index/26)) + string(rune('A'+index%26))
}
